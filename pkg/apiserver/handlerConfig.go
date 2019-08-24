package apiserver

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

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
			services.Controller.BootConfigs = append(services.Controller.BootConfigs, newBoot)
		}

		json.NewEncoder(w).Encode(rsp)
	}
}
