package apiserver

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/plunder-app/plunder/pkg/services"
)

// This package provides the capability to retrieve information about mac addresses (from DHCP) that plunder has seen or allocated

// Retrieve the plunder server dhcp configuration
func getDHCP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rsp Response

	// Find the deployment ID
	id := mux.Vars(r)["id"]

	if id == "leases" {

		jsonData, err := json.Marshal(services.Controller.GetLeases())
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			rsp.FriendlyError = "Error retrieving allocated leases"
			rsp.Error = err.Error()
		} else {
			rsp.Payload = jsonData
		}
	}
	// Are we updating the deployment "global"
	if id == "unleased" {

		jsonData, err := json.Marshal(services.Controller.GetUnLeased())

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			rsp.FriendlyError = "Error retrieving allocated leases"
			rsp.Error = err.Error()
		} else {
			rsp.Payload = jsonData
		}
	}

	json.NewEncoder(w).Encode(rsp)
}
