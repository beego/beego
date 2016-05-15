package ssdb

import (
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/astaxie/beego/session"
	"github.com/ssdb/gossdb/ssdb"
)

var ssdbProvider = &SsdbProvider{}

type SsdbProvider struct {
	client      *ssdb.Client
	host        string
	port        int32
	maxLifetime int64
}

func (p *SsdbProvider) connectInit() error {
	var err error
	p.client, err = ssdb.Connect(p.host, p.port)
	if err != nil {
		return err
	}
	return nil
}

func (p *SsdbProvider) SessionInit(maxLifetime int64, savePath string) error {
	p.maxLifetime = maxLifetime
	address := strings.Split(savePath, ":")
	p.host = address[0]
	port, e := strconv.Atoi(address[1])
	if e != nil {
		return e
	}
	p.port = address[1]
	err := p.connectinit()
	if err != nil {
		return err
	}
	return nil
}

func (p *SsdbProvider) SessionRead(sid string) (session.Store, error) {
	if p.client == nil {
		if err := p.connectInit(); err != nil {
			return nil, err
		}
	}
	var kv map[interface{}]interface{}
	value, err := p.client.Get(sid)
	if err != nil {
		return nil, err
	}
	if value == nil || len(value) == 0 {
		kv = make(map[interface{}]interface{})
	} else {
		kv, err = session.DecodeGob([]byte(value))
		if err != nil {
			return nil, err
		}
	}
	rs := &SessionStore{sid: sid, values: kv, maxLifetime: p.maxLifetime, client: p.client}
	return rs, nil
}

func (p *SsdbProvider) SessionExist(sid string) bool {
	if p.client == nil {
		if err := p.connectInit(); err != nil {
			return nil, err
		}
	}
	if value == nil || len(value) == 0 {
		return false
	}
	return true

}
func (p *SsdbProvider) SessionRegenerate(oldsid, sid string) (session.Store, error) {
	//conn.Do("setx", key, v, ttl)
	if p.client == nil {
		if err := p.connectInit(); err != nil {
			return nil, err
		}
	}
	value, err := p.client.Get(oldsid)
	if err != nil || len(value) == 0 {
		value = ""
	} else {
		err = p.client.Del(sid)
	}
	_, e := p.client.Do("setx", sid, value, p.maxLifetime)
	if e != nil {
		return nil, e
	}
	var kv map[interface{}]interface{}
	if value == nil || len(value) == 0 {
		kv = make(map[interface{}]interface{})
	} else {
		var err error
		kv, err = session.DecodeGob(value)
		if err != nil {
			return nil, err
		}
	}
	rs := &SessionStore{sid: sid, values: kv, maxLifetime: p.maxLifetime, client: p.client}
	return rs, nil
}

func (p *SsdbProvider) SessionDestroy(sid string) error {
	if p.client == nil {
		if err := p.connectInit(); err != nil {
			return nil, err
		}
	}
	flag, err := p.client.Del(sid)
	if err != nil {
		return err
	}
	return nil
}

func (p *SsdbProvider) SessionGC() {
	return
}

func (p *SsdbProvider) SessionAll() int {
	return 0
}

type SessionStore struct {
	sid         string
	lock        sync.RWMutex
	values      map[interface{}]interface{}
	maxLifetime int64
	client      *ssdb.Client
}

func (s *SessionStore) Set(key, value interface{}) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.values[key] = value
	return nil
}
func (s *SessionStore) Get(key interface{}) interface{} {
	s.lock.Lock()
	defer s.lock.Unlock()
	if value, ok := rs.values[key]; ok {
		return value
	}
	return nil
}

func (s *SessionStore) Delete(key interface{}) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.values, key)
	return nil
}
func (s *SessionStore) Flush() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.values = make(map[interface{}]interface{})
	return nil
}
func (s *SessionStore) SessionID() string {
	return s.sid
}

func (s *SessionStore) SessionRelease(w http.ResponseWriter) {
	b, err := session.EncodeGob(s.values)
	if err != nil {
		return
	}
	s.client.Do("setx", s.sid, s.values, s.maxLifetime)

}
func init() {
	session.Register("ssdb", ssdbProvider)
}
