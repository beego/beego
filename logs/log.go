package logs

import (
	"fmt"
	"sync"
)

const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelCritical
)

type loggerType func() LoggerInterface

type LoggerInterface interface {
	Init(config string) error
	WriteMsg(msg string, level int) error
	Destroy()
	Flush()
}

var adapters = make(map[string]loggerType)

// Register makes a log provide available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, log loggerType) {
	if log == nil {
		panic("logs: Register provide is nil")
	}
	if _, dup := adapters[name]; dup {
		panic("logs: Register called twice for provider " + name)
	}
	adapters[name] = log
}

type BeeLogger struct {
	lock    sync.Mutex
	level   int
	msg     chan *logMsg
	outputs map[string]LoggerInterface
}

type logMsg struct {
	level int
	msg   string
}

// config need to be correct JSON as string: {"interval":360}
func NewLogger(channellen int64) *BeeLogger {
	bl := new(BeeLogger)
	bl.msg = make(chan *logMsg, channellen)
	bl.outputs = make(map[string]LoggerInterface)
	//bl.SetLogger("console", "") // default output to console
	go bl.StartLogger()
	return bl
}

func (bl *BeeLogger) SetLogger(adaptername string, config string) error {
	bl.lock.Lock()
	defer bl.lock.Unlock()
	if log, ok := adapters[adaptername]; ok {
		lg := log()
		lg.Init(config)
		bl.outputs[adaptername] = lg
		return nil
	} else {
		return fmt.Errorf("logs: unknown adaptername %q (forgotten Register?)", adaptername)
	}
}

func (bl *BeeLogger) DelLogger(adaptername string) error {
	bl.lock.Lock()
	defer bl.lock.Unlock()
	if lg, ok := bl.outputs[adaptername]; ok {
		lg.Destroy()
		delete(bl.outputs, adaptername)
		return nil
	} else {
		return fmt.Errorf("logs: unknown adaptername %q (forgotten Register?)", adaptername)
	}
}

func (bl *BeeLogger) writerMsg(loglevel int, msg string) error {
	if bl.level > loglevel {
		return nil
	}
	lm := new(logMsg)
	lm.level = loglevel
	lm.msg = msg
	bl.msg <- lm
	return nil
}

func (bl *BeeLogger) SetLevel(l int) {
	bl.level = l
}

func (bl *BeeLogger) StartLogger() {
	for {
		select {
		case bm := <-bl.msg:
			for _, l := range bl.outputs {
				l.WriteMsg(bm.msg, bm.level)
			}
		}
	}
}

func (bl *BeeLogger) Trace(format string, v ...interface{}) {
	msg := fmt.Sprintf("[T] "+format, v...)
	bl.writerMsg(LevelTrace, msg)
}

func (bl *BeeLogger) Debug(format string, v ...interface{}) {
	msg := fmt.Sprintf("[D] "+format, v...)
	bl.writerMsg(LevelDebug, msg)
}

func (bl *BeeLogger) Info(format string, v ...interface{}) {
	msg := fmt.Sprintf("[I] "+format, v...)
	bl.writerMsg(LevelInfo, msg)
}

func (bl *BeeLogger) Warn(format string, v ...interface{}) {
	msg := fmt.Sprintf("[W] "+format, v...)
	bl.writerMsg(LevelWarn, msg)
}

func (bl *BeeLogger) Error(format string, v ...interface{}) {
	msg := fmt.Sprintf("[E] "+format, v...)
	bl.writerMsg(LevelError, msg)
}

func (bl *BeeLogger) Critical(format string, v ...interface{}) {
	msg := fmt.Sprintf("[C] "+format, v...)
	bl.writerMsg(LevelCritical, msg)
}

//flush all chan data
func (bl *BeeLogger) Flush() {
	for _, l := range bl.outputs {
		l.Flush()
	}
}

func (bl *BeeLogger) Close() {
	for {
		if len(bl.msg) > 0 {
			bm := <-bl.msg
			for _, l := range bl.outputs {
				l.WriteMsg(bm.msg, bm.level)
			}
		} else {
			break
		}
	}
	for _, l := range bl.outputs {
		l.Flush()
		l.Destroy()
	}
}
