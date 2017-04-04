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
	"strconv"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v2/tenants"
	"github.com/gophercloud/gophercloud/pagination"

	"github.com/intelsdi-x/snap-plugin-collector-nova/nova/v2"
)

// CachedNovas holds authenticated nova clients for given tenants
type CachedNovas map[string]v2.NovaV2

func authNoTenant(auth gophercloud.AuthOptions) gophercloud.AuthOptions {
	return gophercloud.AuthOptions{
		IdentityEndpoint: auth.IdentityEndpoint,
		UserID:           auth.UserID,
		Username:         auth.Username,
		Password:         auth.Password,
		DomainID:         auth.DomainID,
		DomainName:       auth.DomainName,
		AllowReauth:      auth.AllowReauth,
	}
}

// Get performs lazy initalization of client for given tenant using given auth options.
// If optional parameter givent its used as authenticated client
func (self CachedNovas) Get(auth gophercloud.AuthOptions, tenant string, providers ...*gophercloud.ProviderClient) (v2.NovaV2, error) {
	cached, ok := self[tenant]
	if ok {
		return cached, nil
	}

	var provider *gophercloud.ProviderClient

	if len(providers) > 0 {
		provider = providers[0]
	} else {
		newAuth := authNoTenant(auth)
		newAuth.TenantName = tenant
		provider2, err := openstack.AuthenticatedClient(newAuth)
		if err != nil {
			return v2.NovaV2{}, fmt.Errorf("authentication failed: (%v)", err)
		}
		provider = provider2
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{})

	if err != nil {
		return v2.NovaV2{}, fmt.Errorf("retrieving endpoint failed: (%v)", err)
	}

	self[tenant] = v2.NovaV2{Client: client}

	return self[tenant], nil

}

// collector separates plugin's interface from openstack api.
type collector struct {
	NovaCache CachedNovas
	Keystone  *gophercloud.ServiceClient
	Auth      gophercloud.AuthOptions

	config Config
}

// newCollector creates and initializes instance of collector. Error is returned
// when either authentication or endpoint retrieval failed.
func newCollector(config Config) (collectorInterface, error) {
	auth := gophercloud.AuthOptions{
		IdentityEndpoint: config.Url,
		Username:         config.User,
		Password:         config.Pass,
		TenantName:       config.Tenant,
		DomainID:         config.DomaninID,
		DomainName:       config.DomainName,
		AllowReauth:      true,
	}

	self := &collector{
		NovaCache: CachedNovas{},
		Auth:      auth,
		config:    config,
	}

	provider, err := openstack.AuthenticatedClient(auth)

	if err != nil {
		return nil, fmt.Errorf("cannot authenticate: (%v)", err)
	}

	_, err = self.NovaCache.Get(auth, config.Tenant, provider)
	if err != nil {
		return nil, err
	}

	client, err := openstack.NewIdentityV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		return nil, fmt.Errorf("retrieving identity service failed: (%v)", err)
	}
	self.Keystone = client

	return self, nil
}

