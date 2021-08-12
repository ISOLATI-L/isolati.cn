package controller

import (
	"log"
	"net/http"

	"isolati.cn/global"
)

var robotsHandler = http.FileServer(http.Dir(global.ROOT_PATH + "wwwroot"))

func logReferer(r *http.Request) {
	userAgent := r.UserAgent()
	log.Println("Visited robots.txt: ", userAgent)
	result, err := global.DB.Exec(
		`INSERT INTO robots (RuserAgent) VALUES (?);`,
		userAgent,
	)
	if err != nil {
		log.Println(err.Error())
		return
	}
	affected, err := result.RowsAffected()
	if err != nil {
		log.Println(err.Error())
		return
	}
	if affected == 0 {
		log.Println(result)
	}
}

func handleRobots(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		go logReferer(r)
		robotsHandler.ServeHTTP(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func registerRobotsRoutes() {
	http.HandleFunc("/robots.txt", handleRobots)
}
