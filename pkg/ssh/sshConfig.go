package ssh

import (
	"golang.org/x/crypto/ssh"
)

// Hosts - The array of all hosts once loaded
var Hosts []HostSSHConfig

// HostSSHConfig - The struct of an SSH connection
type HostSSHConfig struct {
	Host         string
	User         string
	Timeout      int
	ClientConfig *ssh.ClientConfig
	Session      *ssh.Session
}

// SetPassword - Turn a password string into an SSH auth method
func SetPassword(password string) []ssh.AuthMethod {
	return []ssh.AuthMethod{ssh.Password(password)}
}
