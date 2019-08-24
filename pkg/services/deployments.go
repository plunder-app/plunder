package services

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/plunder-app/plunder/pkg/utils"
	log "github.com/sirupsen/logrus"
)

// DefaultBootType specifies what a server will default to if no config is found
var DefaultBootType string

// This stores the mapping for a url to the data /macaddress.file => data
var httpPaths map[string]string

func init() {
	// Initialise the paths map
	httpPaths = make(map[string]string)
}

func rebuildConfiguration(updateConfig *DeploymentConfigurationFile) error {

	// If HTTP isn't enabled we can't build the multiplexer for URLs
	if mux == nil {
		return fmt.Errorf("Deployment HTTP Server isn't enabled, so parsing deployments isn't possible")
	}

	log.Debugf("Parsing [%d] Configurations", len(updateConfig.Configs))
	for i := range updateConfig.Configs {
		// inMemipxeConfig is a custom configuration that matches kernel/initrd & cmdline and is 00:11:22:33:44:55.ipxe
		var inMemipxeConfig string

		// inMemipxeConfig is a custom configuration that is specific to the boot type [preseed/kickstart/vsphere] and is 00:11:22:33:44:55.cfg
		var inMemBootConfig string

		// imMemESXiKickstart is a custom configuraton specific to vSphere for it's kickstart
		var imMemESXiKickstart string

		// We need to move all ":" to "-" to make life a little easier for filesystems and internet standards
		dashMac := strings.Replace(updateConfig.Configs[i].MAC, ":", "-", -1)

		// Find the deployment configuration for this host, either custom or inherit from the controller
		bootConfig := findBootConfigForDeployment(updateConfig.Configs[i])

		// If there is no deployment configuration under this name return an error
		if bootConfig == nil {
			errorString := fmt.Errorf("Host [%s] uses unknown config [%s], stopping config update", updateConfig.Configs[i].MAC, updateConfig.Configs[i].ConfigName)
			log.Errorln(errorString)
			return errorString
		}

		// Ensure this entry has the correct mapping
		updateConfig.Configs[i].ConfigBoot = *bootConfig

		// This will populate anything missing from the global configuration
		updateConfig.Configs[i].ConfigHost.PopulateConfiguration(updateConfig.GlobalServerConfig)

		// Look for understood config types
		switch updateConfig.Configs[i].ConfigName {
		case "preseed":
			inMemipxeConfig = utils.IPXEPreeseed(httpAddress, bootConfig.Kernel, bootConfig.Initrd, bootConfig.Cmdline)
			log.Debugf("Generating preseed ipxeConfig for [%s]", dashMac)
			inMemBootConfig = updateConfig.Configs[i].ConfigHost.BuildPreeSeedConfig()

		case "kickstart":
			inMemipxeConfig = utils.IPXEKickstart(httpAddress, bootConfig.Kernel, bootConfig.Initrd, bootConfig.Cmdline)
			log.Debugf("Generating kickstart ipxeConfig for [%s]", dashMac)
			inMemBootConfig = updateConfig.Configs[i].ConfigHost.BuildKickStartConfig()

		case "vsphere":
			inMemipxeConfig = utils.IPXEVSphere(httpAddress, bootConfig.Kernel, bootConfig.Cmdline)
			log.Debugf("Generating vsphere ipxeConfig for [%s]", dashMac)
			inMemBootConfig = updateConfig.Configs[i].ConfigHost.BuildESXiConfig()
			imMemESXiKickstart = updateConfig.Configs[i].ConfigHost.BuildESXiKickStart()

		default:
			log.Debugf("Building configuration for configName [%s]", updateConfig.Configs[i].ConfigBoot.ConfigName)
			inMemipxeConfig = utils.IPXEAnyBoot(httpAddress, bootConfig.Kernel, bootConfig.Initrd, bootConfig.Cmdline)
		}

		// Build the configuration that is passed to iPXE on boot
		if inMemipxeConfig != "" {
			path := fmt.Sprintf("/%s.ipxe", dashMac)
			if _, ok := httpPaths[path]; !ok {
				// Only create the handler if one doesn't exist
				mux.HandleFunc(path, rootHandler)
			}

			httpPaths[path] = inMemipxeConfig
		}

		// Build a boot configuration that is passed to a kernel
		if inMemBootConfig != "" {
			path := fmt.Sprintf("/%s.cfg", dashMac)
			if _, ok := httpPaths[path]; !ok {
				// Only create the handler if one doesn't exist
				mux.HandleFunc(path, rootHandler)
			}
			httpPaths[path] = inMemBootConfig
		}

		// Build a vSphere kickstart configuration that is passed to an installer
		if imMemESXiKickstart != "" {
			path := fmt.Sprintf("/%s.ks", dashMac)
			if _, ok := httpPaths[path]; !ok {
				// Only create the handler if one doesn't exist
				mux.HandleFunc(path, rootHandler)
			}
			httpPaths[path] = imMemESXiKickstart
		}

	}
	if len(updateConfig.Configs) == 0 {
		// No changes, leave as is (with a warning)
		log.Warnln("No deployment configuration, any existing configuration will remain")
	} else {
		// Updated configuration has been parsed, update internal deployment configuration
		log.Infoln("Updating of deployment configuration complete")
		Deployments = *updateConfig
	}

	return nil
}

