// +build small

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
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/core/ctypes"

	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/cdata"
)

var smthErr = fmt.Errorf("smth")

func testingConfig() (cfg1 plugin.ConfigType, cfg2 *cdata.ConfigDataNode) {
	cfg1 = plugin.NewPluginConfigType()
	cfg2 = cdata.NewNode()

	cfg1.AddItem("openstack_user", ctypes.ConfigValueStr{Value: "x"})
	cfg1.AddItem("openstack_pass", ctypes.ConfigValueStr{Value: "x"})
	cfg1.AddItem("openstack_tenant", ctypes.ConfigValueStr{Value: "asdf"})
	cfg1.AddItem("openstack_auth_url", ctypes.ConfigValueStr{Value: "x"})

	cfg2.AddItem("openstack_user", ctypes.ConfigValueStr{Value: "x"})
	cfg2.AddItem("openstack_pass", ctypes.ConfigValueStr{Value: "x"})
	cfg2.AddItem("openstack_tenant", ctypes.ConfigValueStr{Value: "asdf"})
	cfg2.AddItem("openstack_auth_url", ctypes.ConfigValueStr{Value: "x"})

	cfg1.AddItem("allocation_ratio_cores", ctypes.ConfigValueFloat{Value: 3})
	cfg1.AddItem("allocation_ratio_ram", ctypes.ConfigValueFloat{Value: 4})
	cfg1.AddItem("reserved_node_cores", ctypes.ConfigValueFloat{Value: 5})
	cfg1.AddItem("reserved_node_ram_mb", ctypes.ConfigValueFloat{Value: 6})

	cfg2.AddItem("allocation_ratio_cores", ctypes.ConfigValueFloat{Value: 3})
	cfg2.AddItem("allocation_ratio_ram", ctypes.ConfigValueFloat{Value: 4})
	cfg2.AddItem("reserved_node_cores", ctypes.ConfigValueFloat{Value: 5})
	cfg2.AddItem("reserved_node_ram_mb", ctypes.ConfigValueFloat{Value: 6})

	return cfg1, cfg2
}

type collectorMock struct {
	mock.Mock
}

func (self *collectorMock) GetTenants() (map[string]string, error) {
	args := self.Called()
	var r0 map[string]string = nil
	if args.Get(0) != nil {
		r0 = args.Get(0).(map[string]string)
	}
	return r0, args.Error(1)
}

func (self *collectorMock) GetLimitsNames() []string {
	args := self.Called()
	var r0 []string = nil
	if args.Get(0) != nil {
		r0 = args.Get(0).([]string)
	}
	return r0
}

func (self *collectorMock) GetQuotasNames() []string {
	args := self.Called()
	var r0 []string = nil
	if args.Get(0) != nil {
		r0 = args.Get(0).([]string)
	}
	return r0
}
func (self *collectorMock) GetHStatsNames() []string {
	args := self.Called()
	var r0 []string = nil
	if args.Get(0) != nil {
		r0 = args.Get(0).([]string)
	}
	return r0
}
func (self *collectorMock) GetClusterConfigNames() []string {
	args := self.Called()
	var r0 []string = nil
	if args.Get(0) != nil {
		r0 = args.Get(0).([]string)
	}
	return r0
}

func (self *collectorMock) GetLimits(tenant string) (map[string]interface{}, error) {
	args := self.Called(tenant)
	var r0 map[string]interface{} = nil
	if args.Get(0) != nil {
		r0 = args.Get(0).(map[string]interface{})
	}
	return r0, args.Error(1)
}

func (self *collectorMock) GetQuotas(tenant, id string) (map[string]interface{}, error) {
	args := self.Called(tenant, id)
	var r0 map[string]interface{} = nil
	if args.Get(0) != nil {
		r0 = args.Get(0).(map[string]interface{})
	}
	return r0, args.Error(1)
}

func (self *collectorMock) GetHypervisors() (map[string]map[string]interface{}, error) {
	args := self.Called()
	var r0 map[string]map[string]interface{} = nil
	if args.Get(0) != nil {
		r0 = args.Get(0).(map[string]map[string]interface{})
	}
	return r0, args.Error(1)
}

