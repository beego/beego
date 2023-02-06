// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package logs provide a general log interface
// Usage:
//
// import "github.com/beego/beego/v2/core/logs"
//
//	log := NewLogger(10000)
//	log.SetLogger("console", "")
//
//	> the first params stand for how many channel
//
// Use it like this:
//
//	log.Trace("trace")
//	log.Info("info")
//	log.Warn("warning")
//	log.Debug("debug")
//	log.Critical("critical")
package logs

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// RFC5424 log message levels.
const (
	LevelEmergency = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
)

// levelLogLogger is defined to implement log.Logger
// the real log level will be LevelEmergency
const levelLoggerImpl = -1

// Name for adapter with beego official support
const (
	AdapterConsole   = "console"
	AdapterFile      = "file"
	AdapterMultiFile = "multifile"
	AdapterMail      = "smtp"
	AdapterConn      = "conn"
	AdapterEs        = "es"
	AdapterJianLiao  = "jianliao"
	AdapterSlack     = "slack"
	AdapterAliLS     = "alils"
)

// Legacy log level constants to ensure backwards compatibility.
const (
	LevelInfo  = LevelInformational
	LevelTrace = LevelDebug
	LevelWarn  = LevelWarning
)

type newLoggerFunc func() Logger

// Logger defines the behavior of a log provider.
type Logger interface {
	Init(config string) error
	WriteMsg(lm *LogMsg) error
	Destroy()
	Flush()
	SetFormatter(f LogFormatter)
}

var (
	adapters    = make(map[string]newLoggerFunc)
	levelPrefix = [LevelDebug + 1]string{"[M]", "[A]", "[C]", "[E]", "[W]", "[N]", "[I]", "[D]"}
)

// Register makes a log provide available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, log newLoggerFunc) {
	if log == nil {
		panic("logs: Register provide is nil")
	}
	if _, dup := adapters[name]; dup {
		panic("logs: Register called twice for provider " + name)
	}
	adapters[name] = log
}

// BeeLogger is default logger in beego application.
// Can contain several providers and log message into all providers.
type BeeLogger struct {
	lock                sync.Mutex
	init                bool
	enableFuncCallDepth bool
	enableFullFilePath  bool
	asynchronous        bool
	// Whether to discard logs when buffer is full and asynchronous is true
	// No discard by default
	logWithNonBlocking  bool
	wg                  sync.WaitGroup
	level               int
	loggerFuncCallDepth int
	prefix              string
	msgChanLen          int64
	msgChan             chan *LogMsg
	closeChan           chan struct{}
	flushChan           chan struct{}
	outputs             []*nameLogger
	globalFormatter     string
}

const defaultAsyncMsgLen = 1e3

type nameLogger struct {
	Logger
	name string
}

var logMsgPool *sync.Pool

// NewLogger returns a new BeeLogger.
// channelLen: the number of messages in chan(used where asynchronous is true).
// if the buffering chan is full, logger adapters write to file or other way.
func NewLogger(channelLens ...int64) *BeeLogger {
	bl := new(BeeLogger)
	bl.level = LevelDebug
	bl.loggerFuncCallDepth = 3
	bl.msgChanLen = append(channelLens, 0)[0]
	if bl.msgChanLen <= 0 {
		bl.msgChanLen = defaultAsyncMsgLen
	}
	bl.flushChan = make(chan struct{}, 1)
	bl.closeChan = make(chan struct{}, 1)
	bl.setLogger(AdapterConsole)
	return bl
}

// Async sets the log to asynchronous and start the goroutine
func (bl *BeeLogger) Async(msgLen ...int64) *BeeLogger {
	bl.lock.Lock()
	defer bl.lock.Unlock()
	if bl.asynchronous {
		return bl
	}
	bl.asynchronous = true
	if len(msgLen) > 0 && msgLen[0] > 0 {
		bl.msgChanLen = msgLen[0]
	}
	bl.msgChan = make(chan *LogMsg, bl.msgChanLen)
	logMsgPool = &sync.Pool{
		New: func() interface{} {
			return &LogMsg{}
		},
	}
	bl.wg.Add(1)
	go bl.startLogger()
	return bl
}

