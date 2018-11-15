package utils

import (
	"log"
	"net"
)

// FindIPAddress - this will find the address associated with an adapter
func FindIPAddress(addrName string) (string, string, error) {
	var address string
	list, err := net.Interfaces()
	if err != nil {
		return "", "", err
	}

	for _, iface := range list {
		addrs, err := iface.Addrs()
		if err != nil {
			return "", "", err
		}
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					address = ipnet.IP.String()
					// If we're not searching for a specific adapter return the first one
					if addrName == "" {
						return iface.Name, address, nil
					} else
					// If this is the correct adapter return the details
					if iface.Name == addrName {
						return iface.Name, address, nil
					}
				}
			}
		}

	}
	return "", address, nil
}

//ConvertIP -
func ConvertIP(ipAddress string) []byte {
	// net.ParseIP has returned IPv6 sized allocations o_O
	fixIP := net.ParseIP(ipAddress)
	if fixIP == nil {
		log.Fatalf("Couldn't parse the IP address: %s\n", ipAddress)
	}
	if len(fixIP) > 4 {
		return fixIP[len(fixIP)-4:]
	}
	return fixIP
}
