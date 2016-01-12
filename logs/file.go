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
	*MuxWriter
	// The opened file
	Filename string `json:"filename"`

	MaxLines         int `json:"maxlines"`
	maxLinesCurLines int

	// Rotate at size
	MaxSize        int `json:"maxsize"`
	maxSizeCurSize int

	// Rotate daily
	Daily         bool  `json:"daily"`
	MaxDays       int64 `json:"maxdays"`
	dailyOpenDate int

	Rotate bool `json:"rotate"`

	startLock sync.Mutex // Only one log can write to the file

	Level int `json:"level"`

	Perm os.FileMode `json:"perm"`
}

// MuxWriter is an *os.File writer with locker,lock write when rotate
type MuxWriter struct {
	sync.Mutex
	fd *os.File
}

// Write to os.File.
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
		Filename:  "",
		MaxLines:  1000000,
		MaxSize:   1 << 28, //256 MB
		Daily:     true,
		MaxDays:   7,
		Rotate:    true,
		Level:     LevelTrace,
		Perm:      0660,
		MuxWriter: new(MuxWriter),
	}
	// set MuxWriter as Logger's io.Writer
	w.Logger = log.New(w, "", log.Ldate|log.Ltime)
	return w
}

// Init file logger with json config.
// jsonConfig like:
//	{
//	"filename":"logs/beego.log",
//	"maxLines":10000,
//	"maxsize":1<<30,
//	"daily":true,
//	"maxDays":15,
//	"rotate":true,
//  	"perm":0600
//	}
func (w *fileLogWriter) Init(jsonConfig string) error {
	err := json.Unmarshal([]byte(jsonConfig), w)
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
	w.SetFd(fd)
	return w.initFd()
}

func (w *fileLogWriter) doCheck(size int) {
	w.startLock.Lock()
	if w.Rotate {
		if (w.MaxLines > 0 && w.maxLinesCurLines >= w.MaxLines) ||
			(w.MaxSize > 0 && w.maxSizeCurSize >= w.MaxSize) ||
			(w.Daily && time.Now().Day() != w.dailyOpenDate) {
			if err := w.DoRotate(); err != nil {
				fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.Filename, err)
				w.startLock.Unlock()
				return
			}
		}
	}
	w.maxLinesCurLines++
	w.maxSizeCurSize += size
	w.startLock.Unlock()
}

// WriteMsg write logger message into file.
func (w *fileLogWriter) WriteMsg(msg string, level int) error {
	if level > w.Level {
		return nil
	}
	n := 24 + len(msg) // 24 stand for the length "2013/06/23 21:00:22 [T] "
	w.doCheck(n)
	w.Logger.Println(msg)
	return nil
}

func (w *fileLogWriter) createLogFile() (*os.File, error) {
	// Open the log file
	fd, err := os.OpenFile(w.Filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, w.Perm)
	return fd, err
}

func (w *fileLogWriter) initFd() error {
	fd := w.fd
	fInfo, err := fd.Stat()
	if err != nil {
		return fmt.Errorf("get stat err: %s\n", err)
	}
	w.maxSizeCurSize = int(fInfo.Size())
	w.dailyOpenDate = time.Now().Day()
	w.maxLinesCurLines = 0
	if fInfo.Size() > 0 {
		count, err := w.lines()
		if err != nil {
			return err
		}
		w.maxLinesCurLines = count
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
	if err == nil {
		// file exists
		// Find the next available number
		num := 1
		fName := ""
		suffix := filepath.Ext(w.Filename)
		filenameOnly := strings.TrimSuffix(w.Filename, suffix)
		if suffix == "" {
			suffix = ".log"
		}
		for ; err == nil && num <= 999; num++ {
			fName = filenameOnly + fmt.Sprintf(".%s.%03d%s", time.Now().Format("2006-01-02"), num, suffix)
			_, err = os.Lstat(fName)
		}
		// return error if the last file checked still existed
		if err == nil {
			return fmt.Errorf("Rotate: Cannot find free log number to rename %s\n", w.Filename)
		}

		// block Logger's io.Writer
		w.Lock()
		defer w.Unlock()

		fd := w.fd
		fd.Close()

		// close fd before rename
		// Rename the file to its new found name
		err = os.Rename(w.Filename, fName)
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
				fmt.Fprintf(os.Stderr, "Unable to delete old log '%s', error: %+v", path, r)
			}
		}()

		if !info.IsDir() && info.ModTime().Unix() < (time.Now().Unix()-60*60*24*w.MaxDays) {
			if strings.HasPrefix(filepath.Base(path), filepath.Base(w.Filename)) {
				os.Remove(path)
			}
		}
		return
	})
}

// Destroy close the file description, close file writer.
func (w *fileLogWriter) Destroy() {
	w.fd.Close()
}

// Flush flush file logger.
// there are no buffering messages in file logger in memory.
// flush file means sync file from disk.
func (w *fileLogWriter) Flush() {
	w.fd.Sync()
}

func init() {
	Register("file", newFileWriter)
}
