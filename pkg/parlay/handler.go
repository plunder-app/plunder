package parlay

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/plunder-app/plunder/pkg/apiserver"
	"github.com/plunder-app/plunder/pkg/parlay/parlaytypes"
	"github.com/plunder-app/plunder/pkg/services"
	"github.com/plunder-app/plunder/pkg/ssh"
	log "github.com/sirupsen/logrus"
)

var registered bool

// RegisterToAPIServer - will add the endpoints to the API server
func RegisterToAPIServer() {
	// Ensure registration only happens once
	if registered == true {
		return
	}

	// ------------------------------------------------
	//        Parlay API registration
	// ------------------------------------------------

	apiserver.AddDynamicEndpoint("/parlay",
		"/parlay",
		"Create a parlay automation deployment",
		"parlay",
		http.MethodPost,
		postParlay)

	apiserver.AddDynamicEndpoint("/parlay/logs/{id}",
		"/parlay/logs",
		"Retrieve the logs from a parlay deployment",
		"parlayLog",
		http.MethodGet,
		getParlay)

	apiserver.AddDynamicEndpoint("/parlay/logs/{id}",
		"/parlay/logs",
		"Delete the cached logs from a specific parlay deployment",
		"parlayLog",
		http.MethodDelete,
		delParlay)
	registered = true
}

// Retrieve a specific plunder deployment configuration
func postParlay(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rsp apiserver.Response

	if b, err := ioutil.ReadAll(r.Body); err == nil {
		// Parse the treasure map in the POST data
		var m parlaytypes.TreasureMap
		err := json.Unmarshal(b, &m)
		// Unable to parse the JSON payload
		if err != nil {
			rsp.Warning = "Error parsing the parlay actions"
			rsp.Error = err.Error()
		} else {
			// Parsed succesfully, we will deploy this in a go routine and use GET /parlay/MAC to view progress
			//
			err = ssh.ImportHostsFromDeployment(services.Deployments)
			if err != nil {
				rsp.Warning = "Error importing the hosts from deployment"
				rsp.Error = err.Error()
			} else {
				err = DeploySSH(&m, "", true, true)
				if err != nil {
					rsp.Warning = "Error performing the parlay actions"
					rsp.Error = err.Error()
					log.Errorf("%s", err.Error())
				}

			}
		}
	} else {
		rsp.Warning = "Error reading HTTP data"
		rsp.Error = err.Error()

	}

	json.NewEncoder(w).Encode(rsp)
}

// Retrieve a specific parlay automation
func getParlay(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rsp apiserver.Response
	// Find the deployment ID
	id := mux.Vars(r)["id"]

	// We need to revert the mac address back to the correct format (dashes back to colons)
	target := strings.Replace(id, "-", ".", -1)

	// Use the mac address to lookup the deployment
	logs, err := GetTargetLogs(target)
	// If the deployment exists then process the POST data
	if err != nil {
		// RETREIVE the deployment Logs (TODO)
		rsp.Warning = "Error reading Parlay Logs"
		rsp.Error = err.Error()
	} else {
		jsonData, err := json.Marshal(logs)
		if err != nil {

			// RETREIVE the deployment Logs (TODO)
			rsp.Warning = "Error parsing Parlay Logs"
			rsp.Error = err.Error()
		} else {
			rsp.Payload = jsonData
		}
	}

	json.NewEncoder(w).Encode(rsp)
}

// Delete the parlay results from the plunder server
func delParlay(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rsp apiserver.Response
	// Find the deployment ID
	id := mux.Vars(r)["id"]

	// We need to revert the mac address back to the correct format (dashes back to colons)
	target := strings.Replace(id, "-", ".", -1)

	// Use the mac address to lookup the deployment
	err := DeleteTargetLogs(target)
	// If the deployment exists then process the POST data
	if err != nil {

		// RETREIVE the deployment Logs (TODO)
		rsp.Warning = "Error reading deleting logs"
		rsp.Error = err.Error()
	}

	json.NewEncoder(w).Encode(rsp)
}
