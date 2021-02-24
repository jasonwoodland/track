package main

import (
	"log"

	"github.com/gookit/color"
	"github.com/urfave/cli/v2"
)

type Project struct {
	id   int64
	name string
}

func GetProjects() (projects []*Project) {
	rows, err := Db.Query("select id, name from project")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		p := &Project{}
		rows.Scan(&p.id, &p.name)
		projects = append(projects, p)
	}
	return
}

func GetProjectById(id int64) (p *Project) {
	rows, err := Db.Query("select name from project where id = $1", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		p = &Project{
			id: id,
		}
		rows.Scan(&p.name)
	}
	return
}

func GetProjectByName(name string) (p *Project) {
	rows, err := Db.Query("select id from project where name = $1", name)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		p = &Project{
			name: name,
		}
		rows.Scan(&p.id)
	}
	return
}

func (p *Project) GetTask(name string) (t *Task) {
	rows, err := Db.Query("select id from task where project_id = $1 and name = $2", p.id, name)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		t = &Task{
			name:    name,
			project: p,
		}
		rows.Scan(&t.id)
	}
	return
}

func (p *Project) GetTasks() (tasks []*Task) {
	rows, err := Db.Query("select id, name from task where project_id = $1", p.id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		t := &Task{
			project: p,
		}
		rows.Scan(&t.id, &t.name)
		tasks = append(tasks, t)
	}
	return
}

func (p *Project) AddTask(name string) *Task {
	res, err := Db.Exec("insert into task (name, project_id) values ($1, $2)", name, p.id)
	if err != nil {
		log.Fatal(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	return &Task{
		id:      id,
		name:    name,
		project: p,
	}
}

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
				if GetProjectByName(name) != nil {
					color.Printf("Project <magenta>%s</> already exists\n", name)
					return nil
				}
				Db.Exec("insert into project (name) values ($1)", name)
				color.Printf("Added project <magenta>%s</>\n", name)
				return nil
			},
		},
		{
			Name:         "rename",
			Usage:        "Rename a project",
			ArgsUsage:    "old_name new_name",
			BashComplete: ProjectCompletion,
			Action: func(c *cli.Context) error {
				oldName := c.Args().Get(0)
				newName := c.Args().Get(1)
				if oldName == "" || newName == "" {
					cli.ShowSubcommandHelp(c)
					return nil
				}
				if GetProjectByName(oldName) == nil {
					color.Printf("Project <magenta>%s</> doesn't exists\n", oldName)
					return nil
				}
				Db.Exec("update project set name = $1 where name = $2", newName, oldName)
				color.Printf("Renamed project <magenta>%s</> to <magenta>%s</>\n", oldName, newName)
				return nil
			},
		},
		{
			Name:         "remove",
			Aliases:      []string{"rm"},
			Usage:        "Delete a project and all associated tasks",
			ArgsUsage:    "name",
			BashComplete: ProjectCompletion,
			Action: func(c *cli.Context) error {
				name := c.Args().Get(0)
				if name == "" {
					cli.ShowSubcommandHelp(c)
					return nil
				}

				if GetProjectByName(name) == nil {
					color.Printf("Project <magenta>%s</> doesn't exists\n", name)
					return nil
				}

				if !Confirm(color.Sprintf("Delete project <magenta>%s</>?", name), false) {
					return nil
				}

				Db.Exec("delete from project where name = $1", name)
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
		for _, project := range GetProjects() {
			color.Magenta.Println(project.name)
		}
		return nil
	},
}
