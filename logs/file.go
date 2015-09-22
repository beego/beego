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

package logs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// fileLogWriter implements LoggerInterface.
// It writes messages by lines limit, file size limit, or time frequency.
type fileLogWriter struct {
	*log.Logger
	mw *MuxWriter
	// The opened file
	Filename string `json:"filename"`

	Maxlines         int `json:"maxlines"`
	maxlinesCurlines int

	// Rotate at size
	Maxsize        int `json:"maxsize"`
	maxsizeCursize int

	// Rotate daily
	Daily         bool  `json:"daily"`
	Maxdays       int64 `json:"maxdays"`
	dailyOpendate int

	Rotate bool `json:"rotate"`

	startLock sync.Mutex // Only one log can write to the file

	Level int `json:"level"`

	Perm os.FileMode `json:"perm"`
}

// MuxWriter is an *os.File writer with locker.
type MuxWriter struct {
	sync.Mutex
	fd *os.File
}

// write to os.File.
func (l *MuxWriter) Write(b []byte) (int, error) {
	l.Lock()
	defer l.Unlock()
	return l.fd.Write(b)
}

// SetFd set os.File in writer.
func (l *MuxWriter) SetFd(fd *os.File) {
	if l.fd != nil {
		l.fd.Close()
	}
	l.fd = fd
}

// NewFileWriter create a FileLogWriter returning as LoggerInterface.
func newFileWriter() Logger {
	w := &fileLogWriter{
		Filename: "",
		Maxlines: 1000000,
		Maxsize:  1 << 28, //256 MB
		Daily:    true,
		Maxdays:  7,
		Rotate:   true,
		Level:    LevelTrace,
		Perm:     0660,
	}
	// use MuxWriter instead direct use os.File for lock write when rotate
	w.mw = new(MuxWriter)
	// set MuxWriter as Logger's io.Writer
	w.Logger = log.New(w.mw, "", log.Ldate|log.Ltime)
	return w
}

// Init file logger with json config.
// jsonconfig like:
//	{
//	"filename":"logs/beego.log",
//	"maxlines":10000,
//	"maxsize":1<<30,
//	"daily":true,
//	"maxdays":15,
//	"rotate":true,
//  "perm":0600
//	}
func (w *fileLogWriter) Init(jsonconfig string) error {
	err := json.Unmarshal([]byte(jsonconfig), w)
	if err != nil {
		return err
	}
	if len(w.Filename) == 0 {
		return errors.New("jsonconfig must have filename")
	}
	err = w.startLogger()
	return err
}

// start file logger. create log file and set to locker-inside file writer.
func (w *fileLogWriter) startLogger() error {
	fd, err := w.createLogFile()
	if err != nil {
		return err
	}
	w.mw.SetFd(fd)
	return w.initFd()
}

func (w *fileLogWriter) docheck(size int) {
	w.startLock.Lock()
	defer w.startLock.Unlock()
	if w.Rotate && ((w.Maxlines > 0 && w.maxlinesCurlines >= w.Maxlines) ||
		(w.Maxsize > 0 && w.maxsizeCursize >= w.Maxsize) ||
		(w.Daily && time.Now().Day() != w.dailyOpendate)) {
		if err := w.DoRotate(); err != nil {
			fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.Filename, err)
			return
		}
	}
	w.maxlinesCurlines++
	w.maxsizeCursize += size
}

// WriteMsg write logger message into file.
func (w *fileLogWriter) WriteMsg(msg string, level int) error {
	if level > w.Level {
		return nil
	}
	n := 24 + len(msg) // 24 stand for the length "2013/06/23 21:00:22 [T] "
	w.docheck(n)
	w.Logger.Println(msg)
	return nil
}

func (w *fileLogWriter) createLogFile() (*os.File, error) {
	// Open the log file
	fd, err := os.OpenFile(w.Filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, w.Perm)
	return fd, err
}

func (w *fileLogWriter) initFd() error {
	fd := w.mw.fd
	finfo, err := fd.Stat()
	if err != nil {
		return fmt.Errorf("get stat err: %s\n", err)
	}
	w.maxsizeCursize = int(finfo.Size())
	w.dailyOpendate = time.Now().Day()
	w.maxlinesCurlines = 0
	if finfo.Size() > 0 {
		count, err := w.lines()
		if err != nil {
			return err
		}
		w.maxlinesCurlines = count
	}
	return nil
}

func (w *fileLogWriter) lines() (int, error) {
	fd, err := os.Open(w.Filename)
	if err != nil {
		return 0, err
	}
	defer fd.Close()

	buf := make([]byte, 32768) // 32k
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := fd.Read(buf)
		if err != nil && err != io.EOF {
			return count, err
		}

		count += bytes.Count(buf[:c], lineSep)

		if err == io.EOF {
			break
		}
	}

	return count, nil
}

// DoRotate means it need to write file in new file.
// new file name like xx.2013-01-01.2.log
func (w *fileLogWriter) DoRotate() error {
	_, err := os.Lstat(w.Filename)
	if err == nil { // file exists
		// Find the next available number
		num := 1
		fname := ""
		suffix := filepath.Ext(w.Filename)
		filenameOnly := strings.TrimSuffix(w.Filename, suffix)
		if suffix == "" {
			suffix = ".log"
		}
		for ; err == nil && num <= 999; num++ {
			fname = filenameOnly + fmt.Sprintf(".%s.%03d%s", time.Now().Format("2006-01-02"), num, suffix)
			_, err = os.Lstat(fname)
		}
		// return error if the last file checked still existed
		if err == nil {
			return fmt.Errorf("Rotate: Cannot find free log number to rename %s\n", w.Filename)
		}

		// block Logger's io.Writer
		w.mw.Lock()
		defer w.mw.Unlock()

		fd := w.mw.fd
		fd.Close()

		// close fd before rename
		// Rename the file to its newfound home
		err = os.Rename(w.Filename, fname)
		if err != nil {
			return fmt.Errorf("Rotate: %s\n", err)
		}

		// re-start logger
		err = w.startLogger()
		if err != nil {
			return fmt.Errorf("Rotate StartLogger: %s\n", err)
		}

		go w.deleteOldLog()
	}

	return nil
}

func (w *fileLogWriter) deleteOldLog() {
	dir := filepath.Dir(w.Filename)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) (returnErr error) {
		defer func() {
			if r := recover(); r != nil {
				returnErr = fmt.Errorf("Unable to delete old log '%s', error: %+v", path, r)
				fmt.Println(returnErr)
			}
		}()

		if !info.IsDir() && info.ModTime().Unix() < (time.Now().Unix()-60*60*24*w.Maxdays) {
			if strings.HasPrefix(filepath.Base(path), filepath.Base(w.Filename)) {
				os.Remove(path)
			}
		}
		return
	})
}

// Destroy close the file desciption, close file writer.
func (w *fileLogWriter) Destroy() {
	w.mw.fd.Close()
}

// Flush flush file logger.
// there are no buffering messages in file logger in memory.
// flush file means sync file from disk.
func (w *fileLogWriter) Flush() {
	w.mw.fd.Sync()
}

func init() {
	Register("file", newFileWriter)
}
