package apiserver

import (
	"github.com/gorilla/mux"
)

// setAPIEndpoints defines all of the API end points for Plunder
func setAPIEndpoints() *mux.Router {

	// Create a new router
	router := mux.NewRouter()

	// Define the retrieval endpoints for Plunder Server configuration
	router.HandleFunc("/config", getConfig).Methods("GET")

	// Define the creation endpoints for Plunder Server Configuration
	router.HandleFunc("/config", putConfig).Methods("POST")

	// Define the retrieval endpoints for Plunder Deployment configuration
	router.HandleFunc("/deployment", getDeployment).Methods("GET")

	// Define the creation and modification endpoints for Plunder Deployment configuration
	router.HandleFunc("/deployment", putDeployment).Methods("POST")
	router.HandleFunc("/deployment/{ID}", nil).Methods("POST")
	router.HandleFunc("/deployment/{ID}", nil).Methods("DELETE")

	return router
}
