package parlay

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/mitchellh/go-homedir"
)

//Restore provides a checkpoint to resume from
type Restore struct {
	Deployment string   `json:"deployment"` // Name of deployment to restore from
	Action     string   `json:"action"`     // Action to restore from
	Host       string   `json:"host"`       // Single host to start from
	Hosts      []string `json:"hosts"`      // Restart operation on a number of hosts
}

// restore is an interal struct used for execution restoration
var restore Restore

const restoreFile = ".parlay_restore"

// restoreFilePath will build a path where a file will be read/writted
func restoreFilePath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return home + "/" + restoreFile, nil
}

func (r *Restore) createCheckpoint() error {
	// This function will create a checkpoint file that will allow Plunder to restart in the event of failure
	path, err := restoreFilePath()
	if err != nil {
		return err
	}

	// Marshall the struct to a byte array
	b, err := json.Marshal(r)
	if err != nil {
		return err
	}
	// Write the checkpoint file
	err = ioutil.WriteFile(path, b, 0644)

	return err
}

//RestoreFromCheckpoint will attempt to find a restoration checkpoint file
func RestoreFromCheckpoint() *Restore {
	path, err := restoreFilePath()
	if err != nil {
		return nil
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return nil
		}
		var r Restore
		err = json.Unmarshal(b, &r)
		if err != nil {
			return nil
		}
		return &r
	}
	return nil
}