// UpdateDeploymentConfig will read a configuration string and build the iPXE files needed
func UpdateDeploymentConfig(rawDeploymentConfig []byte) error {

	// Separate configuration until everything is processes correctly

	log.Infoln("Updating the Deployment Configuration")
	updateConfig, err := ParseDeployment(rawDeploymentConfig)
	if err != nil {
		return err
	}
	return rebuildConfiguration(updateConfig)

}

// AddDeployment - This function will add a new deployment to the deployment configuration
func AddDeployment(rawDeployment []byte) error {

	var newDeployment DeploymentConfig

	err := json.Unmarshal(rawDeployment, &newDeployment)
	if err != nil {
		return fmt.Errorf("Unable to parse deployment configuration")
	}
	// Find the original deployment via it's mac address
	for i := range Deployments.Configs {
		// Compare this deployment to the one we're looking for
		if Deployments.Configs[i].MAC == newDeployment.MAC {
			return fmt.Errorf("Duplicate entry for MAC address [%s]", newDeployment.MAC)
		}
	}
	// We will now duplicate our configuration
	updateConfig := Deployments
	// We will need to create space to copy the existing configurations over
	updateConfig.Configs = make([]DeploymentConfig, len(Deployments.Configs))
	// Copy our existing configurations into the new configuration
	copy(updateConfig.Configs, Deployments.Configs)
	// Append our new configuration into our new copy
	updateConfig.Configs = append(updateConfig.Configs, newDeployment)

	// Parse the new configuration
	return rebuildConfiguration(&updateConfig)
}

// GetDeployment - This function will add a new deployment to the deployment configuration
func GetDeployment(macAddress string) *DeploymentConfig {
	// Iterate through all the deployments
	for i := range Deployments.Configs {
		if macAddress == Deployments.Configs[i].MAC {
			return &Deployments.Configs[i]
		}
	}
	return nil
}

// UpdateDeployment - This function will add a new deployment to the deployment configuration
func UpdateDeployment(macAddress string, rawDeployment []byte) error {

	var newDeployment DeploymentConfig

	err := json.Unmarshal(rawDeployment, &newDeployment)
	if err != nil {
		return fmt.Errorf("Unable to parse deployment configuration")
	}

	// We will now duplicate our configuration
	updateConfig := Deployments
	// We will need to create space to copy the existing configurations over
	updateConfig.Configs = make([]DeploymentConfig, len(Deployments.Configs))
	// Copy our existing configurations into the new configuration
	copy(updateConfig.Configs, Deployments.Configs)

	// Find the original deployment via it's mac address
	for i := range updateConfig.Configs {
		// Compare this deployment to the one we're looking for
		if updateConfig.Configs[i].MAC == macAddress {
			// Remove the old matching configuration
			updateConfig.Configs = append(updateConfig.Configs[:i], updateConfig.Configs[i+1:]...)
			// Append our new configuration into our new copy
			updateConfig.Configs = append(updateConfig.Configs, newDeployment)

			// Parse the new configuration
			return rebuildConfiguration(&updateConfig)
		}
	}
	return fmt.Errorf("Unable to find existing deployment for MAC address [%s]", macAddress)
}

