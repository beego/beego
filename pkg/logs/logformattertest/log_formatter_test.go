package logformattertest

import (
	"fmt"
	"testing"

	"github.com/astaxie/beego/pkg/common"
	"github.com/astaxie/beego/pkg/logs"
)

func customFormatter(lm *logs.LogMsg) string {
	return fmt.Sprintf("[CUSTOM CONSOLE LOGGING] %s", lm.Msg)
}

func globalFormatter(lm *logs.LogMsg) string {
	return fmt.Sprintf("[GLOBAL] %s", lm.Msg)
}

func TestCustomLoggingFormatter(t *testing.T) {
	// beego.BConfig.Log.AccessLogs = true

	logs.SetLoggerWithOpts("console", []string{`{"color":true}`}, common.SimpleKV{Key: "formatter", Value: customFormatter})

	// Message will be formatted by the customFormatter with colorful text set to true
	logs.Informational("Test message")
}

func TestGlobalLoggingFormatter(t *testing.T) {
	logs.SetGlobalFormatter(globalFormatter)

	logs.SetLogger("console", `{"color":true}`)

	// Message will be formatted by globalFormatter
	logs.Informational("Test message")

}
