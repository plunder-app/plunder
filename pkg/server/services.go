package server

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/plunder-app/plunder/pkg/utils"

	dhcp "github.com/krolaw/dhcp4"
	"github.com/krolaw/dhcp4/conn"
)

// This is needed by other functions to build strings
var httpAddress string

// BootController contains the settings that define how the remote boot will
// behave
type BootController struct {
	AdapterName *string `json:"adapter"` // A physical adapter to bind to e.g. en0, eth0

	// Servers
	EnableDHCP       *bool   `json:"enableDHCP"`     // Enable Server
	DHCPAddress      *string `json:"addressDHCP"`    // Should ideally be the IP of the adapter
	DHCPStartAddress *string `json:"startDHCP"`      // The first available DHCP address
	DHCPLeasePool    *int    `json:"leasePoolDHCP"`  // Size of the IP Address pool
	DHCPGateway      *string `json:"gatewayDHCP"`    // Gatewway to advertise
	DHCPDNS          *string `json:"nameserverDHCP"` // DNS server to advertise

	EnableTFTP  *bool   `json:"enableTFTP"`  // Enable Server
	TFTPAddress *string `json:"addressTFTP"` // Should ideally be the IP of the adapter

	EnableHTTP  *bool   `json:"enableHTTP"`  // Enable Server
	HTTPAddress *string `json:"addressHTTP"` // Should ideally be the IP of the adapter

	// TFTP Configuration
	PXEFileName *string `json:"pxePath"` // undionly.kpxe

	// iPXE file settings - exported
	Kernel  *string `json:"kernelPath"`
	Initrd  *string `json:"initrdPath"`
	Cmdline *string `json:"cmdline"`

	handler *DHCPSettings
}

// DeploymentConfig - contains an accessible "current" configuration
var DeploymentConfig DeploymentConfigurationFile

// DeploymentConfigurationFile - The bootstraps.Configs is used by other packages to manage use case for Mac addresses
type DeploymentConfigurationFile struct {
	GlobalServerConfig HostConfig                 `json:"globalConfig"`
	Deployments        []DeploymentConfigurations `json:"deployments"`
}

// DeploymentConfigurations - is used to parse the files containing all server configurations
type DeploymentConfigurations struct {
	MAC string `json:"mac"`

	// iPXE file settings - exported
	Kernel  string `json:"kernelPath"`
	Initrd  string `json:"initrdPath"`
	Cmdline string `json:"cmdline"`

	Deployment string     `json:"deployment"` // Either preseed or kickstart
	Config     HostConfig `json:"config"`
}

// HostConfig - Defines how a server will be configured by plunder
type HostConfig struct {
	Gateway    string `json:"gateway"`
	IPAddress  string `json:"address"`
	Subnet     string `json:"subnet"`
	NameServer string `json:"nameserver"`
	ServerName string `json:"hostname"`
	NTPServer  string `json:"ntpserver"`
	Adapter    string `json:"adapter"`
	SwapEnable bool   `json:"swapEnabled"`

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

// StartServices - This will start all of the enabled services
func (c *BootController) StartServices(deployment []byte) {
	log.Infof("Starting Remote Boot Services, press CTRL + c to stop")

	if *c.EnableDHCP == true {
		c.handler = &DHCPSettings{}
		c.handler.IP = utils.ConvertIP(*c.DHCPAddress)
		c.handler.Start = utils.ConvertIP(*c.DHCPStartAddress)

		c.handler.LeaseDuration = 2 * time.Hour //TODO, make time modifiable
		c.handler.LeaseRange = *c.DHCPLeasePool
		c.handler.Leases = make(map[int]lease, *c.DHCPLeasePool)

		c.handler.Options = dhcp.Options{
			dhcp.OptionSubnetMask:       []byte{255, 255, 255, 0},
			dhcp.OptionRouter:           []byte(utils.ConvertIP(*c.DHCPGateway)),
			dhcp.OptionDomainNameServer: []byte(utils.ConvertIP(*c.DHCPDNS)),
			dhcp.OptionBootFileName:     []byte(*c.PXEFileName),
		}

		log.Debugf("\nServer IP:\t%s\nAdapter:\t%s\nStart Address:\t%s\nPool Size:\t%d\n", *c.DHCPAddress, *c.AdapterName, *c.DHCPStartAddress, *c.DHCPLeasePool)

		go func() {
			log.Println("RemoteBoot => Starting DHCP")

			newConnection, err := conn.NewUDP4FilterListener(*c.AdapterName, ":67")
			if err != nil {
				log.Fatalf("%v", err)
			}
			//Close the connection when we're tidying up
			defer newConnection.Close()
			err = dhcp.Serve(newConnection, c.handler)
			log.Fatalf("%v", err)
		}()
	}

	if *c.EnableTFTP == true {
		go func() {
			log.Println("RemoteBoot => Starting TFTP")
			log.Debugf("\nServer IP:\t%s\nPXEFile:\t%s\n", *c.TFTPAddress, *c.PXEFileName)

			err := c.serveTFTP()
			log.Fatalf("%v", err)
		}()
	}

	if *c.EnableHTTP == true {
		if c.Kernel == nil || *c.Kernel == "" {
			log.Warn("No Kernel specified in configuration")
		}
		if c.Kernel == nil || *c.Initrd == "" {
			log.Warn("No Initrd specified in configuration")
		}

		httpAddress = *c.HTTPAddress
		httpPaths = make(map[string]string)

		err := UpdateConfiguration(deployment)
		if err != nil {
			log.Fatalf("%v", err)
		}
		go func() {
			log.Println("RemoteBoot => Starting HTTP")
			err := c.serveHTTP()
			log.Fatalf("%v", err)
		}()
	}

	utils.WaitForCtrlC()
}
