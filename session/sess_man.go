package session

import (
	"errors"
	"fmt"
	"net/http"
	"net/textproto"
	"net/url"
	"time"
)

type Manager struct {
	provider    Provider
	config      *ManagerConfig
	idGenerator IdGeneration
}

// NewManager Create new Manager with provider name and json config string.
func NewManager(provideName string, cf *ManagerConfig) (*Manager, error) {
	provider, ok := provides[provideName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
	}
	if cf.Maxlifetime == 0 {
		cf.Maxlifetime = cf.Gclifetime
	}
	if cf.EnableSidInHTTPHeader {
		if cf.SessionNameInHTTPHeader == "" {
			panic(errors.New("SessionNameInHTTPHeader is empty"))
		}
		strMimeHeader := textproto.CanonicalMIMEHeaderKey(cf.SessionNameInHTTPHeader)
		if cf.SessionNameInHTTPHeader != strMimeHeader {
			strErrMsg := "SessionNameInHTTPHeader (" + cf.SessionNameInHTTPHeader + ") has the wrong format, it should be like this : " + strMimeHeader
			panic(errors.New(strErrMsg))
		}
	}
	tg := NewJwtTokenGenerator(cf)
	err := provider.SessionInit(cf.Maxlifetime, cf.ProviderConfig, tg)
	if err != nil {
		return nil, err
	}
	if cf.SessionIDLength == 0 {
		cf.SessionIDLength = 16
	}
	return &Manager{
		provider,
		cf,
		tg,
	}, nil
}

// getRawSessionId retrieves session identifier from HTTP Request.
// First try to retrieve id by reading from cookie, session cookie name is configurable,
// if not exist, then retrieve id from querying parameters.
//
// error is not nil when there is anything wrong.
// sid is empty when need to generate a new session id
// otherwise return an valid session id.
func (manager *Manager) getRawSessionId(r *http.Request) (string, error) {
	cookie, errs := r.Cookie(manager.config.CookieName)
	if errs != nil || cookie.Value == "" {
		var rawSid string
		if manager.config.EnableSidInURLQuery {
			errs := r.ParseForm()
			if errs != nil {
				return "", errs
			}
			rawSid = r.FormValue(manager.config.CookieName)
		}

		// if not found in Cookie / param, then read it from request headers
		if manager.config.EnableSidInHTTPHeader && rawSid == "" {
			sids, isFound := r.Header[manager.config.SessionNameInHTTPHeader]
			if isFound && len(sids) != 0 {
				return sids[0], nil
			}
		}
		return rawSid, nil
	}
	// HTTP Request contains cookie for sessionid info.
	return url.QueryUnescape(cookie.Value)
}

// SessionStart generate or read the session id from http request.
// if session id exists, return SessionStore with this id.
func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session Store, err error) {
	rawSid, errs := manager.getRawSessionId(r)
	if errs != nil {
		return nil, errs
	}

	if rawSid != "" && manager.provider.SessionExist(rawSid) {
		return manager.provider.SessionRead(rawSid)
	}

	// Generate a new session
	rawSid, errs = manager.idGenerator.GenerateSessionID()
	if errs != nil {
		return nil, errs
	}

	session, err = manager.provider.SessionRead(rawSid)
	if err != nil {
		return nil, err
	}
	cookie := &http.Cookie{
		Name:     manager.config.CookieName,
		Value:    url.QueryEscape(rawSid),
		Path:     "/",
		HttpOnly: !manager.config.DisableHTTPOnly,
		Secure:   manager.isSecure(r),
		Domain:   manager.config.Domain,
	}
	if manager.config.CookieLifeTime > 0 {
		cookie.MaxAge = manager.config.CookieLifeTime
		cookie.Expires = time.Now().Add(time.Duration(manager.config.CookieLifeTime) * time.Second)
	}
	if manager.config.EnableSetCookie {
		http.SetCookie(w, cookie)
	}
	r.AddCookie(cookie)

	if manager.config.EnableSidInHTTPHeader {
		r.Header.Set(manager.config.SessionNameInHTTPHeader, rawSid)
		w.Header().Set(manager.config.SessionNameInHTTPHeader, rawSid)
	}

	return
}

// SessionDestroy Destroy session by its id in http request cookie.
func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	if manager.config.EnableSidInHTTPHeader {
		r.Header.Del(manager.config.SessionNameInHTTPHeader)
		w.Header().Del(manager.config.SessionNameInHTTPHeader)
	}

	cookie, err := r.Cookie(manager.config.CookieName)
	if err != nil || cookie.Value == "" {
		return
	}

	sid, _ := url.QueryUnescape(cookie.Value)
	manager.provider.SessionDestroy(sid)
	if manager.config.EnableSetCookie {
		expiration := time.Now()
		cookie = &http.Cookie{Name: manager.config.CookieName,
			Path: "/",
			HttpOnly: !manager.config.DisableHTTPOnly,
			Expires: expiration,
			MaxAge: -1}

		http.SetCookie(w, cookie)
	}
}

// GetSessionStore Get SessionStore by its id.
func (manager *Manager) GetSessionStore(sid string) (sessions Store, err error) {
	sessions, err = manager.provider.SessionRead(sid)
	return
}

// GC Start session gc process.
// it can do gc in times after gc lifetime.
func (manager *Manager) GC() {
	manager.provider.SessionGC()
	time.AfterFunc(time.Duration(manager.config.Gclifetime)*time.Second, func() { manager.GC() })
}

// SessionRegenerateID Regenerate a session id for this SessionStore who's id is saving in http request.
func (manager *Manager) SessionRegenerateID(w http.ResponseWriter, r *http.Request) (session Store) {
	rawSid, err := manager.idGenerator.GenerateSessionID()
	if err != nil {
		return
	}
	cookie, err := r.Cookie(manager.config.CookieName)
	if err != nil || cookie.Value == "" {
		//delete old cookie
		session, _ = manager.provider.SessionRead(rawSid)
		cookie = &http.Cookie{Name: manager.config.CookieName,
			Value: url.QueryEscape(rawSid),
			Path: "/",
			HttpOnly: !manager.config.DisableHTTPOnly,
			Secure: manager.isSecure(r),
			Domain: manager.config.Domain,
		}
	} else {
		oldsid, _ := url.QueryUnescape(cookie.Value)
		session, _ = manager.provider.SessionRegenerate(oldsid, rawSid)
		cookie.Value = url.QueryEscape(rawSid)
		cookie.HttpOnly = true
		cookie.Path = "/"
	}
	if manager.config.CookieLifeTime > 0 {
		cookie.MaxAge = manager.config.CookieLifeTime
		cookie.Expires = time.Now().Add(time.Duration(manager.config.CookieLifeTime) * time.Second)
	}
	if manager.config.EnableSetCookie {
		http.SetCookie(w, cookie)
	}
	r.AddCookie(cookie)

	if manager.config.EnableSidInHTTPHeader {
		r.Header.Set(manager.config.SessionNameInHTTPHeader, rawSid)
		w.Header().Set(manager.config.SessionNameInHTTPHeader, rawSid)
	}
	return
}

// GetActiveSession Get all active sessions count number.
func (manager *Manager) GetActiveSession() int {
	return manager.provider.SessionAll()
}

// SetSecure Set cookie with https.
func (manager *Manager) SetSecure(secure bool) {
	manager.config.Secure = secure
}

// Set cookie with https.
func (manager *Manager) isSecure(req *http.Request) bool {
	if !manager.config.Secure {
		return false
	}
	if req.URL.Scheme != "" {
		return req.URL.Scheme == "https"
	}
	if req.TLS == nil {
		return false
	}
	return true
}
