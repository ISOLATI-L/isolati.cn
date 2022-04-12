package session

import (
	"database/sql"
)

type sessionFromDatabase struct {
	sid string
	db  *sql.DB
}

func newSessionFromDatabase(db *sql.DB, sid string, maxAge int64) *sessionFromDatabase {
	transaction, err := db.Begin()
	if err != nil {
		if transaction != nil {
			transaction.Rollback()
		}
		return nil
	}
	_, err = transaction.Exec(
		`INSERT INTO sessions (Sid, SmaxAge, Sdata)
		VALUES (?, ?, JSON_OBJECT());`,
		sid,
		maxAge,
	)
	if err != nil {
		transaction.Rollback()
		return nil
	}
	transaction.Commit()
	return &sessionFromDatabase{
		sid: sid,
		db:  db,
	}
}

func (si *sessionFromDatabase) getID() string {
	return si.sid
}

func (si *sessionFromDatabase) updateLastAccessedTime() {
	transaction, err := si.db.Begin()
	if err != nil {
		if transaction != nil {
			transaction.Rollback()
		}
		return
	}
	_, err = transaction.Exec(
		`UPDATE sessions SET SlastAccessedTime = CURRENT_TIMESTAMP
		WHERE Sid = ?;`,
		si.sid,
	)
	if err != nil {
		transaction.Rollback()
		return
	}
	transaction.Commit()
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
	transaction, err := fd.db.Begin()
	if err != nil {
		if transaction != nil {
			transaction.Rollback()
		}
		return nil
	}
	row := transaction.QueryRow(
		`SELECT Sid FROM sessions
		WHERE Sid = ?;`,
		sid,
	)
	var Sid string
	err = row.Scan(&Sid)
	if err != nil {
		transaction.Rollback()
		return nil
	}
	if Sid != sid {
		transaction.Rollback()
		return nil
	}
	transaction.Commit()
	return &sessionFromDatabase{
		sid: sid,
		db:  fd.db,
	}
}

func (fd *fromDatabase) set(sid string, key string, value any) error {
	key = "$." + key
	transaction, err := fd.db.Begin()
	if err != nil {
		if transaction != nil {
			transaction.Rollback()
		}
		return nil
	}
	_, err = transaction.Exec(
		`UPDATE sessions SET Sdata = JSON_SET(Sdata, ?, ?)
		WHERE Sid = ?;`,
		key,
		value,
		sid,
	)
	if err != nil {
		transaction.Rollback()
		return err
	}
	transaction.Commit()
	return nil
}

func (fd *fromDatabase) get(sid string, key string) ([]byte, error) {
	key = "$." + key
	transaction, err := fd.db.Begin()
	if err != nil {
		if transaction != nil {
			transaction.Rollback()
		}
		return nil, err
	}
	row := transaction.QueryRow(
		`SELECT JSON_EXTRACT(Sdata, ?) FROM sessions
		WHERE Sid = ?;`,
		key,
		sid,
	)
	result := make([]byte, 0)
	err = row.Scan(&result)
	if err != nil {
		transaction.Rollback()
		return nil, err
	}
	transaction.Commit()
	return result, nil
}

func (fd *fromDatabase) remove(sid string, key string) error {
	key = "$." + key
	transaction, err := fd.db.Begin()
	if err != nil {
		if transaction != nil {
			transaction.Rollback()
		}
		return err
	}
	_, err = transaction.Exec(
		`UPDATE sessions SET Sdata = JSON_REMOVE(Sdata, ?)
		WHERE Sid = ?;`,
		key,
		sid,
	)
	if err != nil {
		transaction.Rollback()
		return err
	}
	transaction.Commit()
	return nil
}

func (fd *fromDatabase) update(sid string) error {
	transaction, err := fd.db.Begin()
	if err != nil {
		if transaction != nil {
			transaction.Rollback()
		}
		return err
	}
	_, err = transaction.Exec(
		`UPDATE sessions SET SlastAccessedTime = CURRENT_TIMESTAMP
		WHERE Sid = ?;`,
		sid,
	)
	if err != nil {
		transaction.Rollback()
		return err
	}
	transaction.Commit()
	return nil
}

func (fd *fromDatabase) destroySession(sid string) error {
	transaction, err := fd.db.Begin()
	if err != nil {
		if transaction != nil {
			transaction.Rollback()
		}
		return err
	}
	_, err = transaction.Exec(
		`DELETE FROM sessions
		WHERE Sid = ?;`,
		sid,
	)
	if err != nil {
		transaction.Rollback()
		return err
	}
	transaction.Commit()
	return nil
}

// 已在数据库设置事件自动清除过期sessions
// 无需在此处进行清除工作
func (fd *fromDatabase) gcSession() bool {
	return false
}
