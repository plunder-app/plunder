package cmd

import (
	"io/ioutil"
	"net"
	"os"

	"github.com/plunder-app/plunder/pkg/utils"

	"github.com/plunder-app/plunder/pkg/server"

	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

//var controller server.BootController
var dhcpSettings server.DHCPSettings

var gateway, dns, startAddress, configPath, deploymentPath, defaultKernel, defaultInitrd, defaultCmdLine *string

var leasecount *int

var anyboot *bool

func init() {

	// Find an example nic to use, that isn't the loopback address
	nicName, nicAddr, err := utils.FindIPAddress("")
	if err != nil {
		log.Warnf("%v", err)
	}
	//
	ip := net.ParseIP(nicAddr)
	ip = ip.To4()
	if ip == nil {
		log.Fatalf("error parsing IP address of adapter [%s]", nicName)
	}

	ip[3]++

	// Prepopulate the flags with the found nic information
	server.Controller.AdapterName = PlunderServer.Flags().String("adapter", nicName, "Name of adapter to use e.g eth0, en0")

	server.Controller.HTTPAddress = PlunderServer.Flags().String("addressHTTP", nicAddr, "Address of HTTP to use, if blank will default to [addressDHCP]")
	server.Controller.TFTPAddress = PlunderServer.Flags().String("addressTFTP", nicAddr, "Address of TFTP to use, if blank will default to [addressDHCP]")

	server.Controller.EnableDHCP = PlunderServer.Flags().Bool("enableDHCP", false, "Enable the DCHP Server")
	server.Controller.EnableTFTP = PlunderServer.Flags().Bool("enableTFTP", false, "Enable the TFTP Server")
	server.Controller.EnableHTTP = PlunderServer.Flags().Bool("enableHTTP", false, "Enable the HTTP Server")

	server.Controller.PXEFileName = PlunderServer.Flags().String("iPXEPath", "undionly.kpxe", "Path to an iPXE bootloader")

	// DHCP Settings
	server.Controller.DHCPConfig.DHCPAddress = PlunderServer.Flags().String("addressDHCP", nicAddr, "Address to advertise leases from, ideally will be the IP address of --adapter")
	server.Controller.DHCPConfig.DHCPGateway = PlunderServer.Flags().String("gateway", nicAddr, "Address of Gateway to use, if blank will default to [addressDHCP]")
	server.Controller.DHCPConfig.DHCPDNS = PlunderServer.Flags().String("dns", nicAddr, "Address of DNS to use, if blank will default to [addressDHCP]")
	server.Controller.DHCPConfig.DHCPLeasePool = PlunderServer.Flags().Int("leasecount", 20, "Amount of leases to advertise")
	server.Controller.DHCPConfig.DHCPStartAddress = PlunderServer.Flags().String("startAddress", ip.String(), "Start advertised address [REQUIRED]")

	//HTTP Settings
	defaultKernel = PlunderServer.Flags().String("kernel", "", "Path to a kernel to set as the *default* kernel")
	defaultInitrd = PlunderServer.Flags().String("initrd", "", "Path to a ramdisk to set as the *default* ramdisk")
	defaultKernel = PlunderServer.Flags().String("cmdline", "", "Additional command line to pass to the *default* kernel")

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
		var deployment []byte
		// If deploymentPath is not blank then the flag has been used
		if *deploymentPath != "" {
			if *anyboot == true {
				log.Errorf("AnyBoot has been enabled, all configuration will be ignored")
			}
			log.Infof("Reading deployment configuration from [%s]", *deploymentPath)
			if _, err := os.Stat(*deploymentPath); !os.IsNotExist(err) {
				deployment, err = ioutil.ReadFile(*deploymentPath)
				if err != nil {
					log.Fatalf("%v", err)
				}
			}
		}

		if *anyboot == true {
			server.AnyBoot = true
		}

		// If configPath is not blank then the flag has been used
		if *configPath != "" {
			log.Infof("Reading configuration from [%s]", *configPath)

			// Check the actual path from the string
			if _, err := os.Stat(*configPath); !os.IsNotExist(err) {
				configFile, err := ioutil.ReadFile(*configPath)
				if err != nil {
					log.Fatalf("%v", err)
				}

				// Read the controller from either a yaml or json format
				err = server.ParseControllerFile(configFile)
				if err != nil {
					log.Fatalf("%v", err)
				}

			} else {
				log.Fatalf("Unable to open [%s]", *configPath)
			}
		}
		if *server.Controller.EnableDHCP == false && *server.Controller.EnableHTTP == false && *server.Controller.EnableTFTP == false {
			log.Fatalln("At least one service is required to be enabled")
		}

		// If we've enabled DHCP, then we need to ensure a start address for the range is populated
		if *server.Controller.EnableDHCP && *server.Controller.DHCPConfig.DHCPStartAddress == "" {
			log.Fatalln("A DHCP Start address is required")
		}

		if *server.Controller.DHCPConfig.DHCPLeasePool == 0 {
			log.Fatalln("At least one available lease is required")
		}

		server.Controller.StartServices(deployment)
		return
	},
}
