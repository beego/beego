package orm

import (
	"sync"
)

const (
	od_CASCADE           = "cascade"
	od_SET_NULL          = "set_null"
	od_SET_DEFAULT       = "set_default"
	od_DO_NOTHING        = "do_nothing"
	defaultStructTagName = "orm"
)

var (
	modelCache = &_modelCache{
		cache:     make(map[string]*modelInfo),
		cacheByFN: make(map[string]*modelInfo),
	}
	supportTag = map[string]int{
		"null":         1,
		"blank":        1,
		"index":        1,
		"unique":       1,
		"pk":           1,
		"auto":         1,
		"auto_now":     1,
		"auto_now_add": 1,
		"size":         2,
		"choices":      2,
		"column":       2,
		"default":      2,
		"rel":          2,
		"reverse":      2,
		"rel_table":    2,
		"rel_through":  2,
		"digits":       2,
		"decimals":     2,
		"on_delete":    2,
		"type":         2,
	}
)

type _modelCache struct {
	sync.RWMutex
	orders    []string
	cache     map[string]*modelInfo
	cacheByFN map[string]*modelInfo
	done      bool
}

func (mc *_modelCache) all() map[string]*modelInfo {
	m := make(map[string]*modelInfo, len(mc.cache))
	for k, v := range mc.cache {
		m[k] = v
	}
	return m
}

func (mc *_modelCache) allOrdered() []*modelInfo {
	m := make([]*modelInfo, 0, len(mc.orders))
	for _, v := range mc.cache {
		m = append(m, v)
	}
	return m
}

func (mc *_modelCache) get(table string) (mi *modelInfo, ok bool) {
	mi, ok = mc.cache[table]
	if ok == false {
		mi, ok = mc.cacheByFN[table]
	}
	return
}

func (mc *_modelCache) set(table string, mi *modelInfo) *modelInfo {
	mii := mc.cache[table]
	mc.cache[table] = mi
	mc.cacheByFN[mi.fullName] = mi
	if mii == nil {
		mc.orders = append(mc.orders, table)
	}
	return mii
}
