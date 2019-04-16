package log

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

//log level, from low to high, more high means more serious
const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

const (
	Ltime  = 1 << iota //time format "2006/01/02 15:04:05"
	Lfile              //file.go:123
	Llevel             //[Trace|Debug|Info...]
)

var LevelName [6]string = [6]string{"Trace", "Debug", "Info", "Warn", "Error", "Fatal"}

const TimeFormat = "2006/01/02 15:04:05"

const maxBufPoolSize = 16

type atomicInt32 int32

func (i *atomicInt32) Set(n int) {
	atomic.StoreInt32((*int32)(i), int32(n))
}

func (i *atomicInt32) Get() int {
	return int(atomic.LoadInt32((*int32)(i)))
}

type Logger struct {
	level atomicInt32
	flag  int

	hMutex  sync.Mutex
	handler Handler

	quit chan struct{}
	msg  chan []byte

	bufMutex sync.Mutex
	bufs     [][]byte

	wg sync.WaitGroup

	closed atomicInt32
}

//new a logger with specified handler and flag
func New(handler Handler, flag int) *Logger {
	var l = new(Logger)

	l.level.Set(LevelInfo)
	l.handler = handler

	l.flag = flag

	l.quit = make(chan struct{})
	l.closed.Set(0)

	l.msg = make(chan []byte, 1024)

	l.bufs = make([][]byte, 0, 16)

	l.wg.Add(1)
	go l.run()

	return l
}

//new a default logger with specified handler and flag: Ltime|Lfile|Llevel
func NewDefault(handler Handler) *Logger {
	return New(handler, Ltime|Lfile|Llevel)
}

func newStdHandler() *StreamHandler {
	h, _ := NewStreamHandler(os.Stdout)
	return h
}

var std = NewDefault(newStdHandler())

func (l *Logger) run() {
	defer l.wg.Done()
	for {
		select {
		case msg := <-l.msg:
			l.hMutex.Lock()
			l.handler.Write(msg)
			l.hMutex.Unlock()
			l.putBuf(msg)
		case <-l.quit:
			//we must log all msg
			if len(l.msg) == 0 {
				return
			}
		}
	}
}

func (l *Logger) popBuf() []byte {
	l.bufMutex.Lock()
	var buf []byte
	if len(l.bufs) == 0 {
		buf = make([]byte, 0, 1024)
	} else {
		buf = l.bufs[len(l.bufs)-1]
		l.bufs = l.bufs[0 : len(l.bufs)-1]
	}
	l.bufMutex.Unlock()

	return buf
}

func (l *Logger) putBuf(buf []byte) {
	l.bufMutex.Lock()
	if len(l.bufs) < maxBufPoolSize {
		buf = buf[0:0]
		l.bufs = append(l.bufs, buf)
	}
	l.bufMutex.Unlock()
}

func (l *Logger) Close() {
	if l.closed.Get() == 1 {
		return
	}
	l.closed.Set(1)

	close(l.quit)

	l.wg.Wait()

	l.quit = nil

	l.handler.Close()
}

//set log level, any log level less than it will not log
func (l *Logger) SetLevel(level int) {
	l.level.Set(level)
}

// name can be in ["trace", "debug", "info", "warn", "error", "fatal"]
func (l *Logger) SetLevelByName(name string) {
	name = strings.ToLower(name)
	switch name {
	case "trace":
		l.SetLevel(LevelTrace)
	case "debug":
		l.SetLevel(LevelDebug)
	case "info":
		l.SetLevel(LevelInfo)
	case "warn":
		l.SetLevel(LevelWarn)
	case "error":
		l.SetLevel(LevelError)
	case "fatal":
		l.SetLevel(LevelFatal)
	}
}

func (l *Logger) SetHandler(h Handler) {
	if l.closed.Get() == 1 {
		return
	}

	l.hMutex.Lock()
	if l.handler != nil {
		l.handler.Close()
	}
	l.handler = h
	l.hMutex.Unlock()
}

