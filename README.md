# snap collector plugin - OpenStack Nova

This plugin monitors openstack nova resources such as limits, quotas and hypervisors.

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license-and-authors)
6. [Acknowledgements](#acknowledgements)

## Getting Started

Plugin retrieves statistics using Nova Compute v2 Rest API.
### System Requirements
* [golang 1.5+](https://golang.org/dl/) - needed only for building

### Operating systems
This plugin should work on any operating system supported by snap.

### Configuration and Usage

Configuration for this plugin is given via global config. Global configuration files are described in snap's documentation. You have to add section "nova-compute" in "collector" section and then specify following options:
-  `"openstack_user"` - user name used to authenticate (ex. `"admin"`)
-  `"openstack_pass"`- password used to authenticate (ex. `"admin"`)
-  `"openstack_tenant"` - tenant name used to authenticate (ex. `"admin"`)
-  `"openstack_auth_url"` - keystone url (ex. `"http://172.16.0.5:5000/v2.0/"`)

These values should correspond to values given in `nova.conf`.
-  `"allocation_ratio_cores"` - oversubscription ratio for vcpus, used to derive some metrics for hypervisors (ex. 1.5)
-  `"allocation_ratio_ram"` - oversubscription ratio for memory, used to derive some metrics for hypervisors (ex. 3)
-  `"reserved_node_cores"` - reserved virtual cores, used to derive some metrics for hypervisors (ex. 2)
-  `"reserved_node_ram_mb"` - reserved virtual memory, used to derive some metrics for hypervisors (ex. 2048)

### Installation
#### Download OpenStack Nova plugin binary:
You can get the pre-built binaries for your OS and architecture at snap's [GitHub Releases](https://github.com/intelsdi-x/snap/releases) page.

#### To build the plugin binary:
Fork https://github.com/intelsdi-x/snap-plugin-collector-nova  
Clone repo into `$GOPATH/src/github.com/intelsdi-x/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-nova.git
```

Build the plugin by running make within the cloned repo:
```
$ make
```
This builds the plugin in `/build/rootfs/`


## Documentation

### Collected Metrics
This plugin has the ability to gather the following metrics:

Namespace | Data Type | Description
----------|-----------|-----------------------
/intel/openstack/nova/cluster/config/allocation_ratio_ram|float64|
/intel/openstack/nova/cluster/config/reserved_node_ram_mb|float64|
/intel/openstack/nova/cluster/config/allocation_ratio_cores|float64|
/intel/openstack/nova/cluster/config/reserved_node_cores|float64|
/intel/openstack/nova/tenant/\<tenant name\>/limits/max_server_meta|int|
/intel/openstack/nova/tenant/\<tenant name\>/limits/max_personality|int|
/intel/openstack/nova/tenant/\<tenant name\>/limits/total_server_groups_used|int|
/intel/openstack/nova/tenant/\<tenant name\>/limits/max_image_meta|int|
/intel/openstack/nova/tenant/\<tenant name\>/limits/max_personality_size|int|
/intel/openstack/nova/tenant/\<tenant name\>/limits/max_server_groups|int|
/intel/openstack/nova/tenant/\<tenant name\>/limits/max_security_group_rules|int|
/intel/openstack/nova/tenant/\<tenant name\>/limits/max_total_keypairs|int|
/intel/openstack/nova/tenant/\<tenant name\>/limits/total_cores_used|int|
/intel/openstack/nova/tenant/\<tenant name\>/limits/total_ram_used|int|
/intel/openstack/nova/tenant/\<tenant name\>/limits/total_instances_used|int|
/intel/openstack/nova/tenant/\<tenant name\>/limits/max_security_groups|int|
/intel/openstack/nova/tenant/\<tenant name\>/limits/total_floating_ips_used|int|
/intel/openstack/nova/tenant/\<tenant name\>/limits/max_total_cores|int|
/intel/openstack/nova/tenant/\<tenant name\>/limits/total_security_groups_used|int|
/intel/openstack/nova/tenant/\<tenant name\>/limits/max_total_floating_ips|int|
/intel/openstack/nova/tenant/\<tenant name\>/limits/max_total_instances|int|
/intel/openstack/nova/tenant/\<tenant name\>/limits/max_total_ram_size|int|
/intel/openstack/nova/tenant/\<tenant name\>/limits/max_server_group_members|int|
/intel/openstack/nova/tenant/\<tenant name\>/quotas/cores|int|
/intel/openstack/nova/tenant/\<tenant name\>/quotas/fixed_ips|int|
/intel/openstack/nova/tenant/\<tenant name\>/quotas/floating_ips|int|
/intel/openstack/nova/tenant/\<tenant name\>/quotas/instances|int|
/intel/openstack/nova/tenant/\<tenant name\>/quotas/key_pairs|int|
/intel/openstack/nova/tenant/\<tenant name\>/quotas/ram|int|
/intel/openstack/nova/tenant/\<tenant name\>/quotas/security_groups|int|
/intel/openstack/nova/hypervisor/\<hypervisor\>/local_gb|int|
/intel/openstack/nova/hypervisor/\<hypervisor\>/local_gb_used|int|
/intel/openstack/nova/hypervisor/\<hypervisor\>/vcpus_used|int|
/intel/openstack/nova/hypervisor/\<hypervisor\>/current_workload|int|
/intel/openstack/nova/hypervisor/\<hypervisor\>/free_disk_gb|int|
/intel/openstack/nova/hypervisor/\<hypervisor\>/free_ram_mb|int|
/intel/openstack/nova/hypervisor/\<hypervisor\>/memory_mb|int|
/intel/openstack/nova/hypervisor/\<hypervisor\>/memory_mb_overcommit_withreserve|float64|
/intel/openstack/nova/hypervisor/\<hypervisor\>/vcpus_overcommit|float64|
/intel/openstack/nova/hypervisor/\<hypervisor\>/vcpus_overcommit_withreserve|float64|
/intel/openstack/nova/hypervisor/\<hypervisor\>/disk_available_least|int|
/intel/openstack/nova/hypervisor/\<hypervisor\>/memory_mb_used|int|
/intel/openstack/nova/hypervisor/\<hypervisor\>/running_vms|int|
/intel/openstack/nova/hypervisor/\<hypervisor\>/hypervisor_version|int|
/intel/openstack/nova/hypervisor/\<hypervisor\>/memory_mb_overcommit|float64|
/intel/openstack/nova/hypervisor/\<hypervisor\>/vcpus|int|


### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release. If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-nova/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-collector-nova/pulls).

## Community Support
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)

## Contributing
We love contributions

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[snap](http://github.com:intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements

* Author: [Lukasz Mroz](https://github.com/lmroz)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.
