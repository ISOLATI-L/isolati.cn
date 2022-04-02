package session

import (
	"database/sql"
)

type sessionFromDatabase struct {
	sid string
	db  *sql.DB
}

func newSessionFromDatabase(db *sql.DB, sid string, maxAge int64) *sessionFromDatabase {
	result, err := db.Exec(
		`INSERT INTO sessions (Sid, SmaxAge, Sdata)
		VALUES (?, ?, JSON_OBJECT());`,
		sid,
		maxAge,
	)
	if err != nil {
		return nil
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return nil
	}
	if affected == 0 {
		return nil
	}
	return &sessionFromDatabase{
		sid: sid,
		db:  db,
	}
}

func (si *sessionFromDatabase) getID() string {
	return si.sid
}

func (si *sessionFromDatabase) updateLastAccessedTime() {
	_, err := si.db.Exec(
		`UPDATE sessions SET SlastAccessedTime = CURRENT_TIMESTAMP
		WHERE Sid = ?;`,
		si.sid,
	)
	if err != nil {
		return
	}
}

type fromDatabase struct {
	db *sql.DB
}

func newFromDatabase(db *sql.DB) *fromDatabase {
	return &fromDatabase{
		db: db,
	}
}

func (fd *fromDatabase) initSession(sid string, maxAge int64) (session, error) {
	newSession := newSessionFromDatabase(fd.db, sid, maxAge)
	return newSession, nil
}

func (fd *fromDatabase) getSession(sid string) session {
	row := fd.db.QueryRow(
		`SELECT Sid FROM sessions
		WHERE Sid = ?;`,
		sid,
	)
	var Sid string
	err := row.Scan(&Sid)
	if err != nil {
		return nil
	}
	if Sid == sid {
		return &sessionFromDatabase{
			sid: sid,
			db:  fd.db,
		}
	} else {
		return nil
	}
}

func (fd *fromDatabase) set(sid string, key string, value any) error {
	key = "$." + key
	_, err := fd.db.Exec(
		`UPDATE sessions SET Sdata = JSON_SET(Sdata, ?, ?)
		WHERE Sid = ?;`,
		key,
		value,
		sid,
	)
	if err != nil {
		return err
	}
	return nil
}

func (fd *fromDatabase) get(sid string, key string) ([]byte, error) {
	key = "$." + key
	row := fd.db.QueryRow(
		`SELECT JSON_EXTRACT(Sdata, ?) FROM sessions
		WHERE Sid = ?;`,
		key,
		sid,
	)
	result := make([]byte, 0)
	err := row.Scan(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *fromDatabase) remove(sid string, key string) error {
	key = "$." + key
	_, err := fd.db.Exec(
		`UPDATE sessions SET Sdata = JSON_REMOVE(Sdata, ?)
		WHERE Sid = ?;`,
		key,
		sid,
	)
	if err != nil {
		return err
	}
	return nil
}

func (fd *fromDatabase) update(sid string) {
	fd.db.Exec(
		`UPDATE sessions SET SlastAccessedTime = CURRENT_TIMESTAMP
		WHERE Sid = ?;`,
		sid,
	)
}

func (fd *fromDatabase) destroySession(sid string) error {
	_, err := fd.db.Exec(
		`DELETE FROM sessions
		WHERE Sid = ?;`,
		sid,
	)
	if err != nil {
		return err
	}
	return nil
}

// 已在数据库设置事件自动清除过期sessions
// 无需在此处进行清除工作
func (fd *fromDatabase) gcSession() bool {
	return false
}
