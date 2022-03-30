package controller

import (
	"database/sql"
	"html/template"
	"io"
	"log"
	"net/http"

	"isolati.cn/db"
	"isolati.cn/global"
	"isolati.cn/session"
)

var loginTemplate = template.New("login")

func showLoginPage(w http.ResponseWriter, r *http.Request) {
	loginTemplate.ExecuteTemplate(w, "layout", layoutMsg{
		PageName: "login",
		ContainerData: sliderContainerData{
			LeftSliderData:  global.LEFT_SLIDER,
			RightSliderData: global.RIGHT_SLIDER,
			ContentData:     nil,
		},
	},
	)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		showLoginPage(w, r)
	case http.MethodPost:
		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("Login Request:", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			break
		}
		log.Println("Login Request:", string(data))
		row := db.DB.QueryRow(
			`SELECT md5password FROM admins
			WHERE md5password=?;`,
			string(data),
		)
		var md5password string
		err = row.Scan(&md5password)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Println("Error Password")
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				log.Println("Query Fail:", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			sid := session.UserSession.BeginSession(w, r)
			err = session.UserSession.Set(sid, "identity", "admin")
			if err != nil {
				log.Println("Set Identity Fail:", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				log.Println("Login Success")
				w.WriteHeader(http.StatusOK)
			}
		}
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
}
