package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"isolati.cn/constant_define"
	"isolati.cn/controller"
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
	// b := make([]byte, 32)
	// for i := 0; i < 100; i++ {
	// 	io.ReadFull(rand.Reader, b)
	// 	log.Println(base64.URLEncoding.EncodeToString(b))
	// }
	var err error
	SQL_Config, err = goconfig.LoadConfigFile(
		constant_define.ROOT_PATH + "SQL.config.ini")
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
	constant_define.DB, err = sql.Open("mysql", connStr)
	constant_define.DB.SetConnMaxLifetime(100)
	constant_define.DB.SetMaxIdleConns(10)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(constant_define.DB)
	ctx := context.Background()
	err = constant_define.DB.PingContext(ctx)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("Connected!")

	constant_define.UserSession = session.NewSessionManager("user", constant_define.DB)

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
		http.FileServer(http.Dir(constant_define.ROOT_PATH+"wwwroot")),
	)
	http.Handle(
		"/js/",
		http.FileServer(http.Dir(constant_define.ROOT_PATH+"wwwroot")),
	)
	http.Handle(
		"/img/",
		http.FileServer(http.Dir(constant_define.ROOT_PATH+"wwwroot")),
	)
	http.Handle(
		"/robots.txt",
		http.FileServer(http.Dir(constant_define.ROOT_PATH+"wwwroot")),
	)
	controller.RegisterRoutes()
	err = server.ListenAndServe()
	// err = server.ListenAndServeTLS(
	// 	constant_define.ROOT_PATH+"isolati.cn.pem",
	// 	constant_define.ROOT_PATH+"isolati.cn.key",
	// )
	if err != nil {
		log.Fatalln(err.Error())
	}
}
