package bootstraps

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// AnyBoot - This flag when set to true will just boot any kernel/initrd/cmdline configuration
var AnyBoot bool

// DeploymentConfig - contains an accessable "current" configuration
var DeploymentConfig DeploymentConfigurationFile

// DeploymentConfigurationFile - The bootstraps.Configs is used by other packages to manage use case for Mac addresses
type DeploymentConfigurationFile struct {
	GlobalServerConfig ServerConfig               `json:"globalConfig"`
	Deployments        []DeploymentConfigurations `json:"deployments"`
}

// DeploymentConfigurations - is used to parse the files containing all server configurations
type DeploymentConfigurations struct {
	MAC        string       `json:"mac"`
	Deployment string       `json:"deployment"` // Either preseed or kickstart
	Config     ServerConfig `json:"config"`
}

// ServerConfig - Defines how a server will be configured by plunder
type ServerConfig struct {
	Gateway    string `json:"gateway"`
	IPAddress  string `json:"address"`
	Subnet     string `json:"subnet"`
	NameServer string `json:"nameserver"`
	ServerName string `json:"hostname"`
	NTPServer  string `json:"ntpserver"`
	Adapter    string `json:"adapter"`

	Username string `json:"username"`
	Password string `json:"password"`

	RepositoryAddress string `json:"repoaddress"`
	// MirrorDirectory is an Ubuntu specific config
	MirrorDirectory string `json:"mirrordir"`

	// SSHKeyPath will typically be loaded from a file ~/.ssh/id_rsa.pub
	SSHKeyPath string `json:"sshkeypath"`

	// Packages to be installed
	Packages string `json:"packages"`
}

// ReadKeyFromFile - will attempt to read an sshkey from a file and populate the struct
func (config *ServerConfig) ReadKeyFromFile(sshKeyPath string) (string, error) {
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
	json.Unmarshal(configFile, &DeploymentConfig)

	if len(DeploymentConfig.Deployments) == 0 {
		log.Warnln("No deployment configurations found")
	}

	for i := range DeploymentConfig.Deployments {
		var newConfig string
		switch DeploymentConfig.Deployments[i].Deployment {
		case "preseed":
			// Build a preseed configuration and write it to disk
			newConfig = DeploymentConfig.Deployments[i].Config.BuildPreeSeedConfig()

		case "kickstart":
			// Build a kickstart configuration and write it to disk
			newConfig = DeploymentConfig.Deployments[i].Config.BuildKickStartConfig()

		default:
			return fmt.Errorf("Unknown deployment method [%s]", DeploymentConfig.Deployments[i].Deployment)
		}

		// We need to move all ":" to "-" to make life a little easier for filesystems and internet standards
		dashMac := strings.Replace(DeploymentConfig.Deployments[i].MAC, ":", "-", -1)

		// Create a filename from the updated name
		filename := fmt.Sprintf("%s.cfg", dashMac)
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		byteCount, err := f.WriteString(newConfig)
		if err != nil {
			return err
		}
		log.Infof("Written %d bytes to file [%s]", byteCount, filename)
		f.Sync()
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
func (config *ServerConfig) PopulateConfiguration() {
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
