
# Plunder

The complete tool for finding **Infrastructure** gold amongst bits of bare-metal!

![Plunder Captain](./image/plunder_captain.png)

## Overview

Plunder is a single-binary server that is all designed in order to make the provisioning of servers, platforms and applications easier. It is deployed as a server that an end user can interact with through it's **Api-server** in order to control and automate the usage. At this time interacting with the api-server is detailed in the source [https://github.com/plunder-app/plunder/blob/master/pkg/apiserver/endpoints.go](https://github.com/plunder-app/plunder/blob/master/pkg/apiserver/endpoints.go), however documentation will be added soon. 

From an end-user interaction a plunder control utility has been created: 

[https://github.com/plunder-app/pldrctl](https://github.com/plunder-app/pldrctl) - provides the capability to query and create deployments and configurations within a plunder instance.

### Services

- `DHCP` - Allocating an IP addressing and pointing to a TFTP server
- `TFTP` - Bootstrapping an Operating system install (uses iPXE)
- `HTTP` - Provides a services where the bootstrap can pull the components needed for the OS install.

An operating system can be easily performed using either **preseed** or **kickstart**, alternatively custom kernels and init ramdisks can be specified to be used based upon Mac address.

### Automation

Further more once the operating system has been provisioned there are usually post-deployment tasks in order to complete an installation. Plunder has the capability to do the following:

- `Remote command execution` - Over SSH (key configured above)
- `Scripting engine` - A JSON/YAML language that also supports plugins to extend the capablities of the automation engine.

A small repository of existing deployment maps has been created [https://github.com/plunder-app/maps](https://github.com/plunder-app/maps)

### Additional features

- `iso support` - Plunder no longer requires a user with elevated privileges to mount an OS ISO in order to read the contents. Plunder can read files directly from the iso file and expose them to an installer through `http`.
- `online updates` - As all configuration to plunder is exposed and managed through an API, it provides the capability of performing most configuration changes with no down time or restarts.
- `in-memory configurations` - Plunder will create all deployment configurations and hold them in memory, meaning that it is stateless and it doesn't leave configuration all over a filesystem
- `VMware deployment support` - Plunder can deploy preseed/kickstart and now vSphere installations.
- `Management of unclaimed devices` - Plunder will watch and keep a pool of devices that aren't being deployed and can force them to reboot/restart until they're needed for deployment.
- `Logging of remote execution` - Plunder can now store all execution logs in-memory until told to clear them.

## Getting Plunder

Prebuilt binaries for Darwin(MacOS)/Linux and Windows can be found on the [releases](https://github.com/plunder-app/plunder/releases) page.

### Building

If you wish to build the code yourself then this can be done simply by running:

```
go get -u github.com/plunder-app/plunder
```
Alternatively clone the repository and either `go build` or `make build`, note that using the makefile will ensure that the current git commit and version number are returned by `plunder version`.

## Usage!

One of the key design concepts was to try to simplify the amount of moving parts required to bootstrap a server, therefore `plunder` aims to be a single tool that you can use. It also aims to simplify the amount of configuration files and configuration work required, it does this by auto-detecting most configuration and producing mainly completed configuration as needed. 

One thing to be aware of is that `plunder` doesn't require replacing anything that already exists in the infrastructure.

The documentation is available [here](./docs/)

### Warning

*NOTE 1* As this provides low-level networking services, only run on a network that is safe to do so. Providing DHCP on a network that already provides DHCP services can lead to un-expected behaviour (and angry network administrators)

*NOTE 2* As DHCP/TFTP and HTTP all bind to low ports < 1024, root access (or sudo) is required to start the plunder services.

# Troubleshooting

PXE booting provides very little feedback when things aren't working, but usually the hand-off is why things wont work i.e. `DHCP` -> `TFTP` boot. Logs from `plunder` should show the hand-off from the CLI.

# Roadmap

- Ability to automate deployments over VMware VMTools

- Windows deployments

- Tidier logging

- Stability enhancements

- Additional plugins

  
