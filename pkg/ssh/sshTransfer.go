package ssh

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pkg/sftp"
)

// ParalellDownload - Allow downloading a file over SFTP from multiple hosts in parallel
func ParalellDownload(hosts []HostSSHConfig, source, destination string, to int) []CommandResult {
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

			if err := host.DownloadFile(source, destination); err != nil {
				res.Error = err
			} else {
				res.Result = "Download completed"
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
				Error:  fmt.Errorf("Download Timed out"),
				Result: "",
			}
			cmdResults = append(cmdResults, failedCommand)

		}
	}
	return cmdResults
}

// DownloadFile -
func (c HostSSHConfig) DownloadFile(source, destination string) error {
	var err error
	c.Connection, err = c.StartConnection()
	if err != nil {
		return err
	}

	// New SFTP client
	sftp, err := sftp.NewClient(c.Connection)
	if err != nil {
		return err
	}
	defer sftp.Close()

	// Open remote source
	sftpSource, err := sftp.Open(source)
	if err != nil {
		return err
	}
	defer sftpSource.Close()

	// Open local destination
	localDestination, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer localDestination.Close()

	//
	_, err = sftpSource.WriteTo(localDestination)
	if err != nil {
		return err
	}

	// An error here isn't cause for alarm, any new transaction should create a new connection
	_ = c.StopConnection()

	return nil
}

// ParalellUpload - Allow uploading a file over SFTP to multiple hosts in parallel
func ParalellUpload(hosts []HostSSHConfig, source, destination string, to int) []CommandResult {
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

			if err := host.UploadFile(source, destination); err != nil {
				res.Error = err
			} else {
				res.Result = "Upload completed"
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
				Error:  fmt.Errorf("Upload Timed out"),
				Result: "",
			}
			cmdResults = append(cmdResults, failedCommand)

		}
	}
	return cmdResults
}

// UploadFile -
func (c HostSSHConfig) UploadFile(source, destination string) error {
	var err error
	c.Connection, err = c.StartConnection()
	if err != nil {
		return err
	}
	// New SFTP client
	sftp, err := sftp.NewClient(c.Connection)
	if err != nil {
		return err
	}
	defer sftp.Close()

	// Open remote source
	sftpDestination, err := sftp.Create(destination)
	if err != nil {
		return err
	}
	defer sftpDestination.Close()

	// Open local destination
	localSource, err := os.Open(source)
	if err != nil {
		return err
	}
	defer localSource.Close()

	// copy source file to destination file
	_, err = io.Copy(sftpDestination, localSource)
	if err != nil {
		return err
	}

	// An error here isn't cause for alarm, any new transaction should create a new connection
	_ = c.StopConnection()

	return nil
}
