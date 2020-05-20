package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/plunder-app/plunder/pkg/apiserver"
)

// RegisterToAPIServer - will add the endpoints to the API server
func RegisterToAPIServer() {

	// ------------------------------------------------
	//        Server configuration API registration
	// ------------------------------------------------

	apiserver.AddDynamicEndpoint("/config",
		"/config",
		"Allows the retrieving of Plunder Server configuration",
		"config",
		http.MethodGet,
		getConfig)

	apiserver.AddDynamicEndpoint("/config",
		"/config",
		"Allows the creation of Plunder Server configuration",
		"config",
		http.MethodPost,
		postConfig)

	apiserver.AddDynamicEndpoint("/config/boot/{id}",
		"/config/boot",
		"Allows the creation of Plunder Server Boot configuration",
		"configBoot",
		http.MethodPost,
		postBootConfig)

	apiserver.AddDynamicEndpoint("/config/boot/{id}",
		"/config/boot",
		"Performs the deletion of Plunder Server Boot configuration",
		"configBoot",
		http.MethodDelete,
		deleteBootConfig)

	// ------------------------------------------------
	//    DHCP configuration API registration
	// ------------------------------------------------

	apiserver.AddDynamicEndpoint("/dhcp/{id}",
		"/dhcp",
		"Allows the retrieval of DHCP information",
		"dhcp",
		http.MethodGet,
		getDHCP)

	// ------------------------------------------------
	//    Deployment configuration API registration
	// ------------------------------------------------

	apiserver.AddDynamicEndpoint("/deployments",
		"/deployments",
		"Allows the retrieving of Plunder Server deployments",
		"deployments",
		http.MethodGet,
		getDeployments)

	apiserver.AddDynamicEndpoint("/deployments",
		"/deployments",
		"Allows the creation of Plunder Server deployments",
		"deployments",
		http.MethodPost,
		postDeployments)

	apiserver.AddDynamicEndpoint("/deployment",
		"/deployment",
		"Allows the creation of a specific Plunder deployment",
		"deployment",
		http.MethodPost,
		postDeployment)

	apiserver.AddDynamicEndpoint("/deployment",
		"/deployment",
		"Allows the patching of a Plunder Server deployment",
		"deployment",
		http.MethodPatch,
		updateDeployment)

	apiserver.AddDynamicEndpoint("/deployment/{id}",
		"/deployment",
		"Allows the retrieval of specific information about a deployment",
		"deploymentID",
		http.MethodGet,
		getSpecificDeployment)

	apiserver.AddDynamicEndpoint("/deployment/{id}",
		"/deployment",
		"Allows the patching of Plunder Server deployments",
		"deploymentID",
		http.MethodPatch,
		updateDeployment)

	apiserver.AddDynamicEndpoint("/deployment/{id}",
		"/deployment",
		"Allows the deletion of a Plunder Server deployment",
		"deploymentID",
		http.MethodDelete,
		deleteDeployment)

	apiserver.AddDynamicEndpoint("/deployment/mac/{id}",
		"/deployment/mac",
		"Allows the deletion of a Plunder Server deployment based upon its MAC address",
		"deploymentMac",
		http.MethodDelete,
		deleteDeploymentMac)

	apiserver.AddDynamicEndpoint("/deployment/address/{id}",
		"/deployment/address",
		"Allows the deletion of a Plunder Server deployment based upon its network address",
		"deploymentAddress",
		http.MethodDelete,
		deleteDeploymentAddress)

}

func getConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rsp apiserver.Response
	jsonData, err := json.Marshal(Controller)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		rsp.Warning = "Error retrieving Server Configuration"
		rsp.Error = err.Error()
	} else {
		rsp.Payload = jsonData
	}
	json.NewEncoder(w).Encode(rsp)
}

// Apply the plunder server global configuration

func postConfig(w http.ResponseWriter, r *http.Request) {
	if b, err := ioutil.ReadAll(r.Body); err == nil {
		var rsp apiserver.Response
		// This function needs to parse both the data and then evaluate the state of running services
		err := ParseControllerData(b)

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			rsp.Warning = "Error updating Server Configuration"
			rsp.Error = err.Error()
		}
		json.NewEncoder(w).Encode(rsp)
		Controller.StartServices(nil)
	}
}

// Apply a Specific Boot Configuration
func postBootConfig(w http.ResponseWriter, r *http.Request) {
	if b, err := ioutil.ReadAll(r.Body); err == nil {
		var rsp apiserver.Response
		// This function needs to parse both the data and then evaluate the state of running services
		var newBoot BootConfig
		err := json.Unmarshal(b, &newBoot)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			rsp.Warning = "Error updating Server Configuration"
			rsp.Error = err.Error()

		} else {
			for x := range Controller.BootConfigs {
				if Controller.BootConfigs[x].ConfigName == newBoot.ConfigName {
					// Found a duplicate
					w.Header().Set("Content-Type", "application/json")
					rsp.Warning = "Error duplicate Server Configuration"
					rsp.Error = fmt.Sprintf("Boot Configuration [%s] already exists", Controller.BootConfigs[x].ConfigName)
					json.NewEncoder(w).Encode(rsp)
					return
				}
			}
			// // Parse the boot configuration (preload ISOs etc.)
			err = newBoot.Parse()
			// err = Controller.ParseBootController()
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				rsp.Warning = "Error updating Server Configuration"
				rsp.Error = err.Error()
			} else {
				// Add the Boot configuration to the controller
				Controller.BootConfigs = append(Controller.BootConfigs, newBoot)
				// Generate the handlers (this can probably GO soon)
				Controller.generateBootTypeHanders()
			}
			// // Parse the boot configuration (preload ISOs etc.)
			// err = Controller.ParseBootController()
			// if err != nil {
			// 	w.Header().Set("Content-Type", "application/json")
			// 	rsp.Warning = "Error updating Server Configuration"
			// 	rsp.Error = err.Error()
			// }
		}

		json.NewEncoder(w).Encode(rsp)
	}
}

