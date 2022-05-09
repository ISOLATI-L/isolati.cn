package controller

import (
	"html/template"
	"net/http"

	"isolati.cn/global"
)

var aboutTemplate = template.New("about")

func showAboutPage(w http.ResponseWriter, r *http.Request) {
	aboutTemplate.ExecuteTemplate(w, "layout", layoutMsg{
		CssFiles: []string{"/css/sliderContainer.css", "/css/about.css"},
		JsFiles:  []string{},
		PageName: "about",
		ContainerData: sliderContainerData{
			ContentData: nil,
		},
	})
}

func handleAbout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		showAboutPage(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func registerAboutRoutes() {
	template.Must(
		aboutTemplate.ParseFiles(
			global.ROOT_PATH+"/wwwroot/layout.html",
			global.ROOT_PATH+"/wwwroot/sliderContainer.html",
			global.ROOT_PATH+"/wwwroot/about.html",
		),
	)
	http.HandleFunc("/about", handleAbout)
}
