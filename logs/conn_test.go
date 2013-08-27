package logs

import (
	"testing"
)

func TestConn(t *testing.T) {
	log := NewLogger(1000)
	log.SetLogger("conn", `{"net":"tcp","addr":":7020"}`)
	log.Info("info")
}
