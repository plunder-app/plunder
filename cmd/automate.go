package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/plunder/pkg/ssh"
)

var deploymentSSH *string

func init() {
	deploymentSSH = plunderAutomateSSH.Flags().String("deployment", "", "Path to a plunder deployment configuration")

	// Automate SSH Flags

	plunderAutomate.AddCommand(plunderAutomateSSH)
	plunderCmd.AddCommand(plunderAutomate)
}

// PlunderAutomate
var plunderAutomate = &cobra.Command{
	Use:   "automate",
	Short: "Automate the deployment of a platform/application",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		cmd.Help()
		return
	},
}

// plunderAutomateSSH
var plunderAutomateSSH = &cobra.Command{
	Use:   "ssh",
	Short: "Automate over ssh",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		// If deploymentPath is not blank then the flag has been used
		log.Infof("%s", *deploymentSSH)
		if *deploymentSSH != "" {
			log.Infof("Reading deployment configuration from [%s]", *deploymentSSH)
			err := ssh.ImportHostsFromDeployment(*deploymentSSH)
			if err != nil {
				log.Fatalf("%v", err)
			}
		}
		log.Infof("Found [%d] ssh configurations", len(ssh.Hosts))

		return
	},
}
