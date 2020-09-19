package logs

import (
	"strconv"
	"testing"
	"time"
)

func TestPatternLogFormatter(t *testing.T) {
	tes := PatternLogFormatter{
		Pattern:    "%F:%n|%w%t>> %m",
		WhenFormat: "2006-01-02",
	}
	when := time.Now()
	lm := &LogMsg{
		Msg:        "message",
		FilePath:   "/User/go/beego/main.go",
		Level:      LevelWarn,
		LineNumber: 10,
		When:       when,
	}
	got := tes.ToString(lm)
	want := lm.FilePath + ":" + strconv.Itoa(lm.LineNumber) + "|" +
		when.Format(tes.WhenFormat) + levelPrefix[lm.Level-1] + ">> " + lm.Msg
	if got != want {
		t.Errorf("want %s, got %s", want, got)
	}
}
