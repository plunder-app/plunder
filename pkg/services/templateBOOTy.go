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
	a.Action = types.WriteImage
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

	a.DestinationDevice = config.DestinationDevice
	a.SourceImage = config.SourceImage
	a.GrowPartition = *config.GrowPartition
	a.LVMRootName = config.LVMRootName
	a.DropToShell = *config.ShellOnFail
	a.NameServer = config.NameServer
	b, _ := json.Marshal(a)
	return string(b)
}
