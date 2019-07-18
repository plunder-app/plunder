# Application Architecture

The purpose of this document is to outline the architecture of the plunder program itself, as it has predominantly been developed by a single developer the logic sometimes is hard to fathom (or to understand after a period of absence)

## Application Server routine

When starting plunder as a server for deployment a number of files are parsed and internal structures populated, below is a step through of the actions that take place. 

### Starting the server (HTTP Enabled)

We will start `plunder` with a *default* json configuration, with the services enabled and pointing to a default ubuntu kernel/initrd. The deployment file has a single server defined in it.

`plunder server --config ./config.json --deployment ./deployment.json`

1. Plunder starts
   - parses flags
   - parses global `config.json` 
2. Plunder will start services enabled in the configuration `controller.StartServices(deployment)` (`cmd\server.go`)
3. If a deployment file is passed then it should be parsed `err := UpdateConfiguration(deployment)` (`pkg\server\services.go`)
   - The parsing of this will generate strings that are mapped to urls that are tracked in a map `httpPaths`
   - The function `UpdateConfiguration(configFile []byte)` (`pkg/server/generator.go`) will generate these in memory by iterating through the file and checking the deployment type.
4. HTTP Server is started with `err := c.serveHTTP()`(`pkg\server\services.go`)
5. This function will create a number of prebuilt PXE boot strings using the kernels etc. from `config.json`, configurations such as `/preboot.ipxe` etc.
6. In the event that new configuration is passed to the server then steps 3 are ran again.

### Client connections

1. A Host starts and proceeds to PXE boot, by doing a DHCP request.
2. The DHCP server defaults to point the `BootFileName` Option `dhcp.OptionBootFileName:[]byte(*c.PXEFileName)`(`pkg/server/services.go`), also checks for the dhcp option `77` saying `iPXE` (which will be false)
3. This is passed of TFTP to the booting host which will start iPXE and re-do a DHCP request
4. This time however the DHCP client will have the option `77` set to `iPXE` which means that it's ready for provisioning. 
5. The DHCP server will look for an existing configuration `deploymentType := FindDeployment(mac)` (`pkg/server/serve_dhcp.go`), which should return `preseed` etc.
6. The DHCP server will then look to see if a specific configuration has been created with `if httpPaths[fmt.Sprintf("%s.ipxe", dashMac)] == ""` (`pkg/server/serve_dhcp.go`), if not it will default to a deployment type
7. If there exists a pre-defined configuration then it will set the DHCP option to that.
