package controller

import (
	"html/template"
	"net/http"

	"isolati.cn/global"
)

var aboutTemplate = template.New("about")

func handleAbout(w http.ResponseWriter, r *http.Request) {
	aboutTemplate.ExecuteTemplate(w, "layout", layoutMsg{
		PageName: "about",
		ContainerData: sliderContainerData{
			LeftSliderData:  global.LEFT_SLIDER,
			RightSliderData: global.RIGHT_SLIDER,
			ContentData:     nil,
		},
	})
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
