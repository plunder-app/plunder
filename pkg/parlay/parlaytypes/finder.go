package parlaytypes

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// This will iterate through a deployment map and build a new deployment map from found deployments
func (m *TreasureMap) FindDeployments(deployment []string) (*TreasureMap, error) {

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

func (d *Deployment) FindHosts(hosts []string) (*Deployment, error) {

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

func (d *Deployment) FindActions(actions []string) ([]Action, error) {
	var newActionList []Action

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
func (m *TreasureMap) FindDeployment(deployment, action, host, logFile string, resume bool) (*TreasureMap, error) {
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
						return nil, fmt.Errorf("No actions have been found, looking for action [%s]", action)
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
						return nil, fmt.Errorf("No host has been found, looking for host [%s]", host)
					}
				}
			}
		}
		// If this is zero it means that no actions have been found
		if len(foundMap.Deployments) == 0 {
			return nil, fmt.Errorf("No deployment has been found, looking for deployment [%s]", deployment)
		}
	} else {
		return nil, fmt.Errorf("No deployment was specified")
	}
	return &foundMap, nil
	//return parlay.DeploySSH(foundMap, logFile, false, false)
}
