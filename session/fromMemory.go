package session

import (
	"log"
	"sync"
	"time"
)

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

func (fm *FromMemory) SetSession(session Session) error {
	fm.sessions[session.GetId()] = session
	return nil
}

func (fm *FromMemory) GetSession(key string) Session {
	return fm.sessions[key]
}

func (fm *FromMemory) DestroySession(sid string) error {
	if fm.sessions[sid] != nil {
		delete(fm.sessions, sid)
	}
	return nil
}

func (fm *FromMemory) GCSession() {
	sessions := fm.sessions
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
			delete(fm.sessions, key)
		}
	}
}
