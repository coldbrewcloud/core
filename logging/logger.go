package logging

import (
	"os"

	_logging "github.com/op/go-logging"
)

type Logger struct {
	logger        *_logging.Logger
	currentFormat string
	currentLevel  Level
}

func (l *Logger) updateLevelFormat() {
	backend := _logging.NewBackendFormatter(
		_logging.NewLogBackend(os.Stdout, "", 0),
		_logging.MustStringFormatter(l.currentFormat))

	leveledBackend := _logging.AddModuleLevel(backend)
	leveledBackend.SetLevel(_logging.Level(l.currentLevel), "")

	l.logger.SetBackend(leveledBackend)
}

func (l *Logger) SetLevel(level Level) {
	l.currentLevel = level
	l.updateLevelFormat()
}

func (l *Logger) SetFormat(format string) {
	l.currentFormat = format
	l.updateLevelFormat()
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	l.logger.Warningf(format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *Logger) Criticalf(format string, args ...interface{}) {
	l.logger.Criticalf(format, args...)
}
