// Beego (http://beego.me/)

// @description beego is an open-source, high-performance web framework for the Go programming language.

// @link        http://github.com/astaxie/beego for the canonical source repository

// @license     http://github.com/astaxie/beego/blob/master/LICENSE

// @authors     astaxie

package utils

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

// SelfPath gets compiled executable file absolute path
func SelfPath() string {
	path, _ := filepath.Abs(os.Args[0])
	return path
}

// SelfDir gets compiled executable file directory
func SelfDir() string {
	return filepath.Dir(SelfPath())
}

// FileExists reports whether the named file or directory exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// Search a file in paths.
// this is often used in search config file in /etc ~/
func SearchFile(filename string, paths ...string) (fullpath string, err error) {
	for _, path := range paths {
		if fullpath = filepath.Join(path, filename); FileExists(fullpath) {
			return
		}
	}
	err = errors.New(fullpath + " not found in paths")
	return
}

// like command grep -E
// for example: GrepFile(`^hello`, "hello.txt")
// \n is striped while read
func GrepFile(patten string, filename string) (lines []string, err error) {
	re, err := regexp.Compile(patten)
	if err != nil {
		return
	}

	fd, err := os.Open(filename)
	if err != nil {
		return
	}
	lines = make([]string, 0)
	reader := bufio.NewReader(fd)
	prefix := ""
	isLongLine := false
	for {
		byteLine, isPrefix, er := reader.ReadLine()
		if er != nil && er != io.EOF {
			return nil, er
		}
		if er == io.EOF {
			break
		}
		line := string(byteLine)
		if isPrefix {
			prefix += line
			continue
		} else {
			isLongLine = true
		}

		line = prefix + line
		if isLongLine {
			prefix = ""
		}
		if re.MatchString(line) {
			lines = append(lines, line)
		}
	}
	return lines, nil
}
