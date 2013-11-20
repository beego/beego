package toolbox

import (
	"os"
	"testing"
)

func TestProcessInput(t *testing.T) {
	ProcessInput("lookup goroutine", os.Stdout)
	ProcessInput("lookup heap", os.Stdout)
	ProcessInput("lookup threadcreate", os.Stdout)
	ProcessInput("lookup block", os.Stdout)
	ProcessInput("gc summary", os.Stdout)
}