// AsyncNonBlockWrite Non-blocking write in asynchronous mode
// Only works if asynchronous write logging is set
func (bl *BeeLogger) AsyncNonBlockWrite() *BeeLogger {
	if !bl.asynchronous {
		return bl
	}
	bl.logWithNonBlocking = true
	return bl
}

// SetLogger provides a given logger adapter into BeeLogger with config string.
// config must in in JSON format like {"interval":360}}
func (bl *BeeLogger) setLogger(adapterName string, configs ...string) error {
	config := append(configs, "{}")[0]
	for _, l := range bl.outputs {
		if l.name == adapterName {
			return fmt.Errorf("logs: duplicate adaptername %q (you have set this logger before)", adapterName)
		}
	}

	logAdapter, ok := adapters[adapterName]
	if !ok {
		return fmt.Errorf("logs: unknown adaptername %q (forgotten Register?)", adapterName)
	}

	lg := logAdapter()

	err := lg.Init(config)
	if err != nil {
		return err
	}

	// Global formatter overrides the default set formatter
	if len(bl.globalFormatter) > 0 {
		fmtr, ok := GetFormatter(bl.globalFormatter)
		if !ok {
			return fmt.Errorf("the formatter with name: %s not found", bl.globalFormatter)
		}
		lg.SetFormatter(fmtr)
	}

	bl.outputs = append(bl.outputs, &nameLogger{name: adapterName, Logger: lg})
	return nil
}

// SetLogger provides a given logger adapter into BeeLogger with config string.
// config must in in JSON format like {"interval":360}}
func (bl *BeeLogger) SetLogger(adapterName string, configs ...string) error {
	bl.lock.Lock()
	defer bl.lock.Unlock()
	if !bl.init {
		bl.outputs = []*nameLogger{}
		bl.init = true
	}
	return bl.setLogger(adapterName, configs...)
}

// DelLogger removes a logger adapter in BeeLogger.
func (bl *BeeLogger) DelLogger(adapterName string) error {
	bl.lock.Lock()
	defer bl.lock.Unlock()
	outputs := make([]*nameLogger, 0, len(bl.outputs))
	for _, lg := range bl.outputs {
		if lg.name == adapterName {
			lg.Destroy()
		} else {
			outputs = append(outputs, lg)
		}
	}
	if len(outputs) == len(bl.outputs) {
		return fmt.Errorf("logs: unknown adaptername %q (forgotten Register?)", adapterName)
	}
	bl.outputs = outputs
	return nil
}

func (bl *BeeLogger) writeToLoggers(lm *LogMsg) {
	for _, l := range bl.outputs {
		err := l.WriteMsg(lm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to WriteMsg to adapter:%v,error:%v\n", l.name, err)
		}
	}
}

func (bl *BeeLogger) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	// writeMsg will always add a '\n' character
	if p[len(p)-1] == '\n' {
		p = p[0 : len(p)-1]
	}
	lm := &LogMsg{
		Msg:   string(p),
		Level: levelLoggerImpl,
		When:  time.Now(),
	}

	// set levelLoggerImpl to ensure all log message will be write out
	err = bl.writeMsg(lm)
	if err == nil {
		return len(p), nil
	}
	return 0, err
}

func (bl *BeeLogger) writeMsg(lm *LogMsg) error {
	if !bl.init {
		bl.lock.Lock()
		bl.setLogger(AdapterConsole)
		bl.lock.Unlock()
	}

	var (
		file string
		line int
		ok   bool
	)

	_, file, line, ok = runtime.Caller(bl.loggerFuncCallDepth)
	if !ok {
		file = "???"
		line = 0
	}
	lm.FilePath = file
	lm.LineNumber = line
	lm.Prefix = bl.prefix

	lm.enableFullFilePath = bl.enableFullFilePath
	lm.enableFuncCallDepth = bl.enableFuncCallDepth

	// set level info in front of filename info
	if lm.Level == levelLoggerImpl {
		// set to emergency to ensure all log will be print out correctly
		lm.Level = LevelEmergency
	}

	if bl.asynchronous {
		logM := logMsgPool.Get().(*LogMsg)
		logM.Level = lm.Level
		logM.Msg = lm.Msg
		logM.When = lm.When
		logM.Args = lm.Args
		logM.FilePath = lm.FilePath
		logM.LineNumber = lm.LineNumber
		logM.Prefix = lm.Prefix

		if bl.outputs != nil {
			if bl.logWithNonBlocking {
				select {
				case bl.msgChan <- lm:
				// discard log when channel is full
				default:
				}
			} else {
				bl.msgChan <- lm
			}
		} else {
			logMsgPool.Put(lm)
		}
	} else {
		bl.writeToLoggers(lm)
	}
	return nil
}

