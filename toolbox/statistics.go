package toolbox

import (
	"fmt"
	"io"
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
	lock        sync.RWMutex
	LengthLimit int //limit the urlmap's length if it's equal to 0 there's no limit
	urlmap      map[string]map[string]*Statistics
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
		if m.LengthLimit > 0 && m.LengthLimit <= len(m.urlmap) {
			return
		}
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
	fmt.Fprintf(rw, "| % -50s| % -10s | % -16s | % -16s | % -16s | % -16s | % -16s |\n", "requestUrl", "method", "times", "used", "max used", "min used", "avg used")
	for k, v := range m.urlmap {
		for kk, vv := range v {
			fmt.Fprintf(rw, "| % -50s| % -10s | % -16d | % -16s | % -16s | % -16s | % -16s |\n", k,
				kk, vv.RequestNum, toS(vv.TotalTime), toS(vv.MaxTime), toS(vv.MinTime), toS(time.Duration(int64(vv.TotalTime)/vv.RequestNum)),
			)
		}
	}
}

var StatisticsMap *UrlMap

func init() {
	StatisticsMap = &UrlMap{
		urlmap: make(map[string]map[string]*Statistics),
	}
}
