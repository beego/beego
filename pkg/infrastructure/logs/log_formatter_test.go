package logs

import (
	"fmt"
	"testing"

	"github.com/astaxie/beego/pkg/infrastructure/utils"
)

func customFormatter(lm *LogMsg) string {
	return fmt.Sprintf("[CUSTOM CONSOLE LOGGING] %s", lm.Msg)
}

func globalFormatter(lm *LogMsg) string {
	return fmt.Sprintf("[GLOBAL] %s", lm.Msg)
}

func TestCustomLoggingFormatter(t *testing.T) {
	// beego.BConfig.Log.AccessLogs = true

	SetLoggerWithOpts("console", []string{`{"color":true}`}, &utils.SimpleKV{Key: "formatter", Value: customFormatter})

	// Message will be formatted by the customFormatter with colorful text set to true
	Informational("Test message")
}

func TestGlobalLoggingFormatter(t *testing.T) {
	SetGlobalFormatter(globalFormatter)

	SetLogger("console", `{"color":true}`)

	// Message will be formatted by globalFormatter
	Informational("Test message")

}
