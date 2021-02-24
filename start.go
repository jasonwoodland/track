package main

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/urfave/cli/v2"
	"time"
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

		if ago, err := time.ParseDuration(c.String("ago")); err == nil {
			startTime = startTime.Add(0 - ago)
		}

		if in, err := time.ParseDuration(c.String("in")); err == nil {
			startTime = startTime.Add(in)
		}

		color.Printf("Running: <magenta>%s</> <blue>%s</> (%s)\n", project.name, task.name, GetHours(task.GetTotal()))
		color.Printf("Started at <green>%s</>\n", startTime.Format("15:04"))

		Db.Exec(
			"insert into frame (task_id, start_time) values ($1, $2)",
			task.id,
			startTime.Format(time.RFC3339),
		)
		return nil
	},
}