// GetTenants returns map of tenant names -> tenant ids
func (self *collector) GetTenants() (map[string]string, error) {
	opts := &tenants.ListOpts{}

	pager := tenants.List(self.Keystone, opts)

	result := map[string]string{}

	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		tenantList, err := tenants.ExtractTenants(page)

		if err != nil {
			return false, err
		}

		for _, t := range tenantList {
			result[t.Name] = t.ID
		}

		return true, nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (self *collector) GetLimitsNames() []string {
	return []string{
		"max_server_meta",
		"max_personality",
		"total_server_groups_used",
		"max_image_meta",
		"max_personality_size",
		"max_server_groups",
		"max_security_group_rules",
		"max_total_keypairs",
		"total_cores_used",
		"total_ram_used",
		"total_instances_used",
		"max_security_groups",
		"total_floating_ips_used",
		"max_total_cores",
		"total_security_groups_used",
		"max_total_floating_ips",
		"max_total_instances",
		"max_total_ram_size",
		"max_server_group_members",
	}
}

func (self *collector) GetQuotasNames() []string {
	return []string{
		"cores",
		"fixed_ips",
		"floating_ips",
		"instances",
		"key_pairs",
		"ram",
		"security_groups",
	}
}

func (self *collector) GetHStatsNames() []string {
	return []string{
		"current_workload",
		"disk_available_least",
		"free_disk_gb",
		"free_ram_mb",
		"local_gb",
		"local_gb_used",
		"memory_mb",
		"memory_mb_used",
		"running_vms",
		"vcpus",
		"vcpus_used",
	}
}

func (self *collector) GetClusterConfigNames() []string {
	return []string{
		"allocation_ratio_ram",
		"reserved_node_ram_mb",
		"allocation_ratio_cores",
		"reserved_node_cores",
	}
}

func toKB(gb int) int {
	if gb < 0 {
		return gb
	} else {
		return gb * 1024 * 1024
	}
}

func (self *collector) GetLimits(tenant string) (map[string]interface{}, error) {
	client, err := self.NovaCache.Get(self.Auth, tenant)
	if err != nil {
		return nil, err
	}

	values, err := client.GetLimits()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"max_server_meta":            values.MaxServerMeta,
		"max_personality":            values.MaxPersonality,
		"total_server_groups_used":   values.TotalServerGroupsUsed,
		"max_image_meta":             values.MaxImageMeta,
		"max_personality_size":       values.MaxPersonalitySize,
		"max_server_groups":          values.MaxServerGroups,
		"max_security_group_rules":   values.MaxSecurityGroupRules,
		"max_total_keypairs":         values.MaxTotalKeypairs,
		"total_cores_used":           values.TotalCoresUsed,
		"total_ram_used":             toKB(values.TotalRAMUsed),
		"total_instances_used":       values.TotalInstancesUsed,
		"max_security_groups":        values.MaxSecurityGroups,
		"total_floating_ips_used":    values.TotalFloatingIPsUsed,
		"max_total_cores":            values.MaxTotalCores,
		"total_security_groups_used": values.TotalSecurityGroupsUsed,
		"max_total_floating_ips":     values.MaxTotalFloatingIPs,
		"max_total_instances":        values.MaxTotalInstances,
		"max_total_ram_size":         toKB(values.MaxTotalRamSize),
		"max_server_group_members":   values.MaxServerGroupMembers,
	}, nil

}

// GetQuotas reads quotas. Besides tenant name it requires class id parameter
// which is usuallt tenant id.
func (self *collector) GetQuotas(tenant, id string) (map[string]interface{}, error) {
	client, err := self.NovaCache.Get(self.Auth, tenant)
	if err != nil {
		return nil, err
	}

	values, err := client.GetQuotas(id)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"cores":           values.Cores,
		"fixed_ips":       values.FixedIps,
		"floating_ips":    values.FloatingIps,
		"instances":       values.Instances,
		"key_pairs":       values.KeyPairs,
		"ram":             toKB(values.Ram),
		"security_groups": values.SecurityGroups,
	}, nil

}

// GetHypervisors reads hypervisors statistics (for all at once). Returns
// map of hypervisors' ids -> (stat name -> stat value)
func (self *collector) GetHypervisors() (map[string]map[string]interface{}, error) {

	client, err := self.NovaCache.Get(self.Auth, self.Auth.TenantName)
	if err != nil {
		return nil, err
	}

	list, err := client.GetHypervisorStatistics()
	if err != nil {
		return nil, err
	}

	result := map[string]map[string]interface{}{}

	for _, values := range list {
		result[strconv.Itoa(values.Id)] = map[string]interface{}{
			"current_workload":                 values.CurrentWorkload,
			"disk_available_least":             values.DiskAvailableLeast,
			"free_disk_gb":                     values.FreeDiskGB,
			"free_ram_mb":                      values.FreeRamMB,
			"hypervisor_version":               values.HypervisorVersion,
			"local_gb":                         values.LocalGB,
			"local_gb_used":                    values.LocalGBUsed,
			"memory_mb":                        values.MemoryMB,
			"memory_mb_used":                   values.MemoryMBUsed,
			"running_vms":                      values.RunningVMs,
			"vcpus":                            values.VCPUs,
			"vcpus_used":                       values.VCPUsUsed,
			"memory_mb_overcommit":             float64(values.MemoryMB) * self.config.RatioRam,
			"memory_mb_overcommit_withreserve": float64(values.MemoryMB)*self.config.RatioRam - self.config.ReservedNRam,
			"vcpus_overcommit":                 float64(values.VCPUs) * self.config.RatioCores,
			"vcpus_overcommit_withreserve":     float64(values.VCPUs)*self.config.RatioCores - self.config.ReservedNCores,
		}
	}

	return result, nil
}

func (self *collector) GetClusterConfig() map[string]interface{} {
	return map[string]interface{}{
		"allocation_ratio_ram":   self.config.RatioRam,
		"reserved_node_ram_mb":   self.config.ReservedNRam,
		"allocation_ratio_cores": self.config.RatioCores,
		"reserved_node_cores":    self.config.ReservedNCores,
	}
}
