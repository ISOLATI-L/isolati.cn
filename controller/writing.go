package controller

import (
	"html/template"
	"net/http"

	"isolati.cn/global"
)

var writingTemplate = template.New("writing")

func showWritingPage(w http.ResponseWriter, r *http.Request) {
	writingTemplate.ExecuteTemplate(w, "layout", layoutMsg{
		PageName: "writing",
		ContainerData: sliderContainerData{
			LeftSliderData:  global.LEFT_SLIDER,
			RightSliderData: global.RIGHT_SLIDER,
			ContentData:     nil,
		},
	})
}

func handleWriting(w http.ResponseWriter, r *http.Request) {
	isAdmin, err := isRequestAdmin(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else if !isAdmin {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		switch r.Method {
		case http.MethodGet:
			showWritingPage(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func registerWritingRoutes() {
	template.Must(
		writingTemplate.ParseFiles(
			global.ROOT_PATH+"/wwwroot/layout.html",
			global.ROOT_PATH+"/wwwroot/sliderContainer.html",
			global.ROOT_PATH+"/wwwroot/writing.html",
		),
	)
	http.HandleFunc("/writing", handleWriting)
}
