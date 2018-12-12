package parlay

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/thebsdbox/plunder/pkg/ssh"
)

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
					for y := range m.Deployments[x].Actions {
						if m.Deployments[x].Actions[y].Name == action {
							// Clear the slice as we will be possibly adding different actions
							foundMap.Deployments[0].Actions = nil

							// If we're not resuming that just add the action that we want to happen
							if resume != true {
								foundMap.Deployments[0].Actions = append(foundMap.Deployments[0].Actions, m.Deployments[x].Actions[y])
							} else {
								// Alternatively add all actions from this point
								foundMap.Deployments[0].Actions = m.Deployments[x].Actions[y:]
							}
							// If this is zero it means that no actions have been found
							if len(foundMap.Deployments[0].Actions) == 0 {
								fmt.Printf("No actions have been found, looking for action [%s]", action)
							}
						}
					}
				}
				// If a host is specified act soley on it
				if host != "" {
					for y := range m.Deployments[x].Hosts {
						if m.Deployments[x].Hosts[y] == host {
							foundMap.Deployments[0].Hosts = append(foundMap.Deployments[0].Hosts, m.Deployments[x].Hosts[y])
						}
					}
					// If this is zero it means that no hosts have been found
					if len(foundMap.Deployments[0].Hosts) == 0 {
						fmt.Printf("No host has been found, looking for host [%s]", host)
					}
				}
			}
		}
		// If this is zero it means that no actions have been found
		if len(foundMap.Deployments) == 0 {
			fmt.Printf("No deployment has been found, looking for deployment [%s]", deployment)
		}
	} else {
		return fmt.Errorf("No deployment was specified")
	}
	return foundMap.DeploySSH(logFile)
}

