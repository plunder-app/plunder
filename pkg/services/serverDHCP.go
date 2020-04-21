package services

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"

	dhcp "github.com/krolaw/dhcp4"
	log "github.com/sirupsen/logrus"
)

// Lease defines a lease that is allocated to a client
type Lease struct {
	MAC    string    `json:"mac"`  // Client's Physical Address
	Expiry time.Time `json:"time"` // When the lease expires
}

// DHCPSettings -
type DHCPSettings struct {
	IP      net.IP       // Server IP to use
	Options dhcp.Options // Options to send to DHCP Clients
	Start   net.IP       // Start of IP range to distribute

	LeaseRange    int           // Number of IPs to distribute (starting from start)
	LeaseDuration time.Duration // Lease period

	Leases   map[int]Lease // Map to keep track of leases
	UnLeased []Lease       // Map to keep track of unleased devices, and when they were seen
}

// Discover - Is the discovering of a DHCP server on the network and the typical result is an lease "offer"
// Request - The Request is typically the acceptance of a DHCP lease
// Release - A Release is the client notifying that server that the lease is no longer required

//ServeDHCP - Is the function that is called when ever plunder recieves DHCP packets.
func (h *DHCPSettings) ServeDHCP(p dhcp.Packet, msgType dhcp.MessageType, options dhcp.Options) (d dhcp.Packet) {
	mac := strings.ToLower(p.CHAddr().String())
	log.Debugf("DCHP Message Type: [%v] from MAC Address [%s]", msgType, mac)

	// Retrieve teh deployment type
	deploymentType := FindDeploymentConfigFromMac(mac)
	// Convert the : in the mac address to dashes to make life easier
	dashMac := strings.Replace(mac, ":", "-", -1)

	// These packets typicallty will be in one of a number of phases:
	switch msgType {
	case dhcp.Discover:

		// Look for an existing license
		free := -1
		for i, v := range h.Leases { // Find previous lease
			if v.MAC == mac {
				free = i
				goto reply
			}
		}

		// Look for a free lease
		if free = h.freeLease(); free == -1 {
			// No leases available
			return
		}
	reply:
		//TODO - work out why this is here
		h.Options[dhcp.OptionVendorClassIdentifier] = h.IP

		// if DHCP option "OptionUserClass" is set to iPXE then we know that it's default booted to the correct bootloader
		if string(options[dhcp.OptionUserClass]) == "iPXE" {
			// This will ensure that the leasing table is kept updated for when a server was last seen
			h.leaseHander(deploymentType, mac)

			// TODO - This can be removed and left in the REQUEST section only

			// if an entry doesnt exist then drop it to a default type, if not then it has its own specific
			if httpPaths[fmt.Sprintf("%s.ipxe", dashMac)] == "" {
				h.Options[dhcp.OptionBootFileName] = []byte("http://" + h.IP.String() + "/" + deploymentType + ".ipxe")
			} else {
				h.Options[dhcp.OptionBootFileName] = []byte("http://" + h.IP.String() + "/" + dashMac + ".ipxe")
			}

		}

		ipLease := dhcp.IPAdd(h.Start, free)
		log.Debugf("Allocated IP [%s] for [%s]", ipLease.String(), mac)

		return dhcp.ReplyPacket(p, dhcp.Offer, h.IP, ipLease, h.LeaseDuration,
			h.Options.SelectOrderOrAll(options[dhcp.OptionParameterRequestList]))

	case dhcp.Request:

		if server, ok := options[dhcp.OptionServerIdentifier]; ok && !net.IP(server).Equal(h.IP) {
			return nil // Message not for this dhcp server
		}
		reqIP := net.IP(options[dhcp.OptionRequestedIPAddress])
		if reqIP == nil {
			reqIP = net.IP(p.CIAddr())
		}

		if len(reqIP) == 4 && !reqIP.Equal(net.IPv4zero) {
			if leaseNum := dhcp.IPRange(h.Start, reqIP) - 1; leaseNum >= 0 && leaseNum < h.LeaseRange {
				if l, exists := h.Leases[leaseNum]; !exists || l.MAC == p.CHAddr().String() {

					// Specify the new lease
					h.Leases[leaseNum] = Lease{
						MAC:    p.CHAddr().String(),
						Expiry: time.Now().Add(h.LeaseDuration),
					}

					// if DHCP option "OptionUserClass" is set to iPXE then we know that it's default booted to the correct bootloader
					if string(options[dhcp.OptionUserClass]) == "iPXE" {
						// Only Print out this notification if it's from the iPXE Boot loader
						log.Infof("Mac address [%s] is assigned a [%s] deployment type", mac, deploymentType)
					}

					// if an entry doesnt exist then drop it to a default type, if not then it has its own specific
					if httpPaths[fmt.Sprintf("/%s.ipxe", dashMac)] == "" {
						h.Options[dhcp.OptionBootFileName] = []byte("http://" + h.IP.String() + "/" + deploymentType + ".ipxe")
					} else {
						h.Options[dhcp.OptionBootFileName] = []byte("http://" + h.IP.String() + "/" + dashMac + ".ipxe")
					}

					return dhcp.ReplyPacket(p, dhcp.ACK, h.IP, reqIP, h.LeaseDuration,
						h.Options.SelectOrderOrAll(options[dhcp.OptionParameterRequestList]))
				}
			}
		}
		return dhcp.ReplyPacket(p, dhcp.NAK, h.IP, nil, 0, nil)

	case dhcp.Release, dhcp.Decline:
		for i, v := range h.Leases {
			if v.MAC == mac {
				log.Debugf("Releasing lease for [%s]", mac)
				delete(h.Leases, i)
				break
			}
		}
	}
	return nil
}

