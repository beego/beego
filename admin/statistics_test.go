package admin

import (
	"testing"
	"time"
)

func TestStatics(t *testing.T) {
	StatisticsMap.AddStatistics("POST", "/api/user", "&admin.user", time.Duration(1000000))
	StatisticsMap.AddStatistics("POST", "/api/user", "&admin.user", time.Duration(1200000))
	StatisticsMap.AddStatistics("GET", "/api/user", "&admin.user", time.Duration(1300000))
	StatisticsMap.AddStatistics("POST", "/api/admin", "&admin.user", time.Duration(1400000))
	StatisticsMap.AddStatistics("POST", "/api/user/astaxie", "&admin.user", time.Duration(1200000))
	StatisticsMap.AddStatistics("POST", "/api/user/xiemengjun", "&admin.user", time.Duration(1300000))
	StatisticsMap.AddStatistics("DELETE", "/api/user", "&admin.user", time.Duration(1400000))
	s := StatisticsMap.GetMap()
}
