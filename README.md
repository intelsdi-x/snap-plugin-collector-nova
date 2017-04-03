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
* [golang 1.6+](https://golang.org/dl/) - needed only for building

### Operating systems
All OSs currently supported by snap:
* Linux/amd64

### Installation
#### Download OpenStack Nova plugin binary:
You can get the pre-built binaries for your OS and architecture at plugins's [GitHub Releases](https://github.com/intelsdi-x/snap-plugin-collector-nova/releases) page. Download the plugins package from the latest release, unzip and store in a path you want `snaptel` to access.

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
This builds the plugin in `./build/`

### Configuration and Usage

* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started).
* Create Global Config, see description in [snap's Global Config](https://github.com/intelsdi-x/snap-plugin-collector-nova/blob/master/README.md#snaps-global-config).
* Load the plugin and create a task, see example in [Examples](https://github.com/intelsdi-x/snap-plugin-collector-nova/blob/master/README.md#examples).

Configuration for this plugin is given via global config. Global configuration files are described in [snap's documentation](https://github.com/intelsdi-x/snap/blob/master/docs/SNAPTELD_CONFIGURATION.md). You have to add section "nova-compute" in "collector" section and then specify following options:
-  `"openstack_user"` - user name used to authenticate (ex. `"admin"`),
-  `"openstack_pass"`- password used to authenticate (ex. `"admin"`),
-  `"openstack_tenant"` - tenant name used to authenticate (ex. `"admin"`),
-  `"openstack_auth_url"` - keystone url (ex. `"http://172.16.0.5:5000/v2.0/"`).
If you're using authentication API in v3 you need to set one of those two configuration options (Note that both keys should be present in config, but only one of them is allowed to be set):
- `"openstack_domain_name"` - domain name
- `"openstack_domain_id"` - domain name

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

## Documentation

### Collected Metrics

List of collected metrics is described in [METRICS.md](https://github.com/intelsdi-x/snap-plugin-collector-nova/blob/master/METRICS.md).

### Examples
Example running nova collector and writing data to a file.

Create or download example global config on snapteld's node:
```
mkdir -p /etc/snap/
curl -sfLO https://github.com/intelsdi-x/snap-plugin-collector-nova/blob/master/examples/cfg/cfg.json -o /etc/snap/snapteld.conf
```

Ensure [snap daemon is running](https://github.com/intelsdi-x/snap#running-snap):
* initd: `sudo service snap-telemetry start`
* systemd: `sudo systemctl start snap-telemetry`
* command line: `sudo snapteld -l 1 -t 0 &`

Download and load snap plugins:
```
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-collector-nova/latest/linux/x86_64/snap-plugin-collector-nova
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-publisher-file/latest/linux/x86_64/snap-plugin-publisher-file
$ snaptel plugin load snap-plugin-collector-nova
$ snaptel plugin load snap-plugin-publisher-file
```

See available metrics for your system

```
$ snaptel metric list
```

Download an [example task file](https://github.com/intelsdi-x/snap-plugin-collector-nova/blob/master/examples/tasks/task.json) and load it:
```
$ curl -sfLO https://github.com/intelsdi-x/snap-plugin-collector-nova/blob/master/examples/tasks/task.json
$ snaptel task create -t task.json
Using task manifest to create task
Task created
ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
Name: Task-02dd7ff4-8106-47e9-8b86-70067cd0a850
State: Running
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


