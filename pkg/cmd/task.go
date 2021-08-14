package cmd

import (
	"log"

	"github.com/gookit/color"
	"github.com/jasonwoodland/track/pkg/completion"
	"github.com/jasonwoodland/track/pkg/db"
	"github.com/jasonwoodland/track/pkg/model"
	"github.com/jasonwoodland/track/pkg/presenter"
	"github.com/urfave/cli/v2"
)

var TaskCmds = &cli.Command{
	Name:  "task",
	Usage: "Manage tasks on a project",
	Subcommands: []*cli.Command{
		{
			Name:         "rename",
			Usage:        "Rename a task",
			ArgsUsage:    "project old_name new_name",
			BashComplete: completion.ProjectTaskCompletion,
			Action: func(c *cli.Context) error {
				if c.Args().Len() != 3 {
					cli.ShowSubcommandHelp(c)
					return nil
				}

				projectName := c.Args().Get(0)
				oldName := c.Args().Get(1)
				newName := c.Args().Get(2)

				project := model.GetProjectByName(projectName)

				if project == nil {
					color.Printf("Project <magenta>%s</> doesn't exist\n", projectName)
					return nil
				}

				if project.GetTask(newName) != nil {
					color.Printf("Task <blue>%s</> already exists on project <magenta>%s</>\n", newName, projectName)
					return nil
				}

				if project.GetTask(oldName) == nil {
					color.Printf("Task <blue>%s</> doesn't exist on project <magenta>%s</>\n", oldName, projectName)
					return nil
				}

				db.Db.Exec("update task set name = $1 where name = $2 and project_id = $3", newName, oldName, project.Id)
				color.Printf("Renamed task <blue>%s</> to <blue>%s</> on project <magenta>%s</>\n", oldName, newName, projectName)
				return nil
			},
		},
		{
			Name:         "remove",
			Aliases:      []string{"rm"},
			Usage:        "Delete a task",
			ArgsUsage:    "project task",
			BashComplete: completion.ProjectTaskCompletion,
			Action: func(c *cli.Context) error {
				if c.Args().Len() != 2 {
					cli.ShowSubcommandHelp(c)
					return nil
				}

				projectName := c.Args().Get(0)
				taskName := c.Args().Get(1)

				project := model.GetProjectByName(projectName)
				if project == nil {
					color.Printf("Project <magenta>%s</> doesn't exists\n", projectName)
					return nil
				}

				task := project.GetTask(taskName)

				if task == nil {
					color.Printf("Task <blue>%s</> doesn't exists on project <magenta>%s</>\n", taskName, projectName)
					return nil
				}

				numFrames := task.GetNumFrames()
				s := "s"
				if numFrames == 1 {
					s = ""
				}

				if !presenter.Confirm(color.Sprintf(
					"Delete task <blue>%s</> and %d frame%s on project <magenta>%s</>?",
					taskName,
					numFrames,
					s,
					projectName,
				), false) {
					return nil
				}

				db.Db.Exec("delete from task where name = $1 and project_id = $2", taskName, project.Id)
				color.Printf(
					"Deleted task <blue>%s</> and %d frame%s on project <magenta>%s</>\n",
					taskName,
					numFrames,
					s,
					projectName,
				)
				return nil
			},
		},
		{
			Name:         "merge",
			Usage:        "Merge a task",
			ArgsUsage:    "from_project from_task to_project to_task",
			BashComplete: completion.ProjectTaskProjectTaskCompletion,
			Action: func(c *cli.Context) error {
				if c.Args().Len() != 4 {
					cli.ShowSubcommandHelp(c)
					return nil
				}

				fromProjectName := c.Args().Get(0)
				fromTaskName := c.Args().Get(1)
				toProjectName := c.Args().Get(2)
				toTaskName := c.Args().Get(3)

				fromProject := model.GetProjectByName(fromProjectName)
				if fromProject == nil {
					color.Printf("Project <magenta>%s</> doesn't exists\n", fromProjectName)
					return nil
				}

				fromTask := fromProject.GetTask(fromTaskName)

				if fromTask == nil {
					color.Printf("Task <blue>%s</> doesn't exists on project <magenta>%s</>\n", fromTaskName, fromProjectName)
					return nil
				}

				toProject := model.GetProjectByName(toProjectName)
				if toProject == nil {
					color.Printf("Project <magenta>%s</> doesn't exists\n", toProjectName)
					return nil
				}

				toTask := toProject.GetTask(toTaskName)

				if toTask == nil {
					color.Printf("Task <blue>%s</> doesn't exists on project <magenta>%s</>\n", toTaskName, toProjectName)
					return nil
				}

				numFrames := fromTask.GetNumFrames()
				s := "s"
				if numFrames == 1 {
					s = ""
				}

				if !presenter.Confirm(color.Sprintf(
					"Merge %d frame%s from <magenta>%s</> <blue>%s</> into <magenta>%s</> <blue>%s</>?",
					numFrames,
					s,
					fromProjectName,
					fromTaskName,
					toProjectName,
					toTaskName,
				), false) {
					return nil
				}

				_, err := db.Db.Exec("update frame set task_id = $1 where task_id = $2", toTask.Id, fromTask.Id)
				if err != nil {
					log.Fatal(err)
				}
				_, err = db.Db.Exec("delete from task where id = $1", fromTask.Id)
				if err != nil {
					log.Fatal(err)
				}

				color.Println("Merged")

				return nil
			},
		},
	},
}
