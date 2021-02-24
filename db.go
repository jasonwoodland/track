package main

import (
	"database/sql"

	"github.com/adrg/xdg"
	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

func openDb() {
	dbFilePath, _ := xdg.DataFile("track-cli/db.sqlite3")
	Db, _ = sql.Open("sqlite3", dbFilePath)
	initDb()
}

func initDb() {
	Db.Exec(`
		create table if not exists project (
			id integer primary key,
			name text
		);
	`)

	Db.Exec("pragma foreign_keys = on")

	Db.Exec(`
		create table if not exists task (
			id integer primary key,
			project_id integer,
			name text,

			foreign key(project_id) references project(id) on delete cascade
		);
	`)
	Db.Exec(`
		create table if not exists frame (
			id integer primary key,
			task_id integer,
			start_time text,
			end_time text,

			foreign key(task_id) references task(id) on delete cascade
		);
	`)
}
