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
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

func testingRead(val interface{}, path string) string {
	for _, el := range strings.Split(path, ".") {
		m, ok := val.(map[string]interface{})
		if !ok {
			return ""
		}
		val, ok = m[el]
		if !ok {
			return ""
		}
	}
	return fmt.Sprint(val)
}

type testingTenant struct {
	Name, Id string
	Modifier string
}

type testingServices struct {
	r    *mux.Router
	serv *httptest.Server

	authorized map[string]*testingTenant
	ctoken     int

	tenants                  []testingTenant
	tenantById, tenantByName map[string]*testingTenant
}

func (self *testingServices) keystoneTokens(w http.ResponseWriter, r *http.Request) {
	var val interface{}
	val = nil
	json.NewDecoder(r.Body).Decode(&val)
	defer r.Body.Close()
	tenantName := testingRead(val, "auth.tenantName")
	tenant := self.tenantByName[tenantName]
	if testingRead(val, "auth.passwordCredentials.username") == "u1" &&
		testingRead(val, "auth.passwordCredentials.password") == "p1" &&
		tenant != nil {
		self.ctoken += 1234
		self.authorized[fmt.Sprint(self.ctoken)] = tenant
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(testingAuthRespOk, self.ctoken, tenant.Id, tenant.Name, self.serv.URL)))
	} else {
		w.WriteHeader(401)
		w.Write([]byte("BAD"))
		return
	}

}

func (self *testingServices) keystoneTenants(w http.ResponseWriter, r *http.Request) {
	if self.authorized[r.Header.Get("X-Auth-Token")] == nil {
		w.WriteHeader(401)
		w.Write([]byte("BAD"))
		return
	}
	w.Header().Add("Content-Type", "application/json")

	sLimit := r.URL.Query().Get("limit")
	limit := 999999999
	if sLimit != "" {
		i, _ := strconv.ParseInt(sLimit, 10, 32)
		limit = int(i)
	}
	sMarker := r.URL.Query().Get("marker")

	entry := `        {
	            "id": "%s",
	            "name": "%s",
	            "description": "A description ...",
	            "enabled": true
	        }`
	str := `{
    "tenants": [
%s
    ],
    "tenants_links": []
}
`
	entries := ""
	delim := ""
	for _, t := range self.tenants {

		if sMarker != "" && t.Id != sMarker && len(entries) == 0 {
			continue
		}

		if t.Id == sMarker {
			sMarker = ""
			continue
		}

		if limit <= 0 {
			break
		}
		limit -= 1

		entries += delim
		delim = ",\n"

		entries += fmt.Sprintf(entry, t.Id, t.Name)

	}

	w.Write([]byte(fmt.Sprintf(str, entries)))
}

func (self *testingServices) novaLimits(w http.ResponseWriter, r *http.Request) {
	if self.authorized[r.Header.Get("X-Auth-Token")] == nil {
		w.WriteHeader(401)
		w.Write([]byte("BAD"))
		return
	}
	id, _ := mux.Vars(r)["tenant_id"]
	tenant := self.tenantById[id]
	if tenant != nil {
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(testingNovaLimits, tenant.Modifier)))
	} else {
		w.WriteHeader(404)
		w.Write([]byte("BAD"))
		return
	}
}

func (self *testingServices) novaQuota(w http.ResponseWriter, r *http.Request) {
	if self.authorized[r.Header.Get("X-Auth-Token")] == nil {
		w.WriteHeader(401)
		w.Write([]byte("BAD"))
		return
	}
	id, _ := mux.Vars(r)["tenant_id"]
	tenant := self.tenantById[id]
	if tenant != nil {
		if mux.Vars(r)["tenant2_id"] != id {
			w.WriteHeader(404)
			w.Write([]byte("BAD"))
		} else {
			w.Header().Add("Content-Type", "application/json")
			w.Write([]byte(fmt.Sprintf(testingQuota, tenant.Modifier)))
			return
		}
	} else {
		w.WriteHeader(404)
		w.Write([]byte("BAD"))
		return
	}
}

func (self *testingServices) novaHypervisor(w http.ResponseWriter, r *http.Request) {
	if self.authorized[r.Header.Get("X-Auth-Token")] == nil {
		w.WriteHeader(401)
		w.Write([]byte("BAD"))
		return
	}
	id, _ := mux.Vars(r)["tenant_id"]
	tenant := self.tenantById[id]
	if tenant != nil {
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(testingHypervisors))
	} else {
		w.WriteHeader(404)
		w.Write([]byte("BAD"))
		return
	}
}

