package cmd

import (
	"fmt"
	"time"

	"github.com/gookit/color"
	"github.com/jasonwoodland/track/pkg/completion"
	"github.com/jasonwoodland/track/pkg/model"
	"github.com/jasonwoodland/track/pkg/util"
	"github.com/urfave/cli/v2"
)

var Status = &cli.Command{
	Name:  "status",
	Usage: "Display status of running task",
	BashComplete: func(c *cli.Context) {
		completion.ShowFlagCompletion(c)
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
			state := model.GetState()
			if !state.Running {
				fmt.Println("Not running\033[J")
				return
			}
			color.Printf(
				"Running: <magenta>%s</> <blue>%s</> (%s, %s total)\033[K\n",
				state.Task.Project.Name,
				state.Task.Name,
				util.GetHours(state.TimeElapsed),
				util.GetHours(state.Task.GetTotal()+state.TimeElapsed),
			)
			color.Printf("Started at <green>%s</> (%s ago)\033[K\n", state.StartTime.Format("15:04"), state.TimeElapsed.Round(time.Second))
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
