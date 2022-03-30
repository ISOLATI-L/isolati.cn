package session

import (
	"isolati.cn/db"
)

var UserSession *SessionManager

func init() {
	UserSession = NewSessionManager(
		"user",
		db.DB,
		DEFAULT_TIME,
		false,
	)
}