// SetLevel sets log message level.
// If message level (such as LevelDebug) is higher than logger level (such as LevelWarning),
// log providers will not be sent the message.
func (bl *BeeLogger) SetLevel(l int) {
	bl.level = l
}

// GetLevel Get Current log message level.
func (bl *BeeLogger) GetLevel() int {
	return bl.level
}

// SetLogFuncCallDepth set log funcCallDepth
func (bl *BeeLogger) SetLogFuncCallDepth(d int) {
	bl.loggerFuncCallDepth = d
}

// GetLogFuncCallDepth return log funcCallDepth for wrapper
func (bl *BeeLogger) GetLogFuncCallDepth() int {
	return bl.loggerFuncCallDepth
}

// EnableFuncCallDepth enable log funcCallDepth
func (bl *BeeLogger) EnableFuncCallDepth(b bool) {
	bl.enableFuncCallDepth = b
}

// set prefix
func (bl *BeeLogger) SetPrefix(s string) {
	bl.prefix = s
}

// start logger chan reading.
// when chan is not empty, write logs.
func (bl *BeeLogger) startLogger() {
	gameOver := false
	for {
		select {
		case bm, ok := <-bl.msgChan:
			// this is a terrible design to have a signal channel that accept two inputs
			// so we only handle the msg if the channel is not closed
			if ok {
				bl.writeToLoggers(bm)
				logMsgPool.Put(bm)
			}
		case <-bl.closeChan:
			bl.flush()
			for _, l := range bl.outputs {
				l.Destroy()
			}
			bl.outputs = nil
			gameOver = true
			bl.wg.Done()
		case <-bl.flushChan:
			bl.flush()
			bl.wg.Done()
		}
		if gameOver {
			break
		}
	}
}

func (bl *BeeLogger) setGlobalFormatter(fmtter string) error {
	bl.globalFormatter = fmtter
	return nil
}

// SetGlobalFormatter sets the global formatter for all log adapters
// don't forget to register the formatter by invoking RegisterFormatter
func SetGlobalFormatter(fmtter string) error {
	return beeLogger.setGlobalFormatter(fmtter)
}

// Emergency Log EMERGENCY level message.
func (bl *BeeLogger) Emergency(format string, v ...interface{}) {
	if LevelEmergency > bl.level {
		return
	}

	lm := &LogMsg{
		Level: LevelEmergency,
		Msg:   format,
		When:  time.Now(),
	}
	if len(v) > 0 {
		lm.Msg = fmt.Sprintf(lm.Msg, v...)
	}

	bl.writeMsg(lm)
}

// Alert Log ALERT level message.
func (bl *BeeLogger) Alert(format string, v ...interface{}) {
	if LevelAlert > bl.level {
		return
	}

	lm := &LogMsg{
		Level: LevelAlert,
		Msg:   format,
		When:  time.Now(),
		Args:  v,
	}
	bl.writeMsg(lm)
}

// Critical Log CRITICAL level message.
func (bl *BeeLogger) Critical(format string, v ...interface{}) {
	if LevelCritical > bl.level {
		return
	}
	lm := &LogMsg{
		Level: LevelCritical,
		Msg:   format,
		When:  time.Now(),
		Args:  v,
	}

	bl.writeMsg(lm)
}

// Error Log ERROR level message.
func (bl *BeeLogger) Error(format string, v ...interface{}) {
	if LevelError > bl.level {
		return
	}
	lm := &LogMsg{
		Level: LevelError,
		Msg:   format,
		When:  time.Now(),
		Args:  v,
	}

	bl.writeMsg(lm)
}

// Warning Log WARNING level message.
func (bl *BeeLogger) Warning(format string, v ...interface{}) {
	if LevelWarn > bl.level {
		return
	}
	lm := &LogMsg{
		Level: LevelWarn,
		Msg:   format,
		When:  time.Now(),
		Args:  v,
	}

	bl.writeMsg(lm)
}

