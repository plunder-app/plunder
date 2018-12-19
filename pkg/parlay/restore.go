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

	return nil
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
		err = json.Unmarshal(b, r)
		if err != nil {
			return nil
		}
		return &r
	}
	return nil
}
