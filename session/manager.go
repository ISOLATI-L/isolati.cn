package session

import (
	"crypto/md5"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

var ErrNoCookies error

func init() {
	ErrNoCookies = errors.New("No Cookies")
}

type SessionManager struct {
	cookieName string
	db         *sql.DB
	storage    provider
	maxAge     int64
	httpOnly   bool
}

func NewSessionManager(cookieName string, db *sql.DB, maxAge int64, httpOnly bool) *SessionManager {
	sessionManager := &SessionManager{
		cookieName: cookieName,
		db:         db,
		storage:    newProvider(),
		maxAge:     maxAge,
		httpOnly:   httpOnly,
	}
	return sessionManager
}

func (m *SessionManager) BeginTransaction() (*sql.Tx, error) {
	transaction, err := m.db.Begin()
	if err != nil {
		if transaction != nil {
			transaction.Rollback()
		}
		return nil, err
	}
	return transaction, nil
}

func (m *SessionManager) GetCookieName() string {
	return m.cookieName
}

func (m *SessionManager) Set(transaction *sql.Tx, sid string, key string, value any) error {
	return m.storage.set(transaction, sid, key, value)
}

func (m *SessionManager) Get(transaction *sql.Tx, sid string, key string) ([]byte, error) {
	return m.storage.get(transaction, sid, key)
}

func (m *SessionManager) Remove(transaction *sql.Tx, sid string, key string) error {
	return m.storage.remove(transaction, sid, key)
}

func (m *SessionManager) BeginSession(transaction *sql.Tx, w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie(m.cookieName)
	if err != nil || len(cookie.Value) == 0 {
		sid := m.randomId(transaction)
		maxAge := m.maxAge
		session, err := m.storage.initSession(transaction, sid, maxAge)
		if err != nil {
			return "", errors.New("Init Session Fail: " + err.Error())
		}
		uid_cookie := &http.Cookie{
			Name:     m.cookieName,
			Value:    url.QueryEscape(sid),
			Path:     "/",
			HttpOnly: m.httpOnly,
			MaxAge:   int(maxAge),
		}
		http.SetCookie(w, uid_cookie)
		return session.getID(), nil
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session, _ := m.getSession(transaction, sid)
		if session == nil {
			maxAge := m.maxAge
			newSession, err := m.storage.initSession(transaction, sid, maxAge)
			if err != nil {
				return "", errors.New("Init Session Fail: " + err.Error())
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
			return newSession.getID(), nil
		}
		session.updateLastAccessedTime(transaction)
		return session.getID(), nil
	}
}

func (m *SessionManager) getSid(transaction *sql.Tx, r *http.Request) (string, error) {
	cookie, err := r.Cookie(m.cookieName)
	if err != nil || len(cookie.Value) == 0 {
		return "", err
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session, err := m.getSession(transaction, sid)
		if session == nil {
			return "", ErrNoCookies
		}
		if err != nil {
			return "", err
		}
		session.updateLastAccessedTime(transaction)
		return session.getID(), nil
	}
}

func (m *SessionManager) SetByRequest(transaction *sql.Tx, r *http.Request, key string, value any) error {
	sid, err := m.getSid(transaction, r)
	if len(sid) == 0 {
		return ErrNoCookies
	} else if err != nil {
		return err
	} else {
		m.storage.update(transaction, sid)
		return m.storage.set(transaction, sid, key, value)
	}
}

func (m *SessionManager) GetByRequest(transaction *sql.Tx, r *http.Request, key string) ([]byte, error) {
	sid, err := m.getSid(transaction, r)
	if len(sid) == 0 {
		return nil, ErrNoCookies
	} else if err != nil {
		return nil, err
	} else {
		m.storage.update(transaction, sid)
		return m.storage.get(transaction, sid, key)
	}
}

func (m *SessionManager) RemoveByRequest(transaction *sql.Tx, r *http.Request, key string) error {
	sid, err := m.getSid(transaction, r)
	if len(sid) == 0 {
		return ErrNoCookies
	} else if err != nil {
		return err
	} else {
		m.storage.update(transaction, sid)
		return m.storage.remove(transaction, sid, key)
	}
}

func (m *SessionManager) getSession(transaction *sql.Tx, sid string) (session, error) {
	return m.storage.getSession(transaction, sid)
}

func (m *SessionManager) IsExists(transaction *sql.Tx, sid string) (bool, error) {
	session, err := m.getSession(transaction, sid)
	return session != nil, err
}

func (m *SessionManager) EndSession(transaction *sql.Tx, w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie(m.cookieName)
	if err != nil || len(cookie.Value) == 0 {
		return err
	}

	sid, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	m.storage.destroySession(transaction, sid)

	cookie2 := http.Cookie{
		MaxAge:  0,
		Name:    m.cookieName,
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(time.Duration(0)),
	}

	http.SetCookie(w, &cookie2)

	return nil
}

func (m *SessionManager) Update(transaction *sql.Tx, w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie(m.cookieName)
	if err != nil {
		return "", err
	}

	sid, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	session, err := m.getSession(transaction, sid)
	if session == nil {
		return "", ErrNoCookies
	}
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	session.updateLastAccessedTime(transaction)

	if m.maxAge != int64(cookie.MaxAge) {
		cookie.MaxAge = int(m.maxAge)
	}
	http.SetCookie(w, cookie)
	return session.getID(), nil
}

func randomId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	randId := fmt.Sprintf("%x", md5.Sum(b))
	return randId
}

func (m *SessionManager) randomId(transaction *sql.Tx) string {
	var randId string
	for {
		randId = randomId()
		session, _ := m.getSession(transaction, randId)
		if session == nil {
			break
		}
	}
	return randId
}
