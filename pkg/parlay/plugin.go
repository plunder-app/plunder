package parlay

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"

	log "github.com/Sirupsen/logrus"
)

// Find plugins
func findPlugins(pluginDir string) []string {
	var plugins []string
	// This function will look for all files in a specified directory (defaults to PWD/plugin)
	filepath.Walk(pluginDir, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			if filepath.Ext(path) == ".plugin" {
				plugins = append(plugins, f.Name())
			}
		}
		return nil
	})
	return plugins
}

//LoadPlugins -
func LoadPlugins() {

	pluginList := findPlugins("./plugin")
	log.Debugf("Found [%d] plugins", len(pluginList))
	for x := range pluginList {
		plug, err := plugin.Open("./plugin/" + pluginList[x])
		if err != nil {
			log.Errorf("Unable to open Plugin [%s], %v", pluginList[x], err)
			continue
		}

		symbol, err := plug.Lookup("ParlayActionList")
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

//LoadPlugins -
func ListPlugins() {

	pluginList := findPlugins("./plugin")
	log.Debugf("Found [%d] plugins", len(pluginList))
	for x := range pluginList {
		plug, err := plugin.Open("./plugin/" + pluginList[x])
		if err != nil {
			log.Errorf("Unable to open Plugin [%s], %v", pluginList[x], err)
			continue
		}

		symbol, err := plug.Lookup("ParlayPluginInfo")
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
