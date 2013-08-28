package logs

import (
	"bufio"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestFile(t *testing.T) {
	log := NewLogger(10000)
	log.SetLogger("file", `{"filename":"test.log"}`)
	log.Trace("test")
	log.Info("info")
	log.Debug("debug")
	log.Warn("warning")
	log.Error("error")
	log.Critical("critical")
	time.Sleep(time.Second * 4)
	f, err := os.Open("test.log")
	if err != nil {
		t.Fatal(err)
	}
	b := bufio.NewReader(f)
	linenum := 0
	for {
		line, _, err := b.ReadLine()
		if err != nil {
			break
		}
		if len(line) > 0 {
			linenum++
		}
	}
	if linenum != 6 {
		t.Fatal(linenum, "not line 6")
	}
	os.Remove("test.log")
}

func TestFile2(t *testing.T) {
	log := NewLogger(10000)
	log.SetLogger("file", `{"filename":"test2.log","level":2}`)
	log.Trace("test")
	log.Info("info")
	log.Debug("debug")
	log.Warn("warning")
	log.Error("error")
	log.Critical("critical")
	time.Sleep(time.Second * 4)
	f, err := os.Open("test2.log")
	if err != nil {
		t.Fatal(err)
	}
	b := bufio.NewReader(f)
	linenum := 0
	for {
		line, _, err := b.ReadLine()
		if err != nil {
			break
		}
		if len(line) > 0 {
			linenum++
		}
	}
	if linenum != 4 {
		t.Fatal(linenum, "not line 4")
	}
	os.Remove("test2.log")
}

func TestFileRotate(t *testing.T) {
	log := NewLogger(10000)
	log.SetLogger("file", `{"filename":"test3.log","maxlines":4}`)
	log.Trace("test")
	log.Info("info")
	log.Debug("debug")
	log.Warn("warning")
	log.Error("error")
	log.Critical("critical")
	time.Sleep(time.Second * 4)
	rotatename := "test3.log" + fmt.Sprintf(".%s.%03d", time.Now().Format("2006-01-02"), 1)
	b, err := exists(rotatename)
	if !b || err != nil {
		t.Fatal("rotate not gen")
	}
	os.Remove(rotatename)
	os.Remove("test3.log")
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func BenchmarkFile(b *testing.B) {
	log := NewLogger(100000)
	log.SetLogger("file", `{"filename":"test4.log"}`)
	for i := 0; i < b.N; i++ {
		log.Trace("trace")
	}
	os.Remove("test4.log")
}