// leaseHandler() will take care of adding and removing leases based upon use-case
func (h *DHCPSettings) leaseHander(deploymentType, mac string) {
	if deploymentType == "" || deploymentType == "autoBoot" || deploymentType == "reboot" {
		// Create a lease for an un-used server (dont by default)
		newUnleased := Lease{
			MAC:    mac,
			Expiry: time.Now(),
		}
		// False by default
		var macFound bool

		// Look through array
		for i := range h.UnLeased {
			if mac == h.UnLeased[i].MAC {
				h.UnLeased[i].Expiry = time.Now()
				// Found this entry
				macFound = true
			}
		}

		// New entry
		if macFound == false {
			// Update the unleased map with this mac address being seen
			h.UnLeased = append(h.UnLeased, newUnleased)
		}
	}

	// If this mac address has no deployment type for whatever reason, ensure a warning message is presented
	if deploymentType == "" {
		log.Warnf("Mac address[%s] is unknown, not returning an address", mac)
	}
}

func (h *DHCPSettings) freeLease() int {
	now := time.Now()
	b := rand.Intn(h.LeaseRange) // Try random first
	for _, v := range [][]int{{b, h.LeaseRange}, {0, b}} {
		for i := v[0]; i < v[1]; i++ {
			if l, ok := h.Leases[i]; !ok || l.Expiry.Before(now) {
				return i
			}
		}
	}
	return -1
}

// GetLeases - This will retrieve all of the allocated leases from the boot controller
func (c *BootController) GetLeases() *[]Lease {
	var l []Lease
	for i := range c.handler.Leases {
		l = append(l, c.handler.Leases[i])
	}
	return &l
}

// GetUnLeased - This will retrieve all of the un-allocated leases from the boot controller
func (c *BootController) GetUnLeased() *[]Lease {
	if c.handler == nil {
		var emptyLease []Lease
		return &emptyLease
	}
	return &c.handler.UnLeased
}

// DelUnLeased - This will retrieve all of the un-allocated leases from the boot controller
func (c *BootController) DelUnLeased(mac string) {
	if len(c.handler.UnLeased) == 0 {
		return
	}
	for i := range c.handler.UnLeased {
		if mac == c.handler.UnLeased[i].MAC {
			c.handler.UnLeased = append(c.handler.UnLeased[:i], c.handler.UnLeased[i+1:]...)
		}
	}
}
