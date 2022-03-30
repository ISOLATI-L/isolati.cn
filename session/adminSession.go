package session

import (
	"isolati.cn/db"
)

var AdminSession *SessionManager

func init() {
	AdminSession = NewSessionManager(
		"user",
		db.DB,
		DEFAULT_TIME,
		false,
	)
}
