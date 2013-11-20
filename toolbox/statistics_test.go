package toolbox

import (
	"os"
	"testing"
	"time"
)

func TestStatics(t *testing.T) {
	StatisticsMap.AddStatistics("POST", "/api/user", "&admin.user", time.Duration(2000))
	StatisticsMap.AddStatistics("POST", "/api/user", "&admin.user", time.Duration(120000))
	StatisticsMap.AddStatistics("GET", "/api/user", "&admin.user", time.Duration(13000))
	StatisticsMap.AddStatistics("POST", "/api/admin", "&admin.user", time.Duration(14000))
	StatisticsMap.AddStatistics("POST", "/api/user/astaxie", "&admin.user", time.Duration(12000))
	StatisticsMap.AddStatistics("POST", "/api/user/xiemengjun", "&admin.user", time.Duration(13000))
	StatisticsMap.AddStatistics("DELETE", "/api/user", "&admin.user", time.Duration(1400))
	StatisticsMap.GetMap(os.Stdout)
}
