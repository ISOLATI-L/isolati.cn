package controller

import (
	"html/template"
	"net/http"

	"isolati.cn/global"
)

var homeTemplate = template.New("home")

func handleHome(w http.ResponseWriter, r *http.Request) {
	homeTemplate.ExecuteTemplate(w, "layout", layoutMsg{
		PageName: "home",
		ContainerData: sliderContainerData{
			LeftSliderData:  global.LEFT_SLIDER,
			RightSliderData: global.RIGHT_SLIDER,
			ContentData:     nil,
		},
	})
}

func registerHomeRoutes() {
	template.Must(
		homeTemplate.ParseFiles(
			global.ROOT_PATH+"/wwwroot/layout.html",
			global.ROOT_PATH+"/wwwroot/sliderContainer.html",
			global.ROOT_PATH+"/wwwroot/home.html",
		),
	)
	http.HandleFunc("/home", handleHome)
}
