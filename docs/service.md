# Service Configuration

## Generating a configuration

A `./plunder config server > config.json` will look at the network configuration of the current machine and build a default configuration file (in json). This file will need opening in your favourite text editor an modifying to ensure that `plunder` works correctly. 

### Modifying the configuration

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
*Example generated configuration above*

By **default** the configuration that is generated will have all of the services disabled (dhcp/tftp/http) and attempting to start plunder will result in an error message saying that no services are being started. 

The `addressTFTP` and `addressHTTP` are still required to be set even if you're not enabling the service, this is because those values will be passed through `DHCP` to a server that is being bootstrapped. So if `TFTP` or `HTTP` services already exist on your network, then modify those values accordingly.

The `pxePath` should point to an iPXE bootloader if needed, however if the file doesn't exist or if the option is blank then `plunder` will fall back to an embedded bootloader. 

The `kernelPath` and `initrdPath` should point to a kernel and init ramdisk on the local filesystem that will be passed to the server once the bootloader has finished.

The `startDHCP` should typically be `addressDHCP` +1 and then the `leasePoolDHCP` defines how many free addresses will be allocated sequentially from the start address.

## Usage
At this point you can start various services and you'll see servers on the network requesting `DHCP` addresses etc.. however in order to do anything we will need to configure the [deployment](./deployment.md).