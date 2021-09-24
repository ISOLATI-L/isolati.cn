package main

import (
	"log"
	"net/http"

	"isolati.cn/controller"
	"isolati.cn/global"
	"isolati.cn/middleware"
)

func main() {
	changePrefixMiddleware := middleware.NewChangePrefixMiddleware(
		nil,
		"isolati.cn",
		"www.",
		"",
	)
	timeoutMiddleware := middleware.NewTimeoutMiddleware(
		&changePrefixMiddleware,
	)
	server := http.Server{
		// Addr: ":8080",
		Addr:    "localhost:8080",
		Handler: &timeoutMiddleware,
	}

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
	controller.RegisterRoutes()

	err := server.ListenAndServe()
	// err = server.ListenAndServeTLS(
	// 	global.ROOT_PATH+"isolati.cn.pem",
	// 	global.ROOT_PATH+"isolati.cn.key",
	// )
	if err != nil {
		log.Fatalln(err.Error())
	}
}
