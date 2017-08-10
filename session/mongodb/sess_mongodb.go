package mongodb

import (
	"net/http"
	"sync"
	"time"

	"github.com/astaxie/beego/session"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	collectionName  = "session"
	mongodbProvider = &Provider{}
)

// Provider mongodb session provider
type Provider struct {
	maxLifetime int64
	savePath    string
	mgoSession  *mgo.Session
}

// SessionInit connect mongodb
func (p *Provider) SessionInit(maxLifetime int64, savePath string) error {
	p.maxLifetime = maxLifetime
	p.savePath = savePath

	// init mongodb session
	if p.mgoSession == nil {
		s, err := mgo.Dial(savePath)
		if err != nil {
			return err
		}
		p.mgoSession = s
	}

	return nil
}

// SessionRead read mongodb seesion by sid
func (p *Provider) SessionRead(sid string) (session.Store, error) {
	mgosession := p.mgoSession.Clone()
	defer mgosession.Close()

	var s bson.M
	change := mgo.Change{
		Update: bson.M{
			"$setOnInsert": bson.M{
				"session_key":    sid,
				"session_data":   nil,
				"session_expire": time.Now().Unix() + p.maxLifetime,
			},
		},
		Upsert: true,
	}
	_, err := mgosession.DB("").C(collectionName).Find(bson.M{"session_key": sid}).Apply(change, &s)
	if err != nil {
		return nil, err
	}

	var kv map[interface{}]interface{}
	if s == nil {
		kv = make(map[interface{}]interface{})
	} else {
		kv, err = session.DecodeGob([]byte(s["session_data"].([]uint8)))
		if err != nil {
			return nil, err
		}

	}

	return &SessionStore{sid: sid, values: kv, maxLifetime: p.maxLifetime, mgoSession: p.mgoSession}, nil
}

// SessionExist check mongodb session exist by id
func (p *Provider) SessionExist(sid string) bool {
	mgosession := p.mgoSession.Clone()
	defer mgosession.Close()

	var s interface{}
	if mgosession.DB("").C(collectionName).Find(bson.M{"session_key": sid}).One(&s) != nil {
		return false
	}
	return true
}

// SessionRegenerate generate new sid for mongodb session
func (p *Provider) SessionRegenerate(oldsid, sid string) (session.Store, error) {
	mgosession := p.mgoSession.Clone()
	defer mgosession.Close()

	var s bson.M
	change := mgo.Change{
		Update: bson.M{
			"$set": bson.M{
				"session_key":    sid,
				"session_expire": time.Now().Unix() + p.maxLifetime,
			},
			"$setOnInsert": bson.M{
				"session_key":    sid,
				"session_data":   nil,
				"session_expire": time.Now().Unix() + p.maxLifetime,
			},
		},
		Upsert: true,
	}
	_, err := mgosession.DB("").C(collectionName).Find(bson.M{"session_key": oldsid}).Apply(change, &s)
	if err != nil {
		return nil, err
	}

	var kv map[interface{}]interface{}
	if s == nil {
		kv = make(map[interface{}]interface{})
	} else {
		kv, err = session.DecodeGob([]byte(s["session_data"].([]uint8)))
		if err != nil {
			return nil, err
		}
	}

	return &SessionStore{sid: sid, values: kv, maxLifetime: p.maxLifetime, mgoSession: p.mgoSession}, nil
}

// SessionDestroy remove mongodb session by sid
func (p *Provider) SessionDestroy(sid string) error {
	mgosession := p.mgoSession.Clone()
	defer mgosession.Close()

	err := mgosession.DB("").C(collectionName).Remove(bson.M{"session_key": sid})
	return err
}

// SessionGC remove all expire mongodb seesion
func (p *Provider) SessionGC() {
	mgosession := p.mgoSession.Clone()
	defer mgosession.Close()

	mgosession.DB("").C(collectionName).RemoveAll(bson.M{"session_expire": bson.M{"$lt": time.Now().Unix()}})
}

// SessionGC return all mongodb session
func (p *Provider) SessionAll() int {
	mgosession := p.mgoSession.Clone()
	defer mgosession.Close()

	count, _ := mgosession.DB("").C(collectionName).Find(nil).Count()

	return count
}

// SessionStore mongodb session store
type SessionStore struct {
	sid         string
	lock        sync.RWMutex
	values      map[interface{}]interface{}
	maxLifetime int64
	mgoSession  *mgo.Session
}

// Set set value in mongodb session
func (s *SessionStore) Set(key, value interface{}) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.values[key] = value
	return nil
}

// Get get value from mongodb session
func (s *SessionStore) Get(key interface{}) interface{} {
	s.lock.Lock()
	defer s.lock.Unlock()
	if value, ok := s.values[key]; ok {
		return value
	}
	return nil
}

// Delete delete value in mongodb session
func (s *SessionStore) Delete(key interface{}) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.values, key)
	return nil
}

// Flush clear all values in mongodb session
func (s *SessionStore) Flush() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.values = make(map[interface{}]interface{})
	return nil
}

// SessionID get mongodb session id
func (s *SessionStore) SessionID() string {
	return s.sid
}

// SessionRelease save session values mongodb
func (s *SessionStore) SessionRelease(w http.ResponseWriter) {
	mgosession := s.mgoSession.Clone()
	defer mgosession.Close()

	b, err := session.EncodeGob(s.values)
	if err != nil {
		return
	}

	mgosession.DB("").C(collectionName).Update(bson.M{"session_key": s.sid}, bson.M{"$set": bson.M{"session_data": b}})
}

func init() {
	session.Register("mongodb", mongodbProvider)
}
