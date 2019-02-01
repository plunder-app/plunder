package main

import (
	"encoding/json"
	"fmt"

	"github.com/thebsdbox/plunder/pkg/parlay/types"
)

const pluginInfo = `This plugin is used to managed kubeadm automation`

// This defines the etcd kubeadm file (should use the kubernetes packages to define at a later point)
const etcdKubeadm = `apiVersion: "kubeadm.k8s.io/%s"
kind: ClusterConfiguration
etcd:
    local:
        serverCertSANs:
        - "%s"
        peerCertSANs:
        - "%s"
        extraArgs:
            initial-cluster: %s=https://%s:2380,%s=https://%s:2380,%s=https://%s:2380
            initial-cluster-state: new
            name: %s
            listen-peer-urls: https://%s:2380
            listen-client-urls: https://%s:2379
            advertise-client-urls: https://%s:2379
            initial-advertise-peer-urls: https://%s:2380`

// This defines the manager kubeadm file (should use the kubernetes packages to define at a later point)

const managerKubeadm = `apiVersion: kubeadm.k8s.io/v1beta1
kind: ClusterConfiguration
kubernetesVersion: %s
apiServer:
  certSANs:
  - "%s"
controlPlaneEndpoint: "%s:%d"
etcd:
    external:
        endpoints:
        - https://%s:2379
        - https://%s:2379
        - https://%s:2379
        caFile: /etc/kubernetes/pki/etcd/ca.crt
        certFile: /etc/kubernetes/pki/apiserver-etcd-client.crt
        keyFile: /etc/kubernetes/pki/apiserver-etcd-client.key`

type etcdMembers struct {
	// Hostnames
	Hostname1 string `json:"hostname1,omitempty"`
	Hostname2 string `json:"hostname2,omitempty"`
	Hostname3 string `json:"hostname3,omitempty"`

	// Addresses
	Address1 string `json:"address1,omitempty"`
	Address2 string `json:"address2,omitempty"`
	Address3 string `json:"address3,omitempty"`

	// Intialise a Certificate Authority
	InitCA bool `json:"initCA,omitempty"`

	// Set kubernetes API version
	APIVersion string `json:"apiversion,omitempty"`
}

type managerMembers struct {
	// ETCD Nodes
	ETCDAddress1 string `json:"etcd01,omitempty"`
	ETCDAddress2 string `json:"etcd02,omitempty"`
	ETCDAddress3 string `json:"etcd03,omitempty"`

	// Version of Kubernetes
	Version string `json:"kubeVersion,omitempty"`

	// Load Balancer details (needed for initialising the first master)
	//loadBalancer

	// Stacked - means ETCD nodes are stacked on managers (false by default)
	Stacked bool `json:"stacked,omitempty"`
}

// Dummy main function
func main() {}

// ParlayActionList - This should return an array of actions
func ParlayActionList() []string {
	return []string{
		"kubeadm/etcd",
		"kubeadm/master"}
}

// ParlayActionDetails - This should return an array of action descriptions
func ParlayActionDetails() []string {
	return []string{
		"This action automates the provisioning of a the first etcd node and certificates for the remaining two nodes",
		"This action handles the configuration of the first master node"}
}

// ParlayPluginInfo - returns information about the plugin
func ParlayPluginInfo() string {
	return pluginInfo
}

// ParlayUsage - Returns the json that matches the specific action
// <- action is a string that defines which action the usage information should be
// <- raw - raw JSON that will be manipulated into a correct struct that matches the action
// -> err is any error that has been generated
func ParlayUsage(action string) (raw json.RawMessage, err error) {

	// This example plugin only has the code for "exampleAction/test" however this switch statement
	// should handle all exposed actions from the plugin
	switch action {
	case "kubeadm/etcd":
		a := etcdMembers{
			Hostname1:  "etcd01.local",
			Hostname2:  "etcd02.local",
			Hostname3:  "etcd03.local",
			InitCA:     true,
			APIVersion: "v1beta1",
			Address1:   "10.0.101",
			Address2:   "10.0.102",
			Address3:   "10.0.103",
		}
		// In order to turn a struct into an map[string]interface we need to turn it into JSON

		return json.Marshal(a)
	default:
		return raw, fmt.Errorf("Action [%s] could not be found", action)
	}
}

// ParlayExec - Parses the action and the data that the action will consume
// <- action a string that details the action to be executed
// <- raw - raw JSON that will be manipulated into a correct struct that matches the action
// -> actions are an array of generated actions that the parser will then execute
// -> err is any error that has been generated
func ParlayExec(action string, raw json.RawMessage) (actions []types.Action, err error) {

	// This example plugin only has the code for "exampleAction/test" however this switch statement
	// should handle all exposed actions from the plugin
	switch action {
	case "kubeadm/etcd":
		var etcdStruct etcdMembers
		// Unmarshall the JSON into the struct
		err = json.Unmarshal(raw, &etcdStruct)
		return etcdStruct.generateActions(), err
	default:
		return
	}
}
