package session

import (
	"log"
	"sync"
	"time"
)

type SessionFromMemory struct {
	sid              string
	lock             sync.Mutex
	lastAccessedTime time.Time
	maxAge           int64
	data             map[interface{}]interface{}
}

func newSessionFromMemory(sid string) *SessionFromMemory {
	return &SessionFromMemory{
		sid:    sid,
		data:   make(map[interface{}]interface{}),
		maxAge: DEFAULT_TIME,
	}
}

func (si *SessionFromMemory) Set(key interface{}, value interface{}) {
	si.lock.Lock()
	defer si.lock.Unlock()
	si.data[key] = value
}

func (si *SessionFromMemory) Get(key interface{}) interface{} {
	return si.data[key]
}

func (si *SessionFromMemory) Remove(key interface{}) error {
	if si.data[key] != nil {
		delete(si.data, key)
	}
	return nil
}

func (si *SessionFromMemory) GetLastAccessedTime() (time.Time, error) {
	return si.lastAccessedTime, nil
}

func (si *SessionFromMemory) UpdateLastAccessedTime() {
	si.lastAccessedTime = time.Now()
}

func (si *SessionFromMemory) GetMaxAge() int64 {
	return si.maxAge
}

func (si *SessionFromMemory) SetMaxAge(age int64) {
	si.maxAge = age
}

func (si *SessionFromMemory) GetId() string {
	return si.sid
}

func (si *SessionFromMemory) Destroy() bool {
	return true
}

type FromMemory struct {
	lock     sync.Mutex
	sessions map[string]Session
}

func newFromMemory() *FromMemory {
	return &FromMemory{
		sessions: make(map[string]Session),
	}
}

func (fm *FromMemory) InitSession(sid string, maxAge int64) (Session, error) {
	fm.lock.Lock()
	defer fm.lock.Unlock()
	newSession := newSessionFromMemory(sid)
	if maxAge != 0 {
		newSession.SetMaxAge(maxAge)
	}
	newSession.UpdateLastAccessedTime()
	fm.sessions[sid] = newSession
	log.Println(newSession)
	return newSession, nil
}

func (fm *FromMemory) GetSession(sid string) Session {
	return fm.sessions[sid]
}

func (fm *FromMemory) DestroySession(sid string) error {
	if fm.sessions[sid] != nil {
		delete(fm.sessions, sid)
	}
	return nil
}

func (fm *FromMemory) GCSession() bool {
	sessions := fm.sessions
	if len(sessions) == 0 {
		return true
	}
	log.Println("xxxxxxxxxxxxxx--gc-session", sessions)
	now := time.Now().Unix()
	for key, value := range sessions {
		time, err := value.GetLastAccessedTime()
		if err != nil {
			continue
		}
		t := time.Unix() + int64(value.GetMaxAge())
		if t < now {
			log.Println("timeout------->", value)
			delete(fm.sessions, key)
		}
	}
	return true
}
