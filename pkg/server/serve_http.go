package server

import (
	"io"
	"net/http"
	"path/filepath"

	"github.com/thebsdbox/plunder/pkg/utils"
)

// These strings container the generated iPXE details that are passed to the bootloader when the correct url is requested
var preseed, kickstart, reboot string

func (c *BootController) serveHTTP() error {
	// if _, err := os.Stat("./plunder.ipxe"); os.IsNotExist(err) {
	// 	log.Println("Auto generating ./plunder.ipxe")
	// 	err = utils.IPXEPreeseed(*c.HTTPAddress, *c.Kernel, *c.Initrd, *c.Cmdline)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	preseed = utils.IPXEPreeseed(*c.HTTPAddress, *c.Kernel, *c.Initrd, *c.Cmdline)
	kickstart = utils.IPXEKickstart(*c.HTTPAddress, *c.Kernel, *c.Initrd, *c.Cmdline)
	reboot = utils.IPXEReboot()
	docroot, err := filepath.Abs("./")
	if err != nil {
		return err
	}

	//httpRoot := http.FileServer(http.Dir(docroot)).ServeHTTP()
	http.Handle("/", http.FileServer(http.Dir(docroot)))
	//http.HandleFunc("/", httpRoot)

	http.HandleFunc("/health", HealthCheckHandler)
	http.HandleFunc("/preseed.ipxe", preseedHandler)
	http.HandleFunc("/reboot.ipxe", rebootHandler)
	http.HandleFunc("/kickstart.ipxe", kickstartHandler)

	return http.ListenAndServe(":80", nil)
}

func preseedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	// Return the preseed content
	io.WriteString(w, preseed)
}

func kickstartHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	// Return the kickstart content
	io.WriteString(w, kickstart)
}

func rebootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	// Return the reboot content
	io.WriteString(w, reboot)
}

// HealthCheckHandler -
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)
}
