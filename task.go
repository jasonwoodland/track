package main

import (
	"log"
	"time"

	"github.com/gookit/color"
	"github.com/urfave/cli/v2"
)

type Task struct {
	id      int64
	name    string
	project *Project
}

func GetTaskById(id int64) (t Task) {
	rows, err := Db.Query("select id, name, project_id from task where id = $1", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		t = Task{
			id: id,
		}
		var projectId int64
		rows.Scan(&t.id, &t.name, &projectId)
		t.project = GetProjectById(projectId)
	}
	return
}

func (t *Task) getNumFrames() (n int) {
	rows, err := Db.Query("select count(*) from frame where task_id = $1", t.id)
	if err != nil {
		log.Fatal(err)
	}
	if rows.Next() {
		rows.Scan(&n)
		return
	}
	return 0
}

func (t *Task) GetFrames() (frames []*Frame) {
	rows, err := Db.Query("select id, start_time, end_time from frame where task_id = $1", t.id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		f := &Frame{
			task: t,
		}
		var startTime, endTime string
		rows.Scan(&f.id, &startTime, &endTime)
		f.startTime, _ = time.Parse(time.RFC3339, startTime)
		f.endTime, _ = time.Parse(time.RFC3339, endTime)
		frames = append(frames, f)
	}
	return
}

func (t *Task) GetTotal() (d time.Duration) {
	rows, err := Db.Query(`
		select
			sum(strftime("%s", end_time) - strftime("%s", start_time)) as total
		from frame
		where task_id = $1
	`, t.id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&d)
		d *= time.Second
	}
	return
}

var TaskCmd = &cli.Command{
	Name:  "task",
	Usage: "Manage tasks on a project",
	Subcommands: []*cli.Command{
		{
			Name:         "rename",
			Usage:        "Rename a task",
			ArgsUsage:    "project old_name new_name",
			BashComplete: ProjectTaskCompletion,
			Action: func(c *cli.Context) error {
				if c.Args().Len() != 3 {
					cli.ShowSubcommandHelp(c)
					return nil
				}

				projectName := c.Args().Get(0)
				oldName := c.Args().Get(1)
				newName := c.Args().Get(2)

				project := GetProjectByName(projectName)

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

				Db.Exec("update task set name = $1 where name = $2 and project_id = $3", newName, oldName, project.id)
				color.Printf("Renamed task <blue>%s</> to <blue>%s</> on project <magenta>%s</>\n", oldName, newName, projectName)
				return nil
			},
		},
		{
			Name:         "remove",
			Aliases:      []string{"rm"},
			Usage:        "Delete a task",
			ArgsUsage:    "project task",
			BashComplete: ProjectTaskCompletion,
			Action: func(c *cli.Context) error {
				if c.Args().Len() != 2 {
					cli.ShowSubcommandHelp(c)
					return nil
				}

				projectName := c.Args().Get(0)
				taskName := c.Args().Get(1)

				project := GetProjectByName(projectName)
				if project == nil {
					color.Printf("Project <magenta>%s</> doesn't exists\n", projectName)
					return nil
				}

				if project.GetTask(taskName) == nil {
					color.Printf("Task <blue>%s</> doesn't exists on project <magenta>%s</>\n", taskName, projectName)
					return nil
				}

				if !Confirm(color.Sprintf("Delete task <blue>%s</> on project <magenta>%s</>?", taskName, projectName), false) {
					return nil
				}

				Db.Exec("delete from project where name = $1 and project_id = $2", taskName, project.id)
				color.Printf("Deleted task <blue>%s</> on project <magenta>%s</>\n", taskName, projectName)
				return nil
			},
		},
	},
}
