package apiserver

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/plunder-app/plunder/pkg/services"
)

// Retrieve the plunder server global configuration

func getConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rsp Response
	jsonData, err := json.Marshal(services.Controller)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		rsp.FriendlyError = "Error retrieving Server Configuration"
		rsp.Error = err.Error()
	} else {
		rsp.Payload = jsonData
	}
	json.NewEncoder(w).Encode(rsp)
}

// Apply the plunder server global configuration

func postConfig(w http.ResponseWriter, r *http.Request) {
	if b, err := ioutil.ReadAll(r.Body); err == nil {
		var rsp Response
		// This function needs to parse both the data and then evaluate the state of running services
		err := services.ParseControllerData(b)

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			rsp.FriendlyError = "Error updating Server Configuration"
			rsp.Error = err.Error()
		}
		json.NewEncoder(w).Encode(rsp)
		services.Controller.StartServices(nil)
	}
}

// Apply a Specific Boot Configuration
func postBootConfig(w http.ResponseWriter, r *http.Request) {
	if b, err := ioutil.ReadAll(r.Body); err == nil {
		var rsp Response
		// This function needs to parse both the data and then evaluate the state of running services
		var newBoot services.BootConfig
		err := json.Unmarshal(b, &newBoot)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			rsp.FriendlyError = "Error updating Server Configuration"
			rsp.Error = err.Error()
		} else {
			// Add the Boot configuration to the controller
			services.Controller.BootConfigs = append(services.Controller.BootConfigs, newBoot)
			// Parse the boot configuration (preload ISOs etc.)
			services.Controller.ParseBootController()
		}

		json.NewEncoder(w).Encode(rsp)
	}
}

// Apply a Specific Boot Configuration
func deleteBootConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Find the deployment ID
	id := mux.Vars(r)["id"]
	var rsp Response

	// We need to revert the mac address back to the correct format (dashes back to colons)
	err := services.Controller.DeleteBootControllerConfig(id)
	if err != nil {

		if err != nil {
			rsp.FriendlyError = "Error updating Deployment Configuration"
			rsp.Error = err.Error()
			rsp.Payload = nil
		}
	}

	json.NewEncoder(w).Encode(rsp)
}
