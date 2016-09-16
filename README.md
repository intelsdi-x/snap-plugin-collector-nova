# snap collector plugin - OpenStack Nova

This plugin monitors openstack nova resources such as limits, quotas and hypervisors.

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Operating systems](#operating-systems)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
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
All OSs currently supported by snap:
* Linux/amd64

### Installation
#### Download OpenStack Nova plugin binary:
You can get the pre-built binaries for your OS and architecture at snap's [GitHub Releases](https://github.com/intelsdi-x/snap/releases) page. Download the plugins package from the latest release, unzip and store in a path you want `snapd` to access.

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

### Configuration and Usage

* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started).
* Create Global Config, see description in [snap's Global Config] (https://github.com/intelsdi-x/snap-plugin-collector-nova/blob/master/README.md#snaps-global-config).
* Load the plugin and create a task, see example in [Examples](https://github.com/intelsdi-x/snap-plugin-collector-nova/blob/master/README.md#examples).

## Documentation

### Collected Metrics

List of collected metrics is described in [METRICS.md](https://github.com/intelsdi-x/snap-plugin-collector-nova/blob/master/METRICS.md).

### snap's Global Config

Configuration for this plugin is given via global config. Global configuration files are described in [snap's documentation](https://github.com/intelsdi-x/snap/blob/master/docs/SNAPD_CONFIGURATION.md). You have to add section "nova-compute" in "collector" section and then specify following options:
-  `"openstack_user"` - user name used to authenticate (ex. `"admin"`),
-  `"openstack_pass"`- password used to authenticate (ex. `"admin"`),
-  `"openstack_tenant"` - tenant name used to authenticate (ex. `"admin"`),
-  `"openstack_auth_url"` - keystone url (ex. `"http://172.16.0.5:5000/v2.0/"`).

These values should correspond to values given in `nova.conf`:
-  `"allocation_ratio_cores"` - oversubscription ratio for vcpus, used to derive some metrics for hypervisors (ex. 1.5),
-  `"allocation_ratio_ram"` - oversubscription ratio for memory, used to derive some metrics for hypervisors (ex. 3),
-  `"reserved_node_cores"` - reserved virtual cores, used to derive some metrics for hypervisors (ex. 2),
-  `"reserved_node_ram_mb"` - reserved virtual memory, used to derive some metrics for hypervisors (ex. 2048).

Example global configuration file for snap-plugin-collector-nova plugin (exemplary file in [examples/cfg/] (examples/cfg/)):
```
{
    "control": {
        "plugins": {
            "collector": {
                "nova-compute": {
                  "all": {
                    "openstack_auth_url": "http://localhost:5000/v2.0/",
                    "openstack_user": "admin",
                    "openstack_pass": "admin",
                    "openstack_tenant": "admin",
                    "allocation_ratio_cores" : 1.5,
                    "allocation_ratio_ram": 3,
                    "reserved_node_cores" : 2,
                    "reserved_node_ram_mb" : 2048
                  }
                }
            },
            "publisher": {},
            "processor": {}
        }
    }
}
```

### Examples

Example running snap-plugin-collector-nova plugin and writing data to a file.

Make sure that your `$SNAP_PATH` is set, if not:
```
$ export SNAP_PATH=<snapDirectoryPath>/build
```
Other paths to files should be set according to your configuration, using a file you should indicate where it is located.

Create Global Config, see example in [examples/cfg/] (examples/cfg/).

In one terminal window, open the snap daemon (in this case with logging set to 1,  trust disabled and global configuration saved in cfg.json ):
```
$ $SNAP_PATH/bin/snapd -l 1 -t 0 --config cfg.json
```
In another terminal window:

Load snap-plugin-collector-nova plugin
```
$ $SNAP_PATH/bin/snapctl plugin load snap-plugin-collector-nova
```
Load file plugin for publishing:
```
$ $SNAP_PATH/bin/snapctl plugin load $SNAP_PATH/plugin/snap-publisher-file
```
See available metrics for your system

```
$ $SNAP_PATH/bin/snapctl metric list
```

Create a task manifest file to use snap-plugin-collector-nova plugin (exemplary files in [examples/tasks/] (examples/tasks/)):
```
{
  "schedule": {
    "interval": "10s",
    "type": "simple"
  },
  "version": 1,
  "workflow": {
    "collect": {
      "config": {},
      "metrics": {
        "/intel/openstack/nova/tenant/admin/limits/max_image_meta": {},
        "/intel/openstack/nova/tenant/admin/limits/max_personality": {},
        "/intel/openstack/nova/tenant/admin/limits/max_personality_size": {},
        "/intel/openstack/nova/tenant/admin/limits/max_security_group_rules": {},
        "/intel/openstack/nova/tenant/admin/limits/max_security_groups": {},
        "/intel/openstack/nova/tenant/admin/limits/max_server_group_members": {},
        "/intel/openstack/nova/tenant/admin/limits/max_server_groups": {},
        "/intel/openstack/nova/tenant/admin/limits/max_server_meta": {},
        "/intel/openstack/nova/tenant/admin/limits/max_total_cores": {},
        "/intel/openstack/nova/tenant/admin/limits/max_total_floating_ips": {},
        "/intel/openstack/nova/tenant/admin/limits/max_total_instances": {},
        "/intel/openstack/nova/tenant/admin/limits/max_total_keypairs": {},
        "/intel/openstack/nova/tenant/admin/limits/max_total_ram_size": {},
        "/intel/openstack/nova/tenant/admin/limits/total_cores_used": {},
        "/intel/openstack/nova/tenant/admin/limits/total_floating_ips_used": {},
        "/intel/openstack/nova/tenant/admin/limits/total_instances_used": {},
        "/intel/openstack/nova/tenant/admin/limits/total_ram_used": {},
        "/intel/openstack/nova/tenant/admin/limits/total_security_groups_used": {},
        "/intel/openstack/nova/tenant/admin/limits/total_server_groups_used": {},
        "/intel/openstack/nova/tenant/admin/quotas/cores": {},
        "/intel/openstack/nova/tenant/admin/quotas/fixed_ips": {},
        "/intel/openstack/nova/tenant/admin/quotas/floating_ips": {},
        "/intel/openstack/nova/tenant/admin/quotas/instances": {},
        "/intel/openstack/nova/tenant/admin/quotas/key_pairs": {},
        "/intel/openstack/nova/tenant/admin/quotas/ram": {},
        "/intel/openstack/nova/tenant/admin/quotas/security_groups": {}
      },
      "process": null,
      "publish": [
                {
                    "plugin_name": "file",
                    "config": {
                        "file": "/tmp/published_nova"
                    }
                }
      ]
    }
  }
}
```

Create a task:
```
$ $SNAP_PATH/bin/snapctl task create -t examples/tasks/task.json
```

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release. If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-nova/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-collector-nova/pulls).

## Community Support
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

And **thank you!** Your contribution, through code and participation, is incredibly important to us.

## License
[snap](http://github.com:intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [Lukasz Mroz](https://github.com/lmroz)


