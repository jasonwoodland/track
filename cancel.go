package main

import (
	"fmt"
	"time"

	"github.com/gookit/color"
	"github.com/urfave/cli/v2"
)

var Cancel = &cli.Command{
	Name:  "cancel",
	Usage: "Cancel a running task",
	Action: func(c *cli.Context) error {
		state := GetState()

		if state.running {
			color.Printf(
				"Cancelled: <magenta>%s</> <blue>%s</> (%s, %s total)\n",
				state.task.project.name,
				state.task.name,
				GetHours(state.timeElapsed),
				GetHours(state.task.GetTotal()),
			)
			color.Printf("Started at <green>%s</> (%s ago)\033[K\n", state.startTime.Format("15:04"), state.timeElapsed.Round(time.Second))

			Db.Exec("delete from frame where end_time is null")
			Db.Exec("delete from task where not exists (select 1 from frame where frame.task_id = task.id)")
		} else {
			fmt.Println("Not runnning")
		}
		return nil
	},
}
