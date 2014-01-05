package session

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"net/http"
	"net/url"
	"sync"
)

var cookiepder = &CookieProvider{}

type CookieSessionStore struct {
	sid    string
	values map[interface{}]interface{} //session data
	lock   sync.RWMutex
}

func (st *CookieSessionStore) Set(key, value interface{}) error {
	st.lock.Lock()
	defer st.lock.Unlock()
	st.values[key] = value
	return nil
}

func (st *CookieSessionStore) Get(key interface{}) interface{} {
	st.lock.RLock()
	defer st.lock.RUnlock()
	if v, ok := st.values[key]; ok {
		return v
	} else {
		return nil
	}
	return nil
}

func (st *CookieSessionStore) Delete(key interface{}) error {
	st.lock.Lock()
	defer st.lock.Unlock()
	delete(st.values, key)
	return nil
}

func (st *CookieSessionStore) Flush() error {
	st.lock.Lock()
	defer st.lock.Unlock()
	st.values = make(map[interface{}]interface{})
	return nil
}

func (st *CookieSessionStore) SessionID() string {
	return st.sid
}

func (st *CookieSessionStore) SessionRelease(w http.ResponseWriter) {
	str, err := encodeCookie(cookiepder.block,
		cookiepder.config.SecurityKey,
		cookiepder.config.SecurityName,
		st.values)
	if err != nil {
		return
	}
	cookie := &http.Cookie{Name: cookiepder.config.CookieName,
		Value:    url.QueryEscape(str),
		Path:     "/",
		HttpOnly: true,
		Secure:   cookiepder.config.Secure}
	http.SetCookie(w, cookie)
	return
}

type cookieConfig struct {
	SecurityKey  string `json:"securityKey"`
	BlockKey     string `json:"blockKey"`
	SecurityName string `json:"securityName"`
	CookieName   string `json:"cookieName"`
	Secure       bool   `json:"secure"`
	Maxage       int    `json:"maxage"`
}

type CookieProvider struct {
	maxlifetime int64
	config      *cookieConfig
	block       cipher.Block
}

func (pder *CookieProvider) SessionInit(maxlifetime int64, config string) error {
	pder.config = &cookieConfig{}
	err := json.Unmarshal([]byte(config), pder.config)
	if err != nil {
		return err
	}
	if pder.config.BlockKey == "" {
		pder.config.BlockKey = string(generateRandomKey(16))
	}
	if pder.config.SecurityName == "" {
		pder.config.SecurityName = string(generateRandomKey(20))
	}
	pder.block, err = aes.NewCipher([]byte(pder.config.BlockKey))
	if err != nil {
		return err
	}
	return nil
}

func (pder *CookieProvider) SessionRead(sid string) (SessionStore, error) {
	kv := make(map[interface{}]interface{})
	kv, _ = decodeCookie(pder.block,
		pder.config.SecurityKey,
		pder.config.SecurityName,
		sid, pder.maxlifetime)
	rs := &CookieSessionStore{sid: sid, values: kv}
	return rs, nil
}

func (pder *CookieProvider) SessionExist(sid string) bool {
	return true
}

func (pder *CookieProvider) SessionRegenerate(oldsid, sid string) (SessionStore, error) {
	return nil, nil
}

func (pder *CookieProvider) SessionDestroy(sid string) error {
	return nil
}

func (pder *CookieProvider) SessionGC() {
	return
}

func (pder *CookieProvider) SessionAll() int {
	return 0
}

func (pder *CookieProvider) SessionUpdate(sid string) error {
	return nil
}

func init() {
	Register("cookie", cookiepder)
}