// DeleteDeploymentMac - This function will delete a deployment based upon it's mac Address
func DeleteDeploymentMac(macAddress string, rawDeployment []byte) error {

	// We will now duplicate our configuration
	updateConfig := Deployments
	// We will need to create space to copy the existing configurations over
	updateConfig.Configs = make([]DeploymentConfig, len(Deployments.Configs))
	// Copy our existing configurations into the new configuration
	copy(updateConfig.Configs, Deployments.Configs)

	// Find the original deployment via it's mac address
	for i := range updateConfig.Configs {
		// Compare this deployment to the one we're looking for
		if updateConfig.Configs[i].MAC == macAddress {
			// Remove the old matching configuration
			updateConfig.Configs = append(updateConfig.Configs[:i], updateConfig.Configs[i+1:]...)
			// Parse the new configuration
			return rebuildConfiguration(&updateConfig)
		}
	}
	return fmt.Errorf("Unable to find existing deployment for Address [%s]", macAddress)

}

// DeleteDeploymentAddress - This function will delete a deployment based upon it's IP Address
func DeleteDeploymentAddress(address string, rawDeployment []byte) error {

	// We will now duplicate our configuration
	updateConfig := Deployments
	// We will need to create space to copy the existing configurations over
	updateConfig.Configs = make([]DeploymentConfig, len(Deployments.Configs))
	// Copy our existing configurations into the new configuration
	copy(updateConfig.Configs, Deployments.Configs)

	// Find the original deployment via it's mac address
	for i := range updateConfig.Configs {
		// Compare this deployment to the one we're looking for
		if updateConfig.Configs[i].ConfigHost.IPAddress == address {
			// Remove the old matching configuration
			updateConfig.Configs = append(updateConfig.Configs[:i], updateConfig.Configs[i+1:]...)
			// Parse the new configuration
			return rebuildConfiguration(&updateConfig)
		}
	}
	return fmt.Errorf("Unable to find existing deployment for Address [%s]", address)

}

// UpdateGlobalDeploymentConfig - This allows updating of the global configuration independently
func UpdateGlobalDeploymentConfig(rawDeployment []byte) error {
	var globalDeploymentConfig HostConfig
	err := json.Unmarshal(rawDeployment, &globalDeploymentConfig)
	if err != nil {
		return fmt.Errorf("Unable to parse deployment configuration")
	}
	// Update the deployments with the new configuration
	Deployments.GlobalServerConfig = globalDeploymentConfig
	return nil
}

//FindDeploymentConfigFromMac - this will return the deployment configuration, allowing the DHCP server to return the correct DHCP options
func FindDeploymentConfigFromMac(mac string) string {

	// AnyBoot will just boot the specified kernel/initrd
	// if AnyBoot == true {
	// 	return "anyboot"
	// }

	if len(Deployments.Configs) == 0 {
		// No configurations have been loaded
		log.Warnln("Attempted to perform Mac Address lookup, however no configurations have been loaded")
		return ""
	}
	for i := range Deployments.Configs {
		log.Debugf("Comparing [%s] to [%s]", mac, strings.ToLower(Deployments.Configs[i].MAC))
		if mac == strings.ToLower(Deployments.Configs[i].MAC) {
			return Deployments.Configs[i].ConfigName
		}
	}
	return DefaultBootType
}
