package main

import (
	"log"
	"time"
)

type State struct {
	running     bool
	task        Task
	startTime   time.Time
	timeElapsed time.Duration
}

func GetState() (s *State) {
	s = &State{}
	rows, err := Db.Query("select task_id, start_time from frame where end_time is null")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		s.running = true
		var taskId int64
		var startTime string
		rows.Scan(&taskId, &startTime)
		s.task = GetTaskById(taskId)
		s.startTime, _ = time.Parse(time.RFC3339, startTime)
		s.timeElapsed = time.Now().Sub(s.startTime)
	}
	return
}
