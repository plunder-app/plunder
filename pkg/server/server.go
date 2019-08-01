package server

import (
	"encoding/json"
	"fmt"

	"github.com/ghodss/yaml"
)

// This is needed by other functions to build strings
var httpAddress string

// Controller contains all the "current" settings for booting servers
var Controller BootController

// Deployments - contains an accessible "current" configuration for all deployments
var Deployments DeploymentConfigurationFile

// ParseControllerFile will read in a byte array and attempt to parse it as yaml or json
func ParseControllerFile(b []byte) error {

	jsonBytes, err := yaml.YAMLToJSON(b)
	if err == nil {
		// If there were no errors then the YAML => JSON was successful, no attempt to unmarshall
		err = json.Unmarshal(jsonBytes, &Controller)
		if err != nil {
			return fmt.Errorf("Unable to parse configuration as either yaml or json")
		}

	} else {
		// Couldn't parse the yaml to JSON
		// Attempt to parse it as JSON
		err = json.Unmarshal(b, &Controller)
		if err != nil {
			return fmt.Errorf("Unable to parse configuration as either yaml or json")
		}
	}
	return nil
}

// ParseDeployment will read in a byte array and attempt to parse it as yaml or json
func ParseDeployment(b []byte) (*DeploymentConfigurationFile, error) {

	var deployment DeploymentConfigurationFile

	jsonBytes, err := yaml.YAMLToJSON(b)
	if err == nil {
		// If there were no errors then the YAML => JSON was successful, no attempt to unmarshall
		err = json.Unmarshal(jsonBytes, &deployment)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse configuration as either yaml or json\n %s", err.Error())
		}

	} else {
		// Couldn't parse the yaml to JSON
		// Attempt to parse it as JSON
		err = json.Unmarshal(b, &deployment)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse configuration as either yaml or json\n %s", err.Error())
		}
	}
	return &deployment, nil
}
