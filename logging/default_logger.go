package logging

var defaultLogger *Logger

func init() {
	defaultLogger = MustGet("__default__")
	defaultLogger.SetFormat(`%{color}%{time:15:04:05.000} %{level:.4s}%{color:reset} %{message}`)
}

func SetLevel(level Level) {
	defaultLogger.currentLevel = level
	defaultLogger.updateLevelFormat()
}

func SetFormat(format string) {
	defaultLogger.currentFormat = format
	defaultLogger.updateLevelFormat()
}

func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

func Warningf(format string, args ...interface{}) {
	defaultLogger.Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}
