package apiserver

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

//Server -
func Server(port int, insecure bool) error {
	e := setAPIEndpoints()
	log.Infof("Starting API server on port %d", port)
	var err error
	address := fmt.Sprintf(":%d", port)
	if insecure == false {
		err = http.ListenAndServeTLS(address, "plunder.pem", "plunder.key", e)
	} else {
		err = http.ListenAndServe(address, e)
	}

	if err != nil {
		return err
	}
	return nil
}
