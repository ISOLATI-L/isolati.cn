package session

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
)

type SessionFromDatabase struct {
	sid string
	db  *sql.DB
}

func newSessionFromDatabase(db *sql.DB, sid string, maxAge int64) *SessionFromDatabase {
	result, err := db.Exec(
		`INSERT INTO sessions (Sid, SmaxAge, Sdata)
		VALUES (?, ?, JSON_OBJECT());`,
		sid,
		maxAge,
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
		sid: sid,
		db:  db,
	}
}

func (si *SessionFromDatabase) UpdateLastAccessedTime() {
	result, err := si.db.Exec(
		`UPDATE sessions SET SlastAccessedTime = CURRENT_TIMESTAMP
		WHERE Sid = ?;`,
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
	newSession := newSessionFromDatabase(fd.db, sid, maxAge)
	// log.Println(newSession)
	return newSession, nil
}

func (fd *FromDatabase) GetSession(sid string) Session {
	row := fd.db.QueryRow(
		`SELECT Sid FROM sessions
		WHERE Sid = ?;`,
		sid,
	)
	var Sid string
	err := row.Scan(&Sid)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	if Sid == sid {
		return &SessionFromDatabase{
			sid: sid,
			db:  fd.db,
		}
	} else {
		return nil
	}
}

func (fd *FromDatabase) Set(sid string, key string, value interface{}) error {
	key = fmt.Sprintf("$.%s", key)
	result, err := fd.db.Exec(
		`UPDATE sessions SET Sdata = JSON_SET(Sdata, ?, ?)
		WHERE Sid = ?;`,
		key,
		value,
		sid,
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
		return errors.New("affected 0 rows")
	}
	return nil
}

func (fd *FromDatabase) Get(sid string, key string) interface{} {
	key = fmt.Sprintf("$.%s", key)
	row := fd.db.QueryRow(
		`SELECT JSON_EXTRACT(Sdata, ?) FROM sessions
		WHERE Sid = ?;`,
		key,
		sid,
	)
	var result interface{}
	err := row.Scan(&result)
	if err != nil {
		log.Println(result)
		return nil
	}
	return result
}

func (fd *FromDatabase) Remove(sid string, key string) error {
	key = fmt.Sprintf("$.%s", key)
	result, err := fd.db.Exec(
		`UPDATE sessions SET Sdata = JSON_REMOVE(Sdata, ?)
		WHERE Sid = ?;`,
		key,
		sid,
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
		return errors.New("affected 0 rows")
	}
	return nil
}

func (fd *FromDatabase) DestroySession(sid string) error {
	result, err := fd.db.Exec(
		`DELETE FROM sessions
		WHERE Sid = ?;`,
		sid,
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
		return errors.New("affected 0 rows")
	}
	return nil
}

// 已在数据库设置事件自动清除过期sessions
// 无需在此处进行清除工作
func (fd *FromDatabase) GCSession() bool {
	return false
}
