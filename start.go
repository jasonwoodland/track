package main

import (
	"fmt"
	"time"

	"github.com/gookit/color"
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
	BashComplete: ProjectTaskCompletion,
	Action: func(c *cli.Context) error {
		startTime := time.Now()

		projectName := c.Args().Get(0)
		taskName := c.Args().Get(1)
		if projectName == "" || taskName == "" {
			cli.ShowSubcommandHelp(c)
			return nil
		}

		state := GetState()
		if state != nil && state.running {
			fmt.Println("Task already running")
			return nil
		}

		project := GetProjectByName(projectName)
		if project == nil {
			color.Printf("Project <magenta>%s</> doesn't exists\n", projectName)
			return nil
		}

		task := project.GetTask(taskName)
		if task == nil {
			color.Printf("Adding task <blue>%s</>\n", taskName)
			task = project.AddTask(taskName)
		}

		ago, err := time.ParseDuration(c.String("ago"))
		if err == nil {
			startTime = startTime.Add(0 - ago)
		}

		if in, err := time.ParseDuration(c.String("in")); err == nil {
			startTime = startTime.Add(in)
		}

		Db.Exec(
			"insert into frame (task_id, start_time) values ($1, $2)",
			task.id,
			startTime.Format(time.RFC3339),
		)

		if c.Bool("watch") {
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

			Cleanup = func() {
				printStatus()
			}

			fmt.Printf("\033[?1049h\033[H")
			for {
				printStatus()
				time.Sleep(time.Second)
				fmt.Printf("\033[H")
			}
		} else {
			if ago != 0 {
				state := GetState()
				color.Printf(
					"Running: <magenta>%s</> <blue>%s</> (%s, %s total)\033[K\n",
					project.name,
					task.name,
					GetHours(state.timeElapsed),
					GetHours(state.task.GetTotal()+state.timeElapsed),
				)
			} else {
				color.Printf("Running: <magenta>%s</> <blue>%s</> (%s)\n", project.name, task.name, GetHours(task.GetTotal()))
			}

			color.Printf("Started at <green>%s</>\n", startTime.Format("15:04"))
		}

		return nil
	},
}
