package cmd

import (
	"fmt"
	"time"

	"github.com/gookit/color"
	"github.com/jasonwoodland/track/pkg/cleanup"
	"github.com/jasonwoodland/track/pkg/completion"
	"github.com/jasonwoodland/track/pkg/db"
	"github.com/jasonwoodland/track/pkg/model"
	"github.com/jasonwoodland/track/pkg/presenter"
	"github.com/jasonwoodland/track/pkg/util"
	"github.com/jasonwoodland/track/pkg/view"
	"github.com/urfave/cli/v2"
)

var Start = &cli.Command{
	Name:      "start",
	Usage:     "Start tracking time for a task",
	ArgsUsage: "project task",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "ago",
			Usage: "Offest the start time with a duration (eg. --ago 5m)",
		},
		&cli.StringFlag{
			Name:  "in",
			Usage: "Start tracking in a given duration (eg. --in 5m)",
		},
		&cli.BoolFlag{
			Name:    "watch",
			Aliases: []string{"w"},
			Usage:   "Output the current status to the screen periodically",
		},
	},
	BashComplete: completion.ProjectTaskCompletion,
	Action: func(c *cli.Context) error {
		startTime := time.Now()

		projectName := c.Args().Get(0)
		taskName := c.Args().Get(1)
		if projectName == "" || taskName == "" {
			cli.ShowSubcommandHelp(c)
			return nil
		}

		project := model.GetProjectByName(projectName)
		if project == nil {
			color.Printf(view.ProjectDoesNotExist, projectName)
			return nil
		}

		task := project.GetTask(taskName)

		state := model.GetState()
		if state != nil && state.Running {
			color.Printf(
				view.AlreadyRunningProjectTaskElapsedTotal,
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
			if task != nil && state.Task.Id == task.Id {
				return nil
			}
			if !presenter.Confirm(view.ConfirmStopRunningTask, true) {
				return nil
			}
		}

		if task == nil {
			color.Printf(view.AddedTask, taskName)
			task = project.AddTask(taskName)
		}

		ago, err := time.ParseDuration(c.String("ago"))
		if err == nil {
			startTime = startTime.Add(0 - ago)
		}

		if in, err := time.ParseDuration(c.String("in")); err == nil {
			startTime = startTime.Add(in)
		}

		db.Db.Exec(
			"insert into frame (task_id, start_time) values ($1, $2)",
			task.Id,
			startTime.Format(time.RFC3339),
		)

		if c.Bool("watch") {
			printStatus := func() {
				state := model.GetState()
				if !state.Running {
					fmt.Println("Not running\033[J")
					return
				}
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
			}

			cleanup.SetCleanupFn(func() {
				printStatus()
			})

			fmt.Printf("\033[?1049h\033[H")
			for {
				printStatus()
				time.Sleep(time.Second)
				fmt.Printf("\033[H")
			}
		} else {
			if ago != 0 {
				state := model.GetState()
				color.Printf(
					view.RunningProjectTaskElapsedTotal,
					project.Name,
					task.Name,
					util.GetHours(state.TimeElapsed),
					util.GetHours(state.Task.GetTotal()),
				)
			} else {
				color.Printf(
					view.RunningProjectTaskTotal,
					project.Name,
					task.Name,
					util.GetHours(task.GetTotal()),
				)
			}

			color.Printf(view.StartedAtTime, startTime.Format("15:04"))
		}

		return nil
	},
}
