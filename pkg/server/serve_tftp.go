package server

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	tftp "github.com/thebsdbox/go-tftp/server"
)

var iPXEData []byte

// HandleWrite : writing is disabled in this service
func HandleWrite(filename string) (w io.Writer, err error) {
	err = errors.New("Server is read only")
	return
}

// HandleRead : read a ROfs file and send over tftp
func HandleRead(filename string) (r io.Reader, err error) {
	r = bytes.NewBuffer(iPXEData)
	return
}

// tftp server
func (c *BootController) serveTFTP() error {

	log.Printf("Opening and caching undionly.kpxe")
	f, err := os.Open(*c.PXEFileName)
	if err != nil {
		log.Printf("Please download the bootloader with the pulliPXE command")
		return err
	}
	// Use bufio.NewReader to get a Reader.
	// ... Then use ioutil.ReadAll to read the entire content.
	r := bufio.NewReader(f)

	iPXEData, err = ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	s := tftp.NewServer("", HandleRead, HandleWrite)
	err = s.Serve(*c.TFTPAddress + ":69")
	if err != nil {
		return err
	}
	return nil
}
