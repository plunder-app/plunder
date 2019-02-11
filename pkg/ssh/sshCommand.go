package ssh

import (
	"fmt"
	"time"
)

// SingleExecute - This will execute a command on a single host
func SingleExecute(cmd string, host HostSSHConfig, to int) CommandResult {
	var configs []HostSSHConfig
	configs = append(configs, host)
	result := ParalellExecute(cmd, configs, to)
	return result[0]
}

//ParalellExecute - This will execute the same command in paralell across multiple hosts
func ParalellExecute(cmd string, hosts []HostSSHConfig, to int) []CommandResult {
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

			if text, err := host.ExecuteCmd(cmd); err != nil {
				// Report any returned values
				res.Error = err
				res.Result = text
			} else {
				res.Result = text
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

// ExecuteCmd -
func (c *HostSSHConfig) ExecuteCmdWithStdin(cmd string) (string, error) {
	if c.Session == nil {
		if _, err := c.StartSession(); err != nil {
			return "", err
		}
	}

	b, err := c.Session.CombinedOutput(cmd)

	return string(b), err
}
