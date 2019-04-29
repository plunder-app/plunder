# Deployment Configuration

## Generating a configuration
A `./plunder config deployment > deployment.json` will create a blank deployment configuration that can be pre-populated in order to create specific deployments.

A configured deployment should resemble something like the example below:

```yaml
{
	"globalConfig": {
		"gateway": "",
		"address": "",
		"subnet": "",
		"nameserver": "",
		"hostname": "",
		"ntpserver": "",
		"adapter": "",
		"swapEnabled": false,
		"username": "",
		"password": "",
		"repoaddress": "",
		"mirrordir": "",
		"sshkeypath": "",
		"packages": ""
	},
	"deployments": [
		{
			"mac": "",
			"kernelPath": "",
			"initrdPath": "",
			"cmdline": "",
			"deployment": "",
			"config": {
				"gateway": "",
				"address": "",
				"subnet": "",
				"nameserver": "",
				"hostname": "",
				"ntpserver": "",
				"adapter": "",
				"swapEnabled": false,
				"username": "",
				"password": "",
				"repoaddress": "",
				"mirrordir": "",
				"sshkeypath": "",
				"packages": ""
			}
		}
	]
}
```

## Configuration overview

The *globalConfig* is the configuration that is inherited by any of the deployment configurations where that information has been omitted, typically a lot of networking information, keys or package information will be shared amongst deployments. 

Placing the same information into an actual deployment will **override** the configuration inherited from the `globalConfig`.

### Shared Configuration overview

- `gateway` - The gateway a server will be configured to use as default router
- `subnet` - The network range server will be configured to use
- `nameserver` - DNS server to resolve hostnames
- `ntpserver` - The address of a timeserver
- `adapter` - Which specific adapter will be configured
- `swapEnabled` - Build the Operating system without swap being created
- `username` - A default user that will be created
- `password` - A password for the above user
- `repoaddress` - The hostname/ip address of the server where the OS packages reside
- `sshkeypath` - The path to an ssh key that will be added to the image for authenticating



### Deployment specific

- `address` - A unique network address that will be added to the server
- `hostname` - A unique hostname to be added to the provisioned server



As mentioned above, a lot of fields can be ignored and the entry from the `globalConfig` will be used.



### Deployments

The deployment contains things that will make a server unique!

- `mac` - The unqique HW mac address of a server to configure

- `kernelPath` - If a specific kernel should be used (for things like LinuxKit)

- `initrdPath` - If a specific init ramdisk should be used

- `cmdline` - Any arguments that should be passed to the kernel ramdisk

  

The `deployment` specifies how the server will be provisioned, there are three options:

- `preseed` Ubuntu/Debian pressed deployment
- `kickstart` CentOS/RHEL deployment
- `reboot` This is for servers that need to be kept on a reboot loop.



The remaining `config` allows updates or overrides to the global confgiguration detailed above.

 

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

## Usage

With configuration for both the services and the deployments completed, they can both be passed to `plunder` in order for servers to be built.

As shown below:

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

## Next Steps
Servers that have their mac addresses in the `deployment` file will be passed the correct bootloader and they will ultimately be provisioned with the networking information as part of the configuration, they also will be provisioned with the credentials and specified ssh key. 

For provisioning applications or a platform details are [here](./provisioning.md).