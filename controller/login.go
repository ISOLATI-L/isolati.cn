package controller

import (
	"html/template"
	"net/http"

	"isolati.cn/global"
)

var loginTemplate = template.New("login")

func showLoginPage(w http.ResponseWriter, r *http.Request) {
	loginTemplate.ExecuteTemplate(w, "layout", layoutMsg{
		PageName: "home",
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
