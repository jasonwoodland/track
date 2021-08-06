package cmd

import (
	"github.com/gookit/color"
	"github.com/jasonwoodland/track/pkg/completion"
	"github.com/jasonwoodland/track/pkg/db"
	"github.com/jasonwoodland/track/pkg/dialog"
	"github.com/jasonwoodland/track/pkg/model"
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

				if project.GetTask(taskName) == nil {
					color.Printf("Task <blue>%s</> doesn't exists on project <magenta>%s</>\n", taskName, projectName)
					return nil
				}

				if !dialog.Confirm(color.Sprintf("Delete task <blue>%s</> on project <magenta>%s</>?", taskName, projectName), false) {
					return nil
				}

				db.Db.Exec("delete from project where name = $1 and project_id = $2", taskName, project.Id)
				color.Printf("Deleted task <blue>%s</> on project <magenta>%s</>\n", taskName, projectName)
				return nil
			},
		},
	},
}
