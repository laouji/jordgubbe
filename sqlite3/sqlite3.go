package sqlite3

import (
	"../config"
	"database/sql"
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

	driver, err := sql.Open("sqlite3", conf.DBPath)
	if err != nil {
		panic(err)
	}

	dbh = &DBH{driver}
}

func GetDBH() *DBH {
	return dbh
}

func (dbh *DBH) LatestId(tableName string) int {
	row := dbh.QueryRow(`SELECT id FROM ` + tableName + ` ORDER BY id DESC LIMIT 1`)

	var id int
	err := row.Scan(&id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return 0
		}
		panic(err)
	}

	return id
}
