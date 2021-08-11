package session

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Provider interface {
	InitSession(sid string, maxAge int64) (Session, error)
	GetSession(sid string) Session
	DestroySession(sid string) error
	GCSession() bool
}

func newProvider(db *sql.DB) Provider {
	if db != nil {
		return newFromDatabase(db)
	} else {
		return newFromMemory()
	}
}

type SessionManager struct {
	cookieName string
	db         *sql.DB
	storage    Provider
	maxAge     int64
	lock       sync.Mutex
}

func NewSessionManager(cookieName string, db *sql.DB) *SessionManager {
	sessionManager := &SessionManager{
		cookieName: cookieName,
		db:         db,
		storage:    newProvider(db),
		maxAge:     DEFAULT_TIME,
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
		session, err := m.storage.InitSession(sid, m.maxAge)
		if err != nil {
			log.Println(err.Error())
			return nil
		}
		maxAge := m.maxAge
		uid_cookie := &http.Cookie{
			Name:     m.cookieName,
			Value:    url.QueryEscape(sid),
			Path:     "/",
			HttpOnly: false,
			MaxAge:   int(maxAge),
		}
		http.SetCookie(w, uid_cookie)
		return session
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session := m.storage.GetSession(sid)
		if session == nil {
			newSession, err := m.storage.InitSession(sid, m.maxAge)
			if err != nil {
				log.Println(err.Error())
				return nil
			}
			maxAge := m.maxAge
			newCookie := http.Cookie{
				Name:     m.cookieName,
				Value:    url.QueryEscape(sid),
				Path:     "/",
				HttpOnly: true,
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

	if m.maxAge != 0 {
		cookie.MaxAge = int(m.maxAge)
	} else {
		cookie.MaxAge = int(session.GetMaxAge())
	}
	http.SetCookie(w, cookie)
	return session
}

func RandomId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	randId := base64.URLEncoding.EncodeToString(b)
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
