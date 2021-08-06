package model

import (
	"log"
	"time"

	"github.com/jasonwoodland/track/pkg/db"
)

type State struct {
	Running     bool
	Task        Task
	StartTime   time.Time
	TimeElapsed time.Duration
}

func GetState() (s *State) {
	s = &State{}
	rows, err := db.Db.Query("select task_id, start_time from frame where end_time is null")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		s.Running = true
		var taskId int64
		var startTime string
		rows.Scan(&taskId, &startTime)
		s.Task = GetTaskById(taskId)
		s.StartTime, _ = time.Parse(time.RFC3339, startTime)
		s.TimeElapsed = time.Now().Sub(s.StartTime)
	}
	return
}
