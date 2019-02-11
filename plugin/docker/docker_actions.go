package main

import (
	"fmt"

	"github.com/thebsdbox/plunder/pkg/parlay/types"
)

func (i *image) generateActions(host string) []types.Action {
	var generatedActions []types.Action
	var a types.Action
	if i.ImageFile != "" {
		a = types.Action{
			// Generate etcd server certificate
			ActionType:   "command",
			Command:      fmt.Sprintf("cat %s | ssh %s sudo docker load", i.ImageFile, host),
			CommandLocal: true,
			Name:         "Upload file to remote docker host",
		}
		generatedActions = append(generatedActions, a)
		return generatedActions
	}
	if i.ImageName != "" {
		a = types.Action{
			// Generate etcd server certificate
			ActionType:   "command",
			Command:      fmt.Sprintf("docker save %s | ssh %s sudo docker load", i.ImageFile, host),
			CommandLocal: true,
			Name:         "Upload file to remote docker host",
		}
		generatedActions = append(generatedActions, a)
		return generatedActions
	}
	return nil
}
