
# Plunder

The complete tool for finding kubernetes gold amongst bits of bare-metal!

Plunder is a single-binary service that provides `DHCP`/`TFTP` and `HTTP` functionality to bootstrap bare-metal (and virtual) servers, it also manages the creation of Ubuntu Preseed configurations to manage the building and configuration of the ubuntu operating system (Other OS's may be added later). 

![Plunder Captain](./image/plunder_captain.png)

## Warning

*NOTE 1* As this provides low-level networking services, only run on a network that is safe to do so. Providing DHCP on a network that already provides DHCP services can lead to un-expected behaviour (and angry network administrators)

*NOTE 2* As DHCP/TFTP and HTTP all bind to low ports < 1024, root access (or sudo) is required to start the plunder services.

## Building

Releases may be provdided in the future, but for now grab the source and `go build` or `make build`.

## Configuration

A `./plunder config > config.json` will look at the network configuration of the current machine and build a default configuration file (in json). This file will need opening in your favourite text editor an modifying to ensure that `plunder` works correctly. 

```
{
	"adapter": "ens160",
	"enableDHCP": false,
	"addressDHCP": "192.168.0.110",
	"startDHCP": "",
	"leasePoolDHCP": 20,
	"gatewayDHCP": "192.168.0.110",
	"nameserverDHCP": "192.168.0.110",
	"enableTFTP": false,
	"addressTFTP": "192.168.0.110",
	"enableHTTP": false,
	"addressHTTP": "192.168.0.110",
	"pxePath": "undionly.kpxe",
	"kernelPath": "",
	"initrdPath": "",
	"cmdline": ""
}
```

### Retreiving bootstrap components

The `./plunder get` command will download the `iPXE` bootloader that is needed by the `TFTP` service in order to bootstrap the OS build. 

### Networking configuration
In a mutli-adapter host (recommended) ensure that the correct adapter is used i.e. `eth0 -> eth1` also that the correct IP addresses are used, the IP addresses are needed for the TFTP and HTTP stages as that is all managed through TCP connections. 

*NOTE* the `startDHCP` field is _required_ and should ideally be `addressDHCP` + 1, the `leasePoolDHCP` will then manage a pool of IP addresses from that start address.

### Services configuration
By default, all services are disabled (this is to make sure you look at the configuration before advertising random network services). The `plunder` application will not start if all services are disabled and will present you with a warning message, in order to enable services change `false` to `true`.

### Kernel / Initrd etc.. 

The plan is to have `plunder` mount and extract the correct kernels and netboot `initrd` (TBD)

#Troubleshooting

PXE booting provides very little feedback when things aren't working, but usually the hand-off is why things wont work i.e. `DHCP` -> `TFTP` boot. Logs from `plunder` should show the hand-off from the CLI.

#Roadmap

Fix the templating of Preseed and Kickstart files, automate the entire process end-to-end. May have `plunder` keep all configurations internally and use http handlers to expose them as urls to the boot loader (TBD)

TL;DR make better.

Created on 2018-11-14 17:30:01
