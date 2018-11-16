package bootstraps

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
)

// ServerConfig - Defines how a server will be configured by plunder
type ServerConfig struct {
	Gateway    string `json:"gateway"`
	IPAddress  string `json:"address"`
	Subnet     string `json:"subnet"`
	NameServer string `json:"nameserver"`

	NTPServer string `json:"ntpserver"`

	Username string `json:"username"`
	Password string `json:"password"`

	RepositoryAddress string `json:"repoaddress"`
	// MirrorDirectory is an Ubuntu specific config
	MirrorDirectory string `json:"mirrordir"`

	// SSHKeyPath will typically be loaded from a file ~/.ssh/id_rsa.pub
	SSHKeyPath string `json:"sshkeypath"`
}

// ConfigFile - is used to parse the files containing all server configurations
type ConfigFile struct {
	MAC        string       `json:"mac"`
	Deployment string       `json:"deployment"` // Either preseed or kickstart
	Config     ServerConfig `json:"config"`
}

// ReadKeyFromFile - will attempt to read an sshkey from a file and populate the struct
func (config *ServerConfig) ReadKeyFromFile(sshKeyPath string) error {
	var buffer []byte
	if _, err := os.Stat(sshKeyPath); !os.IsNotExist(err) {
		buffer, err = ioutil.ReadFile(sshKeyPath)
		if err != nil {
			// Unable to read the file
			return err
		}
	} else {
		// File doesn't exist
		return err
	}
	config.SSHKeyPath = string(buffer)
	return nil
}

// GenerateConfigFiles will read a configuration file and build the iPXE files needed
func GenerateConfigFiles(configFile string) error {
	var configs []ConfigFile

	// Check the actual path from the string
	if _, err := os.Stat(configFile); !os.IsNotExist(err) {
		configFile, err := ioutil.ReadFile(configFile)
		if err != nil {
			return err
		}
		json.Unmarshal(configFile, &configs)
	} else {
		return fmt.Errorf("Unable to open [%s]", configFile)
	}
	for i := range configs {
		var newConfig string
		switch configs[i].Deployment {
		case "preseed":
			// Build a preseed configuration and write it to disk
			newConfig = configs[i].Config.BuildPreeSeedConfig()

		case "kickstart":
			// Build a kickstart configuration and write it to disk
			newConfig = configs[i].Config.BuildKickStartConfig()

		default:
			return fmt.Errorf("Unknown deployment method [%s]", configs[i].Deployment)
		}

		filename := fmt.Sprintf("%s.cfg", configs[i].MAC)
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
