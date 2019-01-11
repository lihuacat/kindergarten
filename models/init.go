package models

import (
	"database/sql"
	"time"

	"github.com/astaxie/beego"
	log "github.com/astaxie/beego/logs"
	_ "github.com/bmizerany/pq"
)

var (
	//	database  string
	//	dbConnStr string
	db *sql.DB
)

func init() {
	var err error
	timeout, err := beego.AppConfig.Int64("session_timeout")
	if err != nil {
		log.Error(err)
		timeout = 1200
	}
	sessionTimeout = time.Duration(timeout) * time.Second
	dbName := beego.AppConfig.String("database")
	dbConnStr := beego.AppConfig.String("db_connect")

	db, err = sql.Open(dbName, dbConnStr)
	if err != nil {
		log.Error("Open DB error:", err)
		return
	}
	initDevTimer()
}
