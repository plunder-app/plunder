package server

import (
	"math/rand"
	"net"
	"time"

	log "github.com/Sirupsen/logrus"
	dhcp "github.com/krolaw/dhcp4"
)

type lease struct {
	nic    string    // Client's Addr
	expiry time.Time // When the lease expires
}

// DHCPSettings -
type DHCPSettings struct {
	IP      net.IP       // Server IP to use
	Options dhcp.Options // Options to send to DHCP Clients
	Start   net.IP       // Start of IP range to distribute

	LeaseRange    int           // Number of IPs to distribute (starting from start)
	LeaseDuration time.Duration // Lease period
	Leases        map[int]lease // Map to keep track of leases
}

//ServeDHCP -
func (h *DHCPSettings) ServeDHCP(p dhcp.Packet, msgType dhcp.MessageType, options dhcp.Options) (d dhcp.Packet) {
	log.Infof("DCHP Message: %v", msgType)
	switch msgType {

	case dhcp.Discover:
		if string(options[77]) != "" {
			if string(options[77]) == "iPXE" {
				h.Options[67] = []byte("http://" + h.IP.String() + "/plunder.ipxe")
			}
		}
		free, nic := -1, p.CHAddr().String()
		for i, v := range h.Leases { // Find previous lease
			if v.nic == nic {
				free = i
				goto reply
			}
		}
		if free = h.freeLease(); free == -1 {
			return
		}
	reply:
		h.Options[60] = h.IP
		ipLease := dhcp.IPAdd(h.Start, free)

		return dhcp.ReplyPacket(p, dhcp.Offer, h.IP, ipLease, h.LeaseDuration,
			h.Options.SelectOrderOrAll(options[dhcp.OptionParameterRequestList]))

	case dhcp.Request:
		log.Debugf("Allocated IP for [%s]", p.CHAddr())

		if server, ok := options[dhcp.OptionServerIdentifier]; ok && !net.IP(server).Equal(h.IP) {
			return nil // Message not for this dhcp server
		}
		reqIP := net.IP(options[dhcp.OptionRequestedIPAddress])
		if reqIP == nil {
			reqIP = net.IP(p.CIAddr())
		}

		if len(reqIP) == 4 && !reqIP.Equal(net.IPv4zero) {
			if leaseNum := dhcp.IPRange(h.Start, reqIP) - 1; leaseNum >= 0 && leaseNum < h.LeaseRange {
				if l, exists := h.Leases[leaseNum]; !exists || l.nic == p.CHAddr().String() {
					h.Leases[leaseNum] = lease{nic: p.CHAddr().String(), expiry: time.Now().Add(h.LeaseDuration)}
					return dhcp.ReplyPacket(p, dhcp.ACK, h.IP, reqIP, h.LeaseDuration,
						h.Options.SelectOrderOrAll(options[dhcp.OptionParameterRequestList]))
				}
			}
		}
		return dhcp.ReplyPacket(p, dhcp.NAK, h.IP, nil, 0, nil)

	case dhcp.Release, dhcp.Decline:
		nic := p.CHAddr().String()
		for i, v := range h.Leases {
			if v.nic == nic {
				delete(h.Leases, i)
				break
			}
		}
	}
	return nil
}

func (h *DHCPSettings) freeLease() int {
	now := time.Now()
	b := rand.Intn(h.LeaseRange) // Try random first
	for _, v := range [][]int{{b, h.LeaseRange}, {0, b}} {
		for i := v[0]; i < v[1]; i++ {
			if l, ok := h.Leases[i]; !ok || l.expiry.Before(now) {
				return i
			}
		}
	}
	return -1
}
