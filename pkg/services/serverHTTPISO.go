package services

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/hooklift/iso9660"
	log "github.com/sirupsen/logrus"
)

// TODO - This currently is inefficient and results in an open/parse of an iso for every file operation.
// github.com/qeedquan/iso9660 may need looking at later on. ( thebsdbox / [1/9/19] )
// Comments are left incase we/I revert

// isoMapper at this point just maps the prefix to the path this may change
var isoMapper map[string]string

// iso9660PathSanitiser will take a "standar" file path and convert it into something that make sense within iso9660 TOC
// The iso9660 constraints:
// - A-Z (uppercase)
// - '_' is the only other character
// - Filename can only be 32 characters (inclucing the terminating semicolon ';')

func iso9660PathSanitiser(unsanitisedPath string) string {
	// Get the filename from the string
	fullFilename := filepath.Base(unsanitisedPath)
	// Get the extension
	extension := filepath.Ext(fullFilename)

	// Remove the extension and leave just the filename
	filename := strings.TrimSuffix(fullFilename, extension)
	// Store the filename and shorten if over 31 characters

	trimmedFilename := filename
	pathLength := len(filename) + len(extension)
	if pathLength > 31 {
		// If the path is too long then we shrink the extension to a seperator and three characters
		if len(extension) > 3 {
			extension = extension[0:4]
		}
		// work out how much of the remaining filename can survive
		trimCount := 31 - len(extension)
		trimmedFilename = filename[0:trimCount]
	}

	rebuiltFileName := strings.ToUpper(fmt.Sprintf("%s%s", trimmedFilename, extension))
	// Find if there is a full stop in the file name
	stopCount := strings.Count(rebuiltFileName, ".")
	var isoFilename string

	switch stopCount {
	case 0:
		// Append one as there is no filepat
		isoFilename = fmt.Sprintf("%s.;1", rebuiltFileName)
	case 1:
		// Not needed, just the semicolon
		isoFilename = fmt.Sprintf("%s;1", rebuiltFileName)
	default:
		// Ensure all other stops are changed to underscores
		isoFilename = fmt.Sprintf("%s;1", strings.Replace(rebuiltFileName, ".", "_", stopCount-1))
	}

	//rebuild the path uppercase
	rebuildPath := strings.ToUpper(fmt.Sprintf("%s/%s", filepath.Dir(unsanitisedPath), isoFilename))

	// strD replacer
	replacer := strings.NewReplacer("+", "_", "-", "_", " ", "_", "~", "_")
	// Format the final output
	isoFormatted := replacer.Replace(rebuildPath)

	return isoFormatted
}

// This takes care of parsing a URL to identify if it should map to an ISO hosted file.

// ISOReader -
func isoReader(w http.ResponseWriter, r *http.Request) {

	// Sanitise the URL, there are a number of steps involved with turning the url into something we can use
	// Remove the beginning slash
	rawURL := strings.TrimLeft(r.URL.String(), "/")

	// Unescape the Http query
	isoURL, err := url.QueryUnescape(rawURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprintf("%s", err.Error()))

		return
	}

	// Split the URL to find the prefix (first part of the URL)
	urlElements := strings.Split(isoURL, "/")
	// Ensure the URL can be parsed
	if len(urlElements) > 1 {
		isoPrefix := urlElements[0]

		isoPath := iso9660PathSanitiser(strings.Replace(isoURL, isoPrefix, "", 1))

		// We now have the ISO prefix to look up files, and the path to look up in the ISO
		//Check for ISO
		log.Debugf("Original URL: %s ISO Path: %s", isoURL, isoPath)
		if _, ok := isoMapper[isoPrefix]; ok {
			file, err := os.Open(isoMapper[isoPrefix])
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				io.WriteString(w, fmt.Sprintf("%s", err.Error()))

				return
			}
			defer file.Close()
			r, err := iso9660.NewReader(file)
			if err != nil {
				fmt.Printf("%v", err)
				return
			}
			for {

				f, err := r.Next()
				if err == io.EOF {
					w.WriteHeader(http.StatusNotFound)
					io.WriteString(w, fmt.Sprintf("Unable to read/find file %s", isoPath))
					return
				}
				if f.Name() == isoPath {
					freader := f.Sys().(io.Reader)
					buf := new(bytes.Buffer)
					buf.ReadFrom(freader)
					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-Type", "application/x-binary")
					io.WriteString(w, buf.String())
					return
				}
			}
			// isoFile, err := isoMapper[isoPrefix].Open(isoPath)
			// if err != nil {
			// 	w.WriteHeader(http.StatusNotFound)
			// 	io.WriteString(w, fmt.Sprintf("%s", err.Error()))
			// 	return
			// }

			// fileStat, err := isoFile.Stat()
			// if err != nil {
			// 	w.WriteHeader(http.StatusNotFound)
			// 	io.WriteString(w, fmt.Sprintf("Unable to stat file on ISO %s", isoPath))
			// 	return
			// }
			// fileBytes = make([]byte, fileStat.Size())
			// _, err = isoFile.Read(fileBytes)
			// if err != nil {
			// 	w.WriteHeader(http.StatusNotFound)
			// 	io.WriteString(w, fmt.Sprintf("Unable to read file on ISO %s", isoPath))
			// 	return
			// }
		}

	} else {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, fmt.Sprintf("Unable to find content ISO Prefix %s", isoURL))
		return
	}
	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, fmt.Sprintf("Unable to find content ISO Prefix %s", isoURL))
	return
}

// OpenISO will open an iso and add it to out ISO Map for reading at a later point
func OpenISO(isoPath, isoPrefix string) error {
	// file, err := os.Open(isoPath)
	// if err != nil {
	// 	return err
	// }

	// f, err := iso9660(isoPath)
	// if err != nil {
	// 	return err
	// }

	if isoMapper == nil {
		// Ensure it is initialised before trying to use it
		isoMapper = make(map[string]string)
	}
	// Add the reader
	isoMapper[isoPrefix] = isoPath

	return nil
}
