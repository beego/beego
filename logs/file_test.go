// Beego (http://beego.me/)
//
// @description beego is an open-source, high-performance web framework for the Go programming language.
//
// @link        http://github.com/astaxie/beego for the canonical source repository
//
// @license     http://github.com/astaxie/beego/blob/master/LICENSE
//
// @authors     astaxie
package logs

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestFile(t *testing.T) {
	log := NewLogger(10000)
	log.SetLogger("file", `{"filename":"test.log"}`)
	log.Debug("debug")
	log.Informational("info")
	log.Notice("notice")
	log.Warning("warning")
	log.Error("error")
	log.Alert("alert")
	log.Critical("critical")
	log.Emergency("emergency")
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
	var expected = LevelDebug + 1
	if linenum != expected {
		t.Fatal(linenum, "not "+strconv.Itoa(expected)+" lines")
	}
	os.Remove("test.log")
}

func TestFile2(t *testing.T) {
	log := NewLogger(10000)
	log.SetLogger("file", fmt.Sprintf(`{"filename":"test2.log","level":%d}`, LevelError))
	log.Debug("debug")
	log.Info("info")
	log.Notice("notice")
	log.Warning("warning")
	log.Error("error")
	log.Alert("alert")
	log.Critical("critical")
	log.Emergency("emergency")
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
	var expected = LevelError + 1
	if linenum != expected {
		t.Fatal(linenum, "not "+strconv.Itoa(expected)+" lines")
	}
	os.Remove("test2.log")
}

func TestFileRotate(t *testing.T) {
	log := NewLogger(10000)
	log.SetLogger("file", `{"filename":"test3.log","maxlines":4}`)
	log.Debug("debug")
	log.Info("info")
	log.Notice("notice")
	log.Warning("warning")
	log.Error("error")
	log.Alert("alert")
	log.Critical("critical")
	log.Emergency("emergency")
	time.Sleep(time.Second * 4)
	rotatename := "test3.log" + fmt.Sprintf(".%s.%03d", time.Now().Format("2006-01-02"), 1)
	b, err := exists(rotatename)
	if !b || err != nil {
		t.Fatal("rotate not generated")
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
		log.Debug("debug")
	}
	os.Remove("test4.log")
}
