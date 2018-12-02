package ssh

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

//Execute -
func Execute(cmd string, hosts []HostSSHConfig, to int) {
	// Run parallel ssh session (max 10)
	results := make(chan string, 10)
	timeout := time.After(time.Duration(to) * time.Second)

	// Execute command on hosts
	for _, host := range hosts {
		go func(host HostSSHConfig) {
			var result string

			if text, err := host.ExecuteCmd(cmd); err != nil {
				result = err.Error()
			} else {
				result = text
			}

			results <- fmt.Sprintf("%s > %s -> %s", host, cmd, result)
		}(host)
	}

	for i := 0; i < len(hosts); i++ {
		select {
		case res := <-results:
			if res != "" {
				fmt.Printf(res)
			}
		case <-timeout:
			//color.Red("Timed out!")
			return
		}
	}
}

//ExecuteSingleCommand -
func ExecuteSingleCommand(cmd string, host HostSSHConfig, to int) {
	// Run parallel ssh session (max 10)
	results := make(chan string, 10)
	timeout := time.After(time.Duration(to) * time.Second)

	// // Execute command on hosts
	// for _, host := range hosts {
	// 	go func(host HostSSHConfig) {
	var result string

	if text, err := host.ExecuteCmd(cmd); err != nil {
		result = err.Error()
	} else {
		result = text
	}

	results <- fmt.Sprintf("%s > %s\n%s\n", host, cmd, result)
	// 	}(host)
	// }

	//for i := 0; i < len(hosts); i++ {
	select {
	case res := <-results:
		if res != "" {
			fmt.Println(res)
		}
	case <-timeout:
		//color.Red("Timed out!")
		return
	}
	//}
}

// StartConnection -
func (c *HostSSHConfig) StartConnection() (*ssh.Client, error) {
	var err error

	host := c.Host
	if !strings.ContainsAny(c.Host, ":") {
		host = host + ":22"
	}
	log.Printf("%v", c)
	c.Connection, err = ssh.Dial("tcp", host, c.ClientConfig)
	if err != nil {
		return nil, err
	}
	return c.Connection, nil
}

// StartSession -
func (c *HostSSHConfig) StartSession() (*ssh.Session, error) {
	var err error
	c.Connection, err = c.StartConnection()
	if err != nil {
		return nil, err
	}
	c.Session, err = c.Connection.NewSession()
	if err != nil {
		return nil, err
	}
	return c.Session, err
}

// StopSession -
func (c *HostSSHConfig) StopSession() {
	if c.Session != nil {
		c.Session.Close()
	}
}

// ExecuteCmd -
func (c *HostSSHConfig) ExecuteCmd(cmd string) (string, error) {
	if c.Session == nil {
		if _, err := c.StartSession(); err != nil {
			return "", err
		}
	}

	var stdoutBuf bytes.Buffer
	c.Session.Stdout = &stdoutBuf
	c.Session.Run(cmd)

	return stdoutBuf.String(), nil
}

// To string
func (c HostSSHConfig) String() string {
	return c.User + "@" + c.Host
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
	return nil
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
		log.Fatal(err)
	}
	return nil
}
