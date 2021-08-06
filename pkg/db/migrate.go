package db

type Migration struct {
	Version int
	Up      func()
}

var migrations = []Migration{
	{
		Version: 0,
		Up: func() {
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
			Db.Exec(`
				create table if not exists setting (
					key text primary key,
					value text
				);
			`)
		},
	},
	{
		Version: 1,
		Up: func() {
			Db.Exec(`
				alter table task add column monthly bool default false;
			`)
		},
	},
}

func migrateDb() {
	for _, m := range migrations {
		if settings.SchemaVersion < m.Version {
			m.Up()
			updateSetting(SchemaVersion, m.Version)
		}
	}
}
