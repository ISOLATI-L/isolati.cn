package session

import (
	"database/sql"
)

type sessionFromDatabase struct {
	sid string
}

func newSessionFromDatabase(transaction *sql.Tx, sid string, maxAge int64) (*sessionFromDatabase, error) {
	_, err := transaction.Exec(
		`INSERT INTO sessions (Sid, SmaxAge, Sdata)
		VALUES (?, ?, JSON_OBJECT());`,
		sid,
		maxAge,
	)
	if err != nil {
		return nil, err
	}
	return &sessionFromDatabase{
		sid: sid,
	}, nil
}

func (si *sessionFromDatabase) getID() string {
	return si.sid
}

func (si *sessionFromDatabase) updateLastAccessedTime(transaction *sql.Tx) error {
	_, err := transaction.Exec(
		`UPDATE sessions SET SlastAccessedTime = CURRENT_TIMESTAMP
		WHERE Sid = ?;`,
		si.sid,
	)
	if err != nil {
		return err
	}
	return nil
}

type fromDatabase struct{}

func newFromDatabase() *fromDatabase {
	return &fromDatabase{}
}

func (fd *fromDatabase) initSession(transaction *sql.Tx, sid string, maxAge int64) (session, error) {
	newSession, err := newSessionFromDatabase(transaction, sid, maxAge)
	return newSession, err
}

func (fd *fromDatabase) getSession(transaction *sql.Tx, sid string) (session, error) {
	row := transaction.QueryRow(
		`SELECT Sid FROM sessions
		WHERE Sid = ?;`,
		sid,
	)
	var Sid string
	err := row.Scan(&Sid)
	if err != nil {
		return nil, err
	}
	if Sid != sid {
		return nil, err
	}
	return &sessionFromDatabase{
		sid: sid,
	}, nil
}

func (fd *fromDatabase) set(transaction *sql.Tx, sid string, key string, value any) error {
	key = "$." + key
	_, err := transaction.Exec(
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

func (fd *fromDatabase) get(transaction *sql.Tx, sid string, key string) ([]byte, error) {
	key = "$." + key
	row := transaction.QueryRow(
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

func (fd *fromDatabase) remove(transaction *sql.Tx, sid string, key string) error {
	key = "$." + key
	_, err := transaction.Exec(
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

func (fd *fromDatabase) update(transaction *sql.Tx, sid string) error {
	_, err := transaction.Exec(
		`UPDATE sessions SET SlastAccessedTime = CURRENT_TIMESTAMP
		WHERE Sid = ?;`,
		sid,
	)
	if err != nil {
		return err
	}
	return nil
}

func (fd *fromDatabase) destroySession(transaction *sql.Tx, sid string) error {
	_, err := transaction.Exec(
		`DELETE FROM sessions
		WHERE Sid = ?;`,
		sid,
	)
	if err != nil {
		return err
	}
	return nil
}
