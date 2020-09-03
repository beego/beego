package ssdb

import (
	"context"
	"net/http"

	"github.com/astaxie/beego/pkg/adapter/session"

	beeSsdb "github.com/astaxie/beego/pkg/infrastructure/session/ssdb"
)

// Provider holds ssdb client and configs
type Provider beeSsdb.Provider

// SessionInit init the ssdb with the config
func (p *Provider) SessionInit(maxLifetime int64, savePath string) error {
	return (*beeSsdb.Provider)(p).SessionInit(context.Background(), maxLifetime, savePath)
}

// SessionRead return a ssdb client session Store
func (p *Provider) SessionRead(sid string) (session.Store, error) {
	s, err := (*beeSsdb.Provider)(p).SessionRead(context.Background(), sid)
	return session.CreateNewToOldStoreAdapter(s), err
}

// SessionExist judged whether sid is exist in session
func (p *Provider) SessionExist(sid string) bool {
	res, _ := (*beeSsdb.Provider)(p).SessionExist(context.Background(), sid)
	return res
}

// SessionRegenerate regenerate session with new sid and delete oldsid
func (p *Provider) SessionRegenerate(oldsid, sid string) (session.Store, error) {
	s, err := (*beeSsdb.Provider)(p).SessionRegenerate(context.Background(), oldsid, sid)
	return session.CreateNewToOldStoreAdapter(s), err
}

// SessionDestroy destroy the sid
func (p *Provider) SessionDestroy(sid string) error {
	return (*beeSsdb.Provider)(p).SessionDestroy(context.Background(), sid)
}

// SessionGC not implemented
func (p *Provider) SessionGC() {
	(*beeSsdb.Provider)(p).SessionGC(context.Background())
}

// SessionAll not implemented
func (p *Provider) SessionAll() int {
	return (*beeSsdb.Provider)(p).SessionAll(context.Background())
}

// SessionStore holds the session information which stored in ssdb
type SessionStore beeSsdb.SessionStore

// Set the key and value
func (s *SessionStore) Set(key, value interface{}) error {
	return (*beeSsdb.SessionStore)(s).Set(context.Background(), key, value)
}

// Get return the value by the key
func (s *SessionStore) Get(key interface{}) interface{} {
	return (*beeSsdb.SessionStore)(s).Get(context.Background(), key)
}

// Delete the key in session store
func (s *SessionStore) Delete(key interface{}) error {
	return (*beeSsdb.SessionStore)(s).Delete(context.Background(), key)
}

// Flush delete all keys and values
func (s *SessionStore) Flush() error {
	return (*beeSsdb.SessionStore)(s).Flush(context.Background())
}

// SessionID return the sessionID
func (s *SessionStore) SessionID() string {
	return (*beeSsdb.SessionStore)(s).SessionID(context.Background())
}

// SessionRelease Store the keyvalues into ssdb
func (s *SessionStore) SessionRelease(w http.ResponseWriter) {
	(*beeSsdb.SessionStore)(s).SessionRelease(context.Background(), w)
}
