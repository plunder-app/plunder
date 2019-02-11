package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/plunder/pkg/parlay"
	"github.com/thebsdbox/plunder/pkg/parlay/plugin"
	"github.com/thebsdbox/plunder/pkg/parlay/types"
	"github.com/thebsdbox/plunder/pkg/ssh"
)

// These flags are used to determine a deployment
var deploymentSSH, mapFile, logFile *string

// These flags are used to determine if a particular deployment, action and specific host need to be used.
var deploymentName, actionName, host *string

// These flags are used for management of plugins
var pluginPath, pluginAction, pluginActions *string

// This flag determines if a singular action should occur or wheter to resume all actions from this point
var resume *bool

// UI Json output only, when this is try the UI selections will just create the associated JSON
var jsonOutput *bool

func init() {

	// Global flags for automation
	logFile = plunderAutomate.PersistentFlags().String("logfile", "", "Path to where plunder will write automation logs")
	mapFile = plunderAutomate.PersistentFlags().String("map", "", "Path to a plunder map")

	// SSH Deployment flags
	deploymentSSH = plunderAutomate.PersistentFlags().String("deployconfig", "", "Path to a plunder deployment configuration")

	// Deployment control flags
	deploymentName = plunderAutomateSSH.Flags().String("deployment", "", "Automate a specific deployment")
	actionName = plunderAutomateSSH.Flags().String("action", "", "Automate a specific action")
	host = plunderAutomateSSH.Flags().String("host", "", "Automate the deployment for a specific host")
	resume = plunderAutomateSSH.Flags().Bool("resume", false, "Resume all actions after the one specified by --action")

	// Plugin Flags
	pluginPath = plunderAutomatePluginUsage.Flags().String("plugin", "", "Path to a specific plugin typically ~./plugin/[X].plugin")
	pluginAction = plunderAutomatePluginUsage.Flags().String("action", "", "Action to retrieve the usage of")
	pluginActions = plunderAutomatePluginActions.Flags().String("plugin", "", "Path to a specific plugin typically ~./plugin/[X].plugin")

	jsonOutput = plunderAutomateUI.Flags().Bool("json", false, "Print the JSON to stdout, no execution of commands")

	plunderAutomatePlugins.AddCommand(plunderAutomatePluginUsage)
	plunderAutomatePlugins.AddCommand(plunderAutomatePluginActions)
	plunderAutomatePlugins.AddCommand(plunderAutomatePluginTest)

	// Automate SSH Flags
	plunderAutomate.AddCommand(plunderAutomateValidate)
	plunderAutomate.AddCommand(plunderAutomateSSH)
	plunderAutomate.AddCommand(plunderAutomatePlugins)
	plunderAutomate.AddCommand(plunderAutomateUI)

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

// plunderAutomatePlugins
var plunderAutomatePlugins = &cobra.Command{
	Use:   "plugin",
	Short: "Automate the deployment of a platform/application",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		parlayplugin.ListPlugins()
		return
	},
}

// plunderAutomatePlugins
var plunderAutomatePluginUsage = &cobra.Command{
	Use:   "usage",
	Short: "Display the usage of a plugin action",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		parlayplugin.UsagePlugin(*pluginPath, *pluginAction)
		return
	},
}

// plunderAutomatePlugins
var plunderAutomatePluginActions = &cobra.Command{
	Use:   "actions",
	Short: "Display the actions of a particular plugin",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		parlayplugin.ListPluginActions(*pluginActions)
		return
	},
}

// plunderAutomatePlugins
var plunderAutomatePluginTest = &cobra.Command{
	Use:   "test",
	Short: "Test the actions of the example plugin",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		test := `{ "name": "Example of test action", "type": "exampleAction/test", "plugin": { "credentials": "AAABBBCCCCDDEEEE", "address": "172.0.0.1" }	}`
		var action types.Action
		_ = json.Unmarshal([]byte(test), &action)

		_, err := parlayplugin.ExecuteActionInPlugin("./plugin/example.plugin", "127.0.0.1", "example/test", action.Plugin)
		if err != nil {
			log.Fatalf("%v", err)
		}
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
		log.SetLevel(log.Level(logLevel))

		if *mapFile != "" {
			log.Infof("Reading deployment configuration from [%s]", *mapFile)
			//var err error
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
						err := parlay.ValidateAction(&deployment.Deployments[x].Actions[y])
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

// plunderAutomateUI
var plunderAutomateUI = &cobra.Command{
	Use:   "ui",
	Short: "Enable the user interface to manage a deployment",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		var newMap *parlay.TreasureMap
		if *mapFile != "" {
			log.Infof("Reading deployment configuration from [%s]", *mapFile)
			//var err error
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
				newMap, err = deployment.StartUI()
				if err != nil {
					log.Fatalf("%v", err)
				}

			}
		}

		if *jsonOutput == true {
			b, _ := json.MarshalIndent(newMap, "", "\t")
			fmt.Printf("%s\n", b)
			return
		}

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

		err := newMap.DeploySSH(*logFile)
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}
