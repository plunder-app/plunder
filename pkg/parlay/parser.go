package parlay

import (
	"fmt"
	"os"
	"os/exec"

	log "github.com/Sirupsen/logrus"

	"github.com/thebsdbox/plunder/pkg/ssh"
)

// DeploySSH - will iterate through a deployment and perform the relevant actions
func (m *TreasureMap) DeploySSH() error {

	if len(ssh.Hosts) == 0 {
		return fmt.Errorf("No hosts credentials have been loaded")
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

		if m.Deployments[x].Parallel == true {
			// Begin parallel work
			for y := range m.Deployments[x].Actions {
				switch m.Deployments[x].Actions[y].ActionType {
				case "upload":
					ssh.ParalellUpload(hosts, m.Deployments[x].Actions[y].Source, m.Deployments[x].Actions[y].Destination, m.Deployments[x].Actions[y].Timeout)
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
							return err
						}
					case "download":
						err = hostConfig.DownloadFile(m.Deployments[x].Actions[y].Source, m.Deployments[x].Actions[y].Destination)
						if err != nil {
							return err
						}
					case "command":
						// Build out a configuration based upon the action
						cr := m.Deployments[x].Actions[y].parseAndExecute(&hostConfig)
						if cr.Error != nil {
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

	if a.CommandSudo != "" {
		// Add sudo to the command
		command = fmt.Sprintf("sudo %s %s", a.CommandSudo, a.Command)
	} else {
		command = a.Command
	}

	if a.CommandLocal == true {
		b, cr.Error = exec.Command(command).Output()
		cr.Result = string(b)
	} else {
		log.Debugf("Executing command [%s] on host [%s]", command, h.Host)
		cr = ssh.SingleExecute(command, *h, a.Timeout)
		if cr.Error != nil {
			log.Fatalf("%v", cr.Error)
			return cr
		}
	}

	// Save the results into a key to be used at another point
	if a.CommandSaveAsKey != "" {
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
