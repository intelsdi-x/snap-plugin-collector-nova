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

package v2

import (
	"net/url"

	"github.com/intelsdi-x/snap-plugin-collector-nova/nova"

	"github.com/gophercloud/gophercloud"
)

type LimitsAbsoluteRespV2 struct {
	MaxServerMeta           int `mapstructure:"maxServerMeta"`
	MaxPersonality          int `mapstructure:"maxPersonality"`
	TotalServerGroupsUsed   int `mapstructure:"totalServerGroupsUsed"`
	MaxImageMeta            int `mapstructure:"maxImageMeta"`
	MaxPersonalitySize      int `mapstructure:"maxPersonalitySize"`
	MaxServerGroups         int `mapstructure:"maxServerGroups"`
	MaxSecurityGroupRules   int `mapstructure:"maxSecurityGroupRules"`
	MaxTotalKeypairs        int `mapstructure:"maxTotalKeypairs"`
	TotalCoresUsed          int `mapstructure:"totalCoresUsed"`
	TotalRAMUsed            int `mapstructure:"totalRAMUsed"`
	TotalInstancesUsed      int `mapstructure:"totalInstancesUsed"`
	MaxSecurityGroups       int `mapstructure:"maxSecurityGroups"`
	TotalFloatingIPsUsed    int `mapstructure:"totalFloatingIpsUsed"`
	MaxTotalCores           int `mapstructure:"maxTotalCores"`
	TotalSecurityGroupsUsed int `mapstructure:"totalSecurityGroupsUsed"`
	MaxTotalFloatingIPs     int `mapstructure:"maxTotalFloatingIps"`
	MaxTotalInstances       int `mapstructure:"maxTotalInstances"`
	MaxTotalRamSize         int `mapstructure:"maxTotalRAMSize"`
	MaxServerGroupMembers   int `mapstructure:"maxServerGroupMembers"`
}

type LimitsLimitsRespV2 struct {
	Absolute LimitsAbsoluteRespV2 `mapstructure:"absolute"`
}

type LimitsRespV2 struct {
	Limits LimitsLimitsRespV2 `mapstructure:"limits"`
}

// Converts internal representation used for deserialization to common API
func (self *LimitsAbsoluteRespV2) ToLimitsAbsolute() nova.LimitsAbsolute {
	return nova.LimitsAbsolute{
		MaxServerMeta:           self.MaxServerMeta,
		MaxPersonality:          self.MaxPersonality,
		TotalServerGroupsUsed:   self.TotalServerGroupsUsed,
		MaxImageMeta:            self.MaxImageMeta,
		MaxPersonalitySize:      self.MaxPersonalitySize,
		MaxServerGroups:         self.MaxServerGroups,
		MaxSecurityGroupRules:   self.MaxSecurityGroupRules,
		MaxTotalKeypairs:        self.MaxTotalKeypairs,
		TotalCoresUsed:          self.TotalCoresUsed,
		TotalRAMUsed:            self.TotalRAMUsed,
		TotalInstancesUsed:      self.TotalInstancesUsed,
		MaxSecurityGroups:       self.MaxSecurityGroups,
		TotalFloatingIPsUsed:    self.TotalFloatingIPsUsed,
		MaxTotalCores:           self.MaxTotalCores,
		TotalSecurityGroupsUsed: self.TotalSecurityGroupsUsed,
		MaxTotalFloatingIPs:     self.MaxTotalFloatingIPs,
		MaxTotalInstances:       self.MaxTotalInstances,
		MaxTotalRamSize:         self.MaxTotalRamSize,
		MaxServerGroupMembers:   self.MaxServerGroupMembers,
	}
}

type QuotaSetRespV2 struct {
	Cores          int `mapstructure:"cores"`
	FixedIps       int `mapstructure:"fixed_ips"`
	FloatingIps    int `mapstructure:"floating_ips"`
	Instances      int `mapstructure:"instances"`
	KeyPairs       int `mapstructure:"key_pairs"`
	Ram            int `mapstructure:"ram"`
	SecurityGroups int `mapstructure:"security_groups"`
}

