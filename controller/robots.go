package controller

import (
	"log"
	"net/http"

	"isolati.cn/db"
	"isolati.cn/global"
)

var robotsHandler = http.FileServer(http.Dir(global.ROOT_PATH + "wwwroot"))

func logReferer(r *http.Request) {
	userAgent := r.UserAgent()
	log.Println("Visited robots.txt: ", userAgent)
	transaction, err := db.DB.Begin()
	if err != nil {
		if transaction != nil {
			transaction.Rollback()
		}
		log.Println(err.Error())
		return
	}
	_, err = transaction.Exec(
		`INSERT INTO robots (RuserAgent) VALUES (?);`,
		userAgent,
	)
	if err != nil {
		transaction.Rollback()
		log.Println(err.Error())
		return
	}
	transaction.Commit()
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
