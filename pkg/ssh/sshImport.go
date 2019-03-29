package ssh

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/plunder-app/plunder/pkg/server"
	"golang.org/x/crypto/ssh"
)

// cachedGlobalKey caches the content of the gloal SSH key to save on excessing file operations
var cachedGlobalKey ssh.AuthMethod

// ImportHostsFromDeployment - This will import a list of hosts from a file
func ImportHostsFromDeployment(config []byte) error {

	var deployment server.DeploymentConfigurationFile

	err := json.Unmarshal(config, &deployment)
	if err != nil {
		return err
	}

	if len(deployment.Deployments) == 0 {
		return fmt.Errorf("No deployment configurations found")
	}

	// Find keys that are in the same places as the public Key
	if deployment.GlobalServerConfig.SSHKeyPath != "" {
		// Find if the private key exists
		cachedGlobalKey, err = findPrivateKey(deployment.GlobalServerConfig.SSHKeyPath)
		if err != nil {
			return err
		}
	}

	// Parse the deployments
	for i := range deployment.Deployments {
		var sshHost HostSSHConfig

		sshHost.Host = deployment.Deployments[i].Config.IPAddress

		if deployment.Deployments[i].Config.Username != "" {
			sshHost.User = deployment.Deployments[i].Config.Username
		} else {
			sshHost.User = deployment.GlobalServerConfig.Username
		}

		// Find additional keys that may exist in the same location
		var keys []ssh.AuthMethod

		if deployment.Deployments[i].Config.SSHKeyPath != "" {
			key, err := findPrivateKey(deployment.Deployments[i].Config.SSHKeyPath)
			if err != nil {
				return err
			}
			keys = append(keys, key)
		} else {
			if cachedGlobalKey != nil {
				keys = append(keys, cachedGlobalKey)
			} else {
				return fmt.Errorf("Host [%s] has no key specified", deployment.Deployments[i].Config.IPAddress)
			}
		}
		sshHost.ClientConfig = &ssh.ClientConfig{User: sshHost.User, Auth: keys, HostKeyCallback: ssh.InsecureIgnoreHostKey()}

		Hosts = append(Hosts, sshHost)
	}

	return nil
}

// readKeyFile - Reads a public key from a file
func findPrivateKey(publicKey string) (ssh.AuthMethod, error) {
	// Typically turn id_rsa.pub -> id_rsa
	privateKey := strings.TrimSuffix(publicKey, filepath.Ext(publicKey))

	b, err := ioutil.ReadFile(privateKey)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(b)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}

// readKeyFile - Reads a public key from a file
func readKeyFile(keyfile string) (ssh.AuthMethod, error) {
	b, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(b)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}

// ReadKeyFiles - will read an array of keys from disk
func ReadKeyFiles(keyFiles []string) ([]ssh.AuthMethod, error) {
	methods := []ssh.AuthMethod{}

	for _, keyname := range keyFiles {
		pkey, err := readKeyFile(keyname)
		if err != nil {
			return nil, err
		}
		if pkey != nil {
			methods = append(methods, pkey)
		}
	}

	return methods, nil
}
