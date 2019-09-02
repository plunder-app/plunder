package services

// TYPE DEFINITIONS Below

// BootController contains the settings that define how the remote boot will
type BootController struct {
	AdapterName *string `json:"adapter"` // A physical adapter to bind to e.g. en0, eth0

	// Servers
	EnableDHCP *bool `json:"enableDHCP"` // Enable Server
	//DHCP Configuration
	DHCPConfig dhcpConfig `json:"dhcpConfig,omitempty"`

	// TFTP / HTTP configuration
	EnableTFTP  *bool   `json:"enableTFTP"`  // Enable Server
	TFTPAddress *string `json:"addressTFTP"` // Should ideally be the IP of the adapter
	EnableHTTP  *bool   `json:"enableHTTP"`  // Enable Server
	HTTPAddress *string `json:"addressHTTP"` // Should ideally be the IP of the adapter

	// TFTP Configuration
	PXEFileName *string `json:"pxePath"` // undionly.kpxe

	// Boot Configuration
	BootConfigs []BootConfig `json:"bootConfigs"` // Array of kernel configurations

	handler *DHCPSettings
}

type dhcpConfig struct {
	DHCPAddress      *string `json:"addressDHCP"`    // Should ideally be the IP of the adapter
	DHCPStartAddress *string `json:"startDHCP"`      // The first available DHCP address
	DHCPLeasePool    *int    `json:"leasePoolDHCP"`  // Size of the IP Address pool
	DHCPGateway      *string `json:"gatewayDHCP"`    // Gatewway to advertise
	DHCPDNS          *string `json:"nameserverDHCP"` // DNS server to advertise
}

// BootConfig defines a named configuration for booting
type BootConfig struct {
	ConfigName string `json:"configName"`

	// iPXE file settings - exported
	Kernel  string `json:"kernelPath"`
	Initrd  string `json:"initrdPath"`
	Cmdline string `json:"cmdline"`

	// ISO Reader settings
	ISOPath   string `json:"isoPath,omitempty"`
	ISOPrefix string `json:"isoPrefix,omitempty"`
}

// DeploymentConfigurationFile - The bootstraps.Configs is used by other packages to manage use case for Mac addresses
type DeploymentConfigurationFile struct {
	GlobalServerConfig HostConfig         `json:"globalConfig"`
	Configs            []DeploymentConfig `json:"deployments"`
}

// DeploymentConfig - is used to parse the files containing all server configurations
type DeploymentConfig struct {
	MAC        string     `json:"mac"`
	ConfigName string     `json:"bootConfigName,omitempty"` // To be discovered in the controller BootConfig array
	ConfigBoot BootConfig `json:"bootConfig,omitempty"`     // Array of kernel configurations
	ConfigHost HostConfig `json:"config"`
}

// HostConfig - Defines how a server will be configured by plunder
type HostConfig struct {

	// Not required for the global configuration
	Adapter    string `json:"adapter,omitempty"`  // Adapter to be configured with networking address
	IPAddress  string `json:"address,omitempty"`  // Allocated IP address for a host (ignored for global)
	ServerName string `json:"hostname,omitempty"` // Hostname to be applied to a server

	// Typically shared details
	Gateway    string `json:"gateway,omitempty"` // Default Gateway
	Subnet     string `json:"subnet,omitempty"`  // Subnet to be used for the host
	NameServer string `json:"nameserver,omitempty"`
	NTPServer  string `json:"ntpserver,omitempty"` // Time Server to be used
	SwapEnable bool   `json:"swapEnabled,omitempty"`

	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`

	RepositoryAddress string `json:"repoaddress,omitempty"`
	// MirrorDirectory is an Ubuntu specific config
	MirrorDirectory string `json:"mirrordir,omitempty"`

	// SSHKeyPath will typically be referenced from a file ~/.ssh/id_rsa.pub
	SSHKeyPath string `json:"sshkeypath,omitempty"`
	// SSHKey is a full SSH Key
	SSHKey string `json:"sshkey,omitempty"`

	// Packages to be installed
	Packages string `json:"packages,omitempty"`
}
