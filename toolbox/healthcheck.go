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

var AdminCheckList map[string]HealthChecker

type HealthChecker interface {
	Check() error
}

func AddHealthCheck(name string, hc HealthChecker) {
	AdminCheckList[name] = hc
}

func init() {
	AdminCheckList = make(map[string]HealthChecker)
}
