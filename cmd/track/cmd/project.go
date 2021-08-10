package cmd

import (
	"github.com/gookit/color"
	"github.com/jasonwoodland/track/pkg/completion"
	"github.com/jasonwoodland/track/pkg/db"
	"github.com/jasonwoodland/track/pkg/model"
	"github.com/jasonwoodland/track/pkg/presenter"
	"github.com/urfave/cli/v2"
)

var ProjectCmds = &cli.Command{
	Name:  "project",
	Usage: "Manage projects",
	Subcommands: []*cli.Command{
		{
			Name:      "add",
			Usage:     "Add a new project",
			ArgsUsage: "name",
			Action: func(c *cli.Context) error {
				name := c.Args().Get(0)
				if name == "" {
					cli.ShowSubcommandHelp(c)
					return nil
				}
				if model.GetProjectByName(name) != nil {
					color.Printf("Project <magenta>%s</> already exists\n", name)
					return nil
				}
				db.Db.Exec("insert into project (name) values ($1)", name)
				color.Printf("Added project <magenta>%s</>\n", name)
				return nil
			},
		},
		{
			Name:         "rename",
			Usage:        "Rename a project",
			ArgsUsage:    "old_name new_name",
			BashComplete: completion.ProjectCompletion,
			Action: func(c *cli.Context) error {
				oldName := c.Args().Get(0)
				newName := c.Args().Get(1)
				if oldName == "" || newName == "" {
					cli.ShowSubcommandHelp(c)
					return nil
				}
				if model.GetProjectByName(oldName) == nil {
					color.Printf("Project <magenta>%s</> doesn't exists\n", oldName)
					return nil
				}
				db.Db.Exec("update project set name = $1 where name = $2", newName, oldName)
				color.Printf("Renamed project <magenta>%s</> to <magenta>%s</>\n", oldName, newName)
				return nil
			},
		},
		{
			Name:         "remove",
			Aliases:      []string{"rm"},
			Usage:        "Delete a project and all associated tasks",
			ArgsUsage:    "name",
			BashComplete: completion.ProjectCompletion,
			Action: func(c *cli.Context) error {
				name := c.Args().Get(0)
				if name == "" {
					cli.ShowSubcommandHelp(c)
					return nil
				}

				if model.GetProjectByName(name) == nil {
					color.Printf("Project <magenta>%s</> doesn't exists\n", name)
					return nil
				}

				if !presenter.Confirm(color.Sprintf("Delete project <magenta>%s</>?", name), false) {
					return nil
				}

				db.Db.Exec("delete from project where name = $1", name)
				color.Printf("Deleted project <magenta>%s</>\n", name)
				return nil
			},
		},
	},
}

var Projects = &cli.Command{
	Name:  "projects",
	Usage: "List projects",
	Action: func(c *cli.Context) error {
		for _, project := range model.GetProjects() {
			color.Magenta.Println(project.Name)
		}
		return nil
	},
}
