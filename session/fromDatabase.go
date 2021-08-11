package session

import (
	"database/sql"
	"log"
	"sync"
	"time"
)

type SessionFromDatabase struct {
	sid              string
	lock             sync.Mutex
	lastAccessedTime time.Time
	maxAge           int64
	data             map[interface{}]interface{}
}

func newSessionFromDatabase(sid string) *SessionFromDatabase {
	return &SessionFromDatabase{
		sid:    sid,
		data:   make(map[interface{}]interface{}),
		maxAge: DEFAULT_TIME,
	}
}

func (si *SessionFromDatabase) Set(key interface{}, value interface{}) {
	si.lock.Lock()
	defer si.lock.Unlock()
	si.data[key] = value
}

func (si *SessionFromDatabase) Get(key interface{}) interface{} {
	return si.data[key]
}

func (si *SessionFromDatabase) Remove(key interface{}) error {
	if si.data[key] != nil {
		delete(si.data, key)
	}
	return nil
}

func (si *SessionFromDatabase) GetLastAccessedTime() time.Time {
	return si.lastAccessedTime
}

func (si *SessionFromDatabase) UpdateLastAccessedTime() {
	si.lastAccessedTime = time.Now()
}

func (si *SessionFromDatabase) GetMaxAge() int64 {
	return si.maxAge
}

func (si *SessionFromDatabase) SetMaxAge(age int64) {
	si.maxAge = age
}

func (si *SessionFromDatabase) GetId() string {
	return si.sid
}

type FromDatabase struct {
	lock     sync.Mutex
	db       *sql.DB
	sessions map[string]Session
}

func newFromDatabase(db *sql.DB) *FromDatabase {
	return &FromDatabase{
		db:       db,
		sessions: make(map[string]Session),
	}
}

func (fd *FromDatabase) InitSession(sid string, maxAge int64) (Session, error) {
	fd.lock.Lock()
	defer fd.lock.Unlock()
	newSession := newSessionFromDatabase(sid)
	if maxAge != 0 {
		newSession.SetMaxAge(maxAge)
	}
	newSession.UpdateLastAccessedTime()
	fd.sessions[sid] = newSession
	log.Println(newSession)
	return newSession, nil
}

func (fd *FromDatabase) SetSession(session Session) error {
	fd.sessions[session.GetId()] = session
	return nil
}

func (fd *FromDatabase) GetSession(key string) Session {
	return fd.sessions[key]
}

func (fd *FromDatabase) DestroySession(sid string) error {
	if fd.sessions[sid] != nil {
		delete(fd.sessions, sid)
	}
	return nil
}

func (fd *FromDatabase) GCSession() {
	sessions := fd.sessions
	if len(sessions) == 0 {
		return
	}
	log.Println("xxxxxxxxxxxxxx--gc-session", sessions)
	now := time.Now().Unix()
	for key, value := range sessions {
		t := (value.GetLastAccessedTime().Unix()) +
			int64(value.GetMaxAge())
		if t < now {
			log.Println("timeout------->", value)
			delete(fd.sessions, key)
		}
	}
}