// DeploySSH - will iterate through a deployment and perform the relevant actions
func (m *TreasureMap) DeploySSH(logFile string) error {
	var logging bool
	var file *os.File
	var err error
	if logFile != "" {
		//enable logging
		logging = true
		file, err = os.Create(logFile)
		if err != nil {
			return err
		}
		defer file.Close()
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
		if logging {
			_, err = file.WriteString(fmt.Sprintf("[%s] Beginning Deployment [%s]\n", time.Now().Format(time.ANSIC), m.Deployments[x].Name))
			if err != nil {
				return err
			}
		}
		if m.Deployments[x].Parallel == true {
			// Begin parallel work
			for y := range m.Deployments[x].Actions {
				switch m.Deployments[x].Actions[y].ActionType {
				case "upload":
					if logging {
						_, err = file.WriteString(fmt.Sprintf("[%s] Uploading file [%s] to Destination [%s] to multiple hosts\n", time.Now().Format(time.ANSIC), m.Deployments[x].Actions[y].Source, m.Deployments[x].Actions[y].Destination))
						if err != nil {
							return err
						}
					}
					results := ssh.ParalellUpload(hosts, m.Deployments[x].Actions[y].Source, m.Deployments[x].Actions[y].Destination, m.Deployments[x].Actions[y].Timeout)
					// Unlikely that this should happen
					if len(results) == 0 {
						return fmt.Errorf("No results have been returned from the parallel execution")
					}
					// Parse the results from the parallel updates
					for i := range results {
						if results[i].Error != nil {
							if logging {
								_, err = file.WriteString(fmt.Sprintf("[%s] Error uploading file [%s] to Destination [%s] to host [%s]\n", time.Now().Format(time.ANSIC), m.Deployments[x].Actions[y].Source, m.Deployments[x].Actions[y].Destination, results[i].Host))
								if err != nil {
									return err
								}
								_, err = file.WriteString(fmt.Sprintf("[%s] [%s]\n", time.Now().Format(time.ANSIC), results[i].Error))
								if err != nil {
									return err
								}
							}
							return fmt.Errorf("Upload task [%s] on host [%s] failed with error [%s]", m.Deployments[x].Actions[y].Name, results[i].Host, results[i].Error)
						}
						if logging {
							_, err = file.WriteString(fmt.Sprintf("[%s] Completed uploading file [%s] to Destination [%s] to host [%s]\n", time.Now().Format(time.ANSIC), m.Deployments[x].Actions[y].Source, m.Deployments[x].Actions[y].Destination, results[i].Host))
							if err != nil {
								return err
							}
						}
						log.Infof("Succesfully uploaded [%s] to [%s] on [%s]", m.Deployments[x].Actions[y].Source, m.Deployments[x].Actions[y].Destination, results[i].Host)
					}
				case "download":

				case "command":

				case "pkg":

				case "key":

				default:
					return fmt.Errorf("Unknown Action [%s]", m.Deployments[x].Actions[y].ActionType)
				}
			}
		} else {
			// Begin host by host deployments as part of each deployment
			for z := range m.Deployments[x].Hosts {

				var hostConfig ssh.HostSSHConfig
				// Find the hosts SSH configuration
				for i := range hosts {
					if hosts[i].Host == m.Deployments[x].Hosts[z] {
						hostConfig = hosts[i]
					}
				}

				for y := range m.Deployments[x].Actions {
					switch m.Deployments[x].Actions[y].ActionType {
					case "upload":
						err = hostConfig.UploadFile(m.Deployments[x].Actions[y].Source, m.Deployments[x].Actions[y].Destination)
						if err != nil {
							return fmt.Errorf("Upload task [%s] on host [%s] failed with error [%s]", m.Deployments[x].Actions[y].Name, hostConfig.Host, err)
						}
						log.Infof("Upload Task [%s] on node [%s] completed successfully", m.Deployments[x].Actions[y].Name, hostConfig.Host)
					case "download":
						err = hostConfig.DownloadFile(m.Deployments[x].Actions[y].Source, m.Deployments[x].Actions[y].Destination)
						if err != nil {
							return err
						}
					case "command":
						// Build out a configuration based upon the action
						cr := m.Deployments[x].Actions[y].parseAndExecute(&hostConfig)
						if cr.Error != nil {
							if logging {
								_, err = file.WriteString(fmt.Sprintf("[%s] Command task [%s] on host [%s] failed with error [%s]\n", time.Now().Format(time.ANSIC), m.Deployments[x].Actions[y].Name, hostConfig.Host, cr.Error))
								if err != nil {
									return err
								}
								_, err = file.WriteString(fmt.Sprintf("------------  Output  ------------\n%s\n---------------------------------\n", cr.Result))
								if err != nil {
									return err
								}
							}
							return fmt.Errorf("Command task [%s] on host [%s] failed with error [%s]\n\t[%s]", m.Deployments[x].Actions[y].Name, hostConfig.Host, cr.Error, cr.Result)
						}
						log.Infof("Command Task [%s] on node [%s] completed successfully", m.Deployments[x].Actions[y].Name, hostConfig.Host)
						log.Debugf("Command Results ->\n%s", cr.Result)
						if logging {
							_, err = file.WriteString(fmt.Sprintf("[%s] Command task [%s] on host [%s] has completed succesfully\n", time.Now().Format(time.ANSIC), m.Deployments[x].Actions[y].Name, hostConfig.Host))
							if err != nil {
								return err
							}
						}
						_, err = file.WriteString(fmt.Sprintf("------------  Output  ------------\n%s\n---------------------------------\n", cr.Result))
						if err != nil {
							return err
						}
					case "pkg":

					case "key":

					default:
						return fmt.Errorf("Unknown Action [%s]", m.Deployments[x].Actions[y].ActionType)
					}
				}
			}
		}
	}

	return nil
}

func (a *Action) parseAndExecute(h *ssh.HostSSHConfig) ssh.CommandResult {
	// This will parse the options passed in the action and execute the required string
	var command string
	var cr ssh.CommandResult
	var b []byte

	// An executable Key takes presedence
	if a.KeyName != "" {
		keycmd := Keys[a.KeyName]
		// Check that the key exists
		if keycmd == "" {
			cr.Error = fmt.Errorf("Unable to find command under key '%s'", a.KeyName)
			return cr
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

	if a.CommandLocal == true {
		b, cr.Error = exec.Command(command).Output()
		if cr.Error != nil {
			return cr
		}
		cr.Result = strings.TrimRight(string(b), "\r\n")
	} else {
		log.Debugf("Executing command [%s] on host [%s]", command, h.Host)
		cr = ssh.SingleExecute(command, *h, a.Timeout)
		cr.Result = strings.TrimRight(cr.Result, "\r\n")
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
