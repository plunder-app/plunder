package cmd

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/plunder-app/plunder/pkg/parlay"
	"github.com/plunder-app/plunder/pkg/parlay/types"
	"github.com/plunder-app/plunder/pkg/services"
	"github.com/plunder-app/plunder/pkg/utils"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// These variables are used to capture input from the CLI
var output, detectNic string
var pretty bool

func init() {
	plunderCmd.AddCommand(plunderConfig)
	plunderConfig.PersistentFlags().StringVarP(&output, "output", "o", "json", "Ouput type, should be either JSON or YAML")
	plunderConfig.PersistentFlags().BoolVarP(&pretty, "pretty", "p", false, "Ouput JSON in a pretty/Human readable format")
	plunderServerConfig.Flags().StringVarP(&detectNic, "detect", "d", "", "Attempt to automatically detect the default settings for a specfic interface")

	plunderConfig.AddCommand(plunderServerConfig)
	plunderConfig.AddCommand(plunderDeploymentConfig)
	plunderConfig.AddCommand(PlunderParlayConfig)

	plunderCmd.AddCommand(plunderGet)

}

// PlunderConfig - This is for intialising a blank or partial configuration
var plunderConfig = &cobra.Command{
	Use:   "config",
	Short: "Initialise a plunder configuration",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		cmd.Help()
		return
	},
}

// PlunderServerConfig - This is for intialising a blank or partial configuration
var plunderServerConfig = &cobra.Command{
	Use:   "server",
	Short: "Initialise a plunder configuration",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		// Indent (or pretty-print) the configuration output
		bc := &services.BootConfig{
			Kernel:     "/kernelPath",
			Initrd:     "/initPath",
			Cmdline:    "cmd=options",
			ConfigName: "default",
		}

		detectServerConfig()

		services.Controller.BootConfigs = append(services.Controller.BootConfigs, *bc)
		err := renderOutput(services.Controller, pretty)
		if err != nil {
			log.Fatalf("%v", err)
		}
		return
	},
}

// PlunderDeploymentConfig - This is for intialising a blank or partial configuration
var plunderDeploymentConfig = &cobra.Command{
	Use:   "deployment",
	Short: "Initialise a server configuration",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		// Create an example Global configuration
		globalConfig := services.HostConfig{
			Gateway:           "192.168.0.1",
			NTPServer:         "192.168.0.1",
			NameServer:        "192.168.0.1",
			Adapter:           "ens192",
			Subnet:            "255.255.255.0",
			Username:          "user",
			Password:          "pass",
			Packages:          "nginx openssh-server",
			RepositoryAddress: "192.168.0.1",
			MirrorDirectory:   "/ubuntu",
			SSHKeyPath:        "/home/deploy/.ssh/id_pub.rsa",
			SSHKey:            "ssh-rsa AABBCCDDEE1122334455",
		}

		// Create an example Host configuration
		hostConfig := services.HostConfig{
			IPAddress:  "192.168.0.2",
			ServerName: "Server01",
		}
		hostDeployConfig := services.DeploymentConfig{
			MAC:        "00:11:22:33:44:55",
			ConfigHost: hostConfig,
			ConfigName: "default",
		}

		configuration := &services.DeploymentConfigurationFile{
			GlobalServerConfig: globalConfig,
		}

		configuration.Configs = append(configuration.Configs, hostDeployConfig)
		// Indent (or pretty-print) the configuration output
		err := renderOutput(configuration, pretty)
		if err != nil {
			log.Fatalf("%v", err)
		}
		return
	},
}

// PlunderParlayConfig - This is for intialising a parlay deployment
var PlunderParlayConfig = &cobra.Command{
	Use:   "parlay",
	Short: "Initialise a parlay configuration",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		parlayActionPackage := types.Action{
			Name:         "Add package",
			ActionType:   "pkg",
			PkgManager:   "apt",
			PkgOperation: "install",
			Packages:     "mysql",
		}

		parlayActionCommand := types.Action{
			Name:             "Run Command",
			ActionType:       "command",
			Command:          "which uptime",
			CommandSudo:      "root",
			CommandSaveAsKey: "cmdKey",
		}
		parlayActionUpload := types.Action{
			Name:        "Upload File",
			ActionType:  "upload",
			Source:      "./my_file",
			Destination: "/tmp/file",
		}

		parlayActionDownload := types.Action{
			Name:        "Download File",
			ActionType:  "download",
			Destination: "./my_file",
			Source:      "/tmp/file",
		}

		parlayActionKey := types.Action{
			Name:       "Execute key",
			ActionType: "command",
			KeyName:    "cmdKey",
		}

		parlayDeployment := parlay.Deployment{
			Name:  "Install MySQL",
			Hosts: []string{"192.168.0.1", "192.168.0.2"},
		}

		parlayDeployment.Actions = append(parlayDeployment.Actions, parlayActionPackage)
		parlayDeployment.Actions = append(parlayDeployment.Actions, parlayActionCommand)
		parlayDeployment.Actions = append(parlayDeployment.Actions, parlayActionUpload)
		parlayDeployment.Actions = append(parlayDeployment.Actions, parlayActionDownload)
		parlayDeployment.Actions = append(parlayDeployment.Actions, parlayActionKey)

		parlayConfig := &parlay.TreasureMap{}
		parlayConfig.Deployments = []parlay.Deployment{}
		parlayConfig.Deployments = append(parlayConfig.Deployments, parlayDeployment)

		// Render the output to screen
		err := renderOutput(parlayConfig, pretty)
		if err != nil {
			log.Fatalf("%v", err)
		}
		return
	},
}

// plunderGet - The Get command will pull any required components (iPXE boot files)
var plunderGet = &cobra.Command{
	Use:   "get",
	Short: "Get any components needed for bootstrapping (internet access required)",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		err := utils.PullPXEBooter()
		if err != nil {
			log.Fatalf("%v", err)
		}
		return
	},
}

func renderOutput(data interface{}, pretty bool) error {
	var d []byte
	var err error
	switch strings.ToLower(output) {
	case "yaml":
		d, err = yaml.Marshal(data)
	case "json":
		if pretty {
			d, err = json.MarshalIndent(data, "", "\t")
		} else {
			d, err = json.Marshal(data)
		}
	default:
		return fmt.Errorf("Unknown output type [%s]", output)
	}
	if err != nil {
		return err
	}
	// Print out the output to STDOUT
	fmt.Printf("%s\n", d)
	return nil
}

func detectServerConfig() error {

	// Find an example nic to use, that isn't the loopback address
	nicName, nicAddr, err := utils.FindIPAddress(detectNic)
	if err != nil {
		return err
	}

	// Attempt to parse th returned IP address and apply simple incrementation to determin DHCP start range
	ip := net.ParseIP(nicAddr)
	ip = ip.To4()
	if ip == nil {
		return fmt.Errorf("error parsing IP address of adapter [%s]", detectNic)
	}
	ip[3]++

	// Prepopulate the flags with the found nic information
	services.Controller.AdapterName = &nicName
	services.Controller.HTTPAddress = &nicAddr
	services.Controller.TFTPAddress = &nicAddr

	*services.Controller.PXEFileName = "undionly.kpxe"

	// DHCP Settings
	services.Controller.DHCPConfig.DHCPAddress = &nicAddr
	services.Controller.DHCPConfig.DHCPGateway = &nicAddr
	services.Controller.DHCPConfig.DHCPDNS = &nicAddr
	*services.Controller.DHCPConfig.DHCPLeasePool = 20
	*services.Controller.DHCPConfig.DHCPStartAddress = ip.String()

	return nil
}
