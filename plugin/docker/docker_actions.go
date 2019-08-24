package main

import (
	"fmt"
	"path"

	"github.com/plunder-app/plunder/pkg/parlay/parlaytypes"
)

func (i *image) generateImageActions(host string) []parlaytypes.Action {
	var generatedActions []parlaytypes.Action
	var a parlaytypes.Action
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
			a = parlaytypes.Action{
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

			a = parlaytypes.Action{
				ActionType:     "command",
				Command:        fmt.Sprintf("%s load", dockerRemoteString),
				CommandPipeCmd: fmt.Sprintf("%s %s", dockerLocalString, i.ImageNames[y]),
				Name:           fmt.Sprintf("Upload container image %s to remote docker host", i.ImageNames[y]),
			}
			generatedActions = append(generatedActions, a)
		}
	}

	return generatedActions
}

func (t *tag) generateTagActions(host string) ([]parlaytypes.Action, error) {

	if len(t.SourceNames) != len(t.TargetNames) {
		return nil, fmt.Errorf("The number of images to retag doesn't match the number of tags")
	}
	var generatedActions []parlaytypes.Action

	// Iterate through all of the images and create retagging actions
	for y := range t.SourceNames {
		// Generate the retag action
		var a = parlaytypes.Action{
			ActionType:  "command",
			Command:     fmt.Sprintf("sudo docker tag %s %s", t.SourceNames[y], t.TargetNames[y]),
			CommandSudo: "root",
			Name:        fmt.Sprintf("Retag %s --> %s", t.SourceNames[y], t.TargetNames[y]),
		}

		generatedActions = append(generatedActions, a)
	}
	return generatedActions, nil
}
