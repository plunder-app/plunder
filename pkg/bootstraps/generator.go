package bootstraps

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// Configs - The bootstraps.Configs is used by other packages to manage use case for Mac addresses
var Configs []ConfigFile

// ServerConfig - Defines how a server will be configured by plunder
type ServerConfig struct {
	Gateway    string `json:"gateway"`
	IPAddress  string `json:"address"`
	Subnet     string `json:"subnet"`
	NameServer string `json:"nameserver"`
	ServerName string `json:"hostname"`
	NTPServer  string `json:"ntpserver"`

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

	// Check the actual path from the string
	if _, err := os.Stat(configFile); !os.IsNotExist(err) {
		configFile, err := ioutil.ReadFile(configFile)
		if err != nil {
			return err
		}
		json.Unmarshal(configFile, &Configs)
	} else {
		return fmt.Errorf("Unable to open [%s]", configFile)
	}

	if len(Configs) == 0 {
		log.Warnln("No deployment configurations found")
	}

	for i := range Configs {
		var newConfig string
		switch Configs[i].Deployment {
		case "preseed":
			// Build a preseed configuration and write it to disk
			newConfig = Configs[i].Config.BuildPreeSeedConfig()

		case "kickstart":
			// Build a kickstart configuration and write it to disk
			newConfig = Configs[i].Config.BuildKickStartConfig()

		default:
			return fmt.Errorf("Unknown deployment method [%s]", Configs[i].Deployment)
		}

		// We need to move all ":" to "-" to make life a little easier for filesystems and internet standards
		dashMac := strings.Replace(Configs[i].MAC, ":", "-", -1)

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
	if len(Configs) == 0 {
		// No configurations have been loaded
		log.Warnln("Attempted to perform Mac Address lookup, however no configurations have been loaded")
		return ""
	}
	for i := range Configs {
		log.Debugf("Comparing [%s] to [%s]", mac, Configs[i].MAC)
		if mac == Configs[i].MAC {
			return Configs[i].Deployment
		}
	}
	return ""
}
