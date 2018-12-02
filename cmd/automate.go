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
		//ssh.ExecuteSingleCommand("pwd", ssh.Hosts[0], 1)
		ssh.Execute("hostname", ssh.Hosts, 10)
		log.Infof("%v", ssh.Hosts[0])
		err := ssh.Hosts[0].UploadFile("test", "/tmp/test")
		if err != nil {
			log.Fatalf("%v", err)
		}
		ssh.Execute("hostname", ssh.Hosts, 10)

		ssh.Hosts[0].DownloadFile("/tmp/test", "test1")
		ssh.Hosts[0].StopSession()

		return
	},
}
