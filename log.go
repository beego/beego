package beego

import (
	"github.com/astaxie/beego/logs"
)

// Log levels to control the logging output.
const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarning
	LevelError
	LevelCritical
)

// SetLogLevel sets the global log level used by the simple
// logger.
func SetLevel(l int) {
	BeeLogger.SetLevel(l)
}

// logger references the used application logger.
var BeeLogger *logs.BeeLogger

func init() {
	BeeLogger = logs.NewLogger(10000)
	BeeLogger.SetLogger("console", "")
}

// SetLogger sets a new logger.
func SetLogger(adaptername string, config string) {
	BeeLogger.SetLogger(adaptername, config)
}

// Trace logs a message at trace level.
func Trace(v ...interface{}) {
	BeeLogger.Trace("%v", v...)
}

// Debug logs a message at debug level.
func Debug(v ...interface{}) {
	BeeLogger.Debug("%v", v...)
}

// Info logs a message at info level.
func Info(v ...interface{}) {
	BeeLogger.Info("%v", v...)
}

// Warning logs a message at warning level.
func Warn(v ...interface{}) {
	BeeLogger.Warn("%v", v...)
}

// Error logs a message at error level.
func Error(v ...interface{}) {
	BeeLogger.Error("%v", v...)
}

// Critical logs a message at critical level.
func Critical(v ...interface{}) {
	BeeLogger.Critical("%v", v...)
}
