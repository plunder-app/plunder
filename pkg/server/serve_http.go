package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/plunder-app/plunder/pkg/utils"
	log "github.com/sirupsen/logrus"
)

// These strings container the generated iPXE details that are passed to the bootloader when the correct url is requested
var preseed, kickstart, anyBoot, reboot string

//
var httpPaths map[string]string

//APIPath returns the URI that is used to interact with the plunder API
func APIPath() string {
	return "/deployment"
}

func (c *BootController) serveHTTP() error {
	// Set this so that other functions that build iPXE files can populate the http server details
	preseed = utils.IPXEPreeseed(*c.HTTPAddress, *c.Kernel, *c.Initrd, *c.Cmdline)
	kickstart = utils.IPXEKickstart(*c.HTTPAddress, *c.Kernel, *c.Initrd, *c.Cmdline)
	anyBoot = utils.IPXEAnyBoot(*c.HTTPAddress, *c.Kernel, *c.Initrd, *c.Cmdline)
	reboot = utils.IPXEReboot()

	docroot, err := filepath.Abs("./")
	if err != nil {
		return err
	}

	//httpRoot := http.FileServer(http.Dir(docroot)).ServeHTTP()
	http.Handle("/", http.FileServer(http.Dir(docroot)))

	http.HandleFunc("/health", HealthCheckHandler)
	http.HandleFunc("/preseed.ipxe", preseedHandler)
	http.HandleFunc("/reboot.ipxe", rebootHandler)
	http.HandleFunc("/kickstart.ipxe", kickstartHandler)
	http.HandleFunc("/anyboot.ipxe", anyBootHandler)

	// Update Endpoints - allow the update of various configuration without restarting
	//http.HandleFunc("/config", kickstartHandler) // TODO
	http.HandleFunc(APIPath(), deploymentHandler)

	return http.ListenAndServe(":80", nil)
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

func anyBootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	// Return the kickstart content
	io.WriteString(w, anyBoot)
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

func deploymentHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	switch r.Method {
	case "GET":
		b, err := json.MarshalIndent(DeploymentConfig, "", "\t")
		if err != nil {
			io.WriteString(w, "<b>Unable to Parse Deployment configuration</b>")
		}
		io.WriteString(w, string(b))
	case "POST":
		if b, err := ioutil.ReadAll(r.Body); err == nil {
			if err != nil {
				errorHTML := fmt.Sprintf("<b>Unable to Parse Deployment configuration</b>\n Error: %s", err.Error())
				io.WriteString(w, errorHTML)
			}
			err := UpdateConfiguration(b)
			if err != nil {
				errorHTML := fmt.Sprintf("<b>Unable to Parse Deployment configuration</b>\n Error: %s", err.Error())
				io.WriteString(w, errorHTML)
			}
		}
	default:
		// Unknown HTTP Method for this endpoint
	}
}
