package apiserver

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	//"github.com/gorilla/mux"
)

// EndPointManager - Contains all of the dynamically created endpoints
var EndPointManager []EndPoint

// EndPoint is the source of truth for handling all of the endpoints exposed through the API Server
// it also provides a mechanism to interact with the apiserver to find/create api endpoints
type EndPoint struct {
	Name         string `json:"name"`
	Path         string `json:"path"`
	FunctionPath string `json:"functionEndpoint"`
	Description  string `json:"description"`
	Method       string `json:"method"`
}

// AddDynamicEndpoint - will add an endpoint to the api server and link it back to a function
func AddDynamicEndpoint(endpointPattern, path, description, name, method string, epFunc http.HandlerFunc) {
	for i := range EndPointManager {
		if EndPointManager[i].Name == name && EndPointManager[i].Method == method {
			log.Warnf("Endpoint [%s] already exists with method [%s]", name, method)
		}
	}
	// First we add the endpoint to the Manager so we can query it
	EndPointManager = append(EndPointManager, EndPoint{
		FunctionPath: endpointPattern,
		Path:         path,
		Description:  description,
		Method:       method,
		Name:         name,
	})
	// Then we add the endpoint to the apiServer
	endpoints.HandleFunc(endpointPattern, epFunc).Methods(method)
}

// GetEndpoint - will return the details for an endpoint
func GetEndpoint(name, method string) *EndPoint {
	for i := range EndPointManager {
		if EndPointManager[i].Name == name && EndPointManager[i].Method == method {
			return &EndPointManager[i]
		}
	}
	return nil
}

// FunctionPath - this will return the api server path for any external caller using the package
func FunctionPath() string {
	return "/api"
}
