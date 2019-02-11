package parlay

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/thebsdbox/plunder/pkg/parlay/plugin"
	"github.com/thebsdbox/plunder/pkg/parlay/types"

	"github.com/thebsdbox/plunder/pkg/ssh"
)

var restore Restore

func (m *TreasureMap) findDeployments(deployment []string) (*TreasureMap, error) {

	var newDeploymentList []Deployment

	for x := range deployment {
		for y := range m.Deployments {
			if m.Deployments[y].Name == deployment[x] {
				newDeploymentList = append(newDeploymentList, m.Deployments[y])
			}
		}
	}
	// If this is zero it means that no deployments have been found
	if len(m.Deployments) == 0 {
		return nil, fmt.Errorf("No Deployment(s) have been found")
	}
	m.Deployments = newDeploymentList
	return m, nil
}

func (d *Deployment) findHosts(hosts []string) (*Deployment, error) {

	var newHostList []string

	for x := range hosts {
		for y := range d.Hosts {
			if d.Hosts[y] == hosts[x] {
				newHostList = append(newHostList, d.Hosts[y])
			}
		}
	}
	// If this is zero it means that no hosts have been found
	if len(d.Hosts) == 0 {
		return nil, fmt.Errorf("No Host(s) have been found")
	}
	d.Hosts = newHostList
	return d, nil
}

func (d *Deployment) findActions(actions []string) ([]types.Action, error) {
	var newActionList []types.Action

	for x := range actions {
		for y := range d.Actions {
			if d.Actions[y].Name == actions[x] {
				newActionList = append(newActionList, d.Actions[y])
			}
		}
	}
	// If this is zero it means that no hosts have been found
	if len(d.Actions) == 0 {
		return nil, fmt.Errorf("No Action(s) have been found")
	}
	return newActionList, nil
}

//FindDeployment - takes a number of flags and builds a new map to be processed
func (m *TreasureMap) FindDeployment(deployment, action, host, logFile string, resume bool) error {
	var foundMap TreasureMap
	if deployment != "" {
		log.Debugf("Looking for deployment [%s]", deployment)
		for x := range m.Deployments {
			if m.Deployments[x].Name == deployment {
				foundMap.Deployments = append(foundMap.Deployments, m.Deployments[x])
				// Find a specific action and add or resume from
				if action != "" {
					// Clear the slice as we will be possibly adding different actions
					foundMap.Deployments[0].Actions = nil
					for y := range m.Deployments[x].Actions {
						if m.Deployments[x].Actions[y].Name == action {
							// If we're not resuming that just add the action that we want to happen
							if resume != true {
								foundMap.Deployments[0].Actions = append(foundMap.Deployments[0].Actions, m.Deployments[x].Actions[y])
							} else {
								// Alternatively add all actions from this point
								foundMap.Deployments[0].Actions = m.Deployments[x].Actions[y:]
							}
						}
					}
					// If this is zero it means that no actions have been found
					if len(foundMap.Deployments[0].Actions) == 0 {
						return fmt.Errorf("No actions have been found, looking for action [%s]", action)
					}
				}
				// If a host is specified act soley on it
				if host != "" {
					// Clear the slice as we will be possibly adding different actions
					foundMap.Deployments[0].Hosts = nil
					for y := range m.Deployments[x].Hosts {

						if m.Deployments[x].Hosts[y] == host {
							foundMap.Deployments[0].Hosts = append(foundMap.Deployments[0].Hosts, m.Deployments[x].Hosts[y])
						}
					}
					// If this is zero it means that no hosts have been found
					if len(foundMap.Deployments[0].Hosts) == 0 {
						return fmt.Errorf("No host has been found, looking for host [%s]", host)
					}
				}
			}
		}
		// If this is zero it means that no actions have been found
		if len(foundMap.Deployments) == 0 {
			return fmt.Errorf("No deployment has been found, looking for deployment [%s]", deployment)
		}
	} else {
		return fmt.Errorf("No deployment was specified")
	}
	return foundMap.DeploySSH(logFile)
}

