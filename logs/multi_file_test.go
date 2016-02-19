package logs

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestMultiFile1(t *testing.T) {
	log := NewLogger(10000)
	err := log.SetLogger("multi_file", `{}`)
	if err == nil {
		t.Fatal("shoud have err")
	}
}

func TestMultiFile2(t *testing.T) {
	debugFileName := "test_file2.debug.log"
	os.Remove(debugFileName)
	defer os.Remove(debugFileName)

	log := NewLogger(10000)
	err := log.SetLogger("multi_file", fmt.Sprintf(
		`{
			"levelfiles": [{
				"levelnames": ["debug"],
				"filename": "%s"
			}]
		}`, debugFileName))
	if err != nil {
		t.Fatal(err)
	}

	log.Debug("debug msg")
	b, err := hasContent(debugFileName)
	if !b || err != nil {
		t.Fatalf("log file[%s] doesn't exist", debugFileName)
	}
}

func TestMultiFile3(t *testing.T) {
	debugFileName := "test_file3.debug.log"
	infoFileName := "test_file3.info.log"
	os.Remove(debugFileName)
	os.Remove(infoFileName)
	defer os.Remove(debugFileName)
	defer os.Remove(infoFileName)

	log := NewLogger(10000)
	err := log.SetLogger("multi_file", fmt.Sprintf(
		`{
			"levelfiles": [{
				"levelnames": ["debug", "trace"],
				"filename": "%s"
			},{
				"levelnames": ["info"],
				"filename": "%s"
			}]
		}`, debugFileName, infoFileName))
	if err != nil {
		t.Fatal(err)
	}

	log.Debug("debug msg")
	log.Trace("trace msg")
	{
		b, err := hasContent(debugFileName)
		if !b || err != nil {
			t.Fatalf("log file[%s] doesn't exist", debugFileName)
		}
	}
	{
		lineNum, err := getFileLineNum(debugFileName)
		if err != nil {
			t.Fatalf("log file[%s] get lines err: %v", debugFileName, err)
		}
		if lineNum != 2 {
			t.Fatalf("log file[%s] line number[%d] should be 2", debugFileName, lineNum)
		}
	}

	log.Info("info msg")
	{
		b, err := hasContent(infoFileName)
		if !b || err != nil {
			t.Fatalf("log file[%s] doesn't exist", infoFileName)
		}
	}
}

func TestMultiFile4(t *testing.T) {
	debugFileName := "test_file4.debug.log"
	infoFileName := "test_file4.info.log"
	os.Remove(debugFileName)
	os.Remove(infoFileName)
	defer os.Remove(debugFileName)
	defer os.Remove(infoFileName)

	log := NewLogger(10000)
	err := log.SetLogger("multi_file", fmt.Sprintf(
		`{
			"levelname": "info",
			"levelfiles": [{
				"levelnames": ["debug"],
				"filename": "%s"
			},{
				"levelnames": ["info"],
				"filename": "%s"
			}]
		}`, debugFileName, infoFileName))
	if err != nil {
		t.Fatal(err)
	}

	log.Debug("debug msg")
	{
		b, err := hasContent(debugFileName)
		if err != nil {
			t.Fatal(err)
		}
		if b {
			t.Fatalf("log file[%s] should not have content", debugFileName)
		}
	}

	log.Info("info msg")
	{
		b, err := hasContent(infoFileName)
		if err != nil {
			t.Fatal(err)
		}
		if !b {
			t.Fatalf("log file[%s] should have content", infoFileName)
		}
	}
}

func hasContent(path string) (bool, error) {
	contentBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return false, err
	}
	return len(contentBytes) > 0, nil
}

func getFileLineNum(path string) (int, error) {
	contentBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return len(strings.Split(string(contentBytes), "\n")) - 1, nil
}
