package apiserver

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

var endpoints *mux.Router

func init() {
	// Initialise a new HTTP Router so that connections can be created before the API server starts, any registered will be added to this router and
	// evaluated once the API Server starts
	endpoints = mux.NewRouter()
}

//StartAPIServer - will parse a configuration file and passed variables and start the API Server
func StartAPIServer(path string, port int, insecure bool) error {
	// Open and Parse the server configuration
	conf, err := OpenServerConfig(path)
	if err != nil {
		log.Warnln(err)
		if insecure == false {
			log.Warningln("Secure server enabled, but no certificates have been loaded [no communication to API server is possible]")
		}
		// Create a blank server config as one wont be returned by the above OpenFile
		conf = &ServerConfig{}
	}
	if port != 0 {
		conf.Port = port
	}

	log.Infof("Starting API server on port %d", conf.Port)
	address := fmt.Sprintf(":%d", conf.Port)

	// Add the apiserver end point
	AddDynamicEndpoint("/api",
		"/api",
		"Endpoint for interacting with the api server",
		"apis",
		http.MethodGet,
		getAPIs)

	// Add the apiserver end point
	AddDynamicEndpoint("/api/{function}/{method}",
		"/api",
		"Endpoint for interacting with the api server",
		"apiFunctions",
		http.MethodGet,
		getAPIFunctionMethod)

	// Begin the start of a secure endpoint (TODO)
	if insecure == false {
		cert, err := conf.RetrieveClientCert()
		if err != nil {
			return err
		}
		key, err := conf.RetrieveKey()
		if err != nil {
			return err
		}
		certPair, err := tls.X509KeyPair(cert, key)
		cfg := &tls.Config{Certificates: []tls.Certificate{certPair}}
		srv := &http.Server{
			TLSConfig:    cfg,
			ReadTimeout:  time.Minute,
			WriteTimeout: time.Minute,
			Addr:         address,
			Handler:      endpoints,
		}

		return srv.ListenAndServeTLS("", "")

	}

	// Start an insecure http server (TODO - warning)
	return http.ListenAndServe(address, endpoints)

}