// Apply a Specific Boot Configuration
func deleteBootConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Find the deployment ID
	id := mux.Vars(r)["id"]
	var rsp apiserver.Response

	// We need to revert the mac address back to the correct format (dashes back to colons)
	err := Controller.DeleteBootControllerConfig(id)
	if err != nil {

		if err != nil {
			rsp.Warning = "Error updating Deployment Configuration"
			rsp.Error = err.Error()
			rsp.Payload = nil
		}
	}

	json.NewEncoder(w).Encode(rsp)
}

func getDeployments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rsp apiserver.Response
	jsonData, err := json.Marshal(Deployments)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		rsp.Warning = "Error retrieving deployment Configuration"
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
		err := UpdateDeploymentConfig(b)
		var rsp apiserver.Response

		if err != nil {
			rsp.Warning = "Error updating Deployment Configuration"
			rsp.Error = err.Error()
			rsp.Payload = nil
		}
		json.NewEncoder(w).Encode(rsp)
	}
}

// Retrieve a specific plunder deployment configuration

func getSpecificDeployment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rsp apiserver.Response
	// Find the deployment ID
	id := mux.Vars(r)["id"]
	// We need to revert the mac address back to the correct format (dashes back to colons)
	mac := strings.Replace(id, "-", ":", -1)

	deployment := GetDeployment(mac)

	if deployment != nil {
		jsonData, err := json.Marshal(deployment)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			rsp.Warning = "Error retrieving deployment Configuration"
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
		err := AddDeployment(b)
		var rsp apiserver.Response

		if err != nil {
			rsp.Warning = "Error updating Deployment Configuration"
			rsp.Error = err.Error()
			rsp.Payload = nil
		}
		json.NewEncoder(w).Encode(rsp)
	}

}

// Retrieve a specific plunder deployment configuration
func updateDeployment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rsp apiserver.Response

	// Find the deployment ID
	id := mux.Vars(r)["id"]

	// Are we updating the deployment "global"
	if id == "global" {
		if b, err := ioutil.ReadAll(r.Body); err == nil {
			err := UpdateGlobalDeploymentConfig(b)

			if err != nil {
				rsp.Warning = "Error updating Global Configuration"
				rsp.Error = err.Error()
				rsp.Payload = nil
			}
		}
	} else {
		// We need to revert the mac address back to the correct format (dashes back to colons)
		mac := strings.Replace(id, "-", ":", -1)

		if b, err := ioutil.ReadAll(r.Body); err == nil {
			err := UpdateDeployment(mac, b)

			if err != nil {
				rsp.Warning = "Error updating Deployment Configuration"
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
	var rsp apiserver.Response

	if b, err := ioutil.ReadAll(r.Body); err == nil {
		// Try the Mac address first

		// We need to revert the mac address back to the correct format (dashes back to colons)
		err := DeleteDeploymentMac(strings.Replace(id, "-", ":", -1), b)
		if err != nil {

			// We need to revert the ip address back to the correct format (dashes back to periods)
			err = DeleteDeploymentAddress(strings.Replace(id, "-", ".", -1), b)
			if err != nil {
				rsp.Warning = "Error updating Deployment Configuration"
				rsp.Error = err.Error()
				rsp.Payload = nil
			}
		}
	}
	json.NewEncoder(w).Encode(rsp)

}

// Retrieve a specific plunder deployment configuration
func deleteDeploymentMac(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Find the deployment ID
	id := mux.Vars(r)["id"]
	var rsp apiserver.Response

	if b, err := ioutil.ReadAll(r.Body); err == nil {
		// We need to revert the mac address back to the correct format (dashes back to colons)
		err := DeleteDeploymentMac(strings.Replace(id, "-", ":", -1), b)
		if err != nil {
			rsp.Warning = "Error updating Deployment Configuration"
			rsp.Error = err.Error()
			rsp.Payload = nil
		}

	}
	json.NewEncoder(w).Encode(rsp)

}

// Retrieve a specific plunder deployment configuration
func deleteDeploymentAddress(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Find the deployment ID
	id := mux.Vars(r)["id"]
	var rsp apiserver.Response

	if b, err := ioutil.ReadAll(r.Body); err == nil {
		// We need to revert the mac address back to the correct format (dashes back to colons)
		err = DeleteDeploymentAddress(strings.Replace(id, "-", ".", -1), b)
		if err != nil {
			rsp.Warning = "Error updating Deployment Configuration"
			rsp.Error = err.Error()
			rsp.Payload = nil
		}

	}
	json.NewEncoder(w).Encode(rsp)

}

func getDHCP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rsp apiserver.Response

	// Find the deployment ID
	id := mux.Vars(r)["id"]

	if id == "leases" {

		jsonData, err := json.Marshal(Controller.GetLeases())
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			rsp.Warning = "Error retrieving allocated leases"
			rsp.Error = err.Error()
		} else {
			rsp.Payload = jsonData
		}
	}
	// Are we updating the deployment "global"
	if id == "unleased" {

		jsonData, err := json.Marshal(Controller.GetUnLeased())

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			rsp.Warning = "Error retrieving allocated leases"
			rsp.Error = err.Error()
		} else {
			rsp.Payload = jsonData
		}
	}

	json.NewEncoder(w).Encode(rsp)
}
