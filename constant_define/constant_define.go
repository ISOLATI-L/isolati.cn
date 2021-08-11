package constant_define

import (
	"database/sql"
	"html/template"

	"isolati.cn/session"
)

const ROOT_PATH = "./"
const SHARE_FILES_PATH = "D:/ISOLATI/"

const LEFT_SLIDER = template.HTML(`
<p>左边栏</p>
<p>TO DO</p>
<p>TO DO</p>
<p>TO DO</p>
<p>TO DO</p>
<p>TO DO</p>
<p>TO DO</p>
<p>TO DO</p>
<p>TO DO</p>
<p>TO DO</p>
<p>TO DO</p>
`)
const RIGHT_SLIDER = template.HTML(`
<p>右边栏</p>
<p>TO DO</p>
<p>TO DO</p>
<p>TO DO</p>
<p>TO DO</p>
<p>TO DO</p>
<p>TO DO</p>
<p>TO DO</p>
<p>TO DO</p>
<p>TO DO</p>
<p>TO DO</p>
`)

var DB *sql.DB

var UserSession *session.SessionManager
