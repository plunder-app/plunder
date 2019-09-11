package cmd

import (
	"io/ioutil"
	"os"

	"github.com/plunder-app/plunder/pkg/apiserver"
	"github.com/plunder-app/plunder/pkg/parlay"
	"github.com/plunder-app/plunder/pkg/services"
	"github.com/plunder-app/plunder/pkg/utils"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

//var controller server.BootController
var dhcpSettings services.DHCPSettings

var apiServerPath, gateway, dns, startAddress, configPath, deploymentPath, defaultKernel, defaultInitrd, defaultCmdLine *string

var leasecount, port *int

var anyboot, insecure *bool

func init() {

	// Prepopulate the flags with the found nic information
	services.Controller.AdapterName = PlunderServer.Flags().String("adapter", "", "Name of adapter to use e.g eth0, en0")

	services.Controller.HTTPAddress = PlunderServer.Flags().String("addressHTTP", "", "Address of HTTP to use, if blank will default to [addressDHCP]")
	services.Controller.TFTPAddress = PlunderServer.Flags().String("addressTFTP", "", "Address of TFTP to use, if blank will default to [addressDHCP]")

	services.Controller.EnableDHCP = PlunderServer.Flags().Bool("enableDHCP", false, "Enable the DCHP Server")
	services.Controller.EnableTFTP = PlunderServer.Flags().Bool("enableTFTP", false, "Enable the TFTP Server")
	services.Controller.EnableHTTP = PlunderServer.Flags().Bool("enableHTTP", false, "Enable the HTTP Server")

	services.Controller.PXEFileName = PlunderServer.Flags().String("iPXEPath", "undionly.kpxe", "Path to an iPXE bootloader")

	// DHCP Settings
	services.Controller.DHCPConfig.DHCPAddress = PlunderServer.Flags().String("addressDHCP", "", "Address to advertise leases from, ideally will be the IP address of --adapter")
	services.Controller.DHCPConfig.DHCPGateway = PlunderServer.Flags().String("gateway", "", "Address of Gateway to use, if blank will default to [addressDHCP]")
	services.Controller.DHCPConfig.DHCPDNS = PlunderServer.Flags().String("dns", "", "Address of DNS to use, if blank will default to [addressDHCP]")
	services.Controller.DHCPConfig.DHCPLeasePool = PlunderServer.Flags().Int("leasecount", 20, "Amount of leases to advertise")
	services.Controller.DHCPConfig.DHCPStartAddress = PlunderServer.Flags().String("startAddress", "", "Start advertised address [REQUIRED]")

	//HTTP Settings
	defaultKernel = PlunderServer.Flags().String("kernel", "", "Path to a kernel to set as the *default* kernel")
	defaultInitrd = PlunderServer.Flags().String("initrd", "", "Path to a ramdisk to set as the *default* ramdisk")
	defaultKernel = PlunderServer.Flags().String("cmdline", "", "Additional command line to pass to the *default* kernel")

	// Config File
	configPath = PlunderServer.Flags().String("config", "", "Path to a plunder server configuration")
	deploymentPath = PlunderServer.Flags().String("deployment", "", "Path to a plunder deployment configuration")
	PlunderServer.Flags().StringVar(&services.DefaultBootType, "defaultBoot", "", "In the event a boot type can't be found default to this")

	// API Server configuration
	port = PlunderServer.Flags().IntP("port", "p", 60443, "Port that the Plunder API server will listen on")
	insecure = PlunderServer.Flags().BoolP("insecure", "i", false, "Start the Plunder API server without encryption")
	apiServerPath = PlunderServer.Flags().String("path", ".plunderserver.yaml", "Path to configuration for the API Server")

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
			// if *anyboot == true {
			// 	log.Errorf("AnyBoot has been enabled, all configuration will be ignored")
			// }
			log.Infof("Reading deployment configuration from [%s]", *deploymentPath)
			if _, err := os.Stat(*deploymentPath); !os.IsNotExist(err) {
				deployment, err = ioutil.ReadFile(*deploymentPath)
				if err != nil {
					log.Fatalf("%v", err)
				}
			}
		}

		// if *anyboot == true {
		// 	services.AnyBoot = true
		// }

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
				err = services.ParseControllerData(configFile)
				if err != nil {
					log.Fatalf("%v", err)
				}

			} else {
				log.Fatalf("Unable to open [%s]", *configPath)
			}
		}

		if *services.Controller.EnableDHCP == false && *services.Controller.EnableHTTP == false && *services.Controller.EnableTFTP == false {
			log.Warnln("All services are currently disabled")
		}

		// If we've enabled DHCP, then we need to ensure a start address for the range is populated
		if *services.Controller.EnableDHCP && *services.Controller.DHCPConfig.DHCPStartAddress == "" {
			log.Fatalln("A DHCP Start address is required")
		}

		if *services.Controller.DHCPConfig.DHCPLeasePool == 0 {
			log.Fatalln("At least one available lease is required")
		}

		services.Controller.StartServices(deployment)

		// Run the API server in a seperate go routine
		go func() {
			err := apiserver.StartAPIServer(*apiServerPath, *port, *insecure)
			if err != nil {
				log.Fatalf("%v", err)
			}
		}()

		// Register the packages to the apiserver
		services.RegisterToAPIServer()
		parlay.RegisterToAPIServer()

		// Sit and wait for a control-C
		utils.WaitForCtrlC()

		return
	},
}