func (self *testingServices) populateTenants() {
	self.tenants = []testingTenant{
		testingTenant{Name: "t1", Id: "ab3", Modifier: "100"},
		testingTenant{Name: "t2", Id: "4df", Modifier: "321"},
	}
	self.tenantById = map[string]*testingTenant{}
	self.tenantByName = map[string]*testingTenant{}
	for i, v := range self.tenants {
		self.tenantById[v.Id] = &self.tenants[i]
		self.tenantByName[v.Name] = &self.tenants[i]
	}
}

type tlog struct {
	h http.Handler
}

func (s *tlog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL, "\n", r.Header)
	s.h.ServeHTTP(w, r)
}

func makeFake() *testingServices {
	self := &testingServices{}
	self.populateTenants()
	self.r = mux.NewRouter()
	self.authorized = map[string]*testingTenant{}
	self.ctoken = 100000

	self.r.HandleFunc("/v2.0/tokens", self.keystoneTokens)
	self.r.HandleFunc("/v2.0/tenants", self.keystoneTenants)
	self.r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf(testingKeystoneRoot, self.serv.URL)))
	})

	self.r.HandleFunc("/nova/{tenant_id}/limits", self.novaLimits)
	self.r.HandleFunc("/nova/{tenant_id}/os-quota-sets/{tenant2_id}", self.novaQuota)
	self.r.HandleFunc("/nova/{tenant_id}/os-hypervisors/detail", self.novaHypervisor)

	self.serv = httptest.NewServer(&tlog{self.r})

	return self
}