func (l *Logger) Output(callDepth int, level int, s string) {
	if l.closed.Get() == 1 {
		// closed
		return
	}

	if l.level.Get() > level {
		// higher level can be logged
		return
	}

	buf := l.popBuf()

	if l.flag&Ltime > 0 {
		now := time.Now().Format(TimeFormat)
		buf = append(buf, '[')
		buf = append(buf, now...)
		buf = append(buf, "] "...)
	}

	if l.flag&Lfile > 0 {
		_, file, line, ok := runtime.Caller(callDepth)
		if !ok {
			file = "???"
			line = 0
		} else {
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					file = file[i+1:]
					break
				}
			}
		}

		buf = append(buf, file...)
		buf = append(buf, ':')

		buf = strconv.AppendInt(buf, int64(line), 10)
		buf = append(buf, ' ')
	}

	if l.flag&Llevel > 0 {
		buf = append(buf, '[')
		buf = append(buf, LevelName[level]...)
		buf = append(buf, "] "...)
	}

	buf = append(buf, s...)

	if s[len(s)-1] != '\n' {
		buf = append(buf, '\n')
	}

	l.msg <- buf
}

//log with Trace level
func (l *Logger) Trace(v ...interface{}) {
	l.Output(2, LevelTrace, fmt.Sprint(v...))
}

//log with Debug level
func (l *Logger) Debug(v ...interface{}) {
	l.Output(2, LevelDebug, fmt.Sprint(v...))
}

//log with info level
func (l *Logger) Info(v ...interface{}) {
	l.Output(2, LevelInfo, fmt.Sprint(v...))
}

//log with warn level
func (l *Logger) Warn(v ...interface{}) {
	l.Output(2, LevelWarn, fmt.Sprint(v...))
}

//log with error level
func (l *Logger) Error(v ...interface{}) {
	l.Output(2, LevelError, fmt.Sprint(v...))
}

//log with fatal level
func (l *Logger) Fatal(v ...interface{}) {
	l.Output(2, LevelFatal, fmt.Sprint(v...))
}

//log with Trace level
func (l *Logger) Tracef(format string, v ...interface{}) {
	l.Output(2, LevelTrace, fmt.Sprintf(format, v...))
}

//log with Debug level
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Output(2, LevelDebug, fmt.Sprintf(format, v...))
}

//log with info level
func (l *Logger) Infof(format string, v ...interface{}) {
	l.Output(2, LevelInfo, fmt.Sprintf(format, v...))
}

//log with warn level
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.Output(2, LevelWarn, fmt.Sprintf(format, v...))
}

//log with error level
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Output(2, LevelError, fmt.Sprintf(format, v...))
}

//log with fatal level
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Output(2, LevelFatal, fmt.Sprintf(format, v...))
}

func SetLevel(level int) {
	std.SetLevel(level)
}

// name can be in ["trace", "debug", "info", "warn", "error", "fatal"]
func SetLevelByName(name string) {
	std.SetLevelByName(name)
}

func SetHandler(h Handler) {
	std.SetHandler(h)
}

func Trace(v ...interface{}) {
	std.Output(2, LevelTrace, fmt.Sprint(v...))
}

func Debug(v ...interface{}) {
	std.Output(2, LevelDebug, fmt.Sprint(v...))
}

func Info(v ...interface{}) {
	std.Output(2, LevelInfo, fmt.Sprint(v...))
}

func Warn(v ...interface{}) {
	std.Output(2, LevelWarn, fmt.Sprint(v...))
}

func Error(v ...interface{}) {
	std.Output(2, LevelError, fmt.Sprint(v...))
}

func Fatal(v ...interface{}) {
	std.Output(2, LevelFatal, fmt.Sprint(v...))
}

func Tracef(format string, v ...interface{}) {
	std.Output(2, LevelTrace, fmt.Sprintf(format, v...))
}

func Debugf(format string, v ...interface{}) {
	std.Output(2, LevelDebug, fmt.Sprintf(format, v...))
}

func Infof(format string, v ...interface{}) {
	std.Output(2, LevelInfo, fmt.Sprintf(format, v...))
}

func Warnf(format string, v ...interface{}) {
	std.Output(2, LevelWarn, fmt.Sprintf(format, v...))
}

func Errorf(format string, v ...interface{}) {
	std.Output(2, LevelError, fmt.Sprintf(format, v...))
}

func Fatalf(format string, v ...interface{}) {
	std.Output(2, LevelFatal, fmt.Sprintf(format, v...))
}
