package apiserver

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

//Server - will parse a configuration file and passed variables and start the API Server
func Server(path string, port int, insecure bool) error {
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

	e := setAPIEndpoints()
	log.Infof("Starting API server on port %d", conf.Port)
	address := fmt.Sprintf(":%d", conf.Port)
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
			Handler:      e,
		}
		err = srv.ListenAndServeTLS("", "")

	} else {
		err = http.ListenAndServe(address, e)
	}

	if err != nil {
		return err
	}
	return nil
}
