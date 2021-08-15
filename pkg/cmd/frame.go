package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gookit/color"
	"github.com/jasonwoodland/track/pkg/completion"
	"github.com/jasonwoodland/track/pkg/db"
	"github.com/jasonwoodland/track/pkg/model"
	"github.com/jasonwoodland/track/pkg/presenter"
	"github.com/jasonwoodland/track/pkg/util"
	"github.com/jasonwoodland/track/pkg/view"
	"github.com/urfave/cli/v2"
)

var FrameCmds = &cli.Command{
	Name:  "frame",
	Usage: "Manage recorded frames for a task",
	Subcommands: []*cli.Command{
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
					color.Printf(view.ProjectDoesNotExist, projectName)
					return nil
				}

				task := project.GetTask(taskName)
				if task == nil {
					color.Printf(view.TaskDoesNotExistForProject, taskName, projectName)
					return nil
				}

				frames := task.GetFrames()
				if frameIndex > len(frames) {
					color.Printf(view.FrameDoesNotExistForProjectTask, frameIndex, projectName, taskName)
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
				color.Printf(view.Project, projectName)
				color.Printf(view.Task, taskName)
				color.Printf(
					view.FrameTimesDuration,
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
					color.Printf(view.ProjectDoesNotExist, projectName)
					return nil
				}

				task := project.GetTask(taskName)
				if task == nil {
					color.Printf(view.TaskDoesNotExistForProject, taskName, projectName)
					return nil
				}

				frames := task.GetFrames()
				if frameIndex > len(frames)-1 {
					color.Printf(view.FrameDoesNotExistForProjectTask, frameIndex, projectName, taskName)
					return nil
				}

				if !presenter.Confirm(color.Sprintf(
					view.ConfirmDeleteFrameTimeProjectTask,
					frames[frameIndex].StartTime.Format("Mon Jan 02 15:04"),
					frames[frameIndex].EndTime.Format("15:04"),
					projectName,
					taskName,
				), false) {
					return nil
				}

				frame := frames[frameIndex]
				db.Db.Exec(
					"delete from frame where id = $1",
					frame.Id,
				)

				color.Println(view.Deleted)

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
					color.Printf(view.ProjectDoesNotExist, projectName)
					return nil
				}

				task := project.GetTask(taskName)
				if task == nil {
					color.Printf(view.TaskDoesNotExistForProject, taskName, projectName)
					return nil
				}

				frames := task.GetFrames()
				if frameIndex > len(frames)-1 {
					color.Printf(view.FrameDoesNotExistForProjectTask, frameIndex, projectName, taskName)
					return nil
				}

				newProject := model.GetProjectByName(newProjectName)
				if newProject == nil {
					color.Printf(view.ProjectDoesNotExist, newProjectName)
					return nil
				}

				newTask := newProject.GetTask(newTaskName)
				if task == nil {
					color.Printf(view.AddedTask, taskName)
					task = project.AddTask(taskName)
				}

				if !presenter.Confirm(
					color.Sprintf(
						view.ConfirmMoveFrameTimesFromToProjectTask,
						frames[frameIndex].StartTime.Format("Mon Jan 02 15:04"),
						frames[frameIndex].EndTime.Format("Mon Jan 02"),
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

				fmt.Println(view.Moved)

				return nil
			},
		},
	},
}
