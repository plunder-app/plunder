package main

import (
	"encoding/json"

	"github.com/thebsdbox/plunder/pkg/parlay"
)

const info = `
This example plugin is used to demonstrate the structure of a plugin
`

// Action defines a custom action
type Action struct {
	Name       string     `json:"name"`
	Type       string     `json:"type"`
	TestAction testAction `json:"testDetails,omitempty"`
}

type testAction struct {
	Credentials string `json:"credentials"`
	Address     string `json:"address"`
}

// ParlayActionList - This should return a list of
func ParlayActionList() []string {
	return []string{"exampleAction/test", "exampleAction/demo", "exampleAction/example"}
}

// ParlayPluginInfo - returns information about the plugin
func ParlayPluginInfo() string {
	return info
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

// ParlayUsage -
func ParlayUsage(action string) string {
	var usageJSON string
	switch action {
	case "exampleAction/test":
		a := Action{
			Name: "Example of test action",
			Type: "exampleAction/test",
			TestAction: testAction{
				Credentials: "AAABBBCCCCDDEEEE",
				Address:     "172.0.0.1",
			},
		}
		b, _ := json.MarshalIndent(a, "", "\t")
		return string(b)
	default:
		return usageJSON
	}
}

// Dummy main function
func main() {}
