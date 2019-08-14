package plunderlogging

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
	return l.file.initFileLogger(path)
}

func (l *Logger) InitJSON() {
	l.json.initJSONLogger()
}

// target - the entity we're affecting
// entry - the results of the operation on the target

// WriteLogEntry will capture what is transpiring and where
func (l *Logger) WriteLogEntry(target, entry string) {
	if l.file.enabled {
		l.file.writeEntry(target, entry)
	}
	if l.json.enabled {
		l.json.writeEntry(target, entry)
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
