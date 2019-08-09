package apiserver

import (
	"fmt"

	"github.com/gorilla/mux"
)

// Expose Endpoints to the outside world

//ConfigAPIPath returns the URI that is used to interact with the plunder Configuration API
func ConfigAPIPath() string {
	return "/config"
}

//DeploymentsAPIPath returns the URI that is used to interact with all Plunder deployments
func DeploymentsAPIPath() string {
	return "/deployments"
}

//DeploymentAPIPath returns the URI that is used to interact with the specific deployments
func DeploymentAPIPath() string {
	return "/deployment"
}

//DHCPAPIPath returns the URI that is used to interact with the plunder Configuration API
func DHCPAPIPath() string {
	return "/dhcp"
}

// setAPIEndpoints defines all of the API end points for Plunder
func setAPIEndpoints() *mux.Router {
	// Create a new router
	router := mux.NewRouter()

	// ------------------------------------
	// General configuration management
	// ------------------------------------

	// Define the retrieval endpoints for Plunder Server configuration
	router.HandleFunc(fmt.Sprintf("%s", ConfigAPIPath()), getConfig).Methods("GET")

	// Define the creation endpoints for Plunder Server Configuration
	router.HandleFunc(fmt.Sprintf("%s", ConfigAPIPath()), postConfig).Methods("POST")

	// Define the retrieval endpoints for Plunder Deployment configuration
	router.HandleFunc(fmt.Sprintf("%s", DeploymentsAPIPath()), getDeployments).Methods("GET")

	// Define the retrieval endpoints for Plunder Deployment configuration
	router.HandleFunc(fmt.Sprintf("%s", DeploymentsAPIPath()), postDeployments).Methods("POST")

	// Define the retrieval endpoints for Plunder Server configuration
	router.HandleFunc(fmt.Sprintf("%s/{id}", DHCPAPIPath()), getDHCP).Methods("GET")

	// ------------------------------------
	// Specific configuration management
	// ------------------------------------

	// Define the creation endpoints for Plunder Server Boot Configuration
	router.HandleFunc(fmt.Sprintf("%s/{id}", ConfigAPIPath()), postBootConfig).Methods("POST")

	// Define the creation and modification endpoints for Plunder Deployment configuration
	router.HandleFunc(fmt.Sprintf("%s", DeploymentAPIPath()), postDeployment).Methods("POST")
	router.HandleFunc(fmt.Sprintf("%s/{id}", DeploymentAPIPath()), getSpecificDeployment).Methods("GET")
	router.HandleFunc(fmt.Sprintf("%s/{id}", DeploymentAPIPath()), updateDeployment).Methods("POST")
	router.HandleFunc(fmt.Sprintf("%s/{id}", DeploymentAPIPath()), deleteDeployment).Methods("DELETE")

	return router
}
