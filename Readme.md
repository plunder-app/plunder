
# Plunder

The complete tool for finding kubernetes gold amongst bits of bare-metal!

![Plunder Captain](./image/plunder_captain.png)

Plunder is a single-binary service that provides the following services:

- `DHCP`
- `TFTP` 
- `HTTP`
- `OS Provisioning`
- `Remote command execution`
- `Scripting engine`

This combined functionality provides the capability to bootstrap bare-metal (and virtual) servers, deploy an operating system through:

- Ubuntu/Debian preseed
- CentOS kickstart (still WIP)

## Getting Plunder

Prebuilt binaries for Darwin(MacOS)/Linux and Windows can be found on the [releases](https://github.com/thebsdbox/plunder/releases) page.

### Building

If you wish to build the code yourself then this can be done simply by running:

```
go get github.com/thebsdbox/plunder
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

Fix the templating of Preseed and Kickstart files, automate the entire process end-to-end. May have `plunder` keep all configurations internally and use http handlers to expose them as urls to the boot loader (TBD)

TL;DR make better.

Created on 2018-11-14 17:30:01
