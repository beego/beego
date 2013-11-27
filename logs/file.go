package logs

import (
	"encoding/json"
	"errors"
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

	startLock sync.Mutex // Only one log can write to the file

	level int
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

func NewFileWriter() LoggerInterface {
	w := &FileLogWriter{
		filename: "",
		maxlines: 1000000,
		maxsize:  1 << 28, //256 MB
		daily:    true,
		maxdays:  7,
		rotate:   true,
		level:    LevelTrace,
	}
	// use MuxWriter instead direct use os.File for lock write when rotate
	w.mw = new(MuxWriter)
	// set MuxWriter as Logger's io.Writer
	w.Logger = log.New(w.mw, "", log.Ldate|log.Ltime)
	return w
}

// jsonconfig like this
//{
//	"filename":"logs/beego.log",
//	"maxlines":10000,
//	"maxsize":1<<30,
//	"daily":true,
//	"maxdays":15,
//	"rotate":true
//}
func (w *FileLogWriter) Init(jsonconfig string) error {
	var m map[string]interface{}
	err := json.Unmarshal([]byte(jsonconfig), &m)
	if err != nil {
		return err
	}
	if fn, ok := m["filename"]; !ok {
		return errors.New("jsonconfig must have filename")
	} else {
		w.filename = fn.(string)
	}
	if ml, ok := m["maxlines"]; ok {
		w.maxlines = int(ml.(float64))
	}
	if ms, ok := m["maxsize"]; ok {
		w.maxsize = int(ms.(float64))
	}
	if dl, ok := m["daily"]; ok {
		w.daily = dl.(bool)
	}
	if md, ok := m["maxdays"]; ok {
		w.maxdays = int64(md.(float64))
	}
	if rt, ok := m["rotate"]; ok {
		w.rotate = rt.(bool)
	}
	if lv, ok := m["level"]; ok {
		w.level = int(lv.(float64))
	}
	err = w.StartLogger()
	return err
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

func (w *FileLogWriter) WriteMsg(msg string, level int) error {
	if level < w.level {
		return nil
	}
	n := 24 + len(msg) // 24 stand for the length "2013/06/23 21:00:22 [T] "
	w.docheck(n)
	w.Logger.Println(msg)
	return nil
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
			return err
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

func (w *FileLogWriter) Destroy() {
	w.mw.fd.Close()
}

func (w *FileLogWriter) Flush() {
	w.mw.fd.Sync()
}

func init() {
	Register("file", NewFileWriter)
}
