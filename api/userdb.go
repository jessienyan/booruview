package api

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var (
	userDB *sql.DB
)

func InitUserDatabase() (err error) {
	userDB, err = sql.Open("sqlite", DatabasePath)
	if err != nil {
		return err
	}
	_, err = userDB.Exec("pragma journal_mode = WAL; pragma synchronous = normal")
	return err
}

func UserDB() *sql.DB {
	return userDB
}
