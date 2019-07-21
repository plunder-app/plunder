package server

import (
	"time"

	"github.com/plunder-app/plunder/pkg/utils"
	log "github.com/sirupsen/logrus"

	dhcp "github.com/krolaw/dhcp4"
	"github.com/krolaw/dhcp4/conn"
)

// find BootConfig will look through a Boot controller for a booting configuration identified through a configuration name
func (c *BootController) findBootConfig(configName string) *bootConfig {
	for i := range c.BootConfigs {
		if *c.BootConfigs[i].ConfigName == configName {
			return &c.BootConfigs[i]
		}
	}

	// No configuration could be found
	return nil
}

// find BootConfig will look through a Boot controller for a booting configuration identified through a configuration name
func (c *BootController) setBootConfig(configName, kernel, initrd, cmdline *string) {
	newConfig := &bootConfig{
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

		go func() {
			log.Println("RemoteBoot => Starting DHCP")

			newConnection, err := conn.NewUDP4FilterListener(*c.AdapterName, ":67")
			if err != nil {
				log.Fatalf("%v", err)
			}
			//Close the connection when we're tidying up
			defer newConnection.Close()
			err = dhcp.Serve(newConnection, c.handler)
			log.Fatalf("%v", err)
		}()
	}

	if *c.EnableTFTP == true {
		go func() {
			log.Println("RemoteBoot => Starting TFTP")
			log.Debugf("\nServer IP:\t%s\nPXEFile:\t%s\n", *c.TFTPAddress, *c.PXEFileName)

			err := c.serveTFTP()
			log.Fatalf("%v", err)
		}()
	}

	if *c.EnableHTTP == true {
		if len(c.BootConfigs) == 0 {
			log.Warn("No Kernel specified in configuration")
		}
		// if c.Kernel == nil || *c.Kernel == "" {
		// 	log.Warn("No Kernel specified in configuration")
		// }
		// if c.Kernel == nil || *c.Initrd == "" {
		// 	log.Warn("No Initrd specified in configuration")
		// }

		httpAddress = *c.HTTPAddress
		httpPaths = make(map[string]string)

		// If a Deployment file is set then update the configuration
		if len(deployment) != 0 {
			err := UpdateConfiguration(deployment)
			if err != nil {
				log.Fatalf("%v", err)
			}
		}

		go func() {
			log.Println("RemoteBoot => Starting HTTP")
			err := c.serveHTTP()
			log.Fatalf("%v", err)
		}()
	}

	utils.WaitForCtrlC()
}