// Notice Log NOTICE level message.
func (bl *BeeLogger) Notice(format string, v ...interface{}) {
	if LevelNotice > bl.level {
		return
	}
	lm := &LogMsg{
		Level: LevelNotice,
		Msg:   format,
		When:  time.Now(),
		Args:  v,
	}

	bl.writeMsg(lm)
}

// Informational Log INFORMATIONAL level message.
func (bl *BeeLogger) Informational(format string, v ...interface{}) {
	if LevelInfo > bl.level {
		return
	}
	lm := &LogMsg{
		Level: LevelInfo,
		Msg:   format,
		When:  time.Now(),
		Args:  v,
	}

	bl.writeMsg(lm)
}

// Debug Log DEBUG level message.
func (bl *BeeLogger) Debug(format string, v ...interface{}) {
	if LevelDebug > bl.level {
		return
	}
	lm := &LogMsg{
		Level: LevelDebug,
		Msg:   format,
		When:  time.Now(),
		Args:  v,
	}

	bl.writeMsg(lm)
}

// Warn Log WARN level message.
// compatibility alias for Warning()
func (bl *BeeLogger) Warn(format string, v ...interface{}) {
	if LevelWarn > bl.level {
		return
	}
	lm := &LogMsg{
		Level: LevelWarn,
		Msg:   format,
		When:  time.Now(),
		Args:  v,
	}

	bl.writeMsg(lm)
}

// Info Log INFO level message.
// compatibility alias for Informational()
func (bl *BeeLogger) Info(format string, v ...interface{}) {
	if LevelInfo > bl.level {
		return
	}
	lm := &LogMsg{
		Level: LevelInfo,
		Msg:   format,
		When:  time.Now(),
		Args:  v,
	}

	bl.writeMsg(lm)
}

// Trace Log TRACE level message.
// compatibility alias for Debug()
func (bl *BeeLogger) Trace(format string, v ...interface{}) {
	if LevelDebug > bl.level {
		return
	}
	lm := &LogMsg{
		Level: LevelDebug,
		Msg:   format,
		When:  time.Now(),
		Args:  v,
	}

	bl.writeMsg(lm)
}

// Flush flush all chan data.
func (bl *BeeLogger) Flush() {
	if bl.asynchronous {
		bl.flushChan <- struct{}{}
		bl.wg.Wait()
		bl.wg.Add(1)
		return
	}
	bl.flush()
}

// Close close logger, flush all chan data and destroy all adapters in BeeLogger.
func (bl *BeeLogger) Close() {
	if bl.asynchronous {
		bl.closeChan <- struct{}{}
		bl.wg.Wait()
		close(bl.msgChan)
	} else {
		bl.flush()
		for _, l := range bl.outputs {
			l.Destroy()
		}
		bl.outputs = nil
	}
	close(bl.flushChan)
	close(bl.closeChan)
}

// Reset close all outputs, and set bl.outputs to nil
func (bl *BeeLogger) Reset() {
	bl.Flush()
	for _, l := range bl.outputs {
		l.Destroy()
	}
	bl.outputs = nil
}

func (bl *BeeLogger) flush() {
	if bl.asynchronous {
		for {
			if len(bl.msgChan) > 0 {
				bm, ok := <-bl.msgChan
				if !ok {
					continue
				}
				bl.writeToLoggers(bm)
				logMsgPool.Put(bm)
				continue
			}
			break
		}
	}
	for _, l := range bl.outputs {
		l.Flush()
	}
}

// beeLogger references the used application logger.
var beeLogger = NewLogger()

// GetBeeLogger returns the default BeeLogger
func GetBeeLogger() *BeeLogger {
	return beeLogger
}

var beeLoggerMap = struct {
	sync.RWMutex
	logs map[string]*log.Logger
}{
	logs: map[string]*log.Logger{},
}

