package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/thebsdbox/plunder/pkg/bootstraps"
	"github.com/thebsdbox/plunder/pkg/utils"

	log "github.com/Sirupsen/logrus"

	"github.com/spf13/cobra"
)

func init() {
	plunderCmd.AddCommand(plunderConfig)
	plunderConfig.AddCommand(plunderServerConfig)
	plunderConfig.AddCommand(plunderDeploymentConfig)

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
		b, err := json.MarshalIndent(controller, "", "\t")
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
		var configuration bootstraps.DeploymentConfigurationFile
		configuration.Deployments = append(configuration.Deployments, bootstraps.DeploymentConfigurations{})
		// Indent (or pretty-print) the configuration output
		b, err := json.MarshalIndent(configuration, "", "\t")
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