func (self *collectorMock) GetClusterConfig() map[string]interface{} {
	args := self.Called()
	var r0 map[string]interface{} = nil
	if args.Get(0) != nil {
		r0 = args.Get(0).(map[string]interface{})
	}
	return r0
}

func TestGetMetricTypes(t *testing.T) {

	Convey("GetMetricTypes", t, func() {

		cfg1, _ := testingConfig()
		m := &collectorMock{}
		orgCollector := testingCollector

		testingCollector = func(config Config) (collectorInterface, error) {
			return m, nil
		}

		sut := NovaPlugin{initializedMutex: new(sync.Mutex)}

		Reset(func() {
			testingCollector = orgCollector
		})

		m.On("GetLimitsNames").Return([]string{
			"limit1",
			"limit2"}, nil)

		m.On("GetQuotasNames").Return([]string{
			"q1",
			"q2"}, nil)

		m.On("GetClusterConfigNames").Return([]string{
			"c1",
			"c2"}, nil)

		Convey("when list of tenants is returned", func() {
			m.On("GetTenants").Return(map[string]string{"asdf": "t_1", "efg": "t_2"}, nil).Once()
			Convey("returns namespace for each tenant", func() {

				m.On("GetHypervisors").Return(map[string]map[string]interface{}{
					"h1": map[string]interface{}{"s100": 5, "s200": 6},
					"h2": map[string]interface{}{"s105": 2, "s205": 3},
				}, nil).Once()

				mts, err := sut.GetMetricTypes(cfg1)
				m.AssertExpectations(t)

				Convey("with no error", func() {
					So(err, ShouldBeNil)
				})

				tenants := map[string]bool{}
				for _, mt := range mts {
					tenants[mt.Namespace().Strings()[4]] = true
				}

				So(tenants["asdf"], ShouldBeTrue)
				So(tenants["efg"], ShouldBeTrue)

				Convey("with hypervisor metrics", func() {

					hvs := map[string]bool{}

					for _, mt := range mts {
						hvs[mt.Namespace().Strings()[3]+":"+mt.Namespace().Strings()[4]] = true
					}

					So(hvs["hypervisor:h1"], ShouldBeTrue)
					So(hvs["hypervisor:h2"], ShouldBeTrue)

				})
				Convey("with limits", func() {

					lms := map[string]bool{}

					for _, mt := range mts {
						lms[strings.Join(mt.Namespace().Strings()[4:], ":")] = true
					}

					So(lms["asdf:limits:limit1"], ShouldBeTrue)
					So(lms["asdf:limits:limit2"], ShouldBeTrue)
					So(lms["efg:limits:limit1"], ShouldBeTrue)
					So(lms["efg:limits:limit2"], ShouldBeTrue)
					So(lms["asdf:limits:limit1"], ShouldBeTrue)
					So(lms["asdf:limits:limit2"], ShouldBeTrue)
					So(lms["efg:limits:limit1"], ShouldBeTrue)
					So(lms["efg:limits:limit2"], ShouldBeTrue)

				})
				Convey("with quotas", func() {

					qts := map[string]bool{}

					for _, mt := range mts {
						qts[strings.Join(mt.Namespace().Strings()[4:], ":")] = true
					}

					So(qts["asdf:quotas:q1"], ShouldBeTrue)
					So(qts["asdf:quotas:q2"], ShouldBeTrue)
					So(qts["efg:quotas:q1"], ShouldBeTrue)
					So(qts["efg:quotas:q2"], ShouldBeTrue)
					So(qts["asdf:quotas:q1"], ShouldBeTrue)
					So(qts["asdf:quotas:q2"], ShouldBeTrue)
					So(qts["efg:quotas:q1"], ShouldBeTrue)
					So(qts["efg:quotas:q2"], ShouldBeTrue)

				})

				Convey("with config values", func() {

					cvs := map[string]bool{}

					for _, mt := range mts {
						cvs[strings.Join(mt.Namespace().Strings()[3:], ":")] = true
					}

					So(cvs["cluster:config:c1"], ShouldBeTrue)
					So(cvs["cluster:config:c2"], ShouldBeTrue)

				})
			})
		})
		Convey("returns error when failed", func() {

			Convey("tenants", func() {
				m.On("GetTenants").Return(nil, smthErr)
				_, err := sut.GetMetricTypes(cfg1)
				So(err, ShouldNotBeNil)
			})

			Convey("hypervisors", func() {

				m.On("GetTenants").Return(map[string]string{"asdf": "t_1", "efg": "t_2"}, nil)

				m.On("GetHypervisors").Return(nil, smthErr)
				_, err := sut.GetMetricTypes(cfg1)
				So(err, ShouldNotBeNil)
			})

		})

	})

}

