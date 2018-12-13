# Deployment Configuration

## Generating a configuration
A `./plunder config deployment > deployment.json` will create a blank deployment configuration that can be pre-populated in order to create specific deployments.

A configured deployment should resemble something like the example below:

```
{
	"globalConfig": {
		"gateway": "192.168.1.1",
		"subnet": "255.255.255.0",
		"nameserver": "8.8.8.8",
		"username": "ubuntu",
		"password": "secret",
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
				"address": "192.168.1.3",
				"hostname": "etcd01",
				[...]
			}
		}
	]
}
```

## Configuration overview

The *globalConfig* is the configuration that is inherited by any of the deployment configurations where that information has been omitted, typically a lot of networking information, keys or package information will be shared amongst deployments. 

Placing the same information into an actual deployment will **override** the configuration inherited from the `globalConfig`.

The `deployment` specifies how the server will be provisioned, there are three options:

- `preseed` Ubuntu/Debian pressed deployment
- `kickstart` CentOS/RHEL deployment
- `reboot` This is for servers that need to be kept on a reboot loop.

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