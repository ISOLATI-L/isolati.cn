package controller

import (
	"log"
	"net/http"
	"regexp"
)

type sliderContainerData struct {
	LeftSliderData  interface{}
	RightSliderData interface{}
	ContentData     interface{}
}

type layoutMsg struct {
	PageName      string
	ContainerData interface{}
}

var pattern *regexp.Regexp

func handlerInit() {
	var err error
	pattern, err = regexp.Compile(`/(.+)$`)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func handleToHome(w http.ResponseWriter, r *http.Request) {
	matches := pattern.FindStringSubmatch(r.URL.Path)
	if len(matches) > 0 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		// http.Redirect(w, r, "/home", http.StatusMovedPermanently)
		http.Redirect(w, r, "/home", http.StatusFound)
	}
}

func RegisterRoutes() {
	http.HandleFunc("/", handleToHome)
	handlerInit()
	registerHomeRoutes()
	registerAboutRoutes()
	registerFilesRoutes()
	registerVideosRoutes()
}
