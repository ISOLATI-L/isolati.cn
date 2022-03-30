package session

import (
	"crypto/md5"
	"crypto/rand"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type SessionManager struct {
	cookieName string
	db         *sql.DB
	storage    provider
	maxAge     int64
	httpOnly   bool
	lock       sync.Mutex
}

func NewSessionManager(cookieName string, db *sql.DB, maxAge int64, httpOnly bool) *SessionManager {
	sessionManager := &SessionManager{
		cookieName: cookieName,
		db:         db,
		storage:    newProvider(db),
		maxAge:     maxAge,
		httpOnly:   httpOnly,
	}
	go sessionManager.gc()
	return sessionManager
}

func (m *SessionManager) GetCookieName() string {
	return m.cookieName
}

func (m *SessionManager) BeginSession(w http.ResponseWriter, r *http.Request) string {
	m.lock.Lock()
	defer m.lock.Unlock()
	cookie, err := r.Cookie(m.cookieName)
	if err != nil || cookie.Value == "" {
		sid := m.randomId()
		maxAge := m.maxAge
		session, err := m.storage.initSession(sid, maxAge)
		if err != nil {
			log.Println(err.Error())
			return ""
		}
		uid_cookie := &http.Cookie{
			Name:     m.cookieName,
			Value:    url.QueryEscape(sid),
			Path:     "/",
			HttpOnly: m.httpOnly,
			MaxAge:   int(maxAge),
		}
		http.SetCookie(w, uid_cookie)
		return session.getID()
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session := m.getSession(sid)
		if session == nil {
			maxAge := m.maxAge
			newSession, err := m.storage.initSession(sid, maxAge)
			if err != nil {
				log.Println(err.Error())
				return ""
			}
			newCookie := http.Cookie{
				Name:     m.cookieName,
				Value:    url.QueryEscape(sid),
				Path:     "/",
				HttpOnly: m.httpOnly,
				MaxAge:   int(maxAge),
				Expires:  time.Now().Add(time.Duration(maxAge)),
			}
			http.SetCookie(w, &newCookie)
			return newSession.getID()
		}
		session.updateLastAccessedTime()
		return session.getID()
	}
}

func (m *SessionManager) getSession(sid string) session {
	return m.storage.getSession(sid)
}

func (m *SessionManager) IsExists(sid string) bool {
	session := m.getSession(sid)
	return session != nil
}

func (m *SessionManager) Set(sid string, key string, value interface{}) error {
	return m.storage.set(sid, key, value)
}

func (m *SessionManager) Get(sid string, key string) interface{} {
	return m.storage.get(sid, key)
}

func (m *SessionManager) Remove(sid string, key string) error {
	return m.storage.remove(sid, key)
}

func (m *SessionManager) EndSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(m.cookieName)
	if err != nil || cookie.Value == "" {
		return
	} else {
		m.lock.Lock()
		defer m.lock.Unlock()

		sid, err := url.QueryUnescape(cookie.Value)
		if err != nil {
			log.Println(err.Error())
			return
		}
		m.storage.destroySession(sid)

		cookie2 := http.Cookie{
			MaxAge:  0,
			Name:    m.cookieName,
			Value:   "",
			Path:    "/",
			Expires: time.Now().Add(time.Duration(0)),
		}

		http.SetCookie(w, &cookie2)
	}
}

func (m *SessionManager) Update(w http.ResponseWriter, r *http.Request) string {
	m.lock.Lock()
	defer m.lock.Unlock()

	cookie, err := r.Cookie(m.cookieName)
	if err != nil {
		return ""
	}

	sid, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		log.Println(err.Error())
		return ""
	}

	session := m.getSession(sid)
	if session == nil {
		return ""
	}
	session.updateLastAccessedTime()

	if m.maxAge != int64(cookie.MaxAge) {
		cookie.MaxAge = int(m.maxAge)
	}
	http.SetCookie(w, cookie)
	return session.getID()
}

func randomId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	randId := fmt.Sprintf("%x", md5.Sum(b))
	return randId
}

func (m *SessionManager) randomId() string {
	var randId string
	for {
		randId = randomId()
		if !m.IsExists(randId) {
			break
		}
	}
	return randId
}

const AGE2 = int(60 * time.Second)

func (m *SessionManager) gc() {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.storage.gcSession() {
		time.AfterFunc(time.Duration(AGE2), m.gc)
	}
}
