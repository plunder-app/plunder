package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/plunder-app/plunder/pkg/bootstraps"

	"github.com/plunder-app/plunder/pkg/utils"

	"github.com/plunder-app/plunder/pkg/server"

	"github.com/spf13/cobra"

	log "github.com/Sirupsen/logrus"
)

var controller server.BootController
var dhcpSettings server.DHCPSettings

var gateway, dns, startAddress, configPath, deploymentPath *string

var leasecount *int

var anyboot *bool

func init() {

	// Find an example nic to use, that isn't the loopback address
	nicName, nicAddr, err := utils.FindIPAddress("")
	if err != nil {
		log.Warnf("%v", err)
	}

	// Prepopulate the flags with the found nic information
	controller.AdapterName = PlunderServer.Flags().String("adapter", nicName, "Name of adapter to use e.g eth0, en0")

	controller.DHCPAddress = PlunderServer.Flags().String("addressDHCP", nicAddr, "Address to advertise leases from, ideally will be the IP address of --adapter")
	controller.HTTPAddress = PlunderServer.Flags().String("addressHTTP", nicAddr, "Address of HTTP to use, if blank will default to [addressDHCP]")
	controller.TFTPAddress = PlunderServer.Flags().String("addressTFTP", nicAddr, "Address of TFTP to use, if blank will default to [addressDHCP]")

	controller.EnableDHCP = PlunderServer.Flags().Bool("enableDHCP", false, "Enable the DCHP Server")
	controller.EnableTFTP = PlunderServer.Flags().Bool("enableTFTP", false, "Enable the TFTP Server")
	controller.EnableHTTP = PlunderServer.Flags().Bool("enableHTTP", false, "Enable the HTTP Server")

	controller.PXEFileName = PlunderServer.Flags().String("iPXEPath", "undionly.kpxe", "Path to an iPXE bootloader")

	// DHCP Settings
	controller.DHCPGateway = PlunderServer.Flags().String("gateway", nicAddr, "Address of Gateway to use, if blank will default to [addressDHCP]")
	controller.DHCPDNS = PlunderServer.Flags().String("dns", nicAddr, "Address of DNS to use, if blank will default to [addressDHCP]")
	controller.DHCPLeasePool = PlunderServer.Flags().Int("leasecount", 20, "Amount of leases to advertise")
	controller.DHCPStartAddress = PlunderServer.Flags().String("startAddress", "", "Start advertised address [REQUIRED]")

	//HTTP Settings
	controller.Kernel = PlunderServer.Flags().String("kernel", "", "Path to a kernel to boot from")
	controller.Initrd = PlunderServer.Flags().String("initrd", "", "Path to an initrd to boot from")
	controller.Cmdline = PlunderServer.Flags().String("cmdline", "", "Additional command line to pass to the kernel")

	// Config File
	configPath = PlunderServer.Flags().String("config", "", "Path to a plunder server configuration")
	deploymentPath = PlunderServer.Flags().String("deployment", "", "Path to a plunder deployment configuration")
	anyboot = PlunderServer.Flags().Bool("anyboot", false, "Should be used without a configuration, this will boot the kernel/initrd")
	plunderCmd.AddCommand(PlunderServer)
}

// PlunderServer - This is for intialising a blank or partial configuration
var PlunderServer = &cobra.Command{
	Use:   "server",
	Short: "Start Plunder Services",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		// If deploymentPath is not blank then the flag has been used
		if *deploymentPath != "" {
			if *anyboot == true {
				log.Errorf("AnyBoot has been enabled, all configuration will be ignored")
			}
			log.Infof("Reading deployment configuration from [%s]", *deploymentPath)
			if _, err := os.Stat(*deploymentPath); !os.IsNotExist(err) {
				deployment, err := ioutil.ReadFile(*deploymentPath)
				if err != nil {
					log.Fatalf("%v", err)
				}
				err = bootstraps.UpdateConfiguration(deployment)
				if err != nil {
					log.Fatalf("%v", err)
				}
			}
		}

		if *anyboot == true {
			bootstraps.AnyBoot = true
		}

		// If configPath is not blank then the flag has been used
		if *configPath != "" {
			log.Infof("Reading server configuration from [%s]", *configPath)

			// Check the actual path from the string
			if _, err := os.Stat(*configPath); !os.IsNotExist(err) {
				configFile, err := ioutil.ReadFile(*configPath)
				if err != nil {
					log.Fatalf("%v", err)
				}
				json.Unmarshal(configFile, &controller)
			} else {
				log.Fatalf("Unable to open [%s]", *configPath)
			}
		}

		if *controller.EnableDHCP == false && *controller.EnableHTTP == false && *controller.EnableTFTP == false {
			log.Fatalln("At least one service is required to be enabled")
		}

		// If we've enabled DHCP, then we need to ensure a start address for the range is populated
		if *controller.EnableDHCP && *controller.DHCPStartAddress == "" {
			log.Fatalln("A DHCP Start address is required")
		}

		if *controller.DHCPLeasePool == 0 {
			log.Fatalln("At least one available lease is required")
		}

		controller.StartServices()
		return
	},
}
