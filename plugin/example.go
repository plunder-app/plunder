package main

const info = `
This example plugin is used to demonstrate the structure of a plugin
`

// ParlayActionList - This should return a list of
func ParlayActionList() []string {
	return []string{"exampleAction/test", "exampleAction/demo", "exampleAction/example"}
}

// ParlayPluginInfo - returns information about the plugin
func ParlayPluginInfo() string {
	return info
}
