package utils

import (
	"fmt"
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
	return "", "", fmt.Errorf("Unknown interface [%s]", addrName)
}

// FindAllIPAddresses - Will return all IP addresses for a server
func FindAllIPAddresses() ([]net.IP, error) {
	var IPS []net.IP
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip != nil {
				IPS = append(IPS, net.IP(ip))
			}
			// process IP address
		}
	}
	return IPS, nil
}

//ConvertIP -
func ConvertIP(ipAddress string) ([]byte, error) {
	// net.ParseIP has returned IPv6 sized allocations o_O
	fixIP := net.ParseIP(ipAddress)
	if fixIP == nil {
		return nil, fmt.Errorf("Couldn't parse the IP address: %s", ipAddress)
	}
	if len(fixIP) > 4 {
		return fixIP[len(fixIP)-4:], nil
	}
	return fixIP, nil
}
