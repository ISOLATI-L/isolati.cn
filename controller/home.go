package controller

import (
	"WEB_ISOLATI/constant_define"
	"html/template"
	"net/http"
)

var homeTemplate = template.New("home")

func handleHome(w http.ResponseWriter, r *http.Request) {
	homeTemplate.ExecuteTemplate(w, "layout", layoutMsg{
		PageName: "home",
		ContainerData: sliderContainerData{
			LeftSliderData:  constant_define.LEFT_SLIDER,
			RightSliderData: constant_define.RIGHT_SLIDER,
			ContentData:     "阿巴阿巴",
		},
	})
}

func registerHomeRoutes() {
	template.Must(
		homeTemplate.ParseFiles(
			constant_define.ROOT_PATH+"/wwwroot/layout.html",
			constant_define.ROOT_PATH+"/wwwroot/sliderContainer.html",
			constant_define.ROOT_PATH+"/wwwroot/home.html",
		),
	)
	http.HandleFunc("/home", handleHome)
}
