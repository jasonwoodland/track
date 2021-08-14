package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/gookit/color"
	"github.com/jasonwoodland/track/pkg/db"
	"github.com/jasonwoodland/track/pkg/model"
	"github.com/jasonwoodland/track/pkg/util"
	"github.com/urfave/cli/v2"
)

var Stop = &cli.Command{
	Name:  "stop",
	Usage: "Stop a running task",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "ago",
			Usage: "Offest the start time with a duration (eg. --ago 5m)",
		},
		&cli.StringFlag{
			Name:  "in",
			Usage: "Start tracking in a given duration (eg. --in 5m)",
		},
	},
	Action: func(c *cli.Context) error {
		state := model.GetState()
		endTime := time.Now()

		if ago, err := time.ParseDuration(c.String("ago")); err == nil {
			endTime = endTime.Add(0 - ago)
			state.TimeElapsed -= ago
		}

		if in, err := time.ParseDuration(c.String("in")); err == nil {
			endTime = endTime.Add(in)
			state.TimeElapsed += in
		}

		res, err := db.Db.Exec(
			"update frame set end_time = $1 where end_time is null",
			endTime.Format(time.RFC3339),
		)
		if err != nil {
			log.Fatal(err)
		}

		if n, _ := res.RowsAffected(); n == 0 {
			fmt.Println("No task started")
		} else {
			color.Printf(
				"Stopped: <magenta>%s</> <blue>%s</> (%s, %s total)\n",
				state.Task.Project.Name,
				state.Task.Name,
				util.GetHours(state.TimeElapsed),
				util.GetHours(state.Task.GetTotal()),
			)
			color.Printf(
				"Finished at <green>%s</> (%s)\n",
				endTime.Format("15:04"),
				state.TimeElapsed.Round(time.Second),
			)
		}

		return nil
	},
}
