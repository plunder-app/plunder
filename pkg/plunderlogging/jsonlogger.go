package plunderlogging

import (
	"fmt"
	"time"
)

// JSONLogger allows parlay to log output to an in-memory jsonStruct
type JSONLogger struct {
	enabled bool
	logger  map[string]JSONLog
}

// JSONLog contains all of the output from a parlay execution
type JSONLog struct {
	State   string         `json:"state"`
	Entries []JSONLogEntry `json:"entries"`
}

// JSONLogEntry contains the details a specific action
type JSONLogEntry struct {
	Created time.Time `json:"created"`
	Entry   string    `json:"entry"`
}

func (j *JSONLogger) initJSONLogger() {
	j.enabled = true
	j.logger = make(map[string]JSONLog)
}

func (j *JSONLogger) writeEntry(target, entry string) {

	// Create new entry
	newEntry := JSONLogEntry{
		Created: time.Now(),
		Entry:   entry,
	}

	// Check if the logger exists
	existingLog, ok := j.logger[target]
	if ok {
		// Update an existing entry

		existingLog.Entries = append(existingLog.Entries, newEntry)
	} else {
		// Create a new logger
		newLog := JSONLog{
			State: "Running",
		}
		// Append the entry to it
		newLog.Entries = append(existingLog.Entries, newEntry)
		// Update the in-memory log store
		j.logger[target] = newLog
	}
}

func (j *JSONLogger) deleteLog(target string) error {
	// Check if the entry exists
	_, ok := j.logger[target]
	if ok {
		// If it does, then we use the in-built function to delete the log entry
		delete(j.logger, target)
	} else {
		// Return a warning
		return fmt.Errorf("In-Memory logging for [%s] either doesn't exist or has already been delted", target)
	}
	return nil
}

func (j *JSONLogger) setLoggingState(target, state string) error {
	// Check if the logger exists
	existingLog, ok := j.logger[target]
	if ok {
		// Update an existing entry
		existingLog.State = state
	} else {
		return fmt.Errorf("In-Memory logging for [%s] either doesn't exist or has already been delted", target)
	}
	return nil
}
