package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

// Static URL for retrieving the bootloader
const iPXEURL = "https://boot.ipxe.org/undionly.kpxe"

// This header is used by all configurations
const iPXEHeader = `#!ipxe
dhcp
echo .
echo .
echo .
echo .
echo +-------------------- Plunder -------------------------------
echo | 
echo |    address.: ${net0/ip}
echo |    mac.....: ${net0/mac}  
echo |    gateway.: ${net0/gateway} 
echo +------------------------------------------------------------
echo .
echo .
echo .
echo .`

//////////////////////////////
//
// Helper Functions
//
//////////////////////////////

// IPXEReboot -
func IPXEReboot() string {
	script := `
echo MAC ADDRESS is set to reboot, plunder will reboot the server in 5 seconds
sleep 5
reboot
`
	return iPXEHeader + script
}

// IPXEAutoBoot -
func IPXEAutoBoot() string {
	script := `
echo Unknown MAC address, PXE boot will keep retrying until configuration changes
:retry_boot
autoboot || goto retry_boot
`
	return iPXEHeader + script
}

// IPXEPreeseed - This will build an iPXE boot script for Debian/Ubuntu
func IPXEPreeseed(webserverAddress, kernel, initrd, cmdline string) string {
	script := `
kernel http://%s/%s auto=true url=http://%s/${mac:hexhyp}.cfg priority=critical %s netcfg/choose_interface=${netX/mac}
initrd http://%s/%s
boot
`
	// Replace the addresses inline
	buildScript := fmt.Sprintf(script, webserverAddress, kernel, webserverAddress, cmdline, webserverAddress, initrd)

	return iPXEHeader + buildScript
}

// IPXEKickstart - This will build an iPXE boot script for RHEL/CentOS
func IPXEKickstart(webserverAddress, kernel, initrd, cmdline string) string {
	script := `
kernel http://%s/%s auto=true url=http://%s/${mac:hexhyp}.cfg priority=critical %s 
initrd http://%s/%s
boot
`
	// Replace the addresses inline
	buildScript := fmt.Sprintf(script, webserverAddress, kernel, webserverAddress, cmdline, webserverAddress, initrd)

	return iPXEHeader + buildScript
}

// IPXEVSphere - This will build an iPXE boot script for VMware vSphere/ESXi
func IPXEVSphere(webserverAddress, kernel, cmdline string) string {
	script := `
kernel http://%s/%s -c http://%s/${mac:hexhyp}.cfg  ks=http://%s/${mac:hexhyp}.ks %s
boot
`
	// Replace the addresses inline
	buildScript := fmt.Sprintf(script, webserverAddress, kernel, webserverAddress, webserverAddress, cmdline)

	return iPXEHeader + buildScript
}

// IPXEAnyBoot - This will build an iPXE boot script for anything wanting to PXE boot
func IPXEAnyBoot(webserverAddress string, kernel, initrd, cmdline string) string {
	script := `
kernel http://%s/%s auto=true url=http://%s/${mac:hexhyp}.cfg %s 
initrd http://%s/%s
boot
`
	// Replace the addresses inline
	buildScript := fmt.Sprintf(script, webserverAddress, kernel, webserverAddress, cmdline, webserverAddress, initrd)

	return iPXEHeader + buildScript
}

// PullPXEBooter - This will attempt to download the iPXE bootloader
func PullPXEBooter() error {
	log.Infoln("Beginning of iPXE download... ")

	// Create the file
	out, err := os.Create("undionly.kpxe")
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(iPXEURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	log.Infoln("Completed")
	return nil
}
