package apiserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/plunder-app/plunder/pkg/parlay"
	"github.com/plunder-app/plunder/pkg/services"
)

// Retrieve a specific plunder deployment configuration
func postParlay(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rsp Response
	// Find the deployment ID
	id := mux.Vars(r)["id"]

	// We need to revert the mac address back to the correct format (dashes back to colons)
	mac := strings.Replace(id, "-", ":", -1)

	// Use the mac address to lookup the deployment
	deployment := services.GetDeployment(mac)

	// If the deployment exists then process the POST data
	if deployment != nil {

		if b, err := ioutil.ReadAll(r.Body); err == nil {
			// Parse the treasure map in the POST data
			var p parlay.TreasureMap
			err := json.Unmarshal(b, &p)
			// Unable to parse the JSON payload
			if err != nil {
				rsp.FriendlyError = "Error parsing the parlay actions"
				rsp.Error = err.Error()
			} else {
				// Parsed succesfully, we will deploy this in a go routine and use GET /parlay/MAC to view progress
				go p.DeploySSH("")
			}
		} else {
			rsp.FriendlyError = "Error reading HTTP data"
			rsp.Error = err.Error()

		}
	} else {
		rsp.FriendlyError = fmt.Sprintf("Unable to find deployment for server [%s]", mac)
	}

	json.NewEncoder(w).Encode(rsp)
}

// Retrieve a specific parlay automation
func getParlay(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rsp Response
	// Find the deployment ID
	id := mux.Vars(r)["id"]

	// We need to revert the mac address back to the correct format (dashes back to colons)
	mac := strings.Replace(id, "-", ":", -1)

	// Use the mac address to lookup the deployment
	deployment := services.GetDeployment(mac)

	// If the deployment exists then process the POST data
	if deployment != nil {

		// RETREIVE the deployment Logs (TODO)

	} else {
		rsp.FriendlyError = fmt.Sprintf("Unable to find deployment for server [%s]", mac)
	}

	json.NewEncoder(w).Encode(rsp)
}

// Delete the parlay results from the plunder server
func delParlay(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rsp Response
	// Find the deployment ID
	id := mux.Vars(r)["id"]

	// We need to revert the mac address back to the correct format (dashes back to colons)
	mac := strings.Replace(id, "-", ":", -1)

	// Use the mac address to lookup the deployment
	deployment := services.GetDeployment(mac)

	// If the deployment exists then process the POST data
	if deployment != nil {

		// DELETE the deployment logs (TODO)

	} else {
		rsp.FriendlyError = fmt.Sprintf("Unable to find deployment for server [%s]", mac)
	}

	json.NewEncoder(w).Encode(rsp)
}
