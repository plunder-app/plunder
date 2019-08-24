package main

import (
	"encoding/json"
	"fmt"

	"github.com/plunder-app/plunder/pkg/parlay/parlaytypes"
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
func ParlayActions(action string, iface interface{}) []parlaytypes.Action {
	var actions []parlaytypes.Action
	a := parlaytypes.Action{
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

	// This example plugin only has the code for "exampleAction/test" however this switch statement
	// should handle all exposed actions from the plugin
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
func ParlayExec(action, host string, raw json.RawMessage) (actions []parlaytypes.Action, err error) {

	var t pluginTestAction
	// Unmarshall the JSON into the struct
	json.Unmarshal(raw, &t)
	// We can now use the fields as part of the struct

	// This example plugin only has the code for "exampleAction/test" however this switch statement
	// should handle all exposed actions from the plugin
	switch action {
	case "exampleAction/test":
		a := parlaytypes.Action{
			Name:       "Echo the address",
			ActionType: "command",
			Command:    fmt.Sprintf("echo %s", t.Address),
		}
		actions = append(actions, a)

		a.Name = "Echo the Credentials"
		a.Command = fmt.Sprintf("echo %s", t.Credentials)
		actions = append(actions, a)

		return
	default:
		return
	}
}
