package admin

import (
	"io"
	"strconv"
	"sync"
	"time"
)

type Statistics struct {
	RequestUrl        string
	RequestController string
	RequestNum        int64
	MinTime           time.Duration
	MaxTime           time.Duration
	TotalTime         time.Duration
}

type UrlMap struct {
	lock   sync.RWMutex
	urlmap map[string]map[string]*Statistics
}

func (m *UrlMap) AddStatistics(requestMethod, requestUrl, requestController string, requesttime time.Duration) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if method, ok := m.urlmap[requestUrl]; ok {
		if s, ok := method[requestMethod]; ok {
			s.RequestNum += 1
			if s.MaxTime < requesttime {
				s.MaxTime = requesttime
			}
			if s.MinTime > requesttime {
				s.MinTime = requesttime
			}
			s.TotalTime += requesttime
		} else {
			nb := &Statistics{
				RequestUrl:        requestUrl,
				RequestController: requestController,
				RequestNum:        1,
				MinTime:           requesttime,
				MaxTime:           requesttime,
				TotalTime:         requesttime,
			}
			m.urlmap[requestUrl][requestMethod] = nb
		}

	} else {
		methodmap := make(map[string]*Statistics)
		nb := &Statistics{
			RequestUrl:        requestUrl,
			RequestController: requestController,
			RequestNum:        1,
			MinTime:           requesttime,
			MaxTime:           requesttime,
			TotalTime:         requesttime,
		}
		methodmap[requestMethod] = nb
		m.urlmap[requestUrl] = methodmap
	}
}

func (m *UrlMap) GetMap(rw io.Writer) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	rw.Write([]byte("requestURL                    avgTime"))
	for k, v := range m.urlmap {
		rw.Write([]byte(k + ""))
		for kk, vv := range v {
			rw.Write([]byte(kk))
			rw.Write([]byte(strconv.FormatInt(vv.RequestNum, 10)))
			rw.Write([]byte(strconv.FormatInt(int64(vv.TotalTime), 10)))
			rw.Write([]byte(strconv.FormatInt(int64(vv.MaxTime), 10)))
			rw.Write([]byte(strconv.FormatInt(int64(vv.MinTime), 10)))
			rw.Write([]byte(strconv.FormatInt(int64(vv.TotalTime)/vv.RequestNum, 10)))
		}
	}
}

var StatisticsMap *UrlMap

func init() {
	StatisticsMap = &UrlMap{
		urlmap: make(map[string]map[string]*Statistics),
	}
}
