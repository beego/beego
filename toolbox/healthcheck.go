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
