/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2016 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package novaplugin

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"

	"reflect"

	"github.com/intelsdi-x/snap-plugin-utilities/config"
	"github.com/intelsdi-x/snap/core"
)

const (
	// Name of plugin
	Name = "nova-compute"
	// Version of plugin
	Version = 2
	// Type of plugin
	Type = plugin.CollectorPluginType
)

type Config struct {
	User   string `c:"openstack_user"`
	Pass   string `c:"openstack_pass"`
	Tenant string `c:"openstack_tenant"`
	Url    string `c:"openstack_auth_url"`

	RatioCores     float64 `c:"allocation_ratio_cores"`
	RatioRam       float64 `c:"allocation_ratio_ram"`
	ReservedNCores float64 `c:"reserved_node_cores"`
	ReservedNRam   float64 `c:"reserved_node_ram_mb"`
}

// ReadConfig deserializes plugin's configuration from metric or global config
// given in cfg. out shoud be pointer to structure. If field of structure has no
// tag it's name is used as config key, if it has named tag "c" with values
// delimited by commas, first value is used as config key. If second value is
// "weak" and field of structure is string, string representation of read value
// is written. "weak" is optional for string fields and currently forbidden for
// other types.
// Returns nil if operation succeeded or relevant error.
func ReadConfig(cfg interface{}, out interface{}) error {
	outStructValue := reflect.ValueOf(out).Elem()
	outStructType := outStructValue.Type()

	for i := 0; i < outStructType.NumField(); i++ {
		field := outStructType.Field(i)
		tags := strings.Split(field.Tag.Get("c"), ",")

		tag := strings.TrimSpace(tags[0])

		if tag == "" {
			tag = field.Name
		}

		value, err := config.GetConfigItem(cfg, tag)

		if err != nil {
			return err
		}

		fieldValue := outStructValue.Field(i)

		if len(tags) > 1 && strings.TrimSpace(tags[1]) == "weak" {
			if field.Type.Kind() != reflect.String {
				return fmt.Errorf("field %s has to be string: %s found",
					field.Name, field.Type)
			}
			fieldValue.SetString(fmt.Sprint(value))

		} else {
			if !reflect.TypeOf(value).ConvertibleTo(field.Type) {
				return fmt.Errorf("cannot assing config attribute %v to field %v: %v is not convertible to %v",
					tag, field.Name, reflect.TypeOf(value), field.Type,
				)
			}
			converted := reflect.ValueOf(value).Convert(field.Type)
			fieldValue.Set(converted)
		}
	}

	return nil
}

type NovaPlugin struct {
	initialized      bool
	initializedMutex *sync.Mutex

	config    Config
	collector collectorInterface
}

// CollectMetrics returns filled mts table with metric values. Limits and quotas
// are colllected once per tenant. All hypervisor related statistics are collected
// in one call. This method also performs lazy initalization of plugin. Error
// is returned if initalization or any of required call failed.
func (self *NovaPlugin) CollectMetrics(mts []plugin.MetricType) ([]plugin.MetricType, error) {
	if len(mts) > 0 {
		err := self.init(mts[0])

		if err != nil {
			return nil, err
		}

	} else {
		return mts, nil
	}

	t := time.Now()

	limitsFor := map[string]bool{}
	quotasFor := map[string]bool{}
	hypervisors := false
	cluster := false

	results := make([]plugin.MetricType, len(mts))
	for _, mt := range mts {
		id, group, subgroup, _ := parseName(mt.Namespace().Strings())
		if group == GROUP_CLUSTER {
			cluster = true
			continue
		}
		if group == GROUP_HYPERVISOR {
			hypervisors = true
		} else {
			if subgroup == SUBGROUP_LIMITS {
				limitsFor[id] = true
			} else {
				quotasFor[id] = true
			}
		}
	}

	cachedLimits := map[string]map[string]interface{}{}
	for tenant, _ := range limitsFor {
		limits, err := self.collector.GetLimits(tenant)
		if err != nil {
			return nil, fmt.Errorf("cannot read limits for %v: (%v)", tenant, err)
		}
		cachedLimits[tenant] = limits
	}

	cachedQuotas := map[string]map[string]interface{}{}
	var tenantIds map[string]string = nil
	for tenant, _ := range quotasFor {
		if tenantIds == nil {
			tenantIds2, err := self.collector.GetTenants()
			if err != nil {
				return nil, fmt.Errorf("cannot get tenants list: (%v)", err)
			}
			tenantIds = tenantIds2
		}

		quotas, err := self.collector.GetQuotas(tenant, tenantIds[tenant])
		if err != nil {
			return nil, fmt.Errorf("cannot read quotas for %v: (%v)", tenant, err)
		}
		cachedQuotas[tenant] = quotas
	}

	cachedHypervisor := map[string]map[string]interface{}{}
	if hypervisors {
		hStats, err := self.collector.GetHypervisors()
		if err != nil {
			return nil, fmt.Errorf("cannot read hypervisors: (%v)", err)
		}
		cachedHypervisor = hStats
	}

	cachedClusterConfig := map[string]interface{}{}
	if cluster {
		cachedClusterConfig = self.collector.GetClusterConfig()
	}

	for i, mt := range mts {
		id, group, subgroup, metric := parseName(mt.Namespace().Strings())
		mt := plugin.MetricType{
			Namespace_: mt.Namespace(),
			Timestamp_: t,
		}
		if group == GROUP_CLUSTER && id == ID_CONFIG {
			mt.Data_ = cachedClusterConfig[metric]

		} else {
			if group == GROUP_HYPERVISOR {
				mt.Data_ = cachedHypervisor[id][metric]
			} else {
				if subgroup == SUBGROUP_LIMITS {
					mt.Data_ = cachedLimits[id][metric]
				} else {
					mt.Data_ = cachedQuotas[id][metric]
				}
			}
		}
		results[i] = mt
	}

	return results, nil
}

