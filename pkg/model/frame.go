package model

import "time"

type Frame struct {
	Id        int64
	Task      *Task
	StartTime time.Time
	EndTime   time.Time
}
