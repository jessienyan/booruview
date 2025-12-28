package api

import "database/sql"

var (
	userDB *sql.DB
)

func InitUserDatabase() (err error) {
	userDB, err = sql.Open("sqlite3", "/database/sqlite.db")
	return err
}

func UserDB() *sql.DB {
	return userDB
}
