# Service Configuration

## Generating a configuration

A `./plunder config server > config.json` will look at the network configuration of the current machine and build a default configuration file (in json). This file will need opening in your favourite text editor an modifying to ensure that `plunder` works correctly. 

### Modifying the configuration

```json
{
        "adapter": "en0",
        "enableDHCP": false,
        "dhcpConfig": {
                "addressDHCP": "192.168.0.142",
                "startDHCP": "192.168.0.143",
                "leasePoolDHCP": 20,
                "gatewayDHCP": "192.168.0.142",
                "nameserverDHCP": "192.168.0.142"
        },
        "enableTFTP": false,
        "addressTFTP": "192.168.0.142",
        "enableHTTP": false,
        "addressHTTP": "192.168.0.142",
        "pxePath": "undionly.kpxe",
        "bootConfigs": [
                {
                        "configName": "default",
                        "kernelPath": "/kernelPath",
                        "initrdPath": "/initPath",
                        "cmdline": "cmd=options",
                        "isoPrefix": "ubuntu",
                        "isoPath": "/path/to/iso"
                }
        ]
}
```

*Example generated configuration above*

### Sections

By **default** the configuration that is generated will have all of the services disabled (dhcp/tftp/http) and attempting to start plunder will result in an error message saying that no services are being started. 

#### Services

The `enable<service>` will ensure that a particular functionality is enabled within Plunder.

The `addressTFTP` and `addressHTTP` are still required to be set even if you're not enabling the service, this is because those values will be passed through `DHCP` to a server that is being bootstrapped. So if `TFTP` or `HTTP` services already exist on your network, then modify those values accordingly.

#### DHCP


The `dhcpConfig` section details all of the configuration for the running DHCP server such as the  `startDHCP` setting which should typically be `addressDHCP` +1 and then the `leasePoolDHCP` defines how many free addresses will be allocated sequentially from the start address.

#### Boot Configurations

The boot configurations are an array of configurations that define various remote booting configurations and are referenced via the `configName`.

The `kernelPath` and `initrdPath` should point to a kernel and init ramdisk on the local filesystem that will be passed to the server once the bootloader has finished.

Finally, the `isoPrefix` (determines the beginning and unique path to contents) and the `isoPath` allow OS installation content to be read from within an ISO file. 

e.g.

`plunderAddress/isoPrefix/path/to/file`

#### Additional

The `pxePath` should point to an iPXE bootloader if needed, however if the file doesn't exist or if the option is blank then `plunder` will fall back to an embedded bootloader. 

## Usage
At this point you can start various services and you'll see servers on the network requesting `DHCP` addresses etc.. however in order to do anything we will need to configure the [deployment](./deployment.md).