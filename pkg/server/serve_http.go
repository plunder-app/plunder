package server

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/thebsdbox/plunder/pkg/utils"

	log "github.com/Sirupsen/logrus"
)

func (c *BootController) serveHTTP() error {
	if _, err := os.Stat("./plunder.ipxe"); os.IsNotExist(err) {
		log.Println("Auto generating ./plunder.ipxe")
		err = utils.GenerateiPXEScript(*c.HTTPAddress, *c.Kernel, *c.Initrd, *c.Cmdline)
		if err != nil {
			return err
		}
	}

	docroot, err := filepath.Abs("./")
	if err != nil {
		return err
	}

	httpHandler := http.FileServer(http.Dir(docroot))

	return http.ListenAndServe(":80", httpHandler)
}
