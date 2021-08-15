package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/gookit/color"
	"github.com/jasonwoodland/track/pkg/db"
	"github.com/jasonwoodland/track/pkg/model"
	"github.com/jasonwoodland/track/pkg/util"
	"github.com/jasonwoodland/track/pkg/view"
	"github.com/urfave/cli/v2"
)

var Shift = &cli.Command{
	Name:      "shift",
	Usage:     "Shift the start time of a running task",
	ArgsUsage: "-- duration",
	Action: func(c *cli.Context) error {
		if c.Args().Len() != 1 {
			cli.ShowSubcommandHelp(c)
			return nil
		}

		state := model.GetState()
		var duration time.Duration

		duration, err := time.ParseDuration(c.Args().Get(0))
		if err != nil {
			color.Println("Bad duration: %s", c.Args().Get(0))
			return nil
		}

		newStartTime := state.StartTime.Add(0 - duration)

		res, err := db.Db.Exec(
			"update frame set start_time = $1 where end_time is null",
			newStartTime.Format(time.RFC3339),
		)
		if err != nil {
			log.Fatal(err)
		}
		if n, _ := res.RowsAffected(); n == 0 {
			fmt.Println(view.NotRunning)
			return nil
		}

		state = model.GetState()

		color.Printf(
			view.RunningProjectTaskElapsedTotal,
			state.Task.Project.Name,
			state.Task.Name,
			util.GetHours(state.TimeElapsed),
			util.GetHours(state.Task.GetTotal()),
		)
		color.Printf(
			view.StartedAtTimeElapsed,
			state.StartTime.Format("15:04"),
			state.TimeElapsed.Round(time.Second),
		)

		return nil
	},
}
