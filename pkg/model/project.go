package model

import (
	"log"

	"github.com/jasonwoodland/track/pkg/db"
)

type Project struct {
	Id   int64
	Name string
}

func GetProjects() (projects []*Project) {
	rows, err := db.Db.Query("select id, name from project")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		p := &Project{}
		rows.Scan(&p.Id, &p.Name)
		projects = append(projects, p)
	}
	return
}

func GetProjectById(id int64) (p *Project) {
	rows, err := db.Db.Query("select name from project where id = $1", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		p = &Project{
			Id: id,
		}
		rows.Scan(&p.Name)
	}
	return
}

func GetProjectByName(name string) (p *Project) {
	rows, err := db.Db.Query("select id from project where name = $1", name)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		p = &Project{
			Name: name,
		}
		rows.Scan(&p.Id)
	}
	return
}

func (p *Project) GetTask(name string) (t *Task) {
	rows, err := db.Db.Query("select id from task where project_id = $1 and name = $2", p.Id, name)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		t = &Task{
			Name:    name,
			Project: p,
		}
		rows.Scan(&t.Id)
	}
	return
}

func (p *Project) GetTasks() (tasks []*Task) {
	rows, err := db.Db.Query("select id, name from task where project_id = $1", p.Id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		t := &Task{
			Project: p,
		}
		rows.Scan(&t.Id, &t.Name)
		tasks = append(tasks, t)
	}
	return
}

func (p *Project) AddTask(name string) *Task {
	res, err := db.Db.Exec("insert into task (name, project_id) values ($1, $2)", name, p.Id)
	if err != nil {
		log.Fatal(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	return &Task{
		Id:      id,
		Name:    name,
		Project: p,
	}
}
