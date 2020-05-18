package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/dustin/go-humanize"
)

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer interface
// and we can pass this into io.TeeReader() which will report progress on each write cycle.
type WriteCounter struct {
	Total uint64
}

var data []byte

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	return n, nil
}

func tickerProgress(byteCounter uint64) {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(byteCounter))
	fmt.Println("")
}

func imageHandler(w http.ResponseWriter, r *http.Request) {

	log.Infof("Incoming image from [%s]", r.RemoteAddr)

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("BootyImage")
	if handler != nil {
		log.Infof("Beginning to recieve image [%s]", handler.Filename)
	}

	if err != nil {
		log.Errorf("%v", err)
		return
	}
	defer file.Close()

	out, err := os.OpenFile(handler.Filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer out.Close()

	// Create our progress reporter and pass it to be used alongside our writer
	ticker := time.NewTicker(500 * time.Millisecond)
	counter := &WriteCounter{}

	go func() {
		for ; true; <-ticker.C {
			tickerProgress(counter.Total)
		}
	}()
	if _, err = io.Copy(out, io.TeeReader(file, counter)); err != nil {
		log.Errorf("%v", err)
	}

	log.Infof("Written of image [%s] to disk", handler.Filename)
	ticker.Stop()

	w.WriteHeader(http.StatusOK)
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(data)
}

// Serve will start the webserver for BOOTy images
func (c *BootController) serveImageHTTP() error {

	fs := http.FileServer(http.Dir("./images"))
	http.HandleFunc("/image", imageHandler)
	http.Handle("/images/", http.StripPrefix("/images/", fs))
	log.Println("Plunder OS Image Services --> Starting HTTP :3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
