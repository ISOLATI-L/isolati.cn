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
	storage    Provider
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
	go sessionManager.GC()
	return sessionManager
}

func (m *SessionManager) GetCookieName() string {
	return m.cookieName
}

func (m *SessionManager) BeginSession(w http.ResponseWriter, r *http.Request) Session {
	m.lock.Lock()
	defer m.lock.Unlock()
	cookie, err := r.Cookie(m.cookieName)
	if err != nil || cookie.Value == "" {
		sid := m.randomId()
		maxAge := m.maxAge
		session, err := m.storage.InitSession(sid, maxAge)
		if err != nil {
			log.Println(err.Error())
			return nil
		}
		uid_cookie := &http.Cookie{
			Name:     m.cookieName,
			Value:    url.QueryEscape(sid),
			Path:     "/",
			HttpOnly: m.httpOnly,
			MaxAge:   int(maxAge),
		}
		http.SetCookie(w, uid_cookie)
		return session
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session := m.storage.GetSession(sid)
		if session == nil {
			maxAge := m.maxAge
			newSession, err := m.storage.InitSession(sid, maxAge)
			if err != nil {
				log.Println(err.Error())
				return nil
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
			return newSession
		}
		return session
	}
}

func (m *SessionManager) GetSessionById(sid string) Session {
	return m.storage.GetSession(sid)
}

func (m *SessionManager) MemoryIsExists(sid string) bool {
	session := m.storage.GetSession(sid)
	return session != nil
}

func (m *SessionManager) Destroy(w http.ResponseWriter, r *http.Request) {
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
		m.storage.DestroySession(sid)

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

func (m *SessionManager) Update(w http.ResponseWriter, r *http.Request) Session {
	m.lock.Lock()
	defer m.lock.Unlock()

	cookie, err := r.Cookie(m.cookieName)
	if err != nil {
		return nil
	}

	sid, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	session := m.storage.GetSession(sid)
	if session == nil {
		return nil
	}
	session.UpdateLastAccessedTime()

	if m.maxAge != int64(cookie.MaxAge) {
		cookie.MaxAge = int(m.maxAge)
	}
	http.SetCookie(w, cookie)
	return session
}

func RandomId() string {
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
		randId = RandomId()
		if !m.MemoryIsExists(randId) {
			break
		}
	}
	return randId
}

const AGE2 = int(60 * time.Second)

func (m *SessionManager) GC() {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.storage.GCSession() {
		time.AfterFunc(time.Duration(AGE2), m.GC)
	}
}
