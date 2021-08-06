package cmd

import (
	"fmt"
	"time"

	"github.com/gookit/color"
	"github.com/jasonwoodland/track/pkg/db"
	"github.com/jasonwoodland/track/pkg/model"
	"github.com/jasonwoodland/track/pkg/util"
	"github.com/urfave/cli/v2"
)

var Cancel = &cli.Command{
	Name:  "cancel",
	Usage: "Cancel a running task",
	Action: func(c *cli.Context) error {
		state := model.GetState()

		if state.Running {
			color.Printf(
				"Cancelled: <magenta>%s</> <blue>%s</> (%s, %s total)\n",
				state.Task.Project.Name,
				state.Task.Name,
				util.GetHours(state.TimeElapsed),
				util.GetHours(state.Task.GetTotal()),
			)
			color.Printf("Started at <green>%s</> (%s ago)\033[K\n", state.StartTime.Format("15:04"), state.TimeElapsed.Round(time.Second))

			db.Db.Exec("delete from frame where end_time is null")
			db.Db.Exec("delete from task where not exists (select 1 from frame where frame.task_id = task.id)")
		} else {
			fmt.Println("Not runnning")
		}
		return nil
	},
}
