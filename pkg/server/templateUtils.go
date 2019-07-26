package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/plunder-app/plunder/pkg/utils"
	log "github.com/sirupsen/logrus"
)

// AnyBoot - This flag when set to true will just boot any kernel/initrd/cmdline configuration
var AnyBoot bool

// ReadKeyFromFile - will attempt to read an sshkey from a file and populate the struct
func (config *HostConfig) ReadKeyFromFile(sshKeyPath string) (string, error) {
	var buffer []byte
	if _, err := os.Stat(sshKeyPath); !os.IsNotExist(err) {
		buffer, err = ioutil.ReadFile(sshKeyPath)
		if err != nil {
			// Unable to read the file
			return "", err
		}
	} else {
		// File doesn't exist
		return "", err
	}

	// TrimRight will remove the carriage return from the end of the buffer
	singleLine := strings.TrimRight(string(buffer), "\r\n")
	return singleLine, nil
}

// UpdateControllerConfig will read a configuration string and build the iPXE files needed
func UpdateControllerConfig(configFile []byte) error {

	// Separate configuration until everything is processes correctly
	var updateConfig DeploymentConfigurationFile

	log.Infoln("Updating the Deployment Configuration")
	err := json.Unmarshal(configFile, &updateConfig)
	if err != nil {
		return err
	}

	log.Debugf("Parsing [%d] Configurations", len(updateConfig.Configs))
	for i := range updateConfig.Configs {
		var newConfig, ipxeConfig string

		// We need to move all ":" to "-" to make life a little easier for filesystems and internet standards
		dashMac := strings.Replace(updateConfig.Configs[i].MAC, ":", "-", -1)

		// Find the deployment configuration for this host, either custom or inherit from the controller
		bootConfig := findBootConfigForDeployment(updateConfig.Configs[i])

		// If there is no deployment configuration under this name return an error
		if bootConfig == nil {
			errorString := fmt.Errorf("Host [%s] uses unknown config [%s], stopping config update", updateConfig.Configs[i].MAC, updateConfig.Configs[i].ConfigName)
			return errorString
		}

		// Ensure this entry has the correct mapping
		updateConfig.Configs[i].ConfigBoot = *bootConfig

		// This will populate anything missing from the global configuration
		updateConfig.Configs[i].ConfigHost.PopulateConfiguration(updateConfig.GlobalServerConfig)

		// Look for understood config types
		switch updateConfig.Configs[i].ConfigName {
		case "preseed":
			ipxeConfig = utils.IPXEPreeseed(httpAddress, bootConfig.Kernel, bootConfig.Initrd, bootConfig.Cmdline)
			log.Debugf("Generating preseed ipxeConfig for [%s]", dashMac)
			newConfig = updateConfig.Configs[i].ConfigHost.BuildPreeSeedConfig()

		case "kickstart":
			ipxeConfig = utils.IPXEKickstart(httpAddress, bootConfig.Kernel, bootConfig.Initrd, bootConfig.Cmdline)
			log.Debugf("Generating kickstart ipxeConfig for [%s]", dashMac)
			newConfig = updateConfig.Configs[i].ConfigHost.BuildKickStartConfig()

		default:
			log.Debugf("Building configuration for configName [%s]", updateConfig.Configs[i].ConfigBoot.ConfigName)
			ipxeConfig = utils.IPXEAnyBoot(httpAddress, bootConfig.Kernel, bootConfig.Initrd, bootConfig.Cmdline)
		}

		// If we've specified an iPXE configuration then we add it
		if ipxeConfig != "" {
			path := fmt.Sprintf("/%s.ipxe", dashMac)
			// Only add handler if one didn't exist before
			if httpPaths[path] == "" {
				http.HandleFunc(path, rootHandler)
			}
			httpPaths[path] = ipxeConfig

		}
		if newConfig != "" {
			path := fmt.Sprintf("/%s.cfg", dashMac)
			// Only add handler if one didn't exist before
			if httpPaths[path] == "" {
				http.HandleFunc(path, rootHandler)
			}
			httpPaths[path] = newConfig
		}
	}
	if len(updateConfig.Configs) == 0 {
		// No changes, leave as is (with a warning)
		log.Warnln("No deployment configuration, any existing configuration will remain")
	} else {
		// Updated configuration has been parsed, update internal deployment configuration
		log.Infoln("Updating of deployment configuration complete")
		Deployments = updateConfig
	}

	return nil
}

//FindDeployment - this will return the deployment configuration, allowing the DHCP server to return the correct DHCP options
func FindDeployment(mac string) string {

	// AnyBoot will just boot the specified kernel/initrd
	if AnyBoot == true {
		return "anyboot"
	}

	if len(Deployments.Configs) == 0 {
		// No configurations have been loaded
		log.Warnln("Attempted to perform Mac Address lookup, however no configurations have been loaded")
		return ""
	}
	for i := range Deployments.Configs {
		log.Debugf("Comparing [%s] to [%s]", mac, Deployments.Configs[i].MAC)
		if mac == Deployments.Configs[i].MAC {
			return Deployments.Configs[i].ConfigName
		}
	}
	return ""
}

// PopulateConfiguration - This will read a deployment configuration and attempt to fill any missing fields from the global config
func (config *HostConfig) PopulateConfiguration(globalConfig HostConfig) {
	// NETWORK CONFIGURATION

	// Inherit the global Gateway
	if config.Gateway == "" {
		config.Gateway = globalConfig.Gateway
	}

	// Inherit the global Subnet
	if config.Subnet == "" {
		config.Subnet = globalConfig.Subnet
	}

	// Inherit the global Name Server (DNS)
	if config.NameServer == "" {
		config.NameServer = globalConfig.NameServer
	}

	if config.Adapter == "" {
		config.Adapter = globalConfig.Adapter
	}

	// REPOSITORY CONFIGURATION

	// Inherit the global Repository address
	if config.RepositoryAddress == "" {
		config.RepositoryAddress = globalConfig.RepositoryAddress
	}

	// Inherit the global Repository Mirror directory (typically /ubuntu)
	if config.MirrorDirectory == "" {
		config.MirrorDirectory = globalConfig.MirrorDirectory
	}

	// USER CONFIGURATION

	// Inherit the global Username
	if config.Username == "" {
		config.Username = globalConfig.Username
	}

	// Inherit the global Password
	if config.Password == "" {
		config.Password = globalConfig.Password
	}

	// Inherit the global SSH Key Path
	if config.SSHKeyPath == "" {
		config.SSHKeyPath = globalConfig.SSHKeyPath
	}

	// Package Configuration

	// Inherit the global package selection
	if config.Packages == "" {
		config.Packages = globalConfig.Packages
	}
}
