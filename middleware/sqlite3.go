package middleware

import (
	"database/sql"
	"github.com/laouji/jordgubbe/config"
	_ "github.com/mattn/go-sqlite3"
)

var (
	dbh *DBH
)

type DBH struct {
	*sql.DB
}

func Init() {
	conf := config.LoadConfig()

	handle, err := sql.Open("sqlite3", conf.DBPath)
	if err != nil {
		panic(err)
	}

	dbh = &DBH{handle}
}

func GetDBH() *DBH {
	if dbh == nil {
		Init()
	}
	return dbh
}
