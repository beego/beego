package logs

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego/logs"
)

func logLevelName2Int(levelName string) (logLevel int, err error) {

	switch strings.ToLower(levelName) {
	case "emergency":
		logLevel = logs.LevelEmergency
	case "alert":
		logLevel = logs.LevelAlert
	case "critical":
		logLevel = logs.LevelCritical
	case "error":
		logLevel = logs.LevelError
	case "warning", "warn":
		logLevel = logs.LevelWarning
	case "notice":
		logLevel = logs.LevelNotice
	case "info", "informational":
		logLevel = logs.LevelInformational
	case "debug":
		logLevel = logs.LevelDebug
	case "trace":
		logLevel = logs.LevelTrace
	default:
		err = fmt.Errorf("level name[%s] is invalid", levelName)
	}

	return logLevel, err
}
