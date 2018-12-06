package parlay

import (
	"fmt"

	"github.com/thebsdbox/plunder/pkg/ssh"
)

// DeploySSH - will iterate through a deployment and perform the relevant actions
func (deployment *Deployment) DeploySSH() error {

	if len(ssh.Hosts) == 0 {
		return fmt.Errorf("No hosts credentials have been loaded")
	}

	// Build new hosts list from imported SSH servers and compare that we have required credentials
	hosts, err := ssh.FindHosts(deployment.Hosts)
	if err != nil {
		return err
	}

	if deployment.Parallel == true {
		for i := range deployment.Actions {
			switch deployment.Actions[i].ActionType {
			case "upload":
				ssh.ParalellUpload(hosts, deployment.Actions[i].Source, deployment.Actions[i].Destination, deployment.Actions[i].Timeout)
			case "download":

			case "command":

			case "pkg":

			case "key":

			default:
				return fmt.Errorf("Unknown Action [%s]", deployment.Actions[i].ActionType)
			}
		}
	}

	return nil
}