func makeMts(cfg *cdata.ConfigDataNode, ns ...string) []plugin.MetricType {
	res := make([]plugin.MetricType, len(ns))
	for i, v := range ns {
		res[i].Config_ = cfg
		rns := []string{"intel", "openstack", "nova"}
		res[i].Namespace_ = core.NewNamespace(append(rns, strings.Split(v, "/")...)...)
	}

	return res
}

func TestCollectMetrics(t *testing.T) {

	Convey("CollectMetrics", t, func() {

		_, cfg2 := testingConfig()
		m := &collectorMock{}
		orgCollector := testingCollector

		testingCollector = func(config Config) (collectorInterface, error) {
			return m, nil
		}

		sut := NovaPlugin{initializedMutex: new(sync.Mutex)}

		Reset(func() {
			testingCollector = orgCollector
		})

		m.On("GetHypervisors").Return(map[string]map[string]interface{}{
			"h1": map[string]interface{}{"s100": 5, "s200": 6},
			"h2": map[string]interface{}{"s105": 2, "s205": 3},
		}, nil)
		m.On("GetLimits", "efg").Return(map[string]interface{}{
			"limit1": 3,
			"limit2": 5,
		}, nil)
		m.On("GetLimits", "asdf").Return(map[string]interface{}{
			"limit1": 13,
			"limit2": 15,
		}, nil)
		m.On("GetQuotas", "asdf", "t_1").Return(map[string]interface{}{
			"q1": 33,
			"q2": 35,
		}, nil)
		m.On("GetQuotas", "efg", "t_2").Return(map[string]interface{}{
			"q1": 113,
			"q2": 115,
		}, nil)
		m.On("GetClusterConfig").Return(map[string]interface{}{
			"c1": 444,
			"c2": 555,
		}, nil)

		m.On("GetTenants").Return(map[string]string{"asdf": "t_1", "efg": "t_2"}, nil)

		Convey("refreshes list of tentants when per tenant quotas are requested", func() {

			sut.CollectMetrics(makeMts(cfg2, "tenant/efg/quotas/q2"))

			m.AssertCalled(t, "GetTenants")
		})

		Convey("doesn't refresh list of tenants when no per tenant stat is requested", func() {

			sut.CollectMetrics(makeMts(cfg2, "tenant/efg/limits/limit2", "hypervisor/h2/s105"))

			m.AssertNotCalled(t, "GetTenants")

		})

		Convey("performs appropriate calls for limits", func() {

			sut.CollectMetrics(makeMts(cfg2, "tenant/efg/limits/limit2", "tenant/efg/limits/limit1"))
			m.AssertNumberOfCalls(t, "GetLimits", 1)

		})

		Convey("performs appropriate calls for quotas", func() {

			sut.CollectMetrics(makeMts(cfg2, "tenant/efg/quotas/q2", "tenant/efg/quotas/q1"))
			m.AssertNumberOfCalls(t, "GetQuotas", 1)

		})

		Convey("requests hypervisors stats", func() {

			Convey("once when we want these stats", func() {

				sut.CollectMetrics(makeMts(cfg2, "tenant/efg/limits/limit2",
					"hypervisor/h1/s100", "hypervisor/h2/s105"))
				m.AssertNumberOfCalls(t, "GetHypervisors", 1)

			})

			Convey("not when we don't", func() {

				sut.CollectMetrics(makeMts(cfg2, "tenant/efg/quotas/q2"))
				m.AssertNumberOfCalls(t, "GetHypervisors", 0)

			})
		})

		Convey("requests config derived stats", func() {

			Convey("once when we want these stats", func() {

				sut.CollectMetrics(makeMts(cfg2, "tenant/efg/limits/limit2",
					"cluster/config/c1", "cluster/config/c2"))
				m.AssertNumberOfCalls(t, "GetClusterConfig", 1)

			})

			Convey("not when we don't", func() {

				sut.CollectMetrics(makeMts(cfg2, "tenant/efg/quotas/q2"))
				m.AssertNumberOfCalls(t, "GetClusterConfig", 0)

			})

		})

		Convey("correct values are present in results", func() {
			mts, _ := sut.CollectMetrics(makeMts(cfg2,
				"tenant/asdf/limits/limit1",
				"tenant/asdf/limits/limit2",
				"tenant/asdf/quotas/q1",
				"tenant/asdf/quotas/q2",
				"tenant/efg/limits/limit1",
				"tenant/efg/limits/limit2",
				"tenant/efg/quotas/q1",
				"tenant/efg/quotas/q2",
				"hypervisor/h1/s100",
				"hypervisor/h1/s200",
				"hypervisor/h2/s105",
				"hypervisor/h2/s205",
				"cluster/config/c1",
				"cluster/config/c2"))

			dut := map[string]interface{}{}
			for _, v := range mts {
				dut[v.Namespace().String()] = v.Data()
			}

			So(dut["/intel/openstack/nova/tenant/asdf/limits/limit1"], ShouldEqual, 13)
			So(dut["/intel/openstack/nova/tenant/asdf/limits/limit2"], ShouldEqual, 15)
			So(dut["/intel/openstack/nova/tenant/asdf/quotas/q1"], ShouldEqual, 33)
			So(dut["/intel/openstack/nova/tenant/asdf/quotas/q2"], ShouldEqual, 35)
			So(dut["/intel/openstack/nova/tenant/efg/limits/limit1"], ShouldEqual, 3)
			So(dut["/intel/openstack/nova/tenant/efg/limits/limit2"], ShouldEqual, 5)
			So(dut["/intel/openstack/nova/tenant/efg/quotas/q1"], ShouldEqual, 113)
			So(dut["/intel/openstack/nova/tenant/efg/quotas/q2"], ShouldEqual, 115)
			So(dut["/intel/openstack/nova/hypervisor/h1/s100"], ShouldEqual, 5)
			So(dut["/intel/openstack/nova/hypervisor/h1/s200"], ShouldEqual, 6)
			So(dut["/intel/openstack/nova/hypervisor/h2/s105"], ShouldEqual, 2)
			So(dut["/intel/openstack/nova/hypervisor/h2/s205"], ShouldEqual, 3)
			So(dut["/intel/openstack/nova/cluster/config/c1"], ShouldEqual, 444)
			So(dut["/intel/openstack/nova/cluster/config/c2"], ShouldEqual, 555)

		})

		Convey("returns error", func() {
			*m = collectorMock{}

			Convey("when refreshing tenant list failed", func() {

				m.On("GetTenants").Return(nil, smthErr)

				_, dut_err := sut.CollectMetrics(makeMts(cfg2, "tenant/asdf/quotas/foo"))

				So(dut_err, ShouldNotBeNil)

			})

			Convey("when requesting limits", func() {

				m.On("GetLimits", "asdf").Return(nil, smthErr)

				_, dut_err := sut.CollectMetrics(makeMts(cfg2, "tenant/asdf/limits/foo"))

				So(dut_err, ShouldNotBeNil)

			})

			Convey("when requesting quotas", func() {

				m.On("GetTenants").Return(map[string]string{"asdf": "t_1", "efg": "t_2"}, nil)
				m.On("GetQuotas", "asdf", "t_1").Return(nil, smthErr)

				_, dut_err := sut.CollectMetrics(makeMts(cfg2, "tenant/asdf/quotas/foo"))

				So(dut_err, ShouldNotBeNil)

			})

			Convey("when requesting hypervisors stats", func() {

				m.On("GetHypervisors").Return(nil, smthErr)
				_, dut_err := sut.CollectMetrics(makeMts(cfg2, "hypervisor/x/y"))

				So(dut_err, ShouldNotBeNil)

			})

		})

	})

}