func TestCollector(t *testing.T) {

	I := func(s string) int {
		i, _ := strconv.ParseInt(s, 10, 32)
		return int(i)
	}

	keySetS := func(s []string) map[string]bool {
		r := map[string]bool{}
		for _, v := range s {
			r[v] = true
		}
		return r
	}
	keySetMM := func(m map[string]map[string]interface{}) map[string]bool {
		r := map[string]bool{}
		for k, _ := range m {
			r[k] = true
		}
		return r
	}
	keySetM := func(m map[string]interface{}) map[string]bool {
		r := map[string]bool{}
		for k, _ := range m {
			r[k] = true
		}
		return r
	}

	setEq := func(a map[string]bool, b map[string]bool) {
		fail := false
		for k, _ := range a {
			ok := b[k]
			if !ok {
				t.Log(k, "present in DUT but not required")
			}
		}
		for k, _ := range b {
			_, ok := a[k]
			if !ok {
				t.Error("!!!", k, "required but missing in dut")
				fail = true
			}
		}
		if fail {
			t.FailNow()
		}
	}

	Convey("Collector", t, func() {

		mockedOS := makeFake()
		defer mockedOS.serv.Close()
		conf := Config{
			User:   "u1",
			Pass:   "p1",
			Tenant: "t2",
			Url:    mockedOS.serv.URL + "/keystone",
		}

		col, err := newCollector(conf)
		if err != nil {
			t.Fatal(err)
		}
		Convey("returns correct set of limits names", func() {

			required := []string{"max_server_meta",
				"max_personality",
				"total_server_groups_used",
				"max_image_meta",
				"max_personality_size",
				"max_total_keypairs",
				"max_security_group_rules",
				"max_server_groups",
				"total_cores_used",
				"total_ram_used",
				"total_instances_used",
				"max_security_groups",
				"total_floating_ips_used",
				"max_total_cores",
				"max_server_group_members",
				"max_total_floating_ips",
				"total_security_groups_used",
				"max_total_instances",
				"max_total_ram_size"}

			setEq(keySetS(col.GetLimitsNames()), keySetS(required))

		})

		Convey("returns correct limits", func() {

			for name, t := range mockedOS.tenantByName {

				limits, err := col.GetLimits(name)
				So(err, ShouldBeNil)

				So(limits["max_image_meta"], ShouldEqual, I(fmt.Sprintf("%v128", t.Modifier)))
				So(limits["max_personality"], ShouldEqual, I(fmt.Sprintf("%v5", t.Modifier)))
				So(limits["max_personality_size"], ShouldEqual, I(fmt.Sprintf("%v10240", t.Modifier)))
				So(limits["max_security_group_rules"], ShouldEqual, I(fmt.Sprintf("%v20", t.Modifier)))
				So(limits["max_security_groups"], ShouldEqual, I(fmt.Sprintf("%v10", t.Modifier)))
				So(limits["max_server_group_members"], ShouldEqual, I(fmt.Sprintf("%v10", t.Modifier)))
				So(limits["max_server_groups"], ShouldEqual, I(fmt.Sprintf("%v10", t.Modifier)))
				So(limits["max_server_meta"], ShouldEqual, I(fmt.Sprintf("%v128", t.Modifier)))
				So(limits["max_total_cores"], ShouldEqual, I(fmt.Sprintf("%v20", t.Modifier)))
				So(limits["max_total_floating_ips"], ShouldEqual, I(fmt.Sprintf("%v10", t.Modifier)))
				So(limits["max_total_instances"], ShouldEqual, I(fmt.Sprintf("%v10", t.Modifier)))
				So(limits["max_total_keypairs"], ShouldEqual, I(fmt.Sprintf("%v100", t.Modifier)))
				So(limits["max_total_ram_size"], ShouldEqual, 1024*1024*I(fmt.Sprintf("%v51200", t.Modifier)))
				So(limits["total_cores_used"], ShouldEqual, I(fmt.Sprintf("%v1", t.Modifier)))
				So(limits["total_floating_ips_used"], ShouldEqual, I(fmt.Sprintf("%v0", t.Modifier)))
				So(limits["total_instances_used"], ShouldEqual, I(fmt.Sprintf("%v1", t.Modifier)))
				So(limits["total_ram_used"], ShouldEqual, 1024*1024*I(fmt.Sprintf("%v2048", t.Modifier)))
				So(limits["total_security_groups_used"], ShouldEqual, I(fmt.Sprintf("%v1", t.Modifier)))
				So(limits["total_server_groups_used"], ShouldEqual, I(fmt.Sprintf("%v0", t.Modifier)))

			}

		})

		Convey("returns correct set of quotas names", func() {

			required := []string{"cores",
				"fixed_ips",
				"floating_ips",
				"instances",
				"key_pairs",
				"ram",
				"security_groups"}

			setEq(keySetS(col.GetQuotasNames()), keySetS(required))
		})

		Convey("returns correct quotas", func() {

			for name, t := range mockedOS.tenantByName {

				quotas, err := col.GetQuotas(name, t.Id)
				So(err, ShouldBeNil)

				So(quotas["ram"], ShouldEqual, 1024*1024*I(fmt.Sprintf("%v51200", t.Modifier)))
				So(quotas["floating_ips"], ShouldEqual, I(fmt.Sprintf("%v10", t.Modifier)))
				So(quotas["key_pairs"], ShouldEqual, I(fmt.Sprintf("%v120", t.Modifier)))
				So(quotas["instances"], ShouldEqual, I(fmt.Sprintf("%v10", t.Modifier)))
				So(quotas["cores"], ShouldEqual, I(fmt.Sprintf("%v20", t.Modifier)))
				So(quotas["fixed_ips"], ShouldEqual, I(fmt.Sprintf("-%v1", t.Modifier)))
				So(quotas["security_groups"], ShouldEqual, I(fmt.Sprintf("%v10", t.Modifier)))
			}

		})

		Convey("returns correct list of tenants", func() {

			dut, err := col.GetTenants()
			So(err, ShouldBeNil)

			for n, t := range mockedOS.tenantByName {
				So(dut[n], ShouldEqual, t.Id)
			}

		})

		Convey("returns correct list of hypervisors", func() {

			dut, err := col.GetHypervisors()
			So(err, ShouldBeNil)

			setEq(keySetMM(dut), keySetS([]string{"1", "2"}))

		})

		Convey("returns correct stat names for hypervisors", func() {

			dut, err := col.GetHypervisors()
			So(err, ShouldBeNil)
			required := []string{"current_workload", "free_disk_gb", "free_ram_mb",
				"hypervisor_version", "memory_mb", "memory_mb_used",
				"running_vms", "vcpus", "vcpus_used", "memory_mb_overcommit",
				"memory_mb_overcommit_withreserve", "vcpus_overcommit", "vcpus_overcommit_withreserve",
			}
			setEq(keySetM(dut["1"]), keySetS(required))
			setEq(keySetM(dut["2"]), keySetS(required))

		})

		Convey("returns correct stats for hypervisors", func() {

			dut, err := col.GetHypervisors()
			So(err, ShouldBeNil)

			So(dut["1"]["current_workload"], ShouldEqual, 1)
			So(dut["1"]["disk_available_least"], ShouldEqual, 2)
			So(dut["1"]["free_disk_gb"], ShouldEqual, 3)
			So(dut["1"]["free_ram_mb"], ShouldEqual, 4)
			So(dut["1"]["hypervisor_version"], ShouldEqual, 2000000)
			So(dut["1"]["local_gb"], ShouldEqual, 5)
			So(dut["1"]["local_gb_used"], ShouldEqual, 6)
			So(dut["1"]["memory_mb"], ShouldEqual, 7)
			So(dut["1"]["memory_mb_used"], ShouldEqual, 8)
			So(dut["1"]["running_vms"], ShouldEqual, 9)
			So(dut["1"]["vcpus"], ShouldEqual, 10)
			So(dut["1"]["vcpus_used"], ShouldEqual, 11)

			So(dut["2"]["current_workload"], ShouldEqual, 12)
			So(dut["2"]["disk_available_least"], ShouldEqual, 13)
			So(dut["2"]["free_disk_gb"], ShouldEqual, 14)
			So(dut["2"]["free_ram_mb"], ShouldEqual, 15)
			So(dut["2"]["hypervisor_version"], ShouldEqual, 3000000)
			So(dut["2"]["local_gb"], ShouldEqual, 16)
			So(dut["2"]["local_gb_used"], ShouldEqual, 17)
			So(dut["2"]["memory_mb"], ShouldEqual, 18)
			So(dut["2"]["memory_mb_used"], ShouldEqual, 19)
			So(dut["2"]["running_vms"], ShouldEqual, 20)
			So(dut["2"]["vcpus"], ShouldEqual, 21)
			So(dut["2"]["vcpus_used"], ShouldEqual, 22)

		})

	})

}

