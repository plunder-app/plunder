package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Delete the parlay results from the plunder server
func getAPIFunctionMethod(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rsp Response
	// Find the deployment ID
	f := mux.Vars(r)["function"]
	m := mux.Vars(r)["method"]

	ep := GetEndpoint(f, m)
	if ep == nil {
		// RETREIVE the deployment Logs (TODO)
		rsp.Warning = fmt.Sprintf("Unable to find HTTP method [%s] for function [%s]", m, f)
		rsp.Error = "Error looking up in API Server"
	} else {
		jsonData, err := json.Marshal(ep)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			rsp.Warning = "Error retrieving deployment Configuration"
			rsp.Error = err.Error()
		} else {
			rsp.Payload = jsonData
		}
	}

	json.NewEncoder(w).Encode(rsp)
}

// Delete the parlay results from the plunder server
func getAPIs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rsp Response

	jsonData, err := json.Marshal(EndPointManager)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		rsp.Warning = "Error retrieving deployment Configuration"
		rsp.Error = err.Error()
	} else {
		rsp.Payload = jsonData
	}

	json.NewEncoder(w).Encode(rsp)
}
