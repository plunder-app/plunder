package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
)

// Static URL for retrieving the bootloader
const iPXEURL = "https://boot.ipxe.org/undionly.kpxe"

//////////////////////////////
//
// Helper Functions
//
//////////////////////////////

// GenerateiPXEScript - This will build an iPXE boot script
func GenerateiPXEScript(webserverAddress string, kernel string, initrd string, cmdline string) error {
	script := `#!ipxe
dhcp
echo +-------------------- Plunder -------------------------------
echo | hostname: ${hostname}, next-server: ${next-server}
echo | address.: ${net0/ip}
echo | mac.....: ${net0/mac}  
echo | gateway.: ${net0/gateway} 
echo +------------------------------------------------------------
echo .
kernel http://%s/%s auto=true url=http://%s/${mac:hexhyp}.cfg priority=critical %s 
initrd http://%s/%s
boot`
	// Replace the addresses inline
	buildScript := fmt.Sprintf(script, webserverAddress, kernel, webserverAddress, cmdline, webserverAddress, initrd)

	f, err := os.Create("./plunder.ipxe")
	if err != nil {
		return err
	}
	_, err = f.WriteString(buildScript)
	if err != nil {
		return err
	}
	f.Sync()
	return nil
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
	log.Infoln("Completed\n")
	return nil
}
