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

//Response - This is the wrapper for responses back to a client, if any errors are created then the payload isn't guarenteed
type Response struct {
	FriendlyError string `json:"friendlyError,omitempty"`
	Error         string `json:"error,omitempty"`   // when it goes wrong
	Success       string `json:"success,omitempty"` // when it goes correct

	Payload json.RawMessage `json:"payload,omitempty"`
}

// Retrieve the plunder global deployment configuration

func getDeployments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rsp Response
	jsonData, err := json.Marshal(services.Deployments)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		rsp.FriendlyError = "Error retrieving deployment Configuration"
		rsp.Error = err.Error()
	} else {
		rsp.Payload = jsonData
	}
	json.NewEncoder(w).Encode(rsp)
}

// Apply the plunder global deployment configuration

func postDeployments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if b, err := ioutil.ReadAll(r.Body); err == nil {
		err := services.UpdateDeploymentConfig(b)
		var rsp Response

		if err != nil {
			rsp.FriendlyError = "Error updating Deployment Configuration"
			rsp.Error = err.Error()
			rsp.Payload = nil
		}
		json.NewEncoder(w).Encode(rsp)
	}
}

// Retrieve a specific plunder deployment configuration

func getSpecificDeployment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rsp Response
	// Find the deployment ID
	id := mux.Vars(r)["id"]
	// We need to revert the mac address back to the correct format (dashes back to colons)
	mac := strings.Replace(id, "-", ":", -1)

	deployment := services.GetDeployment(mac)

	if deployment != nil {
		jsonData, err := json.Marshal(deployment)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			rsp.FriendlyError = "Error retrieving deployment Configuration"
			rsp.Error = err.Error()
		} else {
			rsp.Payload = jsonData
		}

	} else {
		rsp.Error = fmt.Sprintf("Unable to find %s", mac)
	}

	json.NewEncoder(w).Encode(rsp)

}

// Retrieve a specific plunder deployment configuration
func postDeployment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if b, err := ioutil.ReadAll(r.Body); err == nil {
		err := services.AddDeployment(b)
		var rsp Response

		if err != nil {
			rsp.FriendlyError = "Error updating Deployment Configuration"
			rsp.Error = err.Error()
			rsp.Payload = nil
		}
		json.NewEncoder(w).Encode(rsp)
	}

}

// Retrieve a specific plunder deployment configuration
func updateDeployment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rsp Response

	// Find the deployment ID
	id := mux.Vars(r)["id"]

	// Are we updating the deployment "global"
	if id == "global" {
		if b, err := ioutil.ReadAll(r.Body); err == nil {
			err := services.UpdateGlobalDeploymentConfig(b)

			if err != nil {
				rsp.FriendlyError = "Error updating Global Configuration"
				rsp.Error = err.Error()
				rsp.Payload = nil
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
				rsp.Payload = nil
			}
		}
	}
	json.NewEncoder(w).Encode(rsp)
}

// Retrieve a specific plunder deployment configuration
func deleteDeployment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Find the deployment ID
	id := mux.Vars(r)["id"]
	// We need to revert the mac address back to the correct format (dashes back to colons)
	mac := strings.Replace(id, "-", ":", -1)
	var rsp Response

	if b, err := ioutil.ReadAll(r.Body); err == nil {
		err := services.DeleteDeployment(mac, b)

		if err != nil {
			rsp.FriendlyError = "Error updating Deployment Configuration"
			rsp.Error = err.Error()
			rsp.Payload = nil
		}
	}
	json.NewEncoder(w).Encode(rsp)

}
