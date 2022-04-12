package session

import (
	"encoding/json"
	"errors"
	"time"
)

var ErrNoData error

func init() {
	ErrNoData = errors.New("No Data")
}

type sessionFromMemory struct {
	sid              string
	lastAccessedTime time.Time
	maxAge           int64
	data             map[string][]byte
}

func newSessionFromMemory(sid string, maxAge int64) *sessionFromMemory {
	return &sessionFromMemory{
		sid:    sid,
		data:   make(map[string][]byte),
		maxAge: maxAge,
	}
}

func (si *sessionFromMemory) getID() string {
	return si.sid
}

func (si *sessionFromMemory) set(key string, value any) error {
	var err error
	si.data[key], err = json.Marshal(value)
	return err
}

func (si *sessionFromMemory) get(key string) ([]byte, error) {
	value, ok := si.data[key]
	if ok {
		return value, nil
	} else {
		return nil, ErrNoData
	}
}

func (si *sessionFromMemory) remove(key string) error {
	if si.data[key] != nil {
		delete(si.data, key)
	}
	return nil
}

func (si *sessionFromMemory) getLastAccessedTime() (time.Time, error) {
	return si.lastAccessedTime, nil
}

func (si *sessionFromMemory) updateLastAccessedTime() {
	si.lastAccessedTime = time.Now()
}

type fromMemory struct {
	sessions map[string]session
}

func newFromMemory() *fromMemory {
	return &fromMemory{
		sessions: make(map[string]session),
	}
}

func (fm *fromMemory) initSession(sid string, maxAge int64) (session, error) {
	newSession := newSessionFromMemory(sid, maxAge)
	newSession.updateLastAccessedTime()
	fm.sessions[sid] = newSession
	// log.Println(newSession)
	return newSession, nil
}

func (fm *fromMemory) getSession(sid string) session {
	return fm.sessions[sid]
}

func (fm *fromMemory) set(sid string, key string, value any) error {
	return fm.getSession(sid).(*sessionFromMemory).set(key, value)
}

func (fm *fromMemory) get(sid string, key string) ([]byte, error) {
	return fm.getSession(sid).(*sessionFromMemory).get(key)
}

func (fm *fromMemory) remove(sid string, key string) error {
	return fm.getSession(sid).(*sessionFromMemory).remove(key)
}

func (fm *fromMemory) update(sid string) error {
	fm.getSession(sid).(*sessionFromMemory).lastAccessedTime = time.Now()
	return nil
}

func (fm *fromMemory) destroySession(sid string) error {
	if fm.sessions[sid] != nil {
		delete(fm.sessions, sid)
	}
	return nil
}

func (fm *fromMemory) gcSession() bool {
	sessions := fm.sessions
	if len(sessions) == 0 {
		return true
	}
	// log.Println("xxxxxxxxxxxxxx--gc-session", sessions)
	now := time.Now().Unix()
	for key, value := range sessions {
		time := value.(*sessionFromMemory).lastAccessedTime
		t := time.Unix() + value.(*sessionFromMemory).maxAge
		if t < now {
			// log.Println("timeout------->", value)
			delete(fm.sessions, key)
		}
	}
	return true
}
