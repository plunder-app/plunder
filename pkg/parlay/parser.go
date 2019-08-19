package parlay

import (
	"fmt"
	"time"

	"github.com/plunder-app/plunder/pkg/plunderlogging"

	log "github.com/sirupsen/logrus"

	parlayplugin "github.com/plunder-app/plunder/pkg/parlay/plugin"
	"github.com/plunder-app/plunder/pkg/parlay/types"
	"github.com/plunder-app/plunder/pkg/ssh"
)

// This logger will manage all of the logging for Parlay
var logger plunderlogging.Logger

// GetTargetLogs will retrieve the JSON logs
func GetTargetLogs(target string) (*plunderlogging.JSONLog, error) {
	return logger.GetJSONLogs(target)
}

// DeploySSH - will iterate through a deployment and perform the relevant actions
func (m *TreasureMap) DeploySSH(logFile string, jsonLogging bool) error {

	if logFile != "" {
		//enable logging
		logger.InitLogFile(logFile)

	}
	if jsonLogging {
		logger.InitJSON()
	}
	defer logger.SetLoggingState("", "Finished")

	if len(ssh.Hosts) == 0 {
		log.Warnln("No hosts credentials have been loaded, only commands with commandLocal = true will work")
	}
	if len(m.Deployments) == 0 {
		return fmt.Errorf("No Deployments in parlay map")
	}
	for x := range m.Deployments {
		// Build new hosts list from imported SSH servers and compare that we have required credentials
		hosts, err := ssh.FindHosts(m.Deployments[x].Hosts)
		if err != nil {
			return err
		}

		// Beggining of deployment work
		log.Infof("Beginning Deployment [%s]\n", m.Deployments[x].Name)
		logger.WriteLogEntry("", "", "", fmt.Sprintf("[%s] Beginning Deployment [%s]\n", time.Now().Format(time.ANSIC), m.Deployments[x].Name))

		// Set Restore checkpoint
		restore.Deployment = m.Deployments[x].Name
		restore.Hosts = m.Deployments[x].Hosts

		if m.Deployments[x].Parallel == true {
			// Begin this deployment in parallel across all hosts
			err = parallelDeployment(m.Deployments[x].Actions, hosts, &logger)
			if err != nil {
				return err
			}
		} else {
			// This work will be sequential, one host after the next
			for z := range m.Deployments[x].Hosts {
				var hostConfig ssh.HostSSHConfig
				// Find the hosts SSH configuration
				for i := range hosts {
					if hosts[i].Host == m.Deployments[x].Hosts[z] {
						hostConfig = hosts[i]
					}
				}
				err = sequentialDeployment(m.Deployments[x].Actions, hostConfig, &logger)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Begin host by host deployments as part of each deployment
func sequentialDeployment(action []types.Action, hostConfig ssh.HostSSHConfig, logger *plunderlogging.Logger) error {
	var err error

	for y := range action {
		switch action[y].ActionType {
		case "upload":
			err = hostConfig.UploadFile(action[y].Source, action[y].Destination)
			if err != nil {
				// Set checkpoint
				restore.Action = action[y].Name
				restore.Host = hostConfig.Host
				restore.createCheckpoint()
				logger.WriteLogEntry(hostConfig.Host, action[y].Name, "", err.Error())
				logger.SetLoggingState(hostConfig.Host, "Failed")
				// Return the error
				return fmt.Errorf("Upload task [%s] on host [%s] failed with error [%s]", action[y].Name, hostConfig.Host, err)
			}
			log.Infof("Upload Task [%s] on node [%s] completed successfully", action[y].Name, hostConfig.Host)
		case "download":
			err = hostConfig.DownloadFile(action[y].Source, action[y].Destination)
			if err != nil {
				// Set checkpoint
				restore.Action = action[y].Name
				restore.Host = hostConfig.Host
				restore.createCheckpoint()
				logger.WriteLogEntry(hostConfig.Host, action[y].Name, "", err.Error())
				logger.SetLoggingState(hostConfig.Host, "Failed")
				// Return the error
				return fmt.Errorf("Download task [%s] on host [%s] failed with error [%s]", action[y].Name, hostConfig.Host, err)
			}
			log.Infof("Succesfully Downloaded [%s] to [%s] from [%s]", action[y].Source, action[y].Destination, hostConfig.Host)
		case "command":
			// Build out a configuration based upon the action
			cr := parseAndExecute(action[y], &hostConfig)
			// This will end command execution and print the error
			if cr.Error != nil && action[y].IgnoreFailure == false {
				// Set checkpoint
				restore.Action = action[y].Name
				restore.Host = hostConfig.Host
				restore.createCheckpoint()

				// Output error messages
				logger.WriteLogEntry(hostConfig.Host, action[y].Name, cr.Result, cr.Error.Error())
				logger.SetLoggingState(hostConfig.Host, "Failed")
				// cr.Result is ommited here TODO
				return fmt.Errorf("Command task [%s] on host [%s] failed with error [%s]", action[y].Name, hostConfig.Host, cr.Error)
			}

			// if there is an error and we're set to ignore it then process accordingly
			if cr.Error != nil && action[y].IgnoreFailure == true {
				log.Warnf("Command Task [%s] on node [%s] failed (execution will continute)", action[y].Name, hostConfig.Host)
				log.Debugf("Command Results ->\n%s", cr.Result)
				logger.WriteLogEntry(hostConfig.Host, action[y].Name, cr.Result, cr.Error.Error())

				//logger.WriteLogEntry(hostConfig.Host, fmt.Sprintf("Command task [%s] on host [%s] has failed (execution will continute)\n", action[y].Name, hostConfig.Host))
			}

			// No error, task was completed correctly
			if cr.Error == nil {
				// Output success Messages
				log.Infof("Command Task [%s] on node [%s] completed successfully", action[y].Name, hostConfig.Host)
				log.Debugf("Command Results ->\n%s", cr.Result)
				//logger.WriteLogEntry(hostConfig.Host, fmt.Sprintf("Command task [%s] on host [%s] has completed succesfully\n", action[y].Name, hostConfig.Host))
				//logger.WriteLogEntry(hostConfig.Host, fmt.Sprintf("Command task [%s] Output [%s]\n", action[y].Name, cr.Result))
				logger.WriteLogEntry(hostConfig.Host, action[y].Name, cr.Result, "")

			}
		case "pkg":

		case "key":

		default:
			// Set checkpoint (the actiontype may be modified or spelling issue)
			restore.Action = action[y].Name
			restore.Host = hostConfig.Host
			restore.createCheckpoint()
			pluginActions, err := parlayplugin.ExecuteAction(action[y].ActionType, hostConfig.Host, action[y].Plugin)
			if err != nil {
				return err
			}
			log.Debugf("About to execute [%d] actions", len(pluginActions))
			err = sequentialDeployment(pluginActions, hostConfig, logger)
			if err != nil {
				return err
			}
		}
	}
	logger.SetLoggingState(hostConfig.Host, "Completed")
	return nil
}

// Peform all of the actions in parallel on all hosts in the host array
// this function will make use of the parallel ssh calls
func parallelDeployment(action []types.Action, hosts []ssh.HostSSHConfig, logger *plunderlogging.Logger) error {
	for y := range action {
		switch action[y].ActionType {
		case "upload":

			//TODO - Remove or repurpose GENERAL output
			logger.WriteLogEntry("upload", fmt.Sprintf("Uploading file [%s] to Destination [%s] to multiple hosts\n", action[y].Source, action[y].Destination), "", "")

			results := ssh.ParalellUpload(hosts, action[y].Source, action[y].Destination, action[y].Timeout)
			// TODO - Unlikely that this should happen
			if len(results) == 0 {
				return fmt.Errorf("No results have been returned from the parallel execution")
			}
			// Parse the results from the parallel updates
			for i := range results {
				if results[i].Error != nil {
					// Set checkpoint
					restore.Action = action[y].Name
					restore.createCheckpoint()
					logger.WriteLogEntry(results[i].Host, action[y].Name, "", results[i].Error.Error())
					logger.SetLoggingState(results[i].Host, "Failed")

					//logger.WriteLogEntry("", fmt.Sprintf("[%s] Error uploading file [%s] to Destination [%s] to host [%s]\n", time.Now().Format(time.ANSIC), action[y].Source, action[y].Destination, results[i].Host))
					//logger.WriteLogEntry("", fmt.Sprintf("[%s] [%s]\n", time.Now().Format(time.ANSIC), results[i].Error()))
					return fmt.Errorf("Upload task [%s] on host [%s] failed with error [%s]", action[y].Name, results[i].Host, results[i].Error)
				}
				logger.WriteLogEntry(results[i].Host, action[y].Name, fmt.Sprintf("Completed uploading file [%s] to path [%s]", action[y].Source, action[y].Destination), results[i].Error.Error())

				//logger.WriteLogEntry("", fmt.Sprintf("[%s] Completed uploading file [%s] to Destination [%s] to host [%s]\n", time.Now().Format(time.ANSIC), action[y].Source, action[y].Destination, results[i].Host))
				log.Infof("Succesfully uploaded [%s] to [%s] on [%s]", action[y].Source, action[y].Destination, results[i].Host)
			}
		case "download":
			logger.WriteLogEntry("download", fmt.Sprintf("Downloading file [%s] to Destination [%s] from multiple hosts\n", action[y].Source, action[y].Destination), "", "")

			results := ssh.ParalellDownload(hosts, action[y].Source, action[y].Destination, action[y].Timeout)
			// Unlikely that this should happen
			if len(results) == 0 {
				return fmt.Errorf("No results have been returned from the parallel execution")
			}
			// Parse the results from the parallel updates
			for i := range results {
				if results[i].Error != nil {
					// Set checkpoint
					restore.Action = action[y].Name
					restore.createCheckpoint()
					logger.WriteLogEntry(results[i].Host, action[y].Name, "", results[i].Error.Error())
					logger.SetLoggingState(results[i].Host, "Failed")

					//logger.WriteLogEntry("", fmt.Sprintf("[%s] Error downloading file [%s] to [%s] to host [%s]\n", time.Now().Format(time.ANSIC), action[y].Source, action[y].Destination, results[i].Host))
					//logger.WriteLogEntry("", fmt.Sprintf("[%s] [%s]\n", time.Now().Format(time.ANSIC), results[i].Error))

					return fmt.Errorf("Download task [%s] on host [%s] failed with error [%s]", action[y].Name, results[i].Host, results[i].Error)
				}
				logger.WriteLogEntry(results[i].Host, action[y].Name, fmt.Sprintf("Completed Downloading file [%s] to path [%s]", action[y].Source, action[y].Destination), results[i].Error.Error())

				//logger.WriteLogEntry("", fmt.Sprintf("[%s] Completed uploading file [%s] to Destination [%s] to host [%s]\n", time.Now().Format(time.ANSIC), action[y].Source, action[y].Destination, results[i].Host))
				log.Infof("Succesfully uploaded [%s] to [%s] on [%s]", action[y].Source, action[y].Destination, results[i].Host)
			}
		case "command":
			logger.WriteLogEntry("command", fmt.Sprintf("Executing command action [%s] to multiple hosts\n", action[y].Name), "", "")
			command, err := buildCommand(action[y])
			if err != nil {
				// Set checkpoint
				restore.Action = action[y].Name
				restore.createCheckpoint()

				return err
			}
			crs := ssh.ParalellExecute(command, action[y].CommandPipeFile, action[y].CommandPipeCmd, hosts, action[y].Timeout)
			var errors bool // This will only be set to true if a command fails
			for x := range crs {
				if crs[x].Error != nil {
					// Set checkpoint
					restore.Action = action[y].Name
					restore.createCheckpoint()
					logger.WriteLogEntry(crs[x].Host, action[y].Name, crs[x].Result, crs[x].Error.Error())
					logger.SetLoggingState(crs[x].Host, "Failed")
					//log.Errorf("Command task [%s] on host [%s] failed with error [%s]\n\t[%s]", action[y].Name, crs[x].Host, crs[x].Result, crs[x].Error.Error())
					errors = true // An error has been found
					//logger.WriteLogEntry("", fmt.Sprintf("------------  Output  ------------\n%s\n----------------------------------\n", crs[x].Result))
					return fmt.Errorf("Command task [%s] on host [%s] failed with error [%s]\n\t[%s]", action[y].Name, crs[x].Host, crs[x].Error, crs[x].Result)
				}
				log.Infof("Command Task [%s] on node [%s] completed successfully", action[y].Name, crs[x].Host)
				logger.WriteLogEntry(crs[x].Host, action[y].Name, crs[x].Result, "")

				//logger.WriteLogEntry("", fmt.Sprintf("[%s] Command task [%s] on host [%s] has completed succesfully\n", time.Now().Format(time.ANSIC), action[y].Name, crs[x].Host))
				//logger.WriteLogEntry("", fmt.Sprintf("------------  Output  ------------\n%s\n----------------------------------\n", crs[x].Result))
			}
			if errors == true {
				return fmt.Errorf("An error was encountered on command Task [%s]", action[y].Name)
			}
		case "pkg":

		case "key":

		default:
			return fmt.Errorf("Unknown Action [%s]", action[y].ActionType)
		}
	}
	return nil
}
