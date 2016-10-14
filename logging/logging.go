package logging

import (
	"sync"

	_logging "github.com/op/go-logging"
)

var (
	loggers       = make(map[string]*Logger)
	loggersLock   = sync.Mutex{}
	defaultLevel  = INFO
	defaultFormat = `%{color}%{time:15:04:05.000} %{level:.4s}%{color:reset} [%{module}] %{message}`
)

func MustGet(name string) *Logger {
	loggersLock.Lock()
	defer loggersLock.Unlock()

	if l := loggers[name]; l != nil {
		return l
	}

	newLogger := &Logger{
		logger:        _logging.MustGetLogger(name),
		currentLevel:  defaultLevel,
		currentFormat: defaultFormat,
	}

	newLogger.updateLevelFormat()

	loggers[name] = newLogger

	return newLogger
}

func SetAllLevel(level Level) {
	loggersLock.Lock()
	defer loggersLock.Unlock()

	defaultLevel = level

	for _, l := range loggers {
		l.SetLevel(level)
	}
}

func SetAllFormat(format string) {
	loggersLock.Lock()
	defer loggersLock.Unlock()

	defaultFormat = format

	for _, l := range loggers {
		l.SetFormat(format)
	}
}
