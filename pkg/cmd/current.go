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

type command int

const (
	add command = iota
	sub
)

var Shift = &cli.Command{
	Name:    "current",
	Aliases: []string{"cur"},
	Usage:   "Adjust the current running task",
	Subcommands: []*cli.Command{
		{
			Name:      "add",
			Usage:     "Add to the running duration",
			ArgsUsage: "duration",
			Action:    actionForCommand(add),
		},
		{
			Name:      "sub",
			Usage:     "Subtract from the running duration",
			ArgsUsage: "duration",
			Action:    actionForCommand(sub),
		},
	},
}

func actionForCommand(cmd command) func(*cli.Context) error {
	return func(c *cli.Context) error {
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

		prevTimeElapsed := state.TimeElapsed
		prevTaskTotal := state.Task.GetTotal()
		prevStartTime := state.StartTime
		var newStartTime time.Time

		if cmd == add {
			newStartTime = state.StartTime.Add(0 - duration)
		} else if cmd == sub {
			newStartTime = state.StartTime.Add(duration)
		}

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
			view.RunningProjectTaskPrevElapsedTotal,
			state.Task.Project.Name,
			state.Task.Name,
			util.GetHours(prevTimeElapsed),
			util.GetHours(state.TimeElapsed),
			util.GetHours(prevTaskTotal),
			util.GetHours(state.Task.GetTotal()),
		)
		color.Printf(
			view.StartedAtPrevTimeElapsed,
			prevStartTime.Format("15:04"),
			state.StartTime.Format("15:04"),
			prevTimeElapsed.Round(time.Second),
			state.TimeElapsed.Round(time.Second),
		)

		return nil
	}
}
