package main

import (
	"encoding/json"
	"fmt"

	"github.com/thebsdbox/plunder/pkg/parlay"
)

const pluginInfo = `This example plugin is used to demonstrate the structure of a plugin`

//Test -
type Test struct {
	test  string
	test1 int
	test2 bool
}

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

// ParlayUsage -
func ParlayUsage(action string) (string, error) {
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
		return string(b), nil
	default:
		return usageJSON, fmt.Errorf("Action [%s] could not be found", action)
	}
}

// ParlayExec -
func ParlayExec(action string, iface interface{}) ([]parlay.Action, error) {
	//fmt.Printf("%v", iface.(Test).test1)
	// bob := iface.(*Test)
	// fmt.Printf("%v\n%v\n", bob, iface)
	fmt.Printf("%s", iface["test1"])
	return nil, nil
}
