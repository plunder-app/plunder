package ssh

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

// StartConnection -
func (c *HostSSHConfig) StartConnection() (*ssh.Client, error) {
	var err error

	host := c.Host
	if !strings.ContainsAny(c.Host, ":") {
		host = host + ":22"
	}

	log.Debugf("Beginning connection to [%s] with user [%s] and timeout [%d]", c.Host, c.User, c.ClientConfig.Timeout)
	c.Connection, err = ssh.Dial("tcp", host, c.ClientConfig)
	if err != nil {
		return nil, err
	}
	return c.Connection, nil
}

// StopConnection -
func (c *HostSSHConfig) StopConnection() error {
	if c.Connection != nil {
		return c.Connection.Close()
	}
	return fmt.Errorf("Connection not established")
}

// StartSession -
func (c *HostSSHConfig) StartSession() (*ssh.Session, error) {
	var err error
	c.Connection, err = c.StartConnection()
	if err != nil {
		return nil, err
	}
	c.Session, err = c.Connection.NewSession()
	if err != nil {
		return nil, err
	}
	return c.Session, err
}

// StopSession -
func (c *HostSSHConfig) StopSession() {
	if c.Session != nil {
		c.Session.Close()
	}
}

// To string
func (c HostSSHConfig) String() string {
	return c.User + "@" + c.Host
}

//FindHosts - This will take an array of hosts and find the matching HostSSH Configuration
func FindHosts(parlayHosts []string) ([]HostSSHConfig, error) {
	var hostArray []HostSSHConfig
	for x := range parlayHosts {
		found := false
		for y := range Hosts {

			//TODO : Probably needs strings.ToLower() (needs testing)
			if parlayHosts[x] == Hosts[y].Host {
				hostArray = append(hostArray, Hosts[y])
				found = true
				continue
			}
		}
		if found == false {
			return nil, fmt.Errorf("Host [%s] has no SSH credentials", parlayHosts[x])
		}
	}
	return hostArray, nil
}
