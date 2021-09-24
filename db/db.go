package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Unknwon/goconfig"
	"isolati.cn/global"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func init() {
	SQL_Config, err := goconfig.LoadConfigFile(
		global.ROOT_PATH + "SQL.config.ini")
	if err != nil {
		log.Fatalln(err.Error())
	}
	server, err := SQL_Config.GetValue("SQL_Config", "server")
	if err != nil {
		log.Fatalln(err.Error())
	}
	port, err := SQL_Config.GetValue("SQL_Config", "port")
	if err != nil {
		log.Fatalln(err.Error())
	}
	user, err := SQL_Config.GetValue("SQL_Config", "user")
	if err != nil {
		log.Fatalln(err.Error())
	}
	password, err := SQL_Config.GetValue("SQL_Config", "password")
	if err != nil {
		log.Fatalln(err.Error())
	}
	database, err := SQL_Config.GetValue("SQL_Config", "database")
	if err != nil {
		log.Fatalln(err.Error())
	}

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		user, password, server, port, database)
	log.Println(connStr)
	DB, err = sql.Open("mysql", connStr)
	DB.SetConnMaxLifetime(100)
	DB.SetMaxIdleConns(10)
	if err != nil {
		log.Fatalln(err.Error())
	}
	ctx := context.Background()
	err = DB.PingContext(ctx)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("Connected!")
}
