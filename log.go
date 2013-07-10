package beego

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type FileLogWriter struct {
	*log.Logger
	mw *MuxWriter
	// The opened file
	filename string

	maxlines          int
	maxlines_curlines int

	// Rotate at size
	maxsize         int
	maxsize_cursize int

	// Rotate daily
	daily          bool
	maxdays        int64
	daily_opendate int

	rotate bool

	startLock sync.Mutex //only one log can writer to the file
}

type MuxWriter struct {
	sync.Mutex
	fd *os.File
}

func (l *MuxWriter) Write(b []byte) (int, error) {
	l.Lock()
	defer l.Unlock()
	return l.fd.Write(b)
}

func (l *MuxWriter) SetFd(fd *os.File) {
	if l.fd != nil {
		l.fd.Close()
	}
	l.fd = fd
}

func NewFileWriter(fname string, rotate bool) *FileLogWriter {
	w := &FileLogWriter{
		filename: fname,
		maxlines: 1000000,
		maxsize:  1 << 28, //256 MB
		daily:    true,
		maxdays:  7,
		rotate:   rotate,
	}
	// use MuxWriter instead direct use os.File for lock write when rotate
	w.mw = new(MuxWriter)
	// set MuxWriter as Logger's io.Writer
	w.Logger = log.New(w.mw, "", log.Ldate|log.Ltime)
	return w
}

// Set rotate at linecount (chainable). Must be called before call StartLogger
func (w *FileLogWriter) SetRotateLines(maxlines int) *FileLogWriter {
	w.maxlines = maxlines
	return w
}

// Set rotate at size (chainable). Must be called before call StartLogger
func (w *FileLogWriter) SetRotateSize(maxsize int) *FileLogWriter {
	w.maxsize = maxsize
	return w
}

// Set rotate daily (chainable). Must be called before call StartLogger
func (w *FileLogWriter) SetRotateDaily(daily bool) *FileLogWriter {
	w.daily = daily
	return w
}

// Set rotate daily's log keep for maxdays, other will delete
func (w *FileLogWriter) SetRotateMaxDays(maxdays int64) *FileLogWriter {
	w.maxdays = maxdays
	return w
}

func (w *FileLogWriter) StartLogger() error {
	fd, err := w.createLogFile()
	if err != nil {
		return err
	}
	w.mw.SetFd(fd)
	err = w.initFd()
	if err != nil {
		return err
	}
	BeeLogger = w
	return nil
}

func (w *FileLogWriter) docheck(size int) {
	w.startLock.Lock()
	defer w.startLock.Unlock()
	if (w.maxlines > 0 && w.maxlines_curlines >= w.maxlines) ||
		(w.maxsize > 0 && w.maxsize_cursize >= w.maxsize) ||
		(w.daily && time.Now().Day() != w.daily_opendate) {
		if err := w.DoRotate(); err != nil {
			fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.filename, err)
			return
		}
	}
	w.maxlines_curlines++
	w.maxsize_cursize += size
}

func (w *FileLogWriter) Printf(format string, v ...interface{}) {
	// Perform the write
	str := fmt.Sprintf(format, v...)
	n := 24 + len(str) // 24 stand for the length "2013/06/23 21:00:22 [T] "

	w.docheck(n)
	w.Logger.Printf(format, v...)
}

func (w *FileLogWriter) createLogFile() (*os.File, error) {
	// Open the log file
	fd, err := os.OpenFile(w.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	return fd, err
}

func (w *FileLogWriter) initFd() error {
	fd := w.mw.fd
	finfo, err := fd.Stat()
	if err != nil {
		return fmt.Errorf("get stat err: %s\n", err)
	}
	w.maxsize_cursize = int(finfo.Size())
	w.daily_opendate = time.Now().Day()
	if finfo.Size() > 0 {
		content, err := ioutil.ReadFile(w.filename)
		if err != nil {
			fmt.Println(err)
		}
		w.maxlines_curlines = len(strings.Split(string(content), "\n"))
	} else {
		w.maxlines_curlines = 0
	}
	return nil
}

func (w *FileLogWriter) DoRotate() error {
	_, err := os.Lstat(w.filename)
	if err == nil { // file exists
		// Find the next available number
		num := 1
		fname := ""
		for ; err == nil && num <= 999; num++ {
			fname = w.filename + fmt.Sprintf(".%s.%03d", time.Now().Format("2006-01-02"), num)
			_, err = os.Lstat(fname)
		}
		// return error if the last file checked still existed
		if err == nil {
			return fmt.Errorf("Rotate: Cannot find free log number to rename %s\n", w.filename)
		}

		// block Logger's io.Writer
		w.mw.Lock()
		defer w.mw.Unlock()

		fd := w.mw.fd
		fd.Close()

		// close fd before rename
		// Rename the file to its newfound home
		err = os.Rename(w.filename, fname)
		if err != nil {
			return fmt.Errorf("Rotate: %s\n", err)
		}

		// re-start logger
		err = w.StartLogger()
		if err != nil {
			return fmt.Errorf("Rotate StartLogger: %s\n", err)
		}

		go w.deleteOldLog()
	}

	return nil
}

func (w *FileLogWriter) deleteOldLog() {
	dir := path.Dir(w.filename)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && info.ModTime().Unix() < (time.Now().Unix()-60*60*24*w.maxdays) {
			if strings.HasPrefix(filepath.Base(path), filepath.Base(w.filename)) {
				os.Remove(path)
			}
		}
		return nil
	})
}

// Log levels to control the logging output.
const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarning
	LevelError
	LevelCritical
)

// logLevel controls the global log level used by the logger.
var level = LevelTrace

// LogLevel returns the global log level and can be used in
// own implementations of the logger interface.
func Level() int {
	return level
}

// SetLogLevel sets the global log level used by the simple
// logger.
func SetLevel(l int) {
	level = l
}

type IBeeLogger interface {
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Flags() int
	Output(calldepth int, s string) error
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
	Prefix() string
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	SetFlags(flag int)
	SetPrefix(prefix string)
}

// logger references the used application logger.
var BeeLogger IBeeLogger

func init() {
	BeeLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
}

// SetLogger sets a new logger.
func SetLogger(l *log.Logger) {
	BeeLogger = l
}

// Trace logs a message at trace level.
func Trace(v ...interface{}) {
	if level <= LevelTrace {
		BeeLogger.Printf("[T] %v\n", v)
	}
}

// Debug logs a message at debug level.
func Debug(v ...interface{}) {
	if level <= LevelDebug {
		BeeLogger.Printf("[D] %v\n", v)
	}
}

// Info logs a message at info level.
func Info(v ...interface{}) {
	if level <= LevelInfo {
		BeeLogger.Printf("[I] %v\n", v)
	}
}

// Warning logs a message at warning level.
func Warn(v ...interface{}) {
	if level <= LevelWarning {
		BeeLogger.Printf("[W] %v\n", v)
	}
}

// Error logs a message at error level.
func Error(v ...interface{}) {
	if level <= LevelError {
		BeeLogger.Printf("[E] %v\n", v)
	}
}

// Critical logs a message at critical level.
func Critical(v ...interface{}) {
	if level <= LevelCritical {
		BeeLogger.Printf("[C] %v\n", v)
	}
}
