package parlay

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/plunder-app/plunder/pkg/parlay/parlaytypes"
	"github.com/plunder-app/plunder/pkg/ssh"
	log "github.com/sirupsen/logrus"
)

func buildCommand(a parlaytypes.Action) (string, error) {
	var command string

	// An executable Key takes presedence
	if a.KeyName != "" {
		keycmd := Keys[a.KeyName]
		// Check that the key exists
		if keycmd == "" {
			return "", fmt.Errorf("Unable to find command under key '%s'", a.KeyName)

		}
		if a.CommandSudo != "" {
			// Add sudo to the Key command
			command = fmt.Sprintf("sudo -n -u %s %s", a.CommandSudo, keycmd)
		} else {
			command = keycmd
		}
	} else {
		// Not using a key, using a shell command
		if a.CommandSudo != "" {
			// Add sudo to the Shell command
			command = fmt.Sprintf("sudo -n -u %s %s", a.CommandSudo, a.Command)
		} else {
			command = a.Command
		}
	}
	return command, nil
}

func parseAndExecute(a parlaytypes.Action, h *ssh.HostSSHConfig) ssh.CommandResult {
	// This will parse the options passed in the action and execute the required string
	var cr ssh.CommandResult
	var b []byte

	command, err := buildCommand(a)
	if err != nil {
		cr.Error = err
		return cr
	}

	if a.CommandLocal == true {
		log.Debugf("Command [%s]", command)
		cmd := exec.Command("bash", "-c", command)
		b, cr.Error = cmd.CombinedOutput()
		if cr.Error != nil {
			return cr
		}
		cr.Result = strings.TrimRight(string(b), "\r\n")
	} else {
		log.Debugf("Executing command [%s] on host [%s]", command, h.Host)
		cr = ssh.SingleExecute(command, a.CommandPipeFile, a.CommandPipeCmd, *h, a.Timeout)

		cr.Result = strings.TrimRight(cr.Result, "\r\n")

		// If the command hasn't returned anything, put a filler in
		if cr.Result == "" {
			cr.Result = "[No Output]"
		}
		if cr.Error != nil {
			return cr
		}
	}

	// Save the results into a key to be used at another point
	if a.CommandSaveAsKey != "" {
		log.Debugf("Adding new results to key [%s]", a.CommandSaveAsKey)
		Keys[a.CommandSaveAsKey] = cr.Result
	}

	// Save the results into a file to be used at another point
	if a.CommandSaveFile != "" {
		var f *os.File
		f, cr.Error = os.Create(a.CommandSaveFile)
		if cr.Error != nil {
			return cr
		}

		defer f.Close()

		_, cr.Error = f.WriteString(cr.Result)
		if cr.Error != nil {
			return cr
		}
		f.Sync()
	}

	return cr
}
