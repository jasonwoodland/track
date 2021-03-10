package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gookit/color"
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
		state := GetState()
		endTime := time.Now()

		if ago, err := time.ParseDuration(c.String("ago")); err == nil {
			endTime = endTime.Add(0 - ago)
			state.timeElapsed -= ago
		}

		if in, err := time.ParseDuration(c.String("in")); err == nil {
			endTime = endTime.Add(in)
			state.timeElapsed += in
		}

		res, err := Db.Exec(
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
				state.task.project.name,
				state.task.name,
				GetHours(state.timeElapsed),
				GetHours(state.task.GetTotal()),
			)
			color.Printf(
				"Finished at <green>%s</> (%s)\n",
				endTime.Format("15:04"),
				state.timeElapsed.Round(time.Second),
			)
		}

		return nil
	},
}
