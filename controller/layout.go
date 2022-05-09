package controller

import (
	"log"
	"net/http"
	"regexp"

	"isolati.cn/global"
)

type sliderContainerData struct {
	ContentData any
}

type layoutMsg struct {
	CssFiles      []string
	JsFiles       []string
	PageName      string
	ContainerData any
}

var pattern *regexp.Regexp

func handlerInit() {
	var err error
	pattern, err = regexp.Compile(`/(.+)$`)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func redirectToHome(w http.ResponseWriter, r *http.Request) {
	matches := pattern.FindStringSubmatch(r.URL.Path)
	if len(matches) > 0 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		// http.Redirect(w, r, "/home", http.StatusMovedPermanently)
		http.Redirect(w, r, "/home", http.StatusFound)
	}
}

func RegisterRoutes() {
	http.Handle(
		"/css/",
		http.FileServer(http.Dir(global.ROOT_PATH+"wwwroot")),
	)
	http.Handle(
		"/js/",
		http.FileServer(http.Dir(global.ROOT_PATH+"wwwroot")),
	)
	http.Handle(
		"/img/",
		http.FileServer(http.Dir(global.ROOT_PATH+"wwwroot")),
	)
	http.Handle(
		"/jquery/",
		http.FileServer(http.Dir(global.ROOT_PATH+"wwwroot")),
	)
	http.Handle(
		"/editormd/",
		http.FileServer(http.Dir(global.ROOT_PATH+"wwwroot")),
	)
	http.HandleFunc("/", redirectToHome)
	handlerInit()
	registerHomeRoutes()
	registerAboutRoutes()
	registerFilesRoutes()
	registerVideosRoutes()
	registerRobotsRoutes()
	registerLoginRoutes()
	registerAdminRoutes()
	registerParagraphsRoutes()
	registerImagesRoutes()
}
