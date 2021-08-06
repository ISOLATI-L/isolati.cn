package controller

import (
	"WEB_ISOLATI/constant_define"
	"html/template"
	"net/http"
)

var aboutTemplate = template.New("about")

func handleAbout(w http.ResponseWriter, r *http.Request) {
	aboutTemplate.ExecuteTemplate(w, "layout", layoutMsg{
		PageName: "about",
		ContainerData: sliderContainerData{
			LeftSliderData:  constant_define.LEFT_SLIDER,
			RightSliderData: constant_define.RIGHT_SLIDER,
			ContentData:     "阿巴阿巴",
		},
	})
}

func registerAboutRoutes() {
	template.Must(
		aboutTemplate.ParseFiles(
			constant_define.ROOT_PATH+"/wwwroot/layout.html",
			constant_define.ROOT_PATH+"/wwwroot/sliderContainer.html",
			constant_define.ROOT_PATH+"/wwwroot/about.html",
		),
	)
	http.HandleFunc("/about", handleAbout)
}
