package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"isolati.cn/controller"
	"isolati.cn/global"
	"isolati.cn/middleware"
	"isolati.cn/session"

	"github.com/Unknwon/goconfig"
	_ "github.com/go-sql-driver/mysql"
)

type anyType = interface{}
type anyArray = []interface{}
type jsonObject = map[string]anyType

var (
	server   string
	port     string
	user     string
	password string
	database string
)

var SQL_Config *goconfig.ConfigFile

func main() {
	var err error
	// var f *os.File
	// f, err = os.OpenFile(
	// 	global.ROOT_PATH+"isolati.cn.log",
	// 	os.O_CREATE|os.O_APPEND|os.O_WRONLY,
	// 	0666,
	// )
	// if err != nil {
	// 	log.Fatalln(err.Error())
	// }
	// log.SetOutput(f)

	// b := make([]byte, 32)
	// for i := 0; i < 100; i++ {
	// 	if _, err := io.ReadFull(rand.Reader, b); err != nil {
	// 		return
	// 	}
	// 	randId := fmt.Sprintf("%x", md5.Sum(b))
	// 	log.Println(randId)
	// }
	SQL_Config, err = goconfig.LoadConfigFile(
		global.ROOT_PATH + "SQL.config.ini")
	if err != nil {
		log.Fatalln(err.Error())
	}
	server, err = SQL_Config.GetValue("SQL_Config", "server")
	if err != nil {
		log.Fatalln(err.Error())
	}
	port, err = SQL_Config.GetValue("SQL_Config", "port")
	if err != nil {
		log.Fatalln(err.Error())
	}
	user, err = SQL_Config.GetValue("SQL_Config", "user")
	if err != nil {
		log.Fatalln(err.Error())
	}
	password, err = SQL_Config.GetValue("SQL_Config", "password")
	if err != nil {
		log.Fatalln(err.Error())
	}
	database, err = SQL_Config.GetValue("SQL_Config", "database")
	if err != nil {
		log.Fatalln(err.Error())
	}

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		user, password, server, port, database)
	log.Println(connStr)
	global.DB, err = sql.Open("mysql", connStr)
	global.DB.SetConnMaxLifetime(100)
	global.DB.SetMaxIdleConns(10)
	if err != nil {
		log.Fatalln(err.Error())
	}
	ctx := context.Background()
	err = global.DB.PingContext(ctx)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("Connected!")

	global.UserSession = session.NewSessionManager(
		"user",
		global.DB,
		session.DEFAULT_TIME,
		false,
	)

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
	http.Handle(
		"/robots.txt",
		http.FileServer(http.Dir(global.ROOT_PATH+"wwwroot")),
	)
	controller.RegisterRoutes()
	err = server.ListenAndServe()
	// err = server.ListenAndServeTLS(
	// 	global.ROOT_PATH+"isolati.cn.pem",
	// 	global.ROOT_PATH+"isolati.cn.key",
	// )
	if err != nil {
		log.Fatalln(err.Error())
	}
}
