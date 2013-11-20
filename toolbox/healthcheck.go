package toolbox

import (
	"fmt"
	"net/http"
)

//type DatabaseCheck struct {
//}

//func (dc *DatabaseCheck) Check() error {
//	if dc.isConnected() {
//		return nil
//	} else {
//		return errors.New("can't connect database")
//	}
//}

//AddHealthCheck("database",&DatabaseCheck{})

var AdminCheckList map[string]HealthChecker

type HealthChecker interface {
	Check() error
}

func AddHealthCheck(name string, hc HealthChecker) {
	AdminCheckList[name] = hc
}

func Healthcheck(rw http.ResponseWriter, req *http.Request) {
	for name, h := range AdminCheckList {
		if err := h.Check(); err != nil {
			fmt.Fprintf(rw, "%s : ok\n", name)
		} else {
			fmt.Fprintf(rw, "%s : %s\n", name, err.Error())
		}
	}
}

func init() {
	AdminCheckList = make(map[string]HealthChecker)
}
