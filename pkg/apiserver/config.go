package apiserver

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/ghodss/yaml"
)

// ClientConfig is the structure of an expected configuration for pldctl
type ClientConfig struct {
	Address    string `json:"address,omitempty"`
	Port       int    `json:"port"`
	ClientCert string `json:"cert"`
}

// ServerConfig is the structure of an expected configuration for pldctl
type ServerConfig struct {
	ClientConfig
	ServerKey string `json:"key"`
}

//OpenClientConfig will open and parse a Plunder server configuration file
func OpenClientConfig(path string) (*ClientConfig, error) {
	var c ClientConfig
	// Create a CA certificate pool and add cert.pem to it
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := yaml.YAMLToJSON(b)
	if err == nil {
		// If there were no errors then the YAML => JSON was successful, no attempt to unmarshall
		err = json.Unmarshal(jsonBytes, &c)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse configuration as either yaml or json")
		}
	} else {
		// Couldn't parse the yaml to JSON
		// Attempt to parse it as JSON
		err = json.Unmarshal(b, &c)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse configuration as either yaml or json")
		}
	}
	return &c, nil
}

//OpenServerConfig will open and parse a Plunder server configuration file
func OpenServerConfig(path string) (*ServerConfig, error) {
	var s ServerConfig
	// Create a CA certificate pool and add cert.pem to it
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := yaml.YAMLToJSON(b)
	if err == nil {
		// If there were no errors then the YAML => JSON was successful, no attempt to unmarshall
		err = json.Unmarshal(jsonBytes, &s)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse configuration as either yaml or json")
		}
	} else {
		// Couldn't parse the yaml to JSON
		// Attempt to parse it as JSON
		err = json.Unmarshal(b, &s)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse configuration as either yaml or json")
		}
	}
	return &s, nil
}

// WriteServerConfig - will write out the server configuration for the API Server
func WriteServerConfig(path, hostname, address string, port int, cert, key []byte) error {
	var s ServerConfig

	// base64 the certificates
	encodedKey := base64.StdEncoding.EncodeToString(key)
	encodedCert := base64.StdEncoding.EncodeToString(cert)

	// Add the encoded certificates to the struct
	s.ClientCert = encodedCert
	s.ServerKey = encodedKey

	// Add the port for automated startup
	s.Port = port

	// Marshall to yaml
	b, err := yaml.Marshal(s)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, b, 0600)
	if err != nil {
		return err
	}
	return nil
}

// WriteClientConfig - will write out the server configuration for the API Server
func WriteClientConfig(path, address string, s *ServerConfig) error {
	var c ClientConfig

	// Add the encoded certificates to the struct
	c.ClientCert = s.ClientCert

	// Add the host information for automated startup
	c.Port = s.Port
	c.Address = address

	// Marshall client configuration to yaml
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, b, 0600)
	if err != nil {
		return err
	}
	return nil

}

//GetServerAddressURL will retrieve a parsed URL
func (c *ClientConfig) GetServerAddressURL() *url.URL {
	var plunderURL url.URL
	plunderURL.Scheme = "https"
	// Build a url
	plunderURL.Host = fmt.Sprintf("%s:%d", c.Address, +c.Port)
	return &plunderURL
}

func retrieveCert(cert string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(cert)
}

// RetrieveKey will decode the base64 certificate
func (s *ServerConfig) RetrieveKey() ([]byte, error) {
	return retrieveCert(s.ServerKey)
}

// RetrieveClientCert will decode the base64 certificate
func (s *ServerConfig) RetrieveClientCert() ([]byte, error) {
	return retrieveCert(s.ClientCert)
}

// RetrieveClientCert will decode the base64 certificate
func (c *ClientConfig) RetrieveClientCert() ([]byte, error) {
	return retrieveCert(c.ClientCert)
}
