package controller

import (
	"encoding/json"
	"html/template"
	"net/http"

	"isolati.cn/global"
	"isolati.cn/session"
)

var adminTemplate = template.New("admin")

func showAdminPage(w http.ResponseWriter, r *http.Request) {
	vIdentity, err := session.UserSession.GetByRequest(r, "identity")
	if err != nil {
		if err == session.ErrNoCookies {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		var identity string
		err = json.Unmarshal(vIdentity, &identity)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			if identity != "admin" {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			} else {
				adminTemplate.ExecuteTemplate(w, "layout", layoutMsg{
					PageName: "admin",
					ContainerData: sliderContainerData{
						LeftSliderData:  global.LEFT_SLIDER,
						RightSliderData: global.RIGHT_SLIDER,
						ContentData:     nil,
					},
				})
			}
		}
	}
}

func handleAdmin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		showAdminPage(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func registerAdminRoutes() {
	template.Must(
		adminTemplate.ParseFiles(
			global.ROOT_PATH+"/wwwroot/layout.html",
			global.ROOT_PATH+"/wwwroot/sliderContainer.html",
			global.ROOT_PATH+"/wwwroot/admin.html",
		),
	)
	http.HandleFunc("/admin", handleAdmin)
}