// GetLogger returns the default BeeLogger
func GetLogger(prefixes ...string) *log.Logger {
	prefix := append(prefixes, "")[0]
	if prefix != "" {
		prefix = fmt.Sprintf(`[%s] `, strings.ToUpper(prefix))
	}
	beeLoggerMap.RLock()
	l, ok := beeLoggerMap.logs[prefix]
	if ok {
		beeLoggerMap.RUnlock()
		return l
	}
	beeLoggerMap.RUnlock()
	beeLoggerMap.Lock()
	defer beeLoggerMap.Unlock()
	l, ok = beeLoggerMap.logs[prefix]
	if !ok {
		l = log.New(beeLogger, prefix, 0)
		beeLoggerMap.logs[prefix] = l
	}
	return l
}

// EnableFullFilePath enables full file path logging. Disabled by default
// e.g "/home/Documents/GitHub/beego/mainapp/" instead of "mainapp"
func EnableFullFilePath(b bool) {
	beeLogger.enableFullFilePath = b
}

// Reset will remove all the adapter
func Reset() {
	beeLogger.Reset()
}

// Async set the beelogger with Async mode and hold msglen messages
func Async(msgLen ...int64) *BeeLogger {
	return beeLogger.Async(msgLen...)
}

// SetLevel sets the global log level used by the simple logger.
func SetLevel(l int) {
	beeLogger.SetLevel(l)
}

// SetPrefix sets the prefix
func SetPrefix(s string) {
	beeLogger.SetPrefix(s)
}

// EnableFuncCallDepth enable log funcCallDepth
func EnableFuncCallDepth(b bool) {
	beeLogger.enableFuncCallDepth = b
}

// SetLogFuncCall set the CallDepth, default is 4
func SetLogFuncCall(b bool) {
	beeLogger.EnableFuncCallDepth(b)
	beeLogger.SetLogFuncCallDepth(3)
}

// SetLogFuncCallDepth set log funcCallDepth
func SetLogFuncCallDepth(d int) {
	beeLogger.loggerFuncCallDepth = d
}

// SetLogger sets a new logger.
func SetLogger(adapter string, config ...string) error {
	return beeLogger.SetLogger(adapter, config...)
}

// Emergency logs a message at emergency level.
func Emergency(f interface{}, v ...interface{}) {
	beeLogger.Emergency(formatPattern(f, v...), v...)
}

// Alert logs a message at alert level.
func Alert(f interface{}, v ...interface{}) {
	beeLogger.Alert(formatPattern(f, v...), v...)
}

// Critical logs a message at critical level.
func Critical(f interface{}, v ...interface{}) {
	beeLogger.Critical(formatPattern(f, v...), v...)
}

// Error logs a message at error level.
func Error(f interface{}, v ...interface{}) {
	beeLogger.Error(formatPattern(f, v...), v...)
}

// Warning logs a message at warning level.
func Warning(f interface{}, v ...interface{}) {
	beeLogger.Warn(formatPattern(f, v...), v...)
}

// Warn compatibility alias for Warning()
func Warn(f interface{}, v ...interface{}) {
	beeLogger.Warn(formatPattern(f, v...), v...)
}

// Notice logs a message at notice level.
func Notice(f interface{}, v ...interface{}) {
	beeLogger.Notice(formatPattern(f, v...), v...)
}

// Informational logs a message at info level.
func Informational(f interface{}, v ...interface{}) {
	beeLogger.Info(formatPattern(f, v...), v...)
}

// Info compatibility alias for Warning()
func Info(f interface{}, v ...interface{}) {
	beeLogger.Info(formatPattern(f, v...), v...)
}

// Debug logs a message at debug level.
func Debug(f interface{}, v ...interface{}) {
	beeLogger.Debug(formatPattern(f, v...), v...)
}

// Trace logs a message at trace level.
// compatibility alias for Warning()
func Trace(f interface{}, v ...interface{}) {
	beeLogger.Trace(formatPattern(f, v...), v...)
}

func formatPattern(f interface{}, v ...interface{}) string {
	var msg string
	switch f.(type) {
	case string:
		msg = f.(string)
		if len(v) == 0 {
			return msg
		}
		if !strings.Contains(msg, "%") {
			// do not contain format char
			msg += strings.Repeat(" %v", len(v))
		}
	default:
		msg = fmt.Sprint(f)
		if len(v) == 0 {
			return msg
		}
		msg += strings.Repeat(" %v", len(v))
	}
	return msg
}
