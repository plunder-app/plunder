package ssh

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/plunder-app/plunder/pkg/services"
	"golang.org/x/crypto/ssh"

	log "github.com/sirupsen/logrus"
)

// cachedGlobalKey caches the content of the gloal SSH key to save on excessing file operations
var cachedGlobalKey ssh.AuthMethod

// cachedUsername caches the content of the gloal Username on the basis that a lot of key based ops will share the same user
var cachedUsername string

// The init function will look for the default key and the default user

func init() {
	u, err := user.Current()
	if err != nil {
		log.Warnf("Failed to find current user, if this is overridden by a deployment configuration this error can be ignored")
	}

	// If the above call hasn't errored, then u shouldn't be nil
	if u != nil {
		cachedUsername = u.Username
	}

	cachedGlobalKey, err = findDefaultKey()
	if err != nil {
		log.Warnf("Failed to find default ssh key, if this is overridden by a deployment configuration this error can be ignored")
	}
}

// AddHost will append additional hosts to the host array that the ssh package will use
func AddHost(address, keypath, username string) error {
	sshHost := HostSSHConfig{
		Host: address,
	}

	// If a username exists use that, alternatively use the cached entry
	if username != "" {
		sshHost.User = username
	} else if cachedUsername != "" {
		sshHost.User = cachedUsername
	} else {
		return fmt.Errorf("No username data for SSH authentication has been entered or loaded")
	}

	// Find additional keys that may exist in the same location
	var keys []ssh.AuthMethod

	if keypath != "" {
		key, err := findPrivateKey(keypath)
		if err != nil {
			return err
		}
		keys = append(keys, key)
	} else {
		if cachedGlobalKey != nil {
			keys = append(keys, cachedGlobalKey)
		} else {
			return fmt.Errorf("Host [%s] has no key specified", address)
		}
	}

	sshHost.ClientConfig = &ssh.ClientConfig{User: sshHost.User, Auth: keys, HostKeyCallback: ssh.InsecureIgnoreHostKey()}

	Hosts = append(Hosts, sshHost)

	return nil
}

// ImportHostsFromDeployment - This will parse a deployment (either file or HTTP post)
func ImportHostsFromDeployment(config []byte) error {

	var deployment services.DeploymentConfigurationFile

	err := json.Unmarshal(config, &deployment)
	if err != nil {
		return err
	}

	if len(deployment.Configs) == 0 {
		return fmt.Errorf("No deployment configurations found")
	}

	// Find keys that are in the same places as the public Key
	if deployment.GlobalServerConfig.SSHKeyPath != "" {
		// Find if the private key exists
		cachedGlobalKey, err = findPrivateKey(deployment.GlobalServerConfig.SSHKeyPath)
		if err != nil {
			return err
		}
	} else {
		log.Debugln("No global configuration has been loaded, will default to local users keys")
	}

	// Find a global username to use, in place of an empty config
	if deployment.GlobalServerConfig.Username != "" {
		cachedUsername = deployment.GlobalServerConfig.Username
	} else {
		log.Debugf("No global configuration has been loaded, default to user [%s]", cachedUsername)
	}

	// Parse the deployments
	for i := range deployment.Configs {
		var sshHost HostSSHConfig

		sshHost.Host = deployment.Configs[i].ConfigHost.IPAddress

		if deployment.Configs[i].ConfigHost.Username != "" {
			sshHost.User = deployment.Configs[i].ConfigHost.Username
		} else {
			sshHost.User = deployment.GlobalServerConfig.Username
		}

		// Find additional keys that may exist in the same location
		var keys []ssh.AuthMethod

		if deployment.Configs[i].ConfigHost.SSHKeyPath != "" {
			key, err := findPrivateKey(deployment.Configs[i].ConfigHost.SSHKeyPath)
			if err != nil {
				return err
			}
			keys = append(keys, key)
		} else {
			if cachedGlobalKey != nil {
				keys = append(keys, cachedGlobalKey)
			} else {
				return fmt.Errorf("Host [%s] has no key specified", deployment.Configs[i].ConfigHost.IPAddress)
			}
		}
		sshHost.ClientConfig = &ssh.ClientConfig{User: sshHost.User, Auth: keys, HostKeyCallback: ssh.InsecureIgnoreHostKey()}

		Hosts = append(Hosts, sshHost)
	}

	return nil
}

// findDefaultKey - This will look in the users $HOME/.ssh/ for a key to add
func findDefaultKey() (ssh.AuthMethod, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
	return findPrivateKey(fmt.Sprintf("%s/.ssh/id_rsa", home))
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
