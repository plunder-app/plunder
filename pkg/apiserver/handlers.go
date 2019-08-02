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

func getConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rsp response
	rsp.Response = services.Controller
	json.NewEncoder(w).Encode(rsp)
}

func getDeployment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rsp response
	rsp.Response = services.Deployments
	json.NewEncoder(w).Encode(rsp)
}

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

func postConfig(w http.ResponseWriter, r *http.Request) {
	if b, err := ioutil.ReadAll(r.Body); err == nil {
		err := services.ParseControllerData(b)
		var rsp response

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			rsp.FriendlyError = "Error updating Server Configuration"
			rsp.Error = err.Error()
			rsp.Response = nil
			json.NewEncoder(w).Encode(rsp)
		}
	}
}

func postDeployment(w http.ResponseWriter, r *http.Request) {

	if b, err := ioutil.ReadAll(r.Body); err == nil {
		err := services.UpdateDeploymentConfig(b)
		var rsp response

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			rsp.FriendlyError = "Error updating Deployment Configuration"
			rsp.Error = err.Error()
			rsp.Response = nil
			json.NewEncoder(w).Encode(rsp)
		}
	}
}
