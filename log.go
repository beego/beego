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
	// The opened file
	filename string

	maxlines          int
	maxlines_curlines int

	// Rotate at size
	maxsize         int
	maxsize_cursize int

	// Rotate daily
	daily          bool
	maxday         int64
	daily_opendate int

	rotate bool

	startLock sync.Mutex //only one log can writer to the file
}

func NewFileWriter(fname string, rotate bool) *FileLogWriter {
	w := &FileLogWriter{
		filename: fname,
		maxlines: 1000000,
		maxsize:  1 << 28, //256 MB
		daily:    true,
		maxday:   7,
		rotate:   rotate,
	}
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

// Set rotate daily's log keep for maxday,other will delete
func (w *FileLogWriter) SetRotateMaxDay(maxday int64) *FileLogWriter {
	w.maxday = maxday
	return w
}

func (w *FileLogWriter) StartLogger() error {
	if err := w.DoRotate(false); err != nil {
		return err
	}
	return nil
}

func (w *FileLogWriter) docheck(size int) {
	w.startLock.Lock()
	defer w.startLock.Unlock()
	if (w.maxlines > 0 && w.maxlines_curlines >= w.maxlines) ||
		(w.maxsize > 0 && w.maxsize_cursize >= w.maxsize) ||
		(w.daily && time.Now().Day() != w.daily_opendate) {
		if err := w.DoRotate(true); err != nil {
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

func (w *FileLogWriter) DoRotate(rotate bool) error {
	if rotate {
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

			// Rename the file to its newfound home
			err = os.Rename(w.filename, fname)
			if err != nil {
				return fmt.Errorf("Rotate: %s\n", err)
			}
			go w.deleteOldLog()
		}
	}

	// Open the log file
	fd, err := os.OpenFile(w.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	w.Logger = log.New(fd, "", log.Ldate|log.Ltime)
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
	BeeLogger = w
	return nil
}

func (w *FileLogWriter) deleteOldLog() {
	dir := path.Dir(w.filename)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && info.ModTime().Unix() < (time.Now().Unix()-60*60*24*w.maxday) {
			os.Remove(path)
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
