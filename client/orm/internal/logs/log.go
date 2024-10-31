package logs

import (
	"io"
	"log"
)

// Log implement the log.Logger
type Log struct {
	*log.Logger
}

func NewLog(out io.Writer) *Log {
	d := new(Log)
	d.Logger = log.New(out, "[ORM]", log.LstdFlags)
	return d
}
