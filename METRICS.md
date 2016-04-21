# snap plugin collector - nova

## Collected Metrics
This plugin has the ability to gather the following metrics:

Namespace | Data Type | Description
----------|-----------|-----------------------
/intel/openstack/nova/cluster/config/allocation_ratio_ram | float64 | The virtual ram to physical ram allocation ratio
/intel/openstack/nova/cluster/config/reserved_node_ram_mb | float64 | Reserved RAM per node - in MB
/intel/openstack/nova/cluster/config/allocation_ratio_cores | float64 |  The overcommit ratio for vCPUs
/intel/openstack/nova/cluster/config/reserved_node_cores | float64 | Reserved cores per node
/intel/openstack/nova/tenant/\<tenant name\>/limits/max_server_meta | int | The maximum number of metadata items associated with a server
/intel/openstack/nova/tenant/\<tenant name\>/limits/max_personality | int | The maximum number of file path/content pairs that can be supplied on server build
/intel/openstack/nova/tenant/\<tenant name\>/limits/total_server_groups_used | int | The total server groups used by this tenant
/intel/openstack/nova/tenant/\<tenant name\>/limits/max_image_meta | int | The maximum number of metadata items associated with an image
/intel/openstack/nova/tenant/\<tenant name\>/limits/max_personality_size | int | The maximum size, in bytes, for each personality file
/intel/openstack/nova/tenant/\<tenant name\>/limits/max_server_groups | int | The maximum number of security group rules allowed for this tenant
/intel/openstack/nova/tenant/\<tenant name\>/limits/max_security_group_rules | int | The maximum number of security groups allowed for this tenant
/intel/openstack/nova/tenant/\<tenant name\>/limits/max_total_keypairs | int| The maximum allowed keypairs allowed for this tenant
/intel/openstack/nova/tenant/\<tenant name\>/limits/total_cores_used | int| Total cores currently in use by tenant
/intel/openstack/nova/tenant/\<tenant name\>/limits/total_ram_used| int | The current RAM used by this tenant shown as gibibyte
/intel/openstack/nova/tenant/\<tenant name\>/limits/total_instances_used | int | The total instances used by this tenant
/intel/openstack/nova/tenant/\<tenant name\>/limits/max_security_groups | int | The maximum number of security groups allowed for this tenant
/intel/openstack/nova/tenant/\<tenant name\>/limits/total_floating_ips_used | int | The total floating IPs used by this tenant
/intel/openstack/nova/tenant/\<tenant name\>/limits/max_total_cores | int | The maximum allowed cores for this tenant
/intel/openstack/nova/tenant/\<tenant name\>/limits/total_security_groups_used | int | The total number of security groups used by this tenant
/intel/openstack/nova/tenant/\<tenant name\>/limits/max_total_floating_ips | int | The maximum allowed floating IPs for this tenant
/intel/openstack/nova/tenant/\<tenant name\>/limits/max_total_instances | int | The maximum number of instances allowed for this tenant
/intel/openstack/nova/tenant/\<tenant name\>/limits/max_total_ram_size| int| The maximum total amount of RAM (MB)
/intel/openstack/nova/tenant/\<tenant name\>/limits/max_server_group_members | int |  The maximum number of server group members for this tenant
/intel/openstack/nova/tenant/\<tenant name\>/quotas/cores| int| The number of allowed instance cores for each tenant
/intel/openstack/nova/tenant/\<tenant name\>/quotas/fixed_ips| int| The number of allowed fixed IP addresses for each tenant, must be equal to or greater than the number of allowed instances
/intel/openstack/nova/tenant/\<tenant name\>/quotas/floating_ips| int| The number of allowed floating IP addresses for each tenant
/intel/openstack/nova/tenant/\<tenant name\>/quotas/instances| int| The number of allowed instances for each tenant
/intel/openstack/nova/tenant/\<tenant name\>/quotas/key_pairs| int| The number of allowed key pairs for each user
/intel/openstack/nova/tenant/\<tenant name\>/quotas/ram| int| The amount of allowed instance RAM, in MB, for each tenant
/intel/openstack/nova/tenant/\<tenant name\>/quotas/security_groups| int| The number of allowed security groups for each tenant
/intel/openstack/nova/hypervisor/\<hypervisor\>/local_gb| int | The total GB of local disk capacity the compute node provides
/intel/openstack/nova/hypervisor/\<hypervisor\>/local_gb_used| int | The amount of disk in GB used on the compute node
/intel/openstack/nova/hypervisor/\<hypervisor\>/current_workload | int | Current workload on the Nova hypervisor
/intel/openstack/nova/hypervisor/\<hypervisor\>/free_disk_gb| int | The calculated amount of disk the compute node has available
/intel/openstack/nova/hypervisor/\<hypervisor\>/free_ram_mb| int | The calculated amount of RAM the compute node has available
/intel/openstack/nova/hypervisor/\<hypervisor\>/memory_mb| int | The total MB of RAM capacity the compute node provides
/intel/openstack/nova/hypervisor/\<hypervisor\>/memory_mb_used | int| The amount of RAM in MB used on the compute node
/intel/openstack/nova/hypervisor/\<hypervisor\>/memory_mb_overcommit | float64 | The multiplication of the amount of RAM in MB used on the compute node and  the virtual ram to physical ram allocation ratio
/intel/openstack/nova/hypervisor/\<hypervisor\>/memory_mb_overcommit_withreserve | float64 | The multiplication of the amount of RAM in MB used on the compute node and  the virtual ram to physical ram allocation ratio, decreased by reserved RAM per node - in MB
/intel/openstack/nova/hypervisor/\<hypervisor\>/vcpus | int | The total number of vCPUs the compute node provides
/intel/openstack/nova/hypervisor/\<hypervisor\>/vcpus_used | int | The number of vCPUs consumed on the compute node
/intel/openstack/nova/hypervisor/\<hypervisor\>/vcpus_overcommit | float64 | The multiplication of the total number of vCPUs the compute node provides and the overcommit ratio for vCPUs
/intel/openstack/nova/hypervisor/\<hypervisor\>/vcpus_overcommit_withreserve | float64 | The multiplication of the total number of vCPUs the compute node provides and the overcommit ratio for vCPUs, decreased by reserved cores per node
/intel/openstack/nova/hypervisor/\<hypervisor\>/disk_available_least | int |  Disk available for the Nova hypervisor shown as gibibyte
/intel/openstack/nova/hypervisor/\<hypervisor\>/running_vms | int| The number of virtual machine instances running on the node
/intel/openstack/nova/hypervisor/\<hypervisor\>/hypervisor_version| int | The hypervisor version
