package services

import (
	"io/ioutil"
	"os"
	"strings"
)

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
