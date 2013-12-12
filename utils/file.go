package utils

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

func SelfPath() string {
	path, _ := filepath.Abs(os.Args[0])
	return path
}

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

// search a file in paths.
// this is offen used in search config file in /etc ~/
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
	for {
		byteLine, isPrefix, er := reader.ReadLine()
		if er != nil && er != io.EOF {
			return nil, er
		}
		line := string(byteLine)
		if isPrefix {
			prefix += line
			continue
		}

		line = prefix + line
		if re.MatchString(line) {
			lines = append(lines, line)
		}
		if er == io.EOF {
			break
		}
	}
	return lines, nil
}
