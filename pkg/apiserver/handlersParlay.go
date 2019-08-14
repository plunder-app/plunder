package apiserver

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/plunder-app/plunder/pkg/services"
)

// Retrieve a specific plunder deployment configuration
func postParlay(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rsp Response
	// Find the deployment ID
	id := mux.Vars(r)["id"]

	if b, err := ioutil.ReadAll(r.Body); err == nil {
		// We need to revert the mac address back to the correct format (dashes back to colons)
		mac := strings.Replace(id, "-", ":", -1)

		err := services.AddDeployment(b)

		if err != nil {
			rsp.FriendlyError = "Error updating Deployment Configuration"
			rsp.Error = err.Error()
			rsp.Payload = nil
		}
	} else {
		rsp.FriendlyError = "Error reading HTTP data"
		rsp.Error = err.Error()

	}

	json.NewEncoder(w).Encode(rsp)

}

// Retrieve a specific parlay automation
func getParlay(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rsp Response
	// Find the deployment ID
	id := mux.Vars(r)["id"]

	if b, err := ioutil.ReadAll(r.Body); err == nil {
		// We need to revert the mac address back to the correct format (dashes back to colons)
		mac := strings.Replace(id, "-", ":", -1)

		err := services.AddDeployment(b)

		if err != nil {
			rsp.FriendlyError = "Error updating Deployment Configuration"
			rsp.Error = err.Error()
			rsp.Payload = nil
		}
	} else {
		rsp.FriendlyError = "Error reading HTTP data"
		rsp.Error = err.Error()

	}

	json.NewEncoder(w).Encode(rsp)

}

// Delete the parlay results from the plunder server
func delParlay(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rsp Response
	// Find the deployment ID
	id := mux.Vars(r)["id"]

	if b, err := ioutil.ReadAll(r.Body); err == nil {
		// We need to revert the mac address back to the correct format (dashes back to colons)
		mac := strings.Replace(id, "-", ":", -1)

		err := services.AddDeployment(b)

		if err != nil {
			rsp.FriendlyError = "Error updating Deployment Configuration"
			rsp.Error = err.Error()
			rsp.Payload = nil
		}
	} else {
		rsp.FriendlyError = "Error reading HTTP data"
		rsp.Error = err.Error()

	}

	json.NewEncoder(w).Encode(rsp)

}
