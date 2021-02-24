package main

import (
	"fmt"
	"log"

	"github.com/urfave/cli/v2"
)

var Cancel = &cli.Command{
	Name:  "cancel",
	Usage: "Cancel a running task",
	Action: func(c *cli.Context) error {
		res, err := Db.Exec("delete from frame where end_time is null")
		Db.Exec("delete from task where not exists (select 1 from frame where frame.task_id = task.id)")
		if err != nil {
			log.Fatal(err)
		}
		n, err := res.RowsAffected()
		if n == 0 {
			fmt.Println("No task started")
		} else {
			fmt.Println("Task cancelled")
		}
		return nil
	},
}
