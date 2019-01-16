package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/plunder/pkg/parlay"
	"github.com/thebsdbox/plunder/pkg/ssh"
)

var deploymentSSH, mapFile, mapFileValidate, logFile *string

// These flags are used to determine if a particular deployment, action and specific host need to be used.
var deploymentName, actionName, host *string

// This flag determines if a singular action should occur or wheter to resume all actions from this point
var resume *bool

func init() {
	// SSH Deployment flags
	logFile = plunderAutomateSSH.Flags().String("logfile", "", "Patht to where plunder will write automation logs")
	deploymentSSH = plunderAutomateSSH.Flags().String("deployconfig", "", "Path to a plunder deployment configuration")
	mapFile = plunderAutomateSSH.Flags().String("map", "", "Path to a plunder map")

	// Deployment control flags
	deploymentName = plunderAutomateSSH.Flags().String("deployment", "", "Automate a specific deployment")
	actionName = plunderAutomateSSH.Flags().String("action", "", "Automate a specific action")
	host = plunderAutomateSSH.Flags().String("host", "", "Automate the deployment for a specific host")
	resume = plunderAutomateSSH.Flags().Bool("resume", false, "Resume all actions after the one specified by --action")

	// Validation flags
	mapFileValidate = plunderAutomateValidate.Flags().String("map", "", "Path to a plunder map")

	// Automate SSH Flags
	plunderAutomate.AddCommand(plunderAutomateValidate)
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
		// log.Infof("%s", *deploymentSSH)
		if *deploymentSSH != "" {
			log.Infof("Reading deployment configuration from [%s]", *deploymentSSH)
			err := ssh.ImportHostsFromDeployment(*deploymentSSH)
			if err != nil {
				cmd.Help()
				log.Fatalf("%v", err)
			}
		} else {
			cmd.Help()
			log.Fatalf("No Deployment information imported")
		}
		log.Infof("Found [%d] ssh configurations", len(ssh.Hosts))

		if *mapFile != "" {
			log.Infof("Reading deployment configuration from [%s]", *mapFile)

			var deployment parlay.TreasureMap
			// // Check the actual path from the string
			if _, err := os.Stat(*mapFile); !os.IsNotExist(err) {
				b, err := ioutil.ReadFile(*mapFile)
				if err != nil {
					log.Fatalf("%v", err)
				}
				err = json.Unmarshal(b, &deployment)
				if err != nil {
					log.Fatalf("%v", err)
				}

				// If a specific deployment is being used then find it's details
				if *deploymentName != "" {
					log.Infof("Looking for deployment [%s]", *deploymentName)

					err = deployment.FindDeployment(*deploymentName, *actionName, *host, *logFile, *resume)
				} else {
					// Parse the entire deployment
					err = deployment.DeploySSH(*logFile)
				}
				if err != nil {
					log.Fatalf("%v", err)
				}
			} else {
				log.Fatalf("%v", err)
			}
		}

		return
	},
}

// plunderAutomateValidate
var plunderAutomateValidate = &cobra.Command{
	Use:   "validate",
	Short: "Validate a deployment map",
	Run: func(cmd *cobra.Command, args []string) {
		if *mapFileValidate != "" {
			log.Infof("Reading deployment configuration from [%s]", *mapFileValidate)
			//var err error
			var deployment parlay.TreasureMap
			// // Check the actual path from the string
			if _, err := os.Stat(*mapFileValidate); !os.IsNotExist(err) {
				b, err := ioutil.ReadFile(*mapFileValidate)
				if err != nil {
					log.Fatalf("%v", err)
				}
				err = json.Unmarshal(b, &deployment)
				if err != nil {
					log.Fatalf("%v", err)
				}
				deploymentCount := len(deployment.Deployments)
				if deploymentCount == 0 {
					log.Fatalf("Zero deployments have been found")
				}
				log.Infof("Validating [%d] deployments", deploymentCount)
				for x := range deployment.Deployments {
					actionCount := len(deployment.Deployments[x].Actions)
					if actionCount == 0 {
						log.Fatalf("Zero deployments have been found")
					}
					log.Infof("Validating [%d] actions", actionCount)
					for y := range deployment.Deployments[x].Actions {
						err := deployment.Deployments[x].Actions[y].ValidateAction()
						if err != nil {
							log.Warnf("Action [%s] Error [%v]", deployment.Deployments[x].Actions[y].Name, err)
						}
					}
				}
			} else {
				log.Fatalf("Unable to open [%s]", *mapFile)
			}
		} else {
			cmd.Help()
			log.Fatalln("No Deployment map specified")
		}
	},
}
