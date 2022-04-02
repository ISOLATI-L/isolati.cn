package controller

import (
	"encoding/json"
	"html/template"
	"net/http"

	"isolati.cn/global"
	"isolati.cn/session"
)

var adminTemplate = template.New("admin")

func isRequestAdmin(r *http.Request) (bool, error) {
	vIdentity, err := session.UserSession.GetByRequest(r, "identity")
	if err != nil {
		if err == session.ErrNoCookies {
			return false, nil
		} else {
			return false, err
		}
	} else {
		var identity string
		err = json.Unmarshal(vIdentity, &identity)
		if err != nil {
			return false, err
		} else {
			return identity == "admin", nil
		}
	}
}

func showAdminPage(w http.ResponseWriter, r *http.Request) {
	adminTemplate.ExecuteTemplate(w, "layout", layoutMsg{
		PageName: "admin",
		ContainerData: sliderContainerData{
			LeftSliderData:  global.LEFT_SLIDER,
			RightSliderData: global.RIGHT_SLIDER,
			ContentData:     nil,
		},
	})
}

func handleAdmin(w http.ResponseWriter, r *http.Request) {
	isAdmin, err := isRequestAdmin(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if !isAdmin {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case "/admin", "/admin/":
			showAdminPage(w, r)
		case "/admin/writing", "/admin/writing/":
			showWritingPage(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	case http.MethodConnect:
		sid := session.UserSession.Update(w, r)
		if len(sid) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusOK)
		}
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
	registerWritingRoutes()
	http.HandleFunc("/admin", handleAdmin)
	http.HandleFunc("/admin/", handleAdmin)
}
