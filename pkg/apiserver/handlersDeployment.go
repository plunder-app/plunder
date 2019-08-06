package apiserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/plunder-app/plunder/pkg/services"
)

type response struct {
	FriendlyError string `json:"friendlyError,omitempty"`
	Error         string `json:"error,omitempty"`

	Response interface{} `json:"response,omitempty"`
}

// Retrieve the plunder global deployment configuration

func getDeployments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rsp response
	rsp.Response = services.Deployments
	json.NewEncoder(w).Encode(rsp)
}

// Apply the plunder global deployment configuration

func postDeployments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if b, err := ioutil.ReadAll(r.Body); err == nil {
		err := services.UpdateDeploymentConfig(b)
		var rsp response

		if err != nil {
			rsp.FriendlyError = "Error updating Deployment Configuration"
			rsp.Error = err.Error()
			rsp.Response = nil
			json.NewEncoder(w).Encode(rsp)
		}
	}
}

// Retrieve a specific plunder deployment configuration

func getSpecificDeployment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var rsp response
	for i := range services.Deployments.Configs {

		if params["id"] == strings.Replace(services.Deployments.Configs[i].MAC, ":", "-", -1) {
			rsp.Response = services.Deployments.Configs[i]
		}
	}
	if rsp.Response == nil {
		rsp.Error = fmt.Sprintf("Unable to find %s", params["id"])
	}
	json.NewEncoder(w).Encode(rsp)

}

// Retrieve a specific plunder deployment configuration
func postDeployment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if b, err := ioutil.ReadAll(r.Body); err == nil {
		err := services.AddDeployment(b)
		var rsp response

		if err != nil {
			rsp.FriendlyError = "Error updating Deployment Configuration"
			rsp.Error = err.Error()
			rsp.Response = nil
			json.NewEncoder(w).Encode(rsp)
		}
	}

}

// Retrieve a specific plunder deployment configuration
func updateDeployment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rsp response

	// Find the deployment ID
	id := mux.Vars(r)["id"]

	// Are we updating the deployment "global"
	if id == "global" {
		if b, err := ioutil.ReadAll(r.Body); err == nil {
			err := services.UpdateGlobalDeploymentConfig(b)
	
			if err != nil {
				rsp.FriendlyError = "Error updating Global Configuration"
				rsp.Error = err.Error()
				rsp.Response = nil
				json.NewEncoder(w).Encode(rsp)
			}
		}
	} else {
	// We need to revert the mac address back to the correct format (dashes back to colons)
	mac := strings.Replace(id, "-", ":", -1)

	if b, err := ioutil.ReadAll(r.Body); err == nil {
		err := services.UpdateDeployment(mac, b)

		if err != nil {
			rsp.FriendlyError = "Error updating Deployment Configuration"
			rsp.Error = err.Error()
			rsp.Response = nil
			json.NewEncoder(w).Encode(rsp)
		}
	}
}
}

// Retrieve a specific plunder deployment configuration
func deleteDeployment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Find the deployment ID
	id := mux.Vars(r)["id"]
	// We need to revert the mac address back to the correct format (dashes back to colons)
	mac := strings.Replace(id, "-", ":", -1)

	if b, err := ioutil.ReadAll(r.Body); err == nil {
		err := services.DeleteDeployment(mac, b)
		var rsp response

		if err != nil {
			rsp.FriendlyError = "Error updating Deployment Configuration"
			rsp.Error = err.Error()
			rsp.Response = nil
			json.NewEncoder(w).Encode(rsp)
		}
	}
}
