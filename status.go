package main

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/urfave/cli/v2"
	"time"
)

var Status = &cli.Command{
	Name:  "status",
	Usage: "Display status of running task",
	BashComplete: func(c *cli.Context) {
		ShowFlagCompletion(c)
	},
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "watch",
			Aliases: []string{"w"},
			Usage:   "Output the current status to the screen periodically",
		},
	},
	Action: func(c *cli.Context) error {
		printStatus := func() {
			state := GetState()
			if !state.running {
				fmt.Println("Not running\033[J")
				return
			}
			color.Printf(
				"Running: <magenta>%s</> <blue>%s</> (%s, %s total)\033[K\n",
				state.task.project.name,
				state.task.name,
				GetHours(state.timeElapsed),
				GetHours(state.task.GetTotal()+state.timeElapsed),
			)
			color.Printf("Started at <green>%s</> (%s ago)\033[K\n", state.startTime.Format("15:04"), state.timeElapsed.Round(time.Second))
		}

		if c.Bool("watch") {
			fmt.Printf("\033[?1049h\033[H")
			for {
				printStatus()
				time.Sleep(time.Second)
				fmt.Printf("\033[H")
			}
		}

		printStatus()

		return nil
	},
}
