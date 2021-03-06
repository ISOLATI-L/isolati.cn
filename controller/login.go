package controller

import (
	"database/sql"
	"encoding/base64"
	"html/template"
	"io"
	"log"
	"net/http"
	"regexp"

	"isolati.cn/db"
	"isolati.cn/global"
	"isolati.cn/session"
)

var loginTemplate = template.New("login")

var refSelector *regexp.Regexp

func showLoginPage(w http.ResponseWriter, r *http.Request) {
	isAdmin, _ := isRequestAdmin(r)
	if isAdmin {
		url := "/admin"
		matches := refSelector.FindStringSubmatch(r.URL.RawQuery)
		if len(matches) > 4 {
			buf, err := base64.URLEncoding.DecodeString(matches[2])
			if err == nil {
				url = string(buf)
			}
		}
		http.Redirect(w, r, url, http.StatusSeeOther)
	} else {
		loginTemplate.ExecuteTemplate(w, "layout", layoutMsg{
			CssFiles: []string{"/css/sliderContainer.css", "/css/login.css"},
			JsFiles:  []string{"/js/login.js", "/js/hash.js", "/js/post.js"},
			PageName: "login",
			ContainerData: sliderContainerData{
				ContentData: nil,
			},
		},
		)
	}
}

func postLoginRequest(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Login Request:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("Login Request:", string(data))
	transaction, err := db.DB.Begin()
	if err != nil {
		if transaction != nil {
			transaction.Rollback()
		}
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	row := transaction.QueryRow(
		`SELECT md5password FROM admins
		WHERE md5password=?;`,
		string(data),
	)
	var md5password string
	err = row.Scan(&md5password)
	if err != nil {
		if err == sql.ErrNoRows {
			transaction.Rollback()
			log.Println("Error Password")
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			transaction.Rollback()
			log.Println("Query Fail:", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
		transaction.Rollback()
		return
	}
	transaction.Commit()

	transaction, err = session.UserSession.BeginTransaction()
	if err != nil {
		log.Println("Begin Transaction Fail:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		transaction.Rollback()
	}
	sid, err := session.UserSession.BeginSession(transaction, w, r)
	if err != nil {
		transaction.Rollback()
		log.Println("Set Identity Fail:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = session.UserSession.Set(transaction, sid, "identity", "admin")
	if err != nil {
		transaction.Rollback()
		log.Println("Set Identity Fail:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	transaction.Commit()
	log.Println("Login Success")
	w.WriteHeader(http.StatusOK)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		showLoginPage(w, r)
	case http.MethodPost:
		postLoginRequest(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func registerLoginRoutes() {
	template.Must(
		loginTemplate.ParseFiles(
			global.ROOT_PATH+"/wwwroot/layout.html",
			global.ROOT_PATH+"/wwwroot/sliderContainer.html",
			global.ROOT_PATH+"/wwwroot/login.html",
		),
	)
	http.HandleFunc("/login", handleLogin)

	var err error
	refSelector, err = regexp.Compile(
		"ref=(\"|%22)(.*)(\"|%22)(&|$)",
	)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
