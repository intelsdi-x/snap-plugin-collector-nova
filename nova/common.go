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

package nova

import (
	"fmt"

	"github.com/gophercloud/gophercloud"

	"github.com/mitchellh/mapstructure"
)

// Nova is abstraction layer for nova-compute service, providing consisten api
// in case of version changes.
type Nova interface {
	GetLimits() (LimitsAbsolute, error)
	GetQuotas() (QuotaSet, error)
	GetHypervisorStatistics() (HypervisorStatistics, error)
}

type QuotaSet struct {
	Cores          int
	FixedIps       int
	FloatingIps    int
	Instances      int
	KeyPairs       int
	Ram            int
	SecurityGroups int
}

type HypervisorStatistics struct {
	Id                 int
	CurrentWorkload    int
	DiskAvailableLeast int
	FreeDiskGB         int
	FreeRamMB          int
	HypervisorVersion  int
	LocalGB            int
	LocalGBUsed        int
	MemoryMB           int
	MemoryMBUsed       int
	RunningVMs         int
	VCPUs              int
	VCPUsUsed          int
}

type LimitsAbsolute struct {
	MaxServerMeta           int
	MaxPersonality          int
	TotalServerGroupsUsed   int
	MaxImageMeta            int
	MaxPersonalitySize      int
	MaxServerGroups         int
	MaxSecurityGroupRules   int
	MaxTotalKeypairs        int
	TotalCoresUsed          int
	TotalRAMUsed            int
	TotalInstancesUsed      int
	MaxSecurityGroups       int
	TotalFloatingIPsUsed    int
	MaxTotalCores           int
	TotalSecurityGroupsUsed int
	MaxTotalFloatingIPs     int
	MaxTotalInstances       int
	MaxTotalRamSize         int
	MaxServerGroupMembers   int
}

// Get performs request to given service and deserializes json response to struct.
// out should be pointer to resulting structure. Returns error if request or
// deserialization failed.
func Get(client *gophercloud.ServiceClient, request string, out interface{}) error {
	url := client.ServiceURL(request)
	var resp interface{}
	_, err := client.Get(url, &resp, nil)

	if err != nil {
		return fmt.Errorf("request failed: (%v)", err)
	}

	err = mapstructure.Decode(resp, out)

	if err != nil {
		return fmt.Errorf("decoding failed: (%v)", err)
	}

	return nil
}
