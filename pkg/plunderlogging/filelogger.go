package plunderlogging

import (
	"fmt"
	"os"
	"sync"
)

// FileLogger allows parlay to log output to a file on the local filesystem
type FileLogger struct {
	enabled bool
	f       *os.File
}

var fileLogging FileLogger

func (l *FileLogger) initFileLogger(logFile string) (err error) {
	l.enabled = true
	l.f, err = os.Create(logFile)
	if err != nil {
		return err
	}
	return nil
}

// This file based logging function may error, but logging should never break the running of a system, so errors are passed to "Debug" logging
func (l *FileLogger) writeEntry(target, entry string) error {
	var fileMutex sync.Mutex
	if l.enabled == true {

		// As this may be called by numerous goroutines, we impose a mutex lock on it
		fileMutex.Lock()
		defer fileMutex.Unlock()

		// TODO - Does this produce readable logging output
		_, err := l.f.WriteString(fmt.Sprintf("Target=%s Entry=%s", target, entry))

		return err
	}
	return nil
}

func (l *FileLogger) setLoggingState(target, state string) error {

	return nil
}
