package apiserver

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/plunder-app/plunder/pkg/services"
)

var response struct {
	Error    string      `json:"error,omitempty"`
	Response interface{} `json:"response,omitempty"`
}

func getConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//response.Error = nil
	response.Response = services.Controller
	json.NewEncoder(w).Encode(response)
}

func getDeployment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//response.Error = nil
	response.Response = services.Deployments
	json.NewEncoder(w).Encode(response)
}

func postConfig(w http.ResponseWriter, r *http.Request) {
	if b, err := ioutil.ReadAll(r.Body); err == nil {
		err := services.ParseControllerFile(b)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			response.Error = "Error updating Server Configuration"
			response.Response = nil
			json.NewEncoder(w).Encode(response)
		}
	}
}

func postDeployment(w http.ResponseWriter, r *http.Request) {

	if b, err := ioutil.ReadAll(r.Body); err == nil {
		err := services.UpdateDeploymentConfig(b)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			response.Error = "Error updating Deployment Configuration"
			response.Response = nil
			json.NewEncoder(w).Encode(response)
		}
	}
}
