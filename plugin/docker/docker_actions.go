package main

import (
	"fmt"
	"path"

	"github.com/thebsdbox/plunder/pkg/parlay/types"
)

func (i *image) generateActions(host string) []types.Action {
	var generatedActions []types.Action
	var a types.Action
	var sshString, dockerRemoteString, dockerLocalString string

	// This should be set to true if sudo (NOPASSWD) is enabled and required on the local host
	if i.DockerLocalSudo == true {
		dockerLocalString = "sudo docker save"
	} else {
		dockerLocalString = "docker save"
	}

	// This should be set to true if sudo (NOPASSWD) is enabled and required on the remote host
	if i.DockerRemoteSudo == true {
		dockerRemoteString = "sudo docker"
	} else {
		dockerRemoteString = "docker"
	}

	if i.DisableSSHSecurity == true {
		sshString = fmt.Sprintf("ssh -o GlobalKnownHostsFile=/dev/null -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no %s@%s ", i.DockerUser, host)
	} else {
		sshString = fmt.Sprintf("ssh %s@%s %s", i.DockerUser, host, dockerRemoteString)
	}

	if i.ImageFile != "" {
		// If we've specified a file (tarball, or tar+gzip) we cat then pipe over SSH to a docker load
		// TODO - Look at using the crypto/ssh library with a stdin pipe
		a = types.Action{
			ActionType:   "command",
			Command:      fmt.Sprintf("cat %s | %s %s load ", i.ImageFile, sshString, dockerRemoteString),
			CommandLocal: true,
			Name:         fmt.Sprintf("Upload container image %s to remote docker host", path.Base(i.ImageFile)),
		}
		generatedActions = append(generatedActions, a)
	} else if i.ImageName != "" {

		// If we've specified a an existing image from the local docker image store then we "save" it (pipe to stdin)
		// then we can cat then pipe over SSH to a docker load
		a = types.Action{
			ActionType:   "command",
			Command:      fmt.Sprintf("%s %s | %s %s load", dockerLocalString, i.ImageName, sshString, dockerRemoteString),
			CommandLocal: true,
			Name:         fmt.Sprintf("Upload container image %s to remote docker host", i.ImageName),
		}
		generatedActions = append(generatedActions, a)
	}

	// If the downloaded tarball contains a bizarre repo/tag then we can rename/(retag) the image locally
	if i.ImageRetag != "" && i.ImageName != "" {
		a = types.Action{
			ActionType:  "command",
			Command:     fmt.Sprintf("%s tag %s %s", dockerRemoteString, i.ImageName, i.ImageRetag),
			CommandSudo: "root",
			Name:        fmt.Sprintf("Retag %s --> %s", i.ImageName, i.ImageRetag),
		}
		generatedActions = append(generatedActions, a)
	}
	return generatedActions
}
