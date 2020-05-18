package services

import (
	"encoding/json"
	"net"

	"github.com/plunder-app/BOOTy/pkg/plunderclient/types"
	"github.com/vishvananda/netlink"
)

//BuildBOOTYconfig - Creates a new presseed configuration using the passed data
func (config *HostConfig) BuildBOOTYconfig() string {
	a := types.BootyConfig{}

	// set the required action
	a.Action = config.BOOTYAction

	// Default to false if not in configuration
	if config.Compressed == nil {
		a.Compressed = false
	} else {
		a.Compressed = *config.Compressed
	}

	// Parse the strings
	subnet := net.ParseIP(config.Subnet)
	ip := net.ParseIP(config.IPAddress)

	// Change into a cidr
	cidr := net.IPNet{
		IP:   ip,
		Mask: subnet.DefaultMask(),
	}
	addr, _ := netlink.ParseAddr(cidr.String())

	// Set configuration
	if addr != nil {
		a.Address = addr.String()
		a.Gateway = config.Gateway
	}

	// READ
	a.DestinationDevice = config.DestinationDevice
	a.SourceImage = config.SourceImage
	// WRITE
	a.DesintationAddress = config.DestinationAddress
	a.SourceDevice = config.SourceDevice

	// Default to false if not in configuration
	if config.GrowPartition == nil {
		a.GrowPartition = 0
	} else {
		a.GrowPartition = *config.GrowPartition
	}
	a.LVMRootName = config.LVMRootName

	// Default to false if not in configuration
	if config.ShellOnFail == nil {
		a.DropToShell = false
	} else {
		a.DropToShell = *config.ShellOnFail
	}

	a.DropToShell = *config.ShellOnFail
	a.NameServer = config.NameServer

	b, _ := json.Marshal(a)
	return string(b)
}
