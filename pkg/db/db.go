package db

import (
	"database/sql"

	"github.com/adrg/xdg"
	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

func OpenDb() {
	dbFilePath, _ := xdg.DataFile("track-cli/db.sqlite3")
	Db, _ = sql.Open("sqlite3", dbFilePath)
	getSettings()
	migrateDb()
}
