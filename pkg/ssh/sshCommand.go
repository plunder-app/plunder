package ssh

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
)

// SingleExecute - This will execute a command on a single host
func SingleExecute(cmd, pipefile, pipecmd string, host HostSSHConfig, to int) CommandResult {
	var configs []HostSSHConfig
	configs = append(configs, host)
	result := ParalellExecute(cmd, pipefile, pipecmd, configs, to)
	return result[0]
}

//ParalellExecute - This will execute the same command in paralell across multiple hosts
func ParalellExecute(cmd, pipefile, pipecmd string, hosts []HostSSHConfig, to int) []CommandResult {
	var cmdResults []CommandResult
	// Run parallel ssh session (max 10)
	results := make(chan CommandResult, 10)
	var d time.Duration

	// Calculate the timeout
	if to == 0 {
		// If no timeout then default to one year (TODO)
		d = time.Duration(8760) * time.Hour
	} else {
		d = time.Duration(to) * time.Second
	}

	// Set the timeout
	timeout := time.After(d)
	// Execute command on hosts
	for _, host := range hosts {
		go func(host HostSSHConfig) {
			res := new(CommandResult)
			res.Host = host.Host

			if pipefile != "" {
				if text, err := host.ExecuteCmdWithStdinFile(cmd, pipefile); err != nil {
					// Report any returned values
					res.Error = err
					res.Result = text
				} else {
					res.Result = text
				}
			} else if pipecmd != "" {
				if text, err := host.ExecuteCmdWithStdinCmd(cmd, pipecmd); err != nil {
					// Report any returned values
					res.Error = err
					res.Result = text
				} else {
					res.Result = text
				}
			} else {
				if text, err := host.ExecuteCmd(cmd); err != nil {
					// Report any returned values
					res.Error = err
					res.Result = text
				} else {
					res.Result = text
				}
			}
			results <- *res
		}(host)
	}

	for i := 0; i < len(hosts); i++ {
		select {
		case res := <-results:
			// Append the results of a succesfull command
			cmdResults = append(cmdResults, res)
		case <-timeout:
			// In the event that a command times out then append the details
			failedCommand := CommandResult{
				Host:   hosts[i].Host,
				Error:  fmt.Errorf("Command Timed out"),
				Result: "",
			}
			cmdResults = append(cmdResults, failedCommand)

		}
	}
	return cmdResults
}

// ExecuteCmd -
func (c *HostSSHConfig) ExecuteCmd(cmd string) (string, error) {
	if c.Session == nil {
		if _, err := c.StartSession(); err != nil {
			return "", err
		}
	}

	b, err := c.Session.CombinedOutput(cmd)

	return string(b), err
}

// ExecuteCmdWithStdinFile -
func (c *HostSSHConfig) ExecuteCmdWithStdinFile(cmd, filePath string) (string, error) {
	if c.Session == nil {
		if _, err := c.StartSession(); err != nil {
			return "", err
		}
	}

	// get a stdin pipe
	si, err := c.Session.StdinPipe()
	if err != nil {
		return "", err
	}

	// get a stdout pipe
	so, err := c.Session.StdoutPipe()
	if err != nil {
		return "", err
	}

	// open file as an io.reader
	// Also resolve the absolute path just incase
	absPath, _ := filepath.Abs(filePath)
	file, err := os.Open(absPath)
	if err != nil {
		return "", err
	}

	// Start a command on our remote session, this should be something that is expecting stdin
	if err := c.Session.Start(cmd); err != nil {
		return "", err
	}

	// do the actual work
	n, err := io.Copy(si, file)
	if err != nil {
		return "", err
	}
	// Close the stdin as we've finised transmitting the data
	si.Close()

	log.Debugf("Copied %d bytes over the stdin pipe", n)
	// wait for process to finishe
	if err := c.Session.Wait(); err != nil {
		return "", err
	}

	// Read all the data from the bu
	var b []byte
	if b, err = ioutil.ReadAll(so); err != nil {
		return "", err

	}
	return string(b), nil

}

// ExecuteCmdWithStdinCmd -
func (c *HostSSHConfig) ExecuteCmdWithStdinCmd(cmd, localCmd string) (string, error) {
	if c.Session == nil {
		if _, err := c.StartSession(); err != nil {
			return "", err
		}
	}

	// get a stdin pipe
	si, err := c.Session.StdinPipe()
	if err != nil {
		return "", err
	}

	// get a stdout pipe
	so, err := c.Session.StdoutPipe()
	if err != nil {
		return "", err
	}

	// Start a command on our remote session, this should be something that is expecting stdin
	if err := c.Session.Start(cmd); err != nil {
		return "", err
	}

	// Start our local command that should be exposing something through stdout
	localExecCmd := exec.Command("bash", "-c", localCmd)
	localso, err := localExecCmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	err = localExecCmd.Start()
	if err != nil {
		return "", err
	}

	// do the actual work
	n, err := io.Copy(si, localso)
	if err != nil {
		return "", err
	}

	// Close the stdin/stdout as we've finised transmitting the data
	si.Close()
	localso.Close()

	log.Debugf("Copied %d bytes over the stdin pipe", n)

	// Wait for local process to finish
	err = localExecCmd.Wait()
	if err != nil {
		return "", err
	}

	// wait for remote process to finish
	if err := c.Session.Wait(); err != nil {
		return "", err
	}

	// Read all the data from the bu
	var b []byte
	if b, err = ioutil.ReadAll(so); err != nil {
		return "", err
	}
	return string(b), nil

}
