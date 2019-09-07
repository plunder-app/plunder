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

//ParlayAPIPath returns the URI that is used to interact with the plunder parlay automation engine
func ParlayAPIPath() string {
	return "/parlay"
}

//LogsHTTPAPIPath returns the URI that is used to interact with the streaming of http logs
func LogsHTTPAPIPath() string {
	return "/logs/http"
}

// setStaticAPIEndpoints defines all of the static API end points for Plunder
func setStaticAPIEndpoints(r *mux.Router) {
	// Create a new router

	// ------------------------------------
	// General configuration management
	// ------------------------------------

	// Define the retrieval endpoints for Plunder Server configuration
	r.HandleFunc(fmt.Sprintf("%s", ConfigAPIPath()), getConfig).Methods("GET")

	// Define the creation endpoints for Plunder Server Configuration
	r.HandleFunc(fmt.Sprintf("%s", ConfigAPIPath()), postConfig).Methods("POST")

	// Define the retrieval endpoints for Plunder Deployment configuration
	r.HandleFunc(fmt.Sprintf("%s", DeploymentsAPIPath()), getDeployments).Methods("GET")

	// Define the retrieval endpoints for Plunder Deployment configuration
	r.HandleFunc(fmt.Sprintf("%s", DeploymentsAPIPath()), postDeployments).Methods("POST")

	// Define the retrieval endpoints for Plunder Server configuration
	r.HandleFunc(fmt.Sprintf("%s/{id}", DHCPAPIPath()), getDHCP).Methods("GET")

	// Define the endpoint for sending commands to a remote host using the parlay engine
	r.HandleFunc(fmt.Sprintf("%s", ParlayAPIPath()), postParlay).Methods("POST")

	// ------------------------------------
	// Specific configuration management
	// ------------------------------------

	// Define the creation endpoints for Plunder Server Boot Configuration
	r.HandleFunc(fmt.Sprintf("%s/{id}", ConfigAPIPath()), postBootConfig).Methods("POST")
	r.HandleFunc(fmt.Sprintf("%s/{id}", ConfigAPIPath()), deleteBootConfig).Methods("DELETE")

	// Define the creation and modification endpoints for Plunder Deployment configuration
	r.HandleFunc(DeploymentAPIPath(), postDeployment).Methods("POST")
	r.HandleFunc(fmt.Sprintf("%s/{id}", DeploymentAPIPath()), getSpecificDeployment).Methods("GET")
	r.HandleFunc(fmt.Sprintf("%s/{id}", DeploymentAPIPath()), updateDeployment).Methods("POST")
	r.HandleFunc(fmt.Sprintf("%s/{id}", DeploymentAPIPath()), deleteDeployment).Methods("DELETE")

	// Delete deployments based upon different criteria
	r.HandleFunc(fmt.Sprintf("%s/mac/{id}", DeploymentAPIPath()), deleteDeploymentMac).Methods("DELETE")
	r.HandleFunc(fmt.Sprintf("%s/address/{id}", DeploymentAPIPath()), deleteDeploymentAddress).Methods("DELETE")

	// Define the endpoint for sending commands to a remote host using the parlay engine
	r.HandleFunc(fmt.Sprintf("%s/logs/{id}", ParlayAPIPath()), getParlay).Methods("GET")
	r.HandleFunc(fmt.Sprintf("%s/logs/{id}", ParlayAPIPath()), delParlay).Methods("DELETE")

	// ------------------------------------
	// Logging endpoints
	// ------------------------------------

	// Define the endpoint for sending commands to a remote host using the parlay engine
	//router.HandleFunc(fmt.Sprintf("%s/{id}", LogsHTTPAPIPath()), handleSSE(loggingCenter)).Methods("GET")
	return
}
