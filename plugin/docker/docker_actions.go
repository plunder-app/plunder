package main

import (
	"fmt"
	"path"

	"github.com/thebsdbox/plunder/pkg/parlay/types"
)

func (i *image) generateActions(host string) []types.Action {
	var generatedActions []types.Action
	var a types.Action
	var sshString string
	if i.DisableSSHSecurity == true {
		sshString = fmt.Sprintf("ssh -o GlobalKnownHostsFile=/dev/null -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no %s@%s sudo docker load", i.DockerUser, host)
	} else {
		sshString = fmt.Sprintf("ssh %s@%s sudo docker load", i.DockerUser, host)
	}

	if i.ImageFile != "" {
		a = types.Action{
			// Generate etcd server certificate
			ActionType:   "command",
			Command:      fmt.Sprintf("cat %s | %s", i.ImageFile, sshString),
			CommandLocal: true,
			Name:         fmt.Sprintf("Upload container image %s to remote docker host", path.Base(i.ImageFile)),
		}
		generatedActions = append(generatedActions, a)
		return generatedActions
	}
	if i.ImageName != "" {
		a = types.Action{
			// Generate etcd server certificate
			ActionType:   "command",
			Command:      fmt.Sprintf("docker save %s | %s", i.ImageFile, sshString),
			CommandLocal: true,
			Name:         fmt.Sprintf("Upload container image %s to remote docker host", i.ImageName),
		}
		generatedActions = append(generatedActions, a)
		return generatedActions
	}
	return nil
}
