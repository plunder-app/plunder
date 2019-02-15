package main

import (
	"fmt"
	"path"

	"github.com/plunder-app/plunder/pkg/parlay/types"
)

func (i *image) generateImageActions(host string) []types.Action {
	var generatedActions []types.Action
	var a types.Action
	var dockerRemoteString, dockerLocalString string

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

	if len(i.ImageFiles) != 0 {
		// If we've specified a file (tarball, or tar+gzip) we cat then pipe over SSH to a docker load

		for y := range i.ImageFiles {
			a = types.Action{
				ActionType:      "command",
				Command:         fmt.Sprintf("%s load ", dockerRemoteString),
				CommandPipeFile: i.ImageFiles[y],
				Name:            fmt.Sprintf("Upload container image %s to remote docker host", path.Base(i.ImageFiles[y])),
			}
			generatedActions = append(generatedActions, a)
		}
	} else if len(i.ImageNames) != 0 {

		// If we've specified a an existing image from the local docker image store then we "save" it (pipe to stdin)
		// then we can cat then pipe over SSH to a docker load
		for y := range i.ImageNames {

			a = types.Action{
				ActionType:     "command",
				Command:        fmt.Sprintf("%s load", dockerRemoteString),
				CommandPipeCmd: fmt.Sprintf("%s %s", dockerLocalString, i.ImageNames[y]),
				Name:           fmt.Sprintf("Upload container image %s to remote docker host", i.ImageNames[y]),
			}
			generatedActions = append(generatedActions, a)
		}
	}

	// If the downloaded tarball contains a bizarre repo/tag then we can rename/(retag) the image locally
	// if i.ImageRetag != "" && i.ImageName != "" {
	// 	a = types.Action{
	// 		ActionType:  "command",
	// 		Command:     fmt.Sprintf("%s tag %s %s", dockerRemoteString, i.ImageName, i.ImageRetag),
	// 		CommandSudo: "root",
	// 		Name:        fmt.Sprintf("Retag %s --> %s", i.ImageName, i.ImageRetag),
	// 	}
	// 	generatedActions = append(generatedActions, a)
	// }
	return generatedActions
}

func (t *tag) generateTagActions(host string) []types.Action {
	var generatedActions []types.Action

	// Generate the retag action
	var a = types.Action{
		ActionType:  "command",
		Command:     fmt.Sprintf("sudo docker tag %s %s", t.SourceName, t.TargetName),
		CommandSudo: "root",
		Name:        fmt.Sprintf("Retag %s --> %s", t.SourceName, t.TargetName),
	}

	generatedActions = append(generatedActions, a)

	return generatedActions
}
