// Package ledis provide session Provider
package ledis

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/ledisdb/ledisdb/config"
	"github.com/ledisdb/ledisdb/ledis"

	"github.com/beego/beego/v2/server/web/session"
)

var (
	ledispder = &Provider{}
	c         *ledis.DB
)

// SessionStore ledis session store
type SessionStore struct {
	sid         string
	lock        sync.RWMutex
	values      map[interface{}]interface{}
	maxlifetime int64
}

// Set value in ledis session
func (ls *SessionStore) Set(ctx context.Context, key, value interface{}) error {
	ls.lock.Lock()
	defer ls.lock.Unlock()
	ls.values[key] = value
	return nil
}

// Get value in ledis session
func (ls *SessionStore) Get(ctx context.Context, key interface{}) interface{} {
	ls.lock.RLock()
	defer ls.lock.RUnlock()
	if v, ok := ls.values[key]; ok {
		return v
	}
	return nil
}

// Delete value in ledis session
func (ls *SessionStore) Delete(ctx context.Context, key interface{}) error {
	ls.lock.Lock()
	defer ls.lock.Unlock()
	delete(ls.values, key)
	return nil
}

// Flush clear all values in ledis session
func (ls *SessionStore) Flush(context.Context) error {
	ls.lock.Lock()
	defer ls.lock.Unlock()
	ls.values = make(map[interface{}]interface{})
	return nil
}

// SessionID get ledis session id
func (ls *SessionStore) SessionID(context.Context) string {
	return ls.sid
}

// SessionRelease save session values to ledis
func (ls *SessionStore) SessionRelease(ctx context.Context, w http.ResponseWriter) {
	ls.lock.RLock()
	values := ls.values
	ls.lock.RUnlock()
	b, err := session.EncodeGob(values)
	if err != nil {
		return
	}
	c.Set([]byte(ls.sid), b)
	c.Expire([]byte(ls.sid), ls.maxlifetime)
}

// SessionReleaseIfPresent save session values to ledis when key is present
// it is not supported now, because ledis has no this feature like SETXX or atomic operation.
// https://github.com/ledisdb/ledisdb/issues/251
// https://github.com/ledisdb/ledisdb/issues/351
func (ls *SessionStore) SessionReleaseIfPresent(ctx context.Context, w http.ResponseWriter) {
	ls.lock.RLock()
	values := ls.values
	ls.lock.RUnlock()
	b, err := session.EncodeGob(values)
	if err != nil {
		return
	}
	r, _ := c.Exists([]byte(ls.sid))
	if r == 1 {
		c.Set([]byte(ls.sid), b)
		c.Expire([]byte(ls.sid), ls.maxlifetime)
	}
}

// Provider ledis session provider
type Provider struct {
	maxlifetime int64
	SavePath    string `json:"save_path"`
	Db          int    `json:"db"`
}

// SessionInit init ledis session
// savepath like ledis server saveDataPath,pool size
// v1.x e.g. 127.0.0.1:6379,100
// v2.x you should pass a json string
// e.g. { "save_path": "my save path", "db": 100}
func (lp *Provider) SessionInit(ctx context.Context, maxlifetime int64, cfgStr string) error {
	var err error
	lp.maxlifetime = maxlifetime
	cfgStr = strings.TrimSpace(cfgStr)
	// we think cfgStr is v2.0, using json to init the session
	if strings.HasPrefix(cfgStr, "{") {
		err = json.Unmarshal([]byte(cfgStr), lp)
	} else {
		err = lp.initOldStyle(cfgStr)
	}

	if err != nil {
		return err
	}

	cfg := new(config.Config)
	cfg.DataDir = lp.SavePath

	var ledisInstance *ledis.Ledis
	ledisInstance, err = ledis.Open(cfg)
	if err != nil {
		return err
	}
	c, err = ledisInstance.Select(lp.Db)
	return err
}

func (lp *Provider) initOldStyle(cfgStr string) error {
	var err error
	configs := strings.Split(cfgStr, ",")
	if len(configs) == 1 {
		lp.SavePath = configs[0]
	} else if len(configs) == 2 {
		lp.SavePath = configs[0]
		lp.Db, err = strconv.Atoi(configs[1])
	}
	return err
}

// SessionRead read ledis session by sid
func (lp *Provider) SessionRead(ctx context.Context, sid string) (session.Store, error) {
	var (
		kv  map[interface{}]interface{}
		err error
	)

	kvs, _ := c.Get([]byte(sid))

	if len(kvs) == 0 {
		kv = make(map[interface{}]interface{})
	} else {
		if kv, err = session.DecodeGob(kvs); err != nil {
			return nil, err
		}
	}

	ls := &SessionStore{sid: sid, values: kv, maxlifetime: lp.maxlifetime}
	return ls, nil
}

// SessionExist check ledis session exist by sid
func (lp *Provider) SessionExist(ctx context.Context, sid string) (bool, error) {
	count, _ := c.Exists([]byte(sid))
	return count != 0, nil
}

// SessionRegenerate generate new sid for ledis session
func (lp *Provider) SessionRegenerate(ctx context.Context, oldsid, sid string) (session.Store, error) {
	count, _ := c.Exists([]byte(sid))
	if count == 0 {
		// oldsid doesn't exist, set the new sid directly
		// ignore error here, since if it returns error
		// the existed value will be 0
		c.Set([]byte(sid), []byte(""))
		c.Expire([]byte(sid), lp.maxlifetime)
	} else {
		data, _ := c.Get([]byte(oldsid))
		c.Set([]byte(sid), data)
		c.Expire([]byte(sid), lp.maxlifetime)
	}
	return lp.SessionRead(context.Background(), sid)
}

// SessionDestroy delete ledis session by id
func (lp *Provider) SessionDestroy(ctx context.Context, sid string) error {
	c.Del([]byte(sid))
	return nil
}

// SessionGC Implement method, no used.
func (lp *Provider) SessionGC(context.Context) {
}

// SessionAll return all active session
func (lp *Provider) SessionAll(context.Context) int {
	return 0
}

func init() {
	session.Register("ledis", ledispder)
}
