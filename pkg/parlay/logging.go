package parlay

import (
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
)

type fileLogger struct {
	enabled bool
	f       *os.File
}

var logging fileLogger

func (l *fileLogger) init(logFile string) (err error) {
	l.enabled = true
	l.f, err = os.Create(logFile)
	if err != nil {
		return err
	}
	return nil
}

// This file based logging function may error, but logging should never break the running of a system, so errors are passed to "Debug" logging
func (l *fileLogger) writeString(logMessage string) {
	var fileMutex sync.Mutex
	if l.enabled == true {

		// As this may be called by numerous goroutines, we impose a mutex lock on it
		fileMutex.Lock()
		defer fileMutex.Unlock()

		_, err := l.f.WriteString(logMessage)
		if err != nil {
			log.Debugf("%v", err)
		}
	}
}
