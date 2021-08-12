package session

import (
	"sync"
	"time"
)

type SessionFromMemory struct {
	sid              string
	lock             sync.Mutex
	lastAccessedTime time.Time
	maxAge           int64
	data             map[string]interface{}
}

func newSessionFromMemory(sid string, maxAge int64) *SessionFromMemory {
	return &SessionFromMemory{
		sid:    sid,
		data:   make(map[string]interface{}),
		maxAge: maxAge,
	}
}

func (si *SessionFromMemory) Set(key string, value interface{}) error {
	si.lock.Lock()
	defer si.lock.Unlock()
	si.data[key] = value
	return nil
}

func (si *SessionFromMemory) Get(key string) interface{} {
	return si.data[key]
}

func (si *SessionFromMemory) Remove(key string) error {
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
	newSession := newSessionFromMemory(sid, maxAge)
	newSession.UpdateLastAccessedTime()
	fm.sessions[sid] = newSession
	// log.Println(newSession)
	return newSession, nil
}

func (fm *FromMemory) GetSession(sid string) Session {
	return fm.sessions[sid]
}

func (fm *FromMemory) Set(sid string, key string, value interface{}) error {
	return fm.GetSession(sid).(*SessionFromMemory).Set(key, value)
}

func (fm *FromMemory) Get(sid string, key string) interface{} {
	return fm.GetSession(sid).(*SessionFromMemory).Get(key)
}

func (fm *FromMemory) Remove(sid string, key string) error {
	return fm.GetSession(sid).(*SessionFromMemory).Remove(key)
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
	// log.Println("xxxxxxxxxxxxxx--gc-session", sessions)
	now := time.Now().Unix()
	for key, value := range sessions {
		time := value.(*SessionFromMemory).lastAccessedTime
		t := time.Unix() + value.(*SessionFromMemory).maxAge
		if t < now {
			// log.Println("timeout------->", value)
			delete(fm.sessions, key)
		}
	}
	return true
}
