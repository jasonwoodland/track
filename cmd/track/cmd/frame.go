package cmd

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gookit/color"
	"github.com/jasonwoodland/track/pkg/completion"
	"github.com/jasonwoodland/track/pkg/db"
	"github.com/jasonwoodland/track/pkg/dialog"
	"github.com/jasonwoodland/track/pkg/model"
	"github.com/jasonwoodland/track/pkg/util"
	"github.com/urfave/cli/v2"
)

var FrameCmds = &cli.Command{
	Name:  "frame",
	Usage: "Manage recorded frames for a task",
	Subcommands: []*cli.Command{
		{
			Name:         "add",
			Usage:        "Add a frame",
			ArgsUsage:    "project task duration",
			BashComplete: completion.ProjectTaskFrameCompletion,
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

				project := model.GetProjectByName(projectName)
				if project == nil {
					color.Printf("Project <magenta>%s</> doesn't exist\n", projectName)
					return nil
				}

				task := project.GetTask(taskName)
				if task == nil {
					color.Printf("Adding task <blue>%s</>\n", taskName)
					task = project.AddTask(taskName)
				}

				startTime := time.Now().Add(0 - duration)
				endTime := time.Now()

				if o, err := time.ParseDuration(c.String("offset")); err == nil {
					startTime = startTime.Add(o)
					endTime = endTime.Add(o)
				}

				db.Db.Exec(
					"insert into frame (task_id, start_time, end_time) values ($1, $2, $3)",
					task.Id,
					startTime.Format(time.RFC3339),
					endTime.Format(time.RFC3339),
				)

				color.Printf(
					"Added: <magenta>%s</> <blue>%s</> (%s, %s total)\033[K\n",
					project.Name,
					task.Name,
					util.GetHours(duration),
					util.GetHours(task.GetTotal()),
				)

				color.Printf(
					"  <gray>[%v]</> <green>%s - %s</> <default>(%s)</>\n",
					task.GetNumFrames()-1,
					startTime.Format("Mon Jan 02 15:04"),
					endTime.Format("15:04"),
					util.GetHours(endTime.Sub(startTime)),
				)
				return nil
			},
		},
		{
			Name:         "edit",
			Usage:        "Edit a frame's start and end times",
			ArgsUsage:    "project task frame",
			BashComplete: completion.ProjectTaskFrameCompletion,
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

				project := model.GetProjectByName(projectName)
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
					frame.StartTime = frame.StartTime.Add(d)
				}

				if d, err := time.ParseDuration(c.String("end")); err == nil {
					frame.EndTime = frame.EndTime.Add(d)
				}

				// TODO 00:00 shown if the frame is currently running.
				color.Printf("Project: <magenta>%s</>\n", projectName)
				color.Printf("  <blue>%s</>\n", taskName)
				color.Printf(
					"    <gray>[%v]</> <green>%s - %s</> (%s)\n",
					frameIndex,
					frame.StartTime.Format("Mon Jan 02 15:04"),
					frame.EndTime.Format("15:04"),
					util.GetHours(frame.EndTime.Sub(frame.StartTime)),
				)

				db.Db.Exec(
					"update frame set start_time = $1, end_time = $2 where id = $3",
					frame.StartTime.Format(time.RFC3339),
					frame.EndTime.Format(time.RFC3339),
					frame.Id,
				)
				return nil
			},
		},
		{
			Name:         "remove",
			Aliases:      []string{"rm"},
			Usage:        "Delete a frame",
			ArgsUsage:    "project task frame",
			BashComplete: completion.ProjectTaskFrameCompletion,
			Action: func(c *cli.Context) error {
				if c.Args().Len() != 3 {
					cli.ShowSubcommandHelp(c)
					return nil
				}

				projectName := c.Args().Get(0)
				taskName := c.Args().Get(1)
				frameIndex, _ := strconv.Atoi(c.Args().Get(2))

				project := model.GetProjectByName(projectName)
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

				if !dialog.Confirm(color.Sprintf("Remove frame <gray>[%v]</> on task <blue>%s</>, on project <magenta>%s</>?", frameIndex, taskName, projectName), false) {
					return nil
				}

				frame := frames[frameIndex]
				db.Db.Exec(
					"delete from frame where id = $1",
					frame.Id,
				)
				color.Printf("Removed frame <gray>[%v]</> on task <blue>%s</>, on project <magenta>%s</>\n", frameIndex, taskName, projectName)
				return nil
			},
		},
		{
			Name:         "move",
			Aliases:      []string{"mv"},
			Usage:        "Move a frame to another project/task",
			ArgsUsage:    "project task frame new_project new_task",
			BashComplete: completion.ProjectTaskFrameProjectTaskCompletion,
			Action: func(c *cli.Context) error {
				if c.Args().Len() != 5 {
					cli.ShowSubcommandHelp(c)
					return nil
				}

				projectName := c.Args().Get(0)
				taskName := c.Args().Get(1)
				frameIndex, _ := strconv.Atoi(c.Args().Get(2))
				newProjectName := c.Args().Get(3)
				newTaskName := c.Args().Get(4)

				project := model.GetProjectByName(projectName)
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

				newProject := model.GetProjectByName(newProjectName)
				if newProject == nil {
					color.Printf("Project <magenta>%s</> doesn't exist\n", newProjectName)
					return nil
				}

				newTask := newProject.GetTask(newTaskName)
				if task == nil {
					color.Printf("Adding task <blue>%s</>\n", taskName)
					task = project.AddTask(taskName)
				}

				if !dialog.Confirm(
					color.Sprintf(
						"Move frame <gray>[%v]</> from <magenta>%s</> <blue>%s</> to <magenta>%s</> <blue>%s</>?",
						frameIndex,
						projectName,
						taskName,
						newProjectName,
						newTaskName,
					),
					false,
				) {
					return nil
				}

				frame := frames[frameIndex]
				db.Db.Exec(
					"update frame set task_id = $1 where id = $2",
					newTask.Id,
					frame.Id,
				)
				fmt.Println("Moved")
				return nil
			},
		},
	},
}
