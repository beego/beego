package ssdb

import (
	"net/http"

	"github.com/astaxie/beego/session"
	"github.com/ssdb/gossdb/ssdb"
)

var ssdbSsdbProvider = &SsdbSsdbProvider{}

type SsdbProvider struct {
}

func (r *SsdbProvider) SessionInit(maxlifetime int64, savePath string) error {
}

func (r *SsdbProvider) SessionRead(sid string) (session.Store, error) {
}

func (r *SsdbProvider) SessionExist(sid string) bool {
}
func (r *SsdbProvider) SessionRegenerate(oldsid, sid string) (session.Store, error) {
}

func (r *SsdbProvider) SessionDestroy(sid string) error {
}

func (r *SsdbProvider) SessionGC() {
	return
}

func (r *SsdbProvider) SessionAll() int {
	return 0
}

type SessionStore struct {
}

func (s *SessionStore) Set(key, value interface{}) error {
}
func (s *SessionStore) Get(key interface{}) interface{} {
}
func (s *SessionStore) Delete(key interface{}) error {
}
func (s *SessionStore) Flush() error {
}
func (s *SessionStore) SessionID() string {
}

func (s *SessionStore) SessionRelease(w http.ResponseWriter) {
}
func init() {
	session.Register("redis", redispder)
}
