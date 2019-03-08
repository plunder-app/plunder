package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/plunder-app/plunder/pkg/utils"
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

// UpdateConfiguration will read a configuration string and build the iPXE files needed
func UpdateConfiguration(configFile []byte) error {

	log.Infoln("Updating the Deployment Configuration")
	err := json.Unmarshal(configFile, &DeploymentConfig)
	if err != nil {
		return err
	}
	if len(DeploymentConfig.Deployments) == 0 {
		log.Warnln("No deployment configurations found")
	}
	for i := range DeploymentConfig.Deployments {
		var newConfig, ipxeConfig string

		// We need to move all ":" to "-" to make life a little easier for filesystems and internet standards
		dashMac := strings.Replace(DeploymentConfig.Deployments[i].MAC, ":", "-", -1)

		switch DeploymentConfig.Deployments[i].Deployment {
		case "preseed":
			// If a kernel and initrd are submitted then Create an .ipxe file
			if DeploymentConfig.Deployments[i].Kernel != "" && DeploymentConfig.Deployments[i].Initrd != "" {
				ipxeConfig = utils.IPXEPreeseed(httpAddress, DeploymentConfig.Deployments[i].Kernel, DeploymentConfig.Deployments[i].Initrd, DeploymentConfig.Deployments[i].Cmdline)
				log.Debugf("Generating ipxeConfig for [%s]", dashMac)
			}
			// Build a preseed configuration and write it to disk
			newConfig = DeploymentConfig.Deployments[i].Config.BuildPreeSeedConfig()

		case "kickstart":
			// If a kernel and initrd are submitted then Create an .ipxe file
			if DeploymentConfig.Deployments[i].Kernel != "" && DeploymentConfig.Deployments[i].Initrd != "" {
				ipxeConfig = utils.IPXEKickstart(httpAddress, DeploymentConfig.Deployments[i].Kernel, DeploymentConfig.Deployments[i].Initrd, DeploymentConfig.Deployments[i].Cmdline)
			}
			// Build a kickstart configuration and write it to disk
			newConfig = DeploymentConfig.Deployments[i].Config.BuildKickStartConfig()

		default:
			return fmt.Errorf("Unknown deployment method [%s]", DeploymentConfig.Deployments[i].Deployment)
		}

		if ipxeConfig != "" {
			path := fmt.Sprintf("/%s.ipxe", dashMac)
			// Only add handler if one didn't exist before
			if httpPaths[path] == "" {
				http.HandleFunc(path, rootHandler)
			}
			httpPaths[path] = ipxeConfig
		}
		path := fmt.Sprintf("/%s.cfg", dashMac)
		// Only add handler if one didn't exist before
		if httpPaths[path] == "" {
			http.HandleFunc(path, rootHandler)
		}
		httpPaths[path] = newConfig
	}
	return nil
}

//FindDeployment - this will return the deployment configuration, allowing the DHCP server to return the correct DHCP options
func FindDeployment(mac string) string {

	// AnyBoot will just boot the specified kernel/initrd
	if AnyBoot == true {
		return "anyboot"
	}

	if len(DeploymentConfig.Deployments) == 0 {
		// No configurations have been loaded
		log.Warnln("Attempted to perform Mac Address lookup, however no configurations have been loaded")
		return ""
	}
	for i := range DeploymentConfig.Deployments {
		log.Debugf("Comparing [%s] to [%s]", mac, DeploymentConfig.Deployments[i].MAC)
		if mac == DeploymentConfig.Deployments[i].MAC {
			return DeploymentConfig.Deployments[i].Deployment
		}
	}
	return ""
}

// PopulateConfiguration - This will read a deployment configuration and attempt to fill any missing fields from the global config
func (config *HostConfig) PopulateConfiguration() {
	// NETWORK CONFIGURATION

	// Inherit the global Gateway
	if config.Gateway == "" {
		config.Gateway = DeploymentConfig.GlobalServerConfig.Gateway
	}

	// Inherit the global Subnet
	if config.Subnet == "" {
		config.Subnet = DeploymentConfig.GlobalServerConfig.Subnet
	}

	// Inherit the global Name Server (DNS)
	if config.NameServer == "" {
		config.NameServer = DeploymentConfig.GlobalServerConfig.NameServer
	}

	if config.Adapter == "" {
		config.Adapter = DeploymentConfig.GlobalServerConfig.Adapter
	}

	// REPOSITORY CONFIGURATION

	// Inherit the global Repository address
	if config.RepositoryAddress == "" {
		config.RepositoryAddress = DeploymentConfig.GlobalServerConfig.RepositoryAddress
	}

	// Inherit the global Repository Mirror directory (typically /ubuntu)
	if config.MirrorDirectory == "" {
		config.MirrorDirectory = DeploymentConfig.GlobalServerConfig.MirrorDirectory
	}

	// USER CONFIGURATION

	// Inherit the global Username
	if config.Username == "" {
		config.Username = DeploymentConfig.GlobalServerConfig.Username
	}

	// Inherit the global Password
	if config.Password == "" {
		config.Password = DeploymentConfig.GlobalServerConfig.Password
	}

	// Inherit the global SSH Key Path
	if config.SSHKeyPath == "" {
		config.SSHKeyPath = DeploymentConfig.GlobalServerConfig.SSHKeyPath
	}

	// Package Configuration

	// Inherit the global package selection
	if config.Packages == "" {
		config.Packages = DeploymentConfig.GlobalServerConfig.Packages
	}
}
