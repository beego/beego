package logs

import (
	"testing"
)

func TestConsole(t *testing.T) {
	log := NewLogger(10000)
	log.SetLogger("console", "")
	log.Trace("trace")
	log.Info("info")
	log.Warn("warning")
	log.Debug("debug")
	log.Critical("critical")
	log2 := NewLogger(100)
	log2.SetLogger("console", `{"level":1}`)
	log.Trace("trace")
	log.Info("info")
	log.Warn("warning")
	log.Debug("debug")
	log.Critical("critical")
}

func BenchmarkConsole(b *testing.B) {
	log := NewLogger(10000)
	log.SetLogger("console", "")
	for i := 0; i < b.N; i++ {
		log.Trace("trace")
	}
}
