package apiserver

import (
	"github.com/gorilla/mux"
)

// setAPIEndpoints defines all of the API end points for Plunder
func setAPIEndpoints() *mux.Router {
	// Create a new router
	router := mux.NewRouter()

	// ------------------------------------
	// Large configuration management
	// ------------------------------------

	// Define the retrieval endpoints for Plunder Server configuration
	router.HandleFunc("/config", getConfig).Methods("GET")

	// Define the creation endpoints for Plunder Server Configuration
	router.HandleFunc("/config", postConfig).Methods("POST")

	// Define the retrieval endpoints for Plunder Deployment configuration
	router.HandleFunc("/deployments", getDeployments).Methods("GET")

	// Define the retrieval endpoints for Plunder Deployment configuration
	router.HandleFunc("/deployments", postDeployments).Methods("POST")

	// ------------------------------------
	// Specific configuration management
	// ------------------------------------

	// Define the creation endpoints for Plunder Server Boot Configuration
	router.HandleFunc("/config/{id}", postBootConfig).Methods("POST")

	// Define the creation and modification endpoints for Plunder Deployment configuration
	router.HandleFunc("/deployment", postDeployment).Methods("POST")
	router.HandleFunc("/deployment/{id}", getSpecificDeployment).Methods("GET")
	router.HandleFunc("/deployment/{id}", updateDeployment).Methods("POST")
	router.HandleFunc("/deployment/{id}", deleteDeployment).Methods("DELETE")

	return router
}
