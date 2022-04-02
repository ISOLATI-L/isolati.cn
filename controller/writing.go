package controller

import (
	"html/template"
	"net/http"

	"isolati.cn/global"
)

var writingTemplate = template.New("writing")

func showWritingPage(w http.ResponseWriter, r *http.Request) {
	writingTemplate.ExecuteTemplate(w, "layout", layoutMsg{
		CssFiles:      []string{"/css/writing.css", "/editormd/css/editormd.min.css"},
		JsFiles:       []string{"/jquery/jquery-3.6.0.min.js", "/editormd/editormd.min.js", "/js/writing.js"},
		PageName:      "writing",
		ContainerData: nil,
	})
}

func registerWritingRoutes() {
	template.Must(
		writingTemplate.ParseFiles(
			global.ROOT_PATH+"/wwwroot/layout.html",
			// global.ROOT_PATH+"/wwwroot/sliderContainer.html",
			global.ROOT_PATH+"/wwwroot/writing.html",
		),
	)
}
