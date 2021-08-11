package session

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"
)

const TABLE string = "sessions"

type SessionFromDatabase struct {
	sid   string
	lock  sync.Mutex
	db    *sql.DB
	table string
}

func newSessionFromDatabase(db *sql.DB, table string, sid string) *SessionFromDatabase {
	result, err := db.Exec(
		`INSERT INTO ? (Sid, Sdata) VALUES (?, JSON_OBJECT());`,
		table,
		sid,
	)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	affected, err := result.RowsAffected()
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	if affected == 0 {
		log.Println(result)
		return nil
	}
	return &SessionFromDatabase{
		sid:   sid,
		db:    db,
		table: table,
	}
}

func (si *SessionFromDatabase) Set(key interface{}, value interface{}) {
	key = fmt.Sprintf("$.%s", key)
	si.lock.Lock()
	defer si.lock.Unlock()
	result, err := si.db.Exec(
		`UPDATE ? SET Sdata = JSON_SET(Sdata, ?, ?)
		WHERE Sid = ?;`,
		si.table,
		key,
		value,
		si.sid,
	)
	if err != nil {
		log.Println(err.Error())
		return
	}
	affected, err := result.RowsAffected()
	if err != nil {
		log.Println(err.Error())
		return
	}
	if affected == 0 {
		log.Println(result)
		return
	}
}

func (si *SessionFromDatabase) Get(key interface{}) interface{} {
	key = fmt.Sprintf("$.%s", key)
	row := si.db.QueryRow(
		`SELECT JSON_EXTRACT(Sdata, ?) FROM ?
		WHERE Sid = ?;`,
		key,
		si.table,
		si.sid,
	)
	var result interface{}
	err := row.Scan(&result)
	if err != nil {
		log.Println(result)
		return nil
	}
	return result
}

func (si *SessionFromDatabase) Remove(key interface{}) error {
	key = fmt.Sprintf("$.%s", key)
	result, err := si.db.Exec(
		`UPDATE ? SET Sdata = JSON_REMOVE(Sdata, ?)
		WHERE Sid = ?;`,
		si.table,
		key,
		si.sid,
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	if affected == 0 {
		log.Println(result)
	}
	return nil
}

func (si *SessionFromDatabase) GetLastAccessedTime() (time.Time, error) {
	row := si.db.QueryRow(
		`SELECT SlastAccessedTime FROM ?
		WHERE Sid = ?;`,
		si.table,
		si.sid,
	)
	var timeStr string
	err := row.Scan(
		&timeStr,
	)
	if err != nil {
		log.Println(err.Error())
		return time.Now(), err
	}
	return time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local)
}

func (si *SessionFromDatabase) UpdateLastAccessedTime() {
	result, err := si.db.Exec(
		`UPDATE ? SET SlastAccessedTime = CURRENT_TIMESTAMP
		WHERE Sid = ?;`,
		si.table,
		si.sid,
	)
	if err != nil {
		log.Println(err.Error())
		return
	}
	affected, err := result.RowsAffected()
	if err != nil {
		log.Println(err.Error())
		return
	}
	if affected == 0 {
		log.Println(result)
	}
}

func (si *SessionFromDatabase) GetMaxAge() int64 {
	row := si.db.QueryRow(
		`SELECT SmaxAge FROM ?
		WHERE Sid = ?;`,
		si.table,
		si.sid,
	)
	var maxAge int64
	err := row.Scan(
		&maxAge,
	)
	if err != nil {
		log.Println(err.Error())
		return 0
	}
	return maxAge
}

func (si *SessionFromDatabase) SetMaxAge(age int64) {
	result, err := si.db.Exec(
		`UPDATE ? SET SmaxAge = ?
		WHERE Sid = ?;`,
		si.table,
		age,
		si.sid,
	)
	if err != nil {
		log.Println(err.Error())
		return
	}
	affected, err := result.RowsAffected()
	if err != nil {
		log.Println(err.Error())
		return
	}
	if affected == 0 {
		log.Println(result)
	}
}

func (si *SessionFromDatabase) GetId() string {
	return si.sid
}

func (si *SessionFromDatabase) Destroy() bool {
	result, err := si.db.Exec(
		`DELETE FROM ?
		WHERE Sid = ?;`,
		si.table,
		si.sid,
	)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	affected, err := result.RowsAffected()
	if err != nil {
		log.Println(err.Error())
		return false
	}
	if affected == 0 {
		log.Println(result)
	}
	return true
}

type FromDatabase struct {
	lock sync.Mutex
	db   *sql.DB
}

func newFromDatabase(db *sql.DB) *FromDatabase {
	return &FromDatabase{
		db: db,
	}
}

func (fd *FromDatabase) InitSession(sid string, maxAge int64) (Session, error) {
	fd.lock.Lock()
	defer fd.lock.Unlock()
	newSession := newSessionFromDatabase(fd.db, TABLE, sid)
	if maxAge != 0 && maxAge != DEFAULT_TIME {
		newSession.SetMaxAge(maxAge)
	}
	log.Println(newSession)
	return newSession, nil
}

func (fd *FromDatabase) GetSession(sid string) Session {
	return &SessionFromDatabase{
		sid:   sid,
		db:    fd.db,
		table: TABLE,
	}
}

func (fd *FromDatabase) DestroySession(sid string) error {
	fd.GetSession(sid).Destroy()
	return nil
}

func (fd *FromDatabase) GCSession() bool {
	return false
}
