package main

import (
	"encoding/json"
	"fmt"

	"github.com/thebsdbox/plunder/pkg/parlay"
)

const pluginInfo = `This example plugin is used to demonstrate the structure of a plugin`

// pluginAction - defines the struct that is unique to the
type pluginTestAction struct {
	Credentials string `json:"credentials"`
	Address     string `json:"address"`
}

// Dummy main function
func main() {}

// ParlayActionList - This should return an array of actions
func ParlayActionList() []string {
	return []string{
		"exampleAction/test",
		"exampleAction/demo",
		"exampleAction/example"}
}

// ParlayActionDetails - This should return an array of action descriptions
func ParlayActionDetails() []string {
	return []string{
		"This action handles the testing part of the example plugin",
		"This action handles the demonstration of the example plugin",
		"This action handles an example of the example plugin!"}
}

// ParlayPluginInfo - returns information about the plugin
func ParlayPluginInfo() string {
	return pluginInfo
}

//ParlayActions -
func ParlayActions(action string, iface interface{}) []parlay.Action {
	var actions []parlay.Action
	a := parlay.Action{
		Command: "example/test",
	}
	actions = append(actions, a)
	return actions
}

// ParlayUsage - Returns the json that matches the specific action
// <- action is a string that defines which action the usage information should be
// <- raw - raw JSON that will be manipulated into a correct struct that matches the action
// -> err is any error that has been generated
func ParlayUsage(action string) (raw json.RawMessage, err error) {
	switch action {
	case "exampleAction/test":
		a := pluginTestAction{
			Credentials: "AAABBBCCCCDDEEEE",
			Address:     "172.0.0.1",
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
func ParlayExec(action string, raw json.RawMessage) (actions []parlay.Action, err error) {

	var t pluginTestAction
	// Unmarshall the JSON into the struct
	json.Unmarshal(raw, &t)
	// We can now use the fields as part of the struct
	fmt.Printf("Address = %s\nCredentials = %s\n", t.Address, t.Credentials)

	return nil, nil
}
