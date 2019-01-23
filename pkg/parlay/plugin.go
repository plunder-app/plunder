package parlay

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"

	log "github.com/Sirupsen/logrus"
)

// Find plugins returns an array of all .plugin files
func findPlugins(pluginDir string) []string {
	var plugins []string
	// This function will look for all files in a specified directory (defaults to PWD/plugin)
	filepath.Walk(pluginDir, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			if filepath.Ext(path) == ".plugin" {
				absPath, _ := filepath.Abs(path)

				plugins = append(plugins, absPath)
			}
		}
		return nil
	})
	return plugins
}

func findFunctionInPlugin(pluginPath, functionName string) (plugin.Symbol, error) {

	plug, err := plugin.Open(pluginPath)
	if err != nil {
		log.Debugf("%v", err)
		return nil, fmt.Errorf("Unable to open Plugin [%s]", pluginPath)

	}

	symbol, err := plug.Lookup(functionName)
	if err != nil {
		log.Debugf("%v", err)
		return nil, fmt.Errorf("Unable to read functions from Plugin [%s]", pluginPath)

	}

	return symbol, nil
}

//LoadPlugins -
func LoadPlugins() {

	pluginList := findPlugins("./plugin")
	log.Debugf("Found [%d] plugins", len(pluginList))
	for x := range pluginList {
		symbol, err := findFunctionInPlugin(pluginList[x], "ParlayPluginList")
		if err != nil {
			log.Errorf("Unable to read functions from Plugin [%s]", pluginList[x])
			continue
		}

		p, ok := symbol.(func() []string)
		if !ok {
			log.Errorf("Unable to read functions from Plugin [%s]", pluginList[x])
			continue
		}

		z := p()
		for y := range z {
			fmt.Println(z[y])
		}
	}
}

//ListPlugins -
func ListPlugins() {

	pluginList := findPlugins("./plugin")
	log.Debugf("Found [%d] plugins", len(pluginList))
	for x := range pluginList {
		symbol, err := findFunctionInPlugin(pluginList[x], "ParlayPluginInfo")
		if err != nil {
			log.Errorf("Unable to read functions from Plugin [%s]", pluginList[x])
			continue
		}

		p, ok := symbol.(func() string)
		if !ok {
			log.Errorf("Unable to read functions from Plugin [%s]", pluginList[x])
			continue
		}

		fmt.Println(p())
	}
}

//UsagePlugin returns the usage of a plugin function
func UsagePlugin(pluginPath, action string) {

	symbol, err := findFunctionInPlugin(pluginPath, "ParlayUsage")
	if err != nil {
		log.Errorf("Unable to read functions from Plugin [%s]", pluginPath)
		return
	}

	p, ok := symbol.(func(string) string)
	if !ok {
		log.Errorf("Unable to read functions from Plugin [%s]", pluginPath)
		return
	}
	result := p(action)
	fmt.Printf("%s\n", result)
}