const (
	testingAuthRespOk = `
{
    "access": {
        "token": {
            "issued_at": "2014-01-30T15:30:58.819584",
            "expires": "2028-01-31T15:30:58Z",
            "id": "%[1]v",
            "tenant": {
                "description": null,
                "enabled": true,
                "id": "%[2]v",
                "name": "%[3]v"
            }
        },
        "serviceCatalog": [
            {
                "endpoints": [
                    {
                        "adminURL": "%[4]v/nova/%[2]v",
                        "region": "RegionOne",
                        "internalURL": "%[4]v/nova/%[2]v",
                        "id": "123",
                        "publicURL": "%[4]v/nova/%[2]v"
                    }
                ],
                "endpoints_links": [],
                "type": "compute",
                "name": "nova"
            }
					]
				}
			}
`
	testingKeystoneRoot = `
{
    "versions": {
        "values": [
            {
                "id": "v2.0",
                "links": [
                    {
                        "href": "%s/v2.0/",
                        "rel": "self"
                    },
                    {
                        "href": "http://docs.openstack.org/",
                        "rel": "describedby",
                        "type": "text/html"
                    }
                ],
                "media-types": [
                    {
                        "base": "application/json",
                        "type": "application/vnd.openstack.identity-v2.0+json"
                    }
                ],
                "status": "stable",
                "updated": "2014-04-17T00:00:00Z"
            }
        ]
    }
}`
	testingNovaLimits = `
	{
    "limits": {
        "rate": [],
        "absolute": {
            "maxServerMeta": %[1]v128,
            "maxPersonality": %[1]v5,
            "totalServerGroupsUsed": %[1]v0,
            "maxImageMeta": %[1]v128,
            "maxPersonalitySize": %[1]v10240,
            "maxTotalKeypairs": %[1]v100,
            "maxSecurityGroupRules": %[1]v20,
            "maxServerGroups": %[1]v10,
            "totalCoresUsed": %[1]v1,
            "totalRAMUsed": %[1]v2048,
            "totalInstancesUsed": %[1]v1,
            "maxSecurityGroups": %[1]v10,
            "totalFloatingIpsUsed": %[1]v0,
            "maxTotalCores": %[1]v20,
            "maxServerGroupMembers": %[1]v10,
            "maxTotalFloatingIps": %[1]v10,
            "totalSecurityGroupsUsed": %[1]v1,
            "maxTotalInstances": %[1]v10,
            "maxTotalRAMSize": %[1]v51200
        }
    }
}
	`
	testingQuota = `
	{
    "quota_set": {
        "cores": %[1]v20,
        "fixed_ips": -%[1]v1,
        "floating_ips": %[1]v10,
        "injected_file_content_bytes": %[1]v10240,
        "injected_file_path_bytes": %[1]v255,
        "injected_files": %[1]v5,
        "instances": %[1]v10,
        "key_pairs": %[1]v120,
        "metadata_items": %[1]v128,
        "ram": %[1]v51200,
        "security_group_rules": %[1]v20,
        "security_groups": %[1]v10
    }
}
	`

	testingHypervisors = `
{
	"hypervisors": [
		{
				"cpu_info": "{\"vendor\": \"Intel\", \"model\": \"SandyBridge\", \"arch\": \"x86_64\", \"features\": [\"pge\", \"avx\", \"clflush\", \"sep\", \"syscall\", \"vme\", \"dtes64\", \"msr\", \"fsgsbase\", \"xsave\", \"vmx\", \"erms\", \"xtpr\", \"cmov\", \"smep\", \"ssse3\", \"est\", \"pat\", \"monitor\", \"smx\", \"pbe\", \"lm\", \"tsc\", \"nx\", \"fxsr\", \"tm\", \"sse4.1\", \"pae\", \"sse4.2\", \"pclmuldq\", \"acpi\", \"tsc-deadline\", \"mmx\", \"osxsave\", \"cx8\", \"mce\", \"de\", \"tm2\", \"ht\", \"dca\", \"lahf_lm\", \"popcnt\", \"mca\", \"pdpe1gb\", \"apic\", \"sse\", \"f16c\", \"pse\", \"ds\", \"invtsc\", \"pni\", \"rdtscp\", \"aes\", \"sse2\", \"ss\", \"ds_cpl\", \"pcid\", \"fpu\", \"cx16\", \"pse36\", \"mtrr\", \"pdcm\", \"rdrand\", \"x2apic\"], \"topology\": {\"cores\": 12, \"threads\": 2, \"sockets\": 2}}",
				"current_workload": 1,
				"disk_available_least": 2,
				"free_disk_gb": 3,
				"free_ram_mb": 4,
				"host_ip": "192.168.20.3",
				"hypervisor_hostname": "node-30.domain.tld",
				"hypervisor_type": "QEMU",
				"hypervisor_version": 2000000,
				"id": 1,
				"local_gb": 5,
				"local_gb_used": 6,
				"memory_mb": 7,
				"memory_mb_used": 8,
				"running_vms": 9,
				"service": {
						"disabled_reason": null,
						"host": "node-30.domain.tld",
						"id": 6
				},
				"state": "up",
				"status": "enabled",
				"vcpus": 10,
				"vcpus_used": 11
		},
		{
				"cpu_info": "{\"vendor\": \"Intel\", \"model\": \"SandyBridge\", \"arch\": \"x86_64\", \"features\": [\"pge\", \"avx\", \"clflush\", \"sep\", \"syscall\", \"vme\", \"dtes64\", \"msr\", \"fsgsbase\", \"xsave\", \"vmx\", \"erms\", \"xtpr\", \"cmov\", \"smep\", \"ssse3\", \"est\", \"pat\", \"monitor\", \"smx\", \"pbe\", \"lm\", \"tsc\", \"nx\", \"fxsr\", \"tm\", \"sse4.1\", \"pae\", \"sse4.2\", \"pclmuldq\", \"acpi\", \"tsc-deadline\", \"mmx\", \"osxsave\", \"cx8\", \"mce\", \"de\", \"tm2\", \"ht\", \"dca\", \"lahf_lm\", \"popcnt\", \"mca\", \"pdpe1gb\", \"apic\", \"sse\", \"f16c\", \"pse\", \"ds\", \"invtsc\", \"pni\", \"rdtscp\", \"aes\", \"sse2\", \"ss\", \"ds_cpl\", \"pcid\", \"fpu\", \"cx16\", \"pse36\", \"mtrr\", \"pdcm\", \"rdrand\", \"x2apic\"], \"topology\": {\"cores\": 12, \"threads\": 2, \"sockets\": 2}}",
				"current_workload": 12,
				"disk_available_least": 13,
				"free_disk_gb": 14,
				"free_ram_mb": 15,
				"host_ip": "192.168.20.4",
				"hypervisor_hostname": "node-31.domain.tld",
				"hypervisor_type": "QEMU",
				"hypervisor_version": 3000000,
				"id": 2,
				"local_gb": 16,
				"local_gb_used": 17,
				"memory_mb": 18,
				"memory_mb_used": 19,
				"running_vms": 20,
				"service": {
						"disabled_reason": null,
						"host": "node-31.domain.tld",
						"id": 7
				},
				"state": "up",
				"status": "enabled",
				"vcpus": 21,
				"vcpus_used": 22
		}
]
}

	`
)