type QuotaRespV2 struct {
	QuotaSet QuotaSetRespV2 `mapstructure:"quota_set"`
}

// Converts internal representation used for deserialization to common API
func (self *QuotaSetRespV2) ToQuotaSet() nova.QuotaSet {
	return nova.QuotaSet{
		Cores:          self.Cores,
		FixedIps:       self.FixedIps,
		FloatingIps:    self.FloatingIps,
		Instances:      self.Instances,
		KeyPairs:       self.KeyPairs,
		Ram:            self.Ram,
		SecurityGroups: self.SecurityGroups,
	}
}

type HypervisorInnerRespV2 struct {
	Id                 int `mapstructure:"id"`
	CurrentWorkload    int `mapstructure:"current_workload"`
	DiskAvailableLeast int `mapstructure:"disk_available_least"`
	FreeDiskGB         int `mapstructure:"free_disk_gb"`
	FreeRamMB          int `mapstructure:"free_ram_mb"`
	HypervisorVersion  int `mapstructure:"hypervisor_version"`
	LocalGB            int `mapstructure:"local_gb"`
	LocalGBUsed        int `mapstructure:"local_gb_used"`
	MemoryMB           int `mapstructure:"memory_mb"`
	MemoryMBUsed       int `mapstructure:"memory_mb_used"`
	RunningVMs         int `mapstructure:"running_vms"`
	VCPUs              int `mapstructure:"vcpus"`
	VCPUsUsed          int `mapstructure:"vcpus_used"`
}

type HypervisorRespV2 struct {
	Hypervisors []HypervisorInnerRespV2 `mapstructure:"hypervisors"`
}

// Converts internal representation used for deserialization to common API
func (self *HypervisorInnerRespV2) ToHypervisorStatistics() nova.HypervisorStatistics {
	return nova.HypervisorStatistics{
		Id:                 self.Id,
		CurrentWorkload:    self.CurrentWorkload,
		DiskAvailableLeast: self.DiskAvailableLeast,
		FreeDiskGB:         self.FreeDiskGB,
		FreeRamMB:          self.FreeRamMB,
		HypervisorVersion:  self.HypervisorVersion,
		LocalGB:            self.LocalGB,
		LocalGBUsed:        self.LocalGBUsed,
		MemoryMB:           self.MemoryMB,
		MemoryMBUsed:       self.MemoryMBUsed,
		RunningVMs:         self.RunningVMs,
		VCPUs:              self.VCPUs,
		VCPUsUsed:          self.VCPUsUsed,
	}
}

type NovaV2 struct {
	Client *gophercloud.ServiceClient
}

//Implements limits retrieval, returns limits or error if rest request failed
func (self *NovaV2) GetLimits() (nova.LimitsAbsolute, error) {
	var out LimitsRespV2
	err := nova.Get(self.Client, "limits", &out)
	if err != nil {
		return nova.LimitsAbsolute{}, err
	}

	return out.Limits.Absolute.ToLimitsAbsolute(), nil
}

//Implements quotas retrieval, returns quotas or error if rest request failed
func (self *NovaV2) GetQuotas(tenantId string) (nova.QuotaSet, error) {
	var out QuotaRespV2
	err := nova.Get(self.Client, "os-quota-sets/"+url.QueryEscape(tenantId), &out)
	if err != nil {
		return nova.QuotaSet{}, err
	}

	return out.QuotaSet.ToQuotaSet(), nil
}

//Implements hypervisor stats retrieval, returns these stats or error if rest request failed
func (self *NovaV2) GetHypervisorStatistics() ([]nova.HypervisorStatistics, error) {
	var out HypervisorRespV2
	err := nova.Get(self.Client, "os-hypervisors/detail", &out)
	if err != nil {
		return nil, err
	}

	result := make([]nova.HypervisorStatistics, len(out.Hypervisors))
	for i, v := range out.Hypervisors {
		result[i] = v.ToHypervisorStatistics()
	}

	return result, nil
}