// DeploySSH - will iterate through a deployment and perform the relevant actions
func (m *TreasureMap) DeploySSH(logFile string) error {
	if logFile != "" {
		//enable logging
		logging.init(logFile)
		defer logging.f.Close()
	}

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
		logging.writeString(fmt.Sprintf("[%s] Beginning Deployment [%s]\n", time.Now().Format(time.ANSIC), m.Deployments[x].Name))

		// Set Restore checkpoint
		restore.Deployment = m.Deployments[x].Name
		restore.Hosts = m.Deployments[x].Hosts

		if m.Deployments[x].Parallel == true {
			// Begin this deployment in parallel across all hosts
			err = parallelDeployment(m.Deployments[x].Actions, hosts)
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
				err = sequentialDeployment(m.Deployments[x].Actions, hostConfig)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Begin host by host deployments as part of each deployment
func sequentialDeployment(action []types.Action, hostConfig ssh.HostSSHConfig) error {
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

				// Return the error
				return fmt.Errorf("Download task [%s] on host [%s] failed with error [%s]", action[y].Name, hostConfig.Host, err)
			}
			log.Infof("Succesfully Downloaded [%s] to [%s] from [%s]", action[y].Source, action[y].Destination, hostConfig.Host)
		case "command":
			// Build out a configuration based upon the action
			cr := parseAndExecute(action[y], &hostConfig)
			if cr.Error != nil {
				// Set checkpoint
				restore.Action = action[y].Name
				restore.Host = hostConfig.Host
				restore.createCheckpoint()

				// Output error messages
				logging.writeString(fmt.Sprintf("[%s] Command task [%s] on host [%s] failed with error [%s]\n", time.Now().Format(time.ANSIC), action[y].Name, hostConfig.Host, cr.Error))
				logging.writeString(fmt.Sprintf("------------  Output  ------------\n%s\n----------------------------------\n", cr.Result))
				return fmt.Errorf("Command task [%s] on host [%s] failed with error [%s]\n\t[%s]", action[y].Name, hostConfig.Host, cr.Error, cr.Result)
			}
			// Output success Messages
			log.Infof("Command Task [%s] on node [%s] completed successfully", action[y].Name, hostConfig.Host)
			log.Debugf("Command Results ->\n%s", cr.Result)
			logging.writeString(fmt.Sprintf("[%s] Command task [%s] on host [%s] has completed succesfully\n", time.Now().Format(time.ANSIC), action[y].Name, hostConfig.Host))
			logging.writeString(fmt.Sprintf("------------  Output  ------------\n%s\n----------------------------------\n", cr.Result))

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
			err = sequentialDeployment(pluginActions, hostConfig)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Peform all of the actions in parallel on all hosts in the host array
// this function will make use of the parallel ssh calls
func parallelDeployment(action []types.Action, hosts []ssh.HostSSHConfig) error {
	for y := range action {
		switch action[y].ActionType {
		case "upload":
			logging.writeString(fmt.Sprintf("[%s] Uploading file [%s] to Destination [%s] to multiple hosts\n", time.Now().Format(time.ANSIC), action[y].Source, action[y].Destination))

			results := ssh.ParalellUpload(hosts, action[y].Source, action[y].Destination, action[y].Timeout)
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

					logging.writeString(fmt.Sprintf("[%s] Error uploading file [%s] to Destination [%s] to host [%s]\n", time.Now().Format(time.ANSIC), action[y].Source, action[y].Destination, results[i].Host))
					logging.writeString(fmt.Sprintf("[%s] [%s]\n", time.Now().Format(time.ANSIC), results[i].Error))
					return fmt.Errorf("Upload task [%s] on host [%s] failed with error [%s]", action[y].Name, results[i].Host, results[i].Error)
				}
				logging.writeString(fmt.Sprintf("[%s] Completed uploading file [%s] to Destination [%s] to host [%s]\n", time.Now().Format(time.ANSIC), action[y].Source, action[y].Destination, results[i].Host))
				log.Infof("Succesfully uploaded [%s] to [%s] on [%s]", action[y].Source, action[y].Destination, results[i].Host)
			}
		case "download":
			logging.writeString(fmt.Sprintf("[%s] Downloading remote file [%s] to [%s] from multiple hosts\n", time.Now().Format(time.ANSIC), action[y].Source, action[y].Destination))

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

					logging.writeString(fmt.Sprintf("[%s] Error downloading file [%s] to [%s] to host [%s]\n", time.Now().Format(time.ANSIC), action[y].Source, action[y].Destination, results[i].Host))
					logging.writeString(fmt.Sprintf("[%s] [%s]\n", time.Now().Format(time.ANSIC), results[i].Error))
					return fmt.Errorf("Download task [%s] on host [%s] failed with error [%s]", action[y].Name, results[i].Host, results[i].Error)
				}
				logging.writeString(fmt.Sprintf("[%s] Completed uploading file [%s] to Destination [%s] to host [%s]\n", time.Now().Format(time.ANSIC), action[y].Source, action[y].Destination, results[i].Host))
				log.Infof("Succesfully uploaded [%s] to [%s] on [%s]", action[y].Source, action[y].Destination, results[i].Host)
			}
		case "command":
			logging.writeString(fmt.Sprintf("[%s] Executing command action [%s] to multiple hosts\n", time.Now().Format(time.ANSIC), action[y].Name))
			command, err := buildCommand(action[y])
			if err != nil {
				// Set checkpoint
				restore.Action = action[y].Name
				restore.createCheckpoint()

				return err
			}
			crs := ssh.ParalellExecute(command, hosts, action[y].Timeout)
			var errors bool // This will only be set to true if a command fails
			for x := range crs {
				if crs[x].Error != nil {
					// Set checkpoint
					restore.Action = action[y].Name
					restore.createCheckpoint()

					log.Errorf("Command task [%s] on host [%s] failed with error [%s]\n\t[%s]", action[y].Name, crs[x].Host, crs[x].Result, crs[x].Error.Error())
					errors = true // An error has been found
					logging.writeString(fmt.Sprintf("------------  Output  ------------\n%s\n----------------------------------\n", crs[x].Result))
					return fmt.Errorf("Command task [%s] on host [%s] failed with error [%s]\n\t[%s]", action[y].Name, crs[x].Host, crs[x].Error, crs[x].Result)
				}
				log.Infof("Command Task [%s] on node [%s] completed successfully", action[y].Name, crs[x].Host)
				logging.writeString(fmt.Sprintf("[%s] Command task [%s] on host [%s] has completed succesfully\n", time.Now().Format(time.ANSIC), action[y].Name, crs[x].Host))
				logging.writeString(fmt.Sprintf("------------  Output  ------------\n%s\n----------------------------------\n", crs[x].Result))
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

func buildCommand(a types.Action) (string, error) {
	var command string

	// An executable Key takes presedence
	if a.KeyName != "" {
		keycmd := Keys[a.KeyName]
		// Check that the key exists
		if keycmd == "" {
			return "", fmt.Errorf("Unable to find command under key '%s'", a.KeyName)

		}
		if a.CommandSudo != "" {
			// Add sudo to the Key command
			command = fmt.Sprintf("sudo -n -u %s %s", a.CommandSudo, keycmd)
		} else {
			command = keycmd
		}
	} else {
		// Not using a key, using a shell command
		if a.CommandSudo != "" {
			// Add sudo to the Shell command
			command = fmt.Sprintf("sudo -n -u %s %s", a.CommandSudo, a.Command)
		} else {
			command = a.Command
		}
	}
	return command, nil
}

func parseAndExecute(a types.Action, h *ssh.HostSSHConfig) ssh.CommandResult {
	// This will parse the options passed in the action and execute the required string
	var cr ssh.CommandResult
	var b []byte

	command, err := buildCommand(a)
	if err != nil {
		cr.Error = err
		return cr
	}

	if a.CommandLocal == true {
		log.Debugf("Command [%s]", command)
		cmd := exec.Command("bash", "-c", command)
		b, cr.Error = cmd.CombinedOutput()
		if cr.Error != nil {
			return cr
		}
		cr.Result = strings.TrimRight(string(b), "\r\n")
	} else {
		log.Debugf("Executing command [%s] on host [%s]", command, h.Host)
		cr = ssh.SingleExecute(command, *h, a.Timeout)

		cr.Result = strings.TrimRight(cr.Result, "\r\n")

		// If the command hasn't returned anything, put a filler in
		if cr.Result == "" {
			cr.Result = "[No Output]"
		}
		if cr.Error != nil {
			return cr
		}
	}

	// Save the results into a key to be used at another point
	if a.CommandSaveAsKey != "" {
		log.Debugf("Adding new results to key [%s]", a.CommandSaveAsKey)
		Keys[a.CommandSaveAsKey] = cr.Result
	}

	// Save the results into a file to be used at another point
	if a.CommandSaveFile != "" {
		var f *os.File
		f, cr.Error = os.Create(a.CommandSaveFile)
		if cr.Error != nil {
			return cr
		}

		defer f.Close()

		_, cr.Error = f.WriteString(cr.Result)
		if cr.Error != nil {
			return cr
		}
		f.Sync()
	}

	return cr
}
