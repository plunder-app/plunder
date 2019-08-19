package plunderlogging

import "fmt"

type Logger struct {
	json JSONLogger
	file FileLogger
}

func (l *Logger) EnableJSONLogging(e bool) {
	l.json.enabled = e
	l.json.initJSONLogger()
}

func (l *Logger) EnableFileLogging(e bool) {
	l.file.enabled = e
}

func (l *Logger) InitLogFile(path string) error {
	if l.file.enabled != true {
		return l.file.initFileLogger(path)
	}
	// Dont re-initialise the file
	return nil

}

func (l *Logger) InitJSON() {
	// Dont re-initialise the json

	if l.json.enabled != true {
		l.json.initJSONLogger()
	}

}

// target - the entity we're affecting
// entry - the results of the operation on the target

// WriteLogEntry will capture what is transpiring and where
func (l *Logger) WriteLogEntry(target, task, entry, err string) {
	if l.file.enabled {
		l.file.writeEntry(target, entry)
	}
	if l.json.enabled {
		l.json.writeEntry(target, task, entry, err)
	}

	// A logging system shouldnt break anything so any errors are just outputed to STDOUT

}

func (l *Logger) SetLoggingState(target, state string) {
	if l.file.enabled {
		l.file.setLoggingState(target, state)
	}
	if l.json.enabled {
		l.json.setLoggingState(target, state)
	}

	// A logging system shouldnt break anything so any errors are just outputed to STDOUT

}

func (l *Logger) GetJSONLogs(target string) (*JSONLog, error) {
	if l.json.logger == nil {
		return nil, fmt.Errorf("JSON Logging hasn't been enabled")
	}
	// Check if the logger exists
	existingLog, ok := l.json.logger[target]
	if ok {
		return existingLog, nil
	}
	return nil, fmt.Errorf("No Logs for Target [%s] exist", target)
}
