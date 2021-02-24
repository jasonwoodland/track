package main

import (
	"log"
	"strconv"
	"time"

	"github.com/gookit/color"
	"github.com/urfave/cli/v2"
)

type Frame struct {
	id        int64
	task      *Task
	startTime time.Time
	endTime   time.Time
}

var FrameCmd = &cli.Command{
	Name:  "frame",
	Usage: "Manage recorded frames for a task",
	Subcommands: []*cli.Command{
		{
			Name:         "add",
			Usage:        "Add a frame",
			ArgsUsage:    "project task duration",
			BashComplete: ProjectTaskFrameCompletion,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "offset",
					Aliases: []string{"o"},
					Usage:   "Duration to offset the start and time by (eg. --offset -5m; add a frame which finished 5 minutes ago)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Args().Len() != 3 {
					cli.ShowSubcommandHelp(c)
					return nil
				}

				projectName := c.Args().Get(0)
				taskName := c.Args().Get(1)
				duration, err := time.ParseDuration(c.Args().Get(2))

				if err != nil {
					log.Fatalf("Bad duration: %s", c.Args().Get(2))
				}

				project := GetProjectByName(projectName)

				if project == nil {
					color.Printf("Project <magenta>%s</> doesn't exist\n", projectName)
					return nil
				}

				task := project.GetTask(taskName)

				if task == nil {
					color.Printf("Task <blue>%s</> doesn't exist on project <magenta>%s</>\n", taskName, projectName)
					return nil
				}

				startTime := time.Now().Add(0 - duration)
				endTime := time.Now()

				if o, err := time.ParseDuration(c.String("offset")); err == nil {
					startTime = startTime.Add(o)
					endTime = endTime.Add(o)
				}

				Db.Exec(
					"insert into frame (task_id, start_time, end_time) values ($1, $2, $3)",
					task.id,
					startTime.Format(time.RFC3339),
					endTime.Format(time.RFC3339),
				)

				color.Printf(
					"Added: <magenta>%s</> <blue>%s</> (%s, %s total)\033[K\n",
					project.name,
					task.name,
					GetHours(duration),
					GetHours(task.GetTotal()+duration),
				)

				color.Printf(
					"  <gray>[%v]</> <green>%s - %s</> <default>(%s)</>\n",
					task.getNumFrames()-1,
					startTime.Format("Mon Jan 02 15:04"),
					endTime.Format("15:04"),
					GetHours(endTime.Sub(startTime)),
				)
				return nil
			},
		},
		{
			Name:         "edit",
			Usage:        "Edit a frame's start and end times",
			ArgsUsage:    "project task frame",
			BashComplete: ProjectTaskFrameCompletion,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "start",
					Aliases: []string{"s"},
					Usage:   "Duration to modify the start time by (eg. --start -5m)",
				},
				&cli.StringFlag{
					Name:    "end",
					Aliases: []string{"e"},
					Usage:   "Duration to modify the end time by (eg. --end -5m)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Args().Len() != 3 {
					cli.ShowSubcommandHelp(c)
					return nil
				}

				projectName := c.Args().Get(0)
				taskName := c.Args().Get(1)
				frameIndex, _ := strconv.Atoi(c.Args().Get(2))

				project := GetProjectByName(projectName)

				if project == nil {
					color.Printf("Project <magenta>%s</> doesn't exist\n", projectName)
					return nil
				}

				task := project.GetTask(taskName)

				if task == nil {
					color.Printf("Task <blue>%s</> doesn't exist on project <magenta>%s</>\n", taskName, projectName)
					return nil
				}

				frames := task.GetFrames()

				if frameIndex > len(frames) {
					color.Printf("Frame <gray>[%v]</> doesn't exist on task <blue>%s</>, on project <magenta>%s</>\n", frameIndex, taskName, projectName)
					return nil
				}

				frame := frames[frameIndex]

				if d, err := time.ParseDuration(c.String("start")); err == nil {
					frame.startTime = frame.startTime.Add(d)
				}

				if d, err := time.ParseDuration(c.String("end")); err == nil {
					frame.endTime = frame.endTime.Add(d)
				}

				// TODO 00:00 shown if the frame is currently running.
				color.Printf("Project: <magenta>%s</>\n", projectName)
				color.Printf("  <blue>%s</>\n", taskName)
				color.Printf(
					"    <gray>[%v]</> <green>%s - %s</> (%s)\n",
					frameIndex,
					frame.startTime.Format("Mon Jan 02 15:04"),
					frame.endTime.Format("15:04"),
					GetHours(frame.endTime.Sub(frame.startTime)),
				)

				Db.Exec(
					"update frame set start_time = $1, end_time = $2 where id = $3",
					frame.startTime.Format(time.RFC3339),
					frame.endTime.Format(time.RFC3339),
					frame.id,
				)
				return nil
			},
		},
		{
			Name:         "remove",
			Aliases:      []string{"rm"},
			Usage:        "Delete a frame",
			ArgsUsage:    "project task frame",
			BashComplete: ProjectTaskFrameCompletion,
			Action: func(c *cli.Context) error {
				if c.Args().Len() != 3 {
					cli.ShowSubcommandHelp(c)
					return nil
				}

				projectName := c.Args().Get(0)
				taskName := c.Args().Get(1)
				frameIndex, _ := strconv.Atoi(c.Args().Get(2))

				project := GetProjectByName(projectName)

				if project == nil {
					color.Printf("Project <magenta>%s</> doesn't exist\n", projectName)
					return nil
				}

				task := project.GetTask(taskName)

				if task == nil {
					color.Printf("Task <blue>%s</> doesn't exist on project <magenta>%s</>\n", taskName, projectName)
					return nil
				}

				frames := task.GetFrames()

				if frameIndex > len(frames)-1 {
					color.Printf("Frame <gray>[%v]</> doesn't exist on <magenta>%s</> <blue>%s</>\n", frameIndex, taskName, projectName)
					return nil
				}

				if !Confirm(color.Sprintf("Remove frame <gray>[%v]</> on task <blue>%s</>, on project <magenta>%s</>?", frameIndex, taskName, projectName), false) {
					return nil
				}

				frame := frames[frameIndex]
				Db.Exec(
					"delete from frame where id = $1",
					frame.id,
				)
				color.Printf("Removed frame <gray>[%v]</> on task <blue>%s</>, on project <magenta>%s</>\n", frameIndex, taskName, projectName)
				return nil
			},
		},
	},
}
