
# Plunder

The complete tool for finding kubernetes gold amongst bits of bare-metal!

Plunder is a single-binary service that provides `DHCP`/`TFTP` and `HTTP` functionality to bootstrap bare-metal (and virtual) servers, it also manages the creation of Ubuntu Preseed configurations to manage the building and configuration of the ubuntu operating system (Other OS's may be added later). 

![Plunder Captain](./image/plunder_captain.png)

## Warning

*NOTE 1* As this provides low-level networking services, only run on a network that is safe to do so. Providing DHCP on a network that already provides DHCP services can lead to un-expected behaviour (and angry network administrators)

*NOTE 2* As DHCP/TFTP and HTTP all bind to low ports < 1024, root access (or sudo) is required to start the plunder services.

## Building

Releases may be provdided in the future, but for now grab the source and `go build` or `make build`.

## Server Configuration

A `./plunder config server > config.json` will look at the network configuration of the current machine and build a default configuration file (in json). This file will need opening in your favourite text editor an modifying to ensure that `plunder` works correctly. 

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

## Deployment Configuration

A `./plunder config deployment > deployment.json` will create a blank deployment configuration that can be pre-populated in order to create specific deployments.

A configured deployment should resemble something like the example below:

```
{
	"globalConfig": {
		"gateway": "192.168.1.1",
		"address": "",
		"subnet": "255.255.255.0",
		"nameserver": "8.8.8.8",
		"hostname": "",
		"ntpserver": "",
		"adapter": "",
		"username": "ubuntu",
		"password": "",
		"repoaddress": "192.168.1.1",
		"mirrordir": "/ubuntu",
		"sshkeypath": "/home/dan/.ssh/id_rsa.pub",
		"packages": "openssh-server iptables libltdl7"
	},
	"deployments": [
		{
			"mac": "00:50:56:a3:64:a2",
			"deployment": "preseed",
			"config": {
				"gateway": "192.168.1.1",
				[...]
			}
		}
	]
}
```

The *globalConfig* is the configuration that is inherited by any of the deployment configurations where that information has been omitted, typically a lot of networking information, keys or package information will be shared amongst deployments. 

### Online updates of deployment configuration
The webserver exposes a `/deployment` end point that can be used to provide an online update of the configuration, this has the following benefits:

- Allows automation of updates, through things like an API call
- Provides no-downtime, stopping and starting the server to load a new configuration can result in a broken installation as the network connection will be broken during restart

*Retrieve the existing configuration*

The currently active configuration can be retrieved through a simple get on the `/deployment` endpoint 

e.g.

`curl -vX <IP ADDRESS>/deployment`

*Updating the configuration*

The configuration can be updated by `POST`ing the configuration JSON to the same URL.

e.g.

`curl -vX POST deploy01/deployment -d @deployment.json --header "Content-Type: application/json"`

### Retreiving bootstrap components (now optional)

The plunder binary has an embedded `iPXE` bootloader, meaning that nothing other than the `plunder` binary is needed in order to deploy servers. However if a newer version is required then the `./plunder get` command will download the `iPXE` bootloader that is needed by the `TFTP` service in order to bootstrap the OS build. 

### Networking configuration
In a mutli-adapter host (recommended) ensure that the correct adapter is used i.e. `eth0 -> eth1` also that the correct IP addresses are used, the IP addresses are needed for the TFTP and HTTP stages as that is all managed through TCP connections. 

*NOTE* the `startDHCP` field is _required_ and should ideally be `addressDHCP` + 1, the `leasePoolDHCP` will then manage a pool of IP addresses from that start address.

### Services configuration
By default, all services are disabled (this is to make sure you look at the configuration before advertising random network services). The `plunder` application will not start if all services are disabled and will present you with a warning message, in order to enable services change `false` to `true`.

### Kernel / Initrd etc.. 

The plan is to have `plunder` mount and extract the correct kernels and netboot `initrd` (TBD)

# Usage

Once the configuration file has been updated the `./plunder server` command will start the required services as shown below:

```
sudo ./plunder server --config ./config.json --deployment ./deployment.json --logLevel 5
[sudo] password for dan: 
INFO[0000] Reading configuration from [./config.json]   
INFO[0000] Starting Remote Boot Services, press CTRL + c to stop 
DEBU[0000] 
Server IP:	192.168.1.1
Adapter:	ens192
Start Address:	192.168.1.2
Pool Size:	100
 
INFO[0000] RemoteBoot => Starting DHCP                  
INFO[0000] RemoteBoot => Starting TFTP                  
DEBU[0000] 
Server IP:	192.168.1.1
PXEFile:	undionly.kpxe
 
INFO[0000] Opening and caching undionly.kpxe            
INFO[0000] RemoteBoot => Starting HTTP                  
INFO[0286] DCHP Message: Discover   
```

# Troubleshooting

PXE booting provides very little feedback when things aren't working, but usually the hand-off is why things wont work i.e. `DHCP` -> `TFTP` boot. Logs from `plunder` should show the hand-off from the CLI.

# Roadmap

Fix the templating of Preseed and Kickstart files, automate the entire process end-to-end. May have `plunder` keep all configurations internally and use http handlers to expose them as urls to the boot loader (TBD)

TL;DR make better.

Created on 2018-11-14 17:30:01
