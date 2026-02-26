package api

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var (
	userDB *sql.DB
)

func InitUserDatabase() (err error) {
	userDB, err = sql.Open("sqlite", "/sqlite.db")
	return err
}

func UserDB() *sql.DB {
	return userDB
}
