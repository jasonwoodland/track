package db

import (
	"log"
	"strconv"
)

type Setting struct {
	Key   string
	Value string
}

type Settings struct {
	SchemaVersion int
}

const (
	SchemaVersion = "SCHEMA_VERSION"
)

var settings Settings

func updateSetting(key string, value interface{}) {
	Db.Exec(`
		insert into setting (key, value) values(?, ?) on conflict (key) do update set value = excluded.value;
	`, key, value)
}

func getSettings() {
	query, err := Db.Query("select * from setting")
	if err != nil {
		log.Fatal(err)
	}

	for query.Next() {
		var setting Setting
		query.Scan(
			&setting.Key,
			&setting.Value,
		)
		switch setting.Key {
		case SchemaVersion:
			settings.SchemaVersion, _ = strconv.Atoi(setting.Value)
		}
	}
}
