package services

import (
	"io"
	"net/http"
	"path/filepath"

	"github.com/plunder-app/plunder/pkg/utils"
	log "github.com/sirupsen/logrus"
)

// These strings container the generated iPXE details that are passed to the bootloader when the correct url is requested
var autoBoot, preseed, kickstart, defaultBoot, vsphere, reboot string

// controller Pointer for the config API endpoint handler
var controller *BootController

var mux *http.ServeMux

func (c *BootController) generateBootTypeHanders() {

	// Find the default configuration
	defaultConfig := findBootConfigForName("default")
	if defaultConfig != nil {
		defaultBoot = utils.IPXEPreeseed(*c.HTTPAddress, defaultConfig.Kernel, defaultConfig.Initrd, defaultConfig.Cmdline)
	} else {
		log.Warnf("Found [%d] configurations and no \"default\" configuration", len(c.BootConfigs))
	}

	// If a preeseed configuration has been configured then add it, and create a HTTP endpoint
	preeseedConfig := findBootConfigForName("preseed")
	if preeseedConfig != nil {
		preseed = utils.IPXEPreeseed(*c.HTTPAddress, preeseedConfig.Kernel, preeseedConfig.Initrd, preeseedConfig.Cmdline)

	}

	// If a kickstart configuration has been configured then add it, and create a HTTP endpoint
	kickstartConfig := findBootConfigForName("kickstart")
	if kickstartConfig != nil {
		kickstart = utils.IPXEPreeseed(*c.HTTPAddress, kickstartConfig.Kernel, kickstartConfig.Initrd, kickstartConfig.Cmdline)
	}

	// If a vsphereConfig configuration has been configured then add it, and create a HTTP endpoint
	vsphereConfig := findBootConfigForName("vsphere")
	if vsphereConfig != nil {
		vsphere = utils.IPXEVSphere(*c.HTTPAddress, vsphereConfig.Kernel, vsphereConfig.Cmdline)
	}
}

func (c *BootController) serveHTTP() error {

	// This function will pre-generate the boot handlers for the various boot types
	c.generateBootTypeHanders()

	autoBoot = utils.IPXEAutoBoot()
	reboot = utils.IPXEReboot()

	docroot, err := filepath.Abs("./")
	if err != nil {
		return err
	}

	// Created only once

	// TOTO - alloew this to be customisable
	mux.Handle("/", http.FileServer(http.Dir(docroot)))

	// Boot handlers
	mux.HandleFunc("/health", HealthCheckHandler)
	mux.HandleFunc("/reboot.ipxe", rebootHandler)
	mux.HandleFunc("/autoBoot.ipxe", autoBootHandler)
	mux.HandleFunc("/default.ipxe", rootHandler)
	mux.HandleFunc("/kickstart.ipxe", kickstartHandler)
	mux.HandleFunc("/preseed.ipxe", preseedHandler)
	mux.HandleFunc("/vsphere.ipxe", vsphereHandler)

	// Set the pointer to the boot config
	controller = c

	return http.ListenAndServe(":80", mux)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugf("Requested URL [%s]", r.RequestURI)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	// Return the preseed content
	log.Debugf("Requested URL [%s]", r.URL.Host)
	io.WriteString(w, httpPaths[r.URL.Path])
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

func vsphereHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	// Return the vsphere content
	io.WriteString(w, vsphere)
}

func defaultBootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	// Return the default boot content
	io.WriteString(w, defaultBoot)
}

func rebootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	// Return the reboot content
	io.WriteString(w, reboot)
}

func autoBootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	// Return the reboot content
	io.WriteString(w, autoBoot)
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
