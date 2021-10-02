package model

import (
	"log"
	"time"

	"github.com/jasonwoodland/track/pkg/db"
)

type Task struct {
	Id      int64
	Name    string
	Project *Project
}

func GetTaskById(id int64) (t Task) {
	rows, err := db.Db.Query("select id, name, project_id from task where id = $1", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		t = Task{
			Id: id,
		}
		var projectId int64
		rows.Scan(&t.Id, &t.Name, &projectId)
		t.Project = GetProjectById(projectId)
	}
	return
}

func (t *Task) GetNumFrames() (n int) {
	rows, err := db.Db.Query("select count(*) from frame where task_id = $1", t.Id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&n)
		return
	}
	return 0
}

func (t *Task) GetFrames() (frames []*Frame) {
	rows, err := db.Db.Query("select id, start_time, end_time from frame where task_id = $1", t.Id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		f := &Frame{
			Task: t,
		}
		var startTime, endTime string
		rows.Scan(&f.Id, &startTime, &endTime)
		f.StartTime, _ = time.Parse(time.RFC3339, startTime)
		f.EndTime, _ = time.Parse(time.RFC3339, endTime)
		frames = append(frames, f)
	}
	return
}

func (t *Task) GetTotal() (d time.Duration) {
	rows, err := db.Db.Query(`
		select
			sum(strftime("%s", coalesce(end_time, datetime('now'))) - strftime("%s", start_time)) as total
		from frame
		where task_id = $1
	`, t.Id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&d)
		d *= time.Second
	}
	return
}
