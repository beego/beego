package admin

import (
	"fmt"
	"net/http"
)

var AdminTaskList map[string]Tasker

type Tasker interface {
	GetStatus() string
	Run() error
}

type Task struct {
	Taskname string
	Spec     Schedule
	Errlist  []map[uint64]string //errtime:errinfo
	ErrLimit int                 //max length for the errlist 0 stand for there' no limit
}

func (t *Task) GetStatus() string {
	return ""
}

func (t *Task) Run() error {
	return nil
}

func (t *Task) SetCron(spec string) {

}

type Schedule struct {
	Second uint64
	Minute uint64
	Hour   uint64
	DOM    uint64
	Month  uint64
	DOW    uint64
}

func StartTask() {

}

func StopTask() {

}

func AddTask(taskname string, t Tasker) {
	AdminTaskList[taskname] = t
}

func TaskStatus(rw http.ResponseWriter, req *http.Request) {
	for tname, t := range AdminTaskList {
		fmt.Fprintf(rw, "%s:%s", tname, t.GetStatus())
	}
}

//to run a Task by http from the querystring taskname
//url like /task?taskname=sendmail
func RunTask(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	taskname := req.Form.Get("taskname")
	if t, ok := AdminTaskList[taskname]; ok {
		err := t.Run()
		if err != nil {
			fmt.Fprintf(rw, "%v", err)
		}
		fmt.Fprintf(rw, "%s run success,Now the Status is %s", t.GetStatus())
	} else {
		fmt.Fprintf(rw, "there's no task which named:%s", taskname)
	}
}

func init() {
	AdminTaskList = make(map[string]Tasker)
}
