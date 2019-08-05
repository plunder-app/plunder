package services

import (
	"time"

	"github.com/plunder-app/plunder/pkg/utils"
	log "github.com/sirupsen/logrus"

	dhcp "github.com/krolaw/dhcp4"
	"github.com/krolaw/dhcp4/conn"
)

var dhcpServer = make(chan bool)
var dhcpError = make(chan error, 1)

var runningDHCP, runningTFTP, runningHTTP bool

// find BootConfig will look through a Boot controller for a booting configuration identified through a configuration name
func findBootConfigForDeployment(deployment DeploymentConfig) *BootConfig {

	// First check is to look inside the deployment configuration for a custom configuration
	if deployment.ConfigBoot.Kernel != "" && deployment.ConfigBoot.Initrd != "" {
		// A Custom Kernel and initrd are specified
		log.Debugf("The server [%s] has a custom bootConfig defined", deployment.MAC)
		return &deployment.ConfigBoot
	}

	// Second check is to find a matching controller configuration to adopt
	for i := range Controller.BootConfigs {
		if Controller.BootConfigs[i].ConfigName == deployment.ConfigName {
			// Set the specific deployment configuration to the controller config
			return &Controller.BootConfigs[i]
		}
	}

	// Either there is no custom kernel/initrd/cmdline or a bootconfig doesn't exist as part of the server configuration
	return nil
}

// find BootConfig will look through a Boot controller for a booting configuration identified through a configuration name
func findBootConfigForName(configName string) *BootConfig {

	// Find a matching controller configuration to return
	for i := range Controller.BootConfigs {
		if Controller.BootConfigs[i].ConfigName == configName {
			return &Controller.BootConfigs[i]
		}
	}

	// No configuration could be found
	return nil
}

// find BootConfig will look through a Boot controller for a booting configuration identified through a configuration name
func (c *BootController) setBootConfig(configName, kernel, initrd, cmdline string) {
	newConfig := &BootConfig{
		ConfigName: configName,
		Kernel:     kernel,
		Initrd:     initrd,
		Cmdline:    cmdline,
	}
	c.BootConfigs = append(c.BootConfigs, *newConfig)
}

// StartServices - This will start all of the enabled services
func (c *BootController) StartServices(deployment []byte) {
	log.Infof("Starting Remote Boot Services, press CTRL + c to stop")

	if *c.EnableDHCP == true {
		c.handler = &DHCPSettings{}
		c.handler.IP = utils.ConvertIP(*c.DHCPConfig.DHCPAddress)
		c.handler.Start = utils.ConvertIP(*c.DHCPConfig.DHCPStartAddress)

		c.handler.LeaseDuration = 2 * time.Hour //TODO, make time modifiable
		c.handler.LeaseRange = *c.DHCPConfig.DHCPLeasePool
		c.handler.Leases = make(map[int]lease, *c.DHCPConfig.DHCPLeasePool)

		c.handler.Options = dhcp.Options{
			dhcp.OptionSubnetMask:       []byte{255, 255, 255, 0},
			dhcp.OptionRouter:           []byte(utils.ConvertIP(*c.DHCPConfig.DHCPGateway)),
			dhcp.OptionDomainNameServer: []byte(utils.ConvertIP(*c.DHCPConfig.DHCPDNS)),
			dhcp.OptionBootFileName:     []byte(*c.PXEFileName),
		}

		log.Debugf("\nServer IP:\t%s\nAdapter:\t%s\nStart Address:\t%s\nPool Size:\t%d\n", *c.DHCPConfig.DHCPAddress, *c.AdapterName, *c.DHCPConfig.DHCPStartAddress, *c.DHCPConfig.DHCPLeasePool)
		log.Println("Plunder Services --> Starting DHCP")

		if runningDHCP == false {
			newConnection, err := conn.NewUDP4FilterListener(*c.AdapterName, ":67")
			if err != nil {
				log.Fatalf("%v", err)
			}
			go func() {
				//Close the connection when we're tidying up
				defer newConnection.Close()
				runningDHCP = true
				dhcpError <- dhcp.Serve(newConnection, c.handler)
				runningDHCP = false

			}()

			go func() {
				select {
				case <-dhcpError:
					log.Infof("%s\n", dhcpError)
				case <-dhcpServer:
					newConnection.Close()
				}
			}()
		}
	} else {
		log.Debugf("Stopping DHCP Server")
		if runningDHCP {
			dhcpServer <- true
			runningDHCP = false
		}

	}

	if *c.EnableTFTP == true {
		go func() {
			log.Println("Plunder Services --> Starting TFTP")
			log.Debugf("\nServer IP:\t%s\nPXEFile:\t%s\n", *c.TFTPAddress, *c.PXEFileName)

			err := c.serveTFTP()
			if err != nil {
				log.Fatalf("%v", err)
			}
		}()
	}

	if *c.EnableHTTP == true {
		if len(c.BootConfigs) == 0 {
			log.Warn("No Boot settings specified in configuration")
		}

		httpAddress = *c.HTTPAddress

		// If a Deployment file is set then update the configuration
		if len(deployment) != 0 {
			err := UpdateDeploymentConfig(deployment)
			if err != nil {
				// Don't quit on error as updated configuration can be uploaded through the API
				log.Errorf("%v", err)
			}
		}

		go func() {
			log.Println("Plunder Services --> Starting HTTP")
			err := c.serveHTTP()
			if err != nil {
				log.Fatalf("%v", err)
			}
		}()
	}

}

// EvaluateServiceConfig - Takes a new configuration and compares the current to new config and then implements the changes
func EvaluateServiceConfig(newConfig *BootController) error {

	return nil
}
