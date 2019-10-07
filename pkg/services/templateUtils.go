package services

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"strings"
)

// ReadKeyFromFile - will attempt to read an sshkey from a file and populate the struct
func (c *HostConfig) ReadKeyFromFile() (string, error) {
	var buffer []byte
	if _, err := os.Stat(c.SSHKeyPath); !os.IsNotExist(err) {
		buffer, err = ioutil.ReadFile(c.SSHKeyPath)
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

// This will attempt to parse an SSH file in the host config and load it as a base64 encoded KEY
func (c *HostConfig) parseSSH() error {
	// If a file is specified then lets read it and base64 the results (as long as a key doesn't already exist)
	if c.SSHKeyPath != "" && c.SSHKey == "" {
		data, err := c.ReadKeyFromFile()
		if err != nil {
			return err
		}
		c.SSHKey = base64.StdEncoding.EncodeToString([]byte(data))
	}
	return nil
}

// PopulateConfiguration - This will read a deployment configuration and attempt to fill any missing fields from the global config
func (c *HostConfig) PopulateConfiguration(globalConfig HostConfig) {
	// NETWORK CONFIGURATION

	// Inherit the global Gateway
	if c.Gateway == "" {
		c.Gateway = globalConfig.Gateway
	}

	// Inherit the global Subnet
	if c.Subnet == "" {
		c.Subnet = globalConfig.Subnet
	}

	// Inherit the global Name Server (DNS)
	if c.NameServer == "" {
		c.NameServer = globalConfig.NameServer
	}

	if c.Adapter == "" {
		c.Adapter = globalConfig.Adapter
	}

	// REPOSITORY CONFIGURATION

	// Inherit the global Repository address
	if c.RepositoryAddress == "" {
		c.RepositoryAddress = globalConfig.RepositoryAddress
	}

	// Inherit the global Repository Mirror directory (typically /ubuntu)
	if c.MirrorDirectory == "" {
		c.MirrorDirectory = globalConfig.MirrorDirectory
	}

	// USER CONFIGURATION

	// Inherit the global Username
	if c.Username == "" {
		c.Username = globalConfig.Username
	}

	// Inherit the global Password
	if c.Password == "" {
		c.Password = globalConfig.Password
	}

	// Inherit the global SSH Key Path
	if c.SSHKeyPath == "" {
		c.SSHKeyPath = globalConfig.SSHKeyPath
	}

	// Package Configuration

	// Inherit the global package selection
	if c.Packages == "" {
		c.Packages = globalConfig.Packages
	}
}
