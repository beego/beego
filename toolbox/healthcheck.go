// Beego (http://beego.me/)

// @description beego is an open-source, high-performance web framework for the Go programming language.

// @link        http://github.com/astaxie/beego for the canonical source repository

// @license     http://github.com/astaxie/beego/blob/master/LICENSE

// @authors     astaxie

package toolbox

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

// health checker map
var AdminCheckList map[string]HealthChecker

// health checker interface
type HealthChecker interface {
	Check() error
}

// add health checker with name string
func AddHealthCheck(name string, hc HealthChecker) {
	AdminCheckList[name] = hc
}

func init() {
	AdminCheckList = make(map[string]HealthChecker)
}