// GetMetricTypes returns list of possible namespaces. Namespaces involving
// limits or quotas are constructed per tenant. Namespaces for hypervisors are
// constructed  in single api call. This  method also performs lazy
// initalization of plugin. Returns error if initalization or any request failed.
func (self *NovaPlugin) GetMetricTypes(cfg plugin.ConfigType) ([]plugin.MetricType, error) {
	err := self.init(cfg)

	if err != nil {
		return nil, err
	}

	names := [][]string{}

	tenants, err := self.collector.GetTenants()

	if err != nil {
		return nil, fmt.Errorf("cannot get tenants list: (%v)", err)
	}

	limitNames := self.collector.GetLimitsNames()
	quotaNames := self.collector.GetQuotasNames()
	configNames := self.collector.GetClusterConfigNames()

	appendNames(&names, ID_CONFIG, GROUP_CLUSTER, "", configNames)

	for tName, _ := range tenants {
		appendNames(&names, tName, GROUP_TENANT, SUBGROUP_LIMITS, limitNames)
		appendNames(&names, tName, GROUP_TENANT, SUBGROUP_QUOTAS, quotaNames)
	}

	hypervisors, err := self.collector.GetHypervisors()

	if err != nil {
		return nil, fmt.Errorf("cannot get hypervisors list: (%v)", err)
	}

	for hName, hVals := range hypervisors {
		for key, _ := range hVals {
			names = append(names, makeName(hName, GROUP_HYPERVISOR, "", key))
		}
	}

	mts := make([]plugin.MetricType, len(names))
	for i, v := range names {
		mts[i].Namespace_ = core.NewNamespace(v...)
	}

	return mts, nil
}

func (self *NovaPlugin) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	c := cpolicy.New()
	return c, nil
}

func (self *NovaPlugin) init(cfg interface{}) error {
	self.initializedMutex.Lock()
	defer self.initializedMutex.Unlock()

	if self.initialized {
		return nil
	}

	err := ReadConfig(cfg, &self.config)

	if err != nil {
		return err
	}

	// testingCollector is a variable that might either be newCollector
	// or fake factory for mocking
	self.collector, err = testingCollector(self.config)

	if err != nil {
		return fmt.Errorf("plugin initalization failed : [%v]", err)
	}

	self.initialized = true

	return nil

}

func New() *NovaPlugin {
	self := new(NovaPlugin)
	self.initializedMutex = new(sync.Mutex)
	return self
}

func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(Name, Version, Type, []string{plugin.SnapGOBContentType}, []string{plugin.SnapGOBContentType})
}

const (
	SUBGROUP_LIMITS = "limits"
	SUBGROUP_QUOTAS = "quotas"

	GROUP_HYPERVISOR = "hypervisor"
	GROUP_TENANT     = "tenant"
	GROUP_CLUSTER    = "cluster"

	ID_CONFIG = "config"
)

var namespacePrefix = []string{"intel", "openstack", "nova"}

func makeName(id, group, subgroup, metric string) []string {
	result := []string{}
	result = append(result, namespacePrefix...)
	result = append(result, group, id)
	if group == GROUP_TENANT {
		result = append(result, subgroup)
	}
	result = append(result, strings.Split(metric, "/")...)

	return result
}

func appendNames(out *[][]string, id, group, subgroup string, metrics []string) {
	for _, m := range metrics {
		*out = append(*out, makeName(id, group, subgroup, m))
	}
}

func parseName(ns []string) (id, group, subgroup, metric string) {
	i := len(namespacePrefix)
	group = ns[i]
	id = ns[i+1]
	if group == GROUP_TENANT {
		subgroup = ns[i+2]
		metric = strings.Join(ns[i+3:], "/")
	} else {
		metric = strings.Join(ns[i+2:], "/")
	}
	return
}

type collectorInterface interface {
	GetTenants() (map[string]string, error)

	GetLimitsNames() []string
	GetQuotasNames() []string
	GetClusterConfigNames() []string
	GetHStatsNames() []string

	GetLimits(tenant string) (map[string]interface{}, error)
	GetQuotas(tenant, id string) (map[string]interface{}, error)
	GetHypervisors() (map[string]map[string]interface{}, error)
	GetClusterConfig() map[string]interface{}
}

//for mocking
var testingCollector func(Config) (collectorInterface, error) = newCollector
