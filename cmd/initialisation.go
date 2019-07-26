package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/plunder-app/plunder/pkg/parlay"
	"github.com/plunder-app/plunder/pkg/parlay/types"
	"github.com/plunder-app/plunder/pkg/server"
	"github.com/plunder-app/plunder/pkg/utils"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

func init() {
	plunderCmd.AddCommand(plunderConfig)
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
		bc := &server.BootConfig{
			Kernel:     "/kernelPath",
			Initrd:     "/initPath",
			Cmdline:    "cmd=options",
			ConfigName: "default",
		}
		server.Controller.BootConfigs = append(server.Controller.BootConfigs, *bc)
		b, err := json.MarshalIndent(server.Controller, "", "\t")
		if err != nil {
			log.Fatalf("%v", err)
		}
		fmt.Printf("\n%s\n", b)
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
		globalConfig := server.HostConfig{
			Gateway:           "192.168.0.1",
			NTPServer:         "192.168.0.1",
			NameServer:        "192.168.0.1",
			Adapter:           "ens192",
			Subnet:            "255.255.255.0",
			Username:          "user",
			Password:          "pass",
			Packages:          "nginx",
			RepositoryAddress: "192.168.0.1",
			MirrorDirectory:   "/ubuntu",
			SSHKeyPath:        "/home/deploy/.ssh/id_pub.rsa",
			SSHKey:            "ssh-rsa AABBCCDDEE1122334455",
		}

		// Create an example Host configuration
		hostConfig := server.HostConfig{
			IPAddress:  "192.168.0.2",
			ServerName: "Server01",
		}
		hostDeployConfig := server.DeploymentConfig{
			MAC:        "00:11:22:33:44:55",
			ConfigHost: hostConfig,
			ConfigName: "default",
		}

		configuration := &server.DeploymentConfigurationFile{
			GlobalServerConfig: globalConfig,
		}

		configuration.Configs = append(configuration.Configs, hostDeployConfig)
		// Indent (or pretty-print) the configuration output
		b, err := json.MarshalIndent(configuration, "", "\t")
		if err != nil {
			log.Fatalf("%v", err)
		}
		fmt.Printf("\n%s\n", b)
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

		// Indent (or pretty-print) the configuration output
		b, err := json.MarshalIndent(parlayConfig, "", "\t")
		if err != nil {
			log.Fatalf("%v", err)
		}
		fmt.Printf("\n%s\n", b)
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
