package services

import (
	"encoding/json"
	"fmt"

	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"
)

// This is needed by other functions to build strings
var httpAddress string

// Controller contains all the "current" settings for booting servers
var Controller BootController

// Deployments - contains an accessible "current" configuration for all deployments
var Deployments DeploymentConfigurationFile

// ParseControllerData will read in a byte array and attempt to parse it as yaml or json
func ParseControllerData(b []byte) error {

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

// ParseBootController - will iterate through the boot controller and see if any changes need applying
// this is mainly for the dynamic loading of ISOs
func (c *BootController) ParseBootController() error {

	for i := range c.BootConfigs {
		// If either the prefix or path are blank then iterate over, both need to be set in order to load the ISO
		if c.BootConfigs[i].ISOPrefix == "" || c.BootConfigs[i].ISOPath == "" {
			log.Debugf("No ISO is being parsed for configuration %s", c.BootConfigs[i].ConfigName)
		} else {

			// Create the prefix
			urlPrefix := fmt.Sprintf("/%s/", c.BootConfigs[i].ISOPrefix)

			// Only create the handler if one doesn't exist
			if _, ok := isoMapper[c.BootConfigs[i].ISOPrefix]; !ok {
				log.Debugf("Adding handler %s", urlPrefix)

				serveMux.HandleFunc(urlPrefix, isoReader)
			}

			// Atempt to open the ISO and add it to the map for usage later
			err := OpenISO(c.BootConfigs[i].ISOPath, c.BootConfigs[i].ISOPrefix)
			if err != nil {
				return err
			}

			fmt.Printf("%v", isoMapper)
			log.Debugf("Updating handler %s for config %s", urlPrefix, c.BootConfigs[i].ConfigName)

		}
	}
	// Parse the boot controllers for new configuration changes
	c.generateBootTypeHanders()
	return nil
}

// DeleteBootControllerConfig - will iterate through the boot controller and see if any changes need applying
// this is mainly for the dynamic loading of ISOs
func (c *BootController) DeleteBootControllerConfig(configName string) error {

	for i := range c.BootConfigs {
		if c.BootConfigs[i].ConfigName == configName {
			// Remove the mapping to an ISO path
			isoMapper[c.BootConfigs[i].ISOPrefix] = ""
			c.BootConfigs = append(c.BootConfigs[:i], c.BootConfigs[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Unable to find boot configuration %s", configName)
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
