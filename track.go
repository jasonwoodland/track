package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
        "strings"
        "github.com/gookit/color"

	"github.com/adrg/xdg"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
)

var db *sql.DB

type project struct {
    id int64
    name string
}

type task struct {
    id int64
    name string
    project *project
}

type frame struct {
    id int64
    task *task
    startTime time.Time
    endTime time.Time
}

type state struct {
    running bool
    task task
    startTime time.Time
    timeElapsed time.Duration
}

func getProjectById(id int64) (p *project) {
    rows, err := db.Query("select name from project where id = $1", id)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    if rows.Next() {
        p = &project{
            id: id,
        }
        rows.Scan(&p.name)
    }
    return
}

func getProjectByName(name string) (p *project) {
    rows, err := db.Query("select id from project where name = $1", name)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    if rows.Next() {
        p = &project{
            name: name,
        }
        rows.Scan(&p.id)
    }
    return
}

func (p *project) getTask(name string) (t *task) {
    rows, err := db.Query("select id from task where project_id = $1 and name = $2", p.id, name)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    if rows.Next() {
        t = &task{
            name: name,
            project: p,
        }
        rows.Scan(&t.id)
    }
    return
}

func (p *project) addTask(name string) *task {
    res, err := db.Exec("insert into task (name, project_id) values ($1, $2)", name, p.id)
    if err != nil {
        log.Fatal(err)
    }
    id, err := res.LastInsertId()
    if err != nil {
        log.Fatal(err)
    }
    return &task{
        id: id,
        name: name,
        project: p,
    }
}

func getProjects() (projects []*project) {
    rows, err := db.Query("select id, name from project")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    for rows.Next() {
        p := &project{}
        rows.Scan(&p.id, &p.name)
        projects = append(projects, p)
    }
    return
}

func getState() (s *state) {
    s = &state{}
    rows, err := db.Query("select task_id, start_time from frame where end_time is null")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    if rows.Next() {
        s.running = true
        var taskId int64
        var startTime string
        rows.Scan(&taskId, &startTime)
        s.task = getTaskById(taskId)
        s.startTime, _ = time.Parse(time.RFC3339, startTime)
        s.timeElapsed = time.Now().Sub(s.startTime)
    }
    return
}

func  getTaskById(id int64) (t task) {
    rows, err := db.Query("select id, name, project_id from task where id = $1", id)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    if rows.Next() {
        t = task{
            id: id,
        }
        var projectId int64
        rows.Scan(&t.id, &t.name, &projectId)
        t.project = getProjectById(projectId)
    }
    return
}

func timeFromShorthand(v string) (t time.Time) {
    layouts := []string{
        "2006",
        "01-02",
        "20060102",
        "2006-01-02",
    }
    for _, l := range layouts {
        if (len(l) == len(v)) {
            t, err := time.Parse(l, v)
            if err != nil {
                log.Fatal(err)
            }
            fmt.Printf("%v", t.Year())
            if t.Year() == 0 {
                t = t.AddDate(time.Now().Year(), 0, 0)
            }
            return t
        }
    }
    log.Fatalf("bad format provided: %s", v)
    return time.Time{}
}

func initDb(db *sql.DB) {
    db.Exec(`
        create table if not exists project (
              id integer primary key,
              name text
            );
          `)
    db.Exec(`
        create table if not exists task (
            id integer primary key,
            project_id integer,
            name text,

            foreign key(project_id) references project(id) on delete cascade
        );
    `)
    db.Exec(`
        create table if not exists frame (
            id integer primary key,
            task_id integer,
            start_time text,
            end_time text,

            foreign key(task_id) references task(id) on delete cascade
        );
    `)
}

type Time time.Time

func (t *Time) Scan(v interface{}) error {
    vt, err := time.Parse(time.RFC3339, string(v.(string)))
    if err != nil {
        return err
    }
    *t = Time(vt)
    return nil
}

func main() {
    dbFilePath, _ := xdg.DataFile("track-cli/db.sqlite3")
    db, _ = sql.Open("sqlite3", dbFilePath)
    initDb(db)
    defer db.Close()

    app := &cli.App{
        Name: "track",
        Usage: "track your time",

        Commands: cli.Commands{
            {
                Name: "start",
                Aliases: []string{"s"},
                Usage: "start tracking time for a task",
                ArgsUsage: "project task",
                Flags: []cli.Flag{
                    &cli.StringFlag{
                        Name: "ago",
                    },
                    &cli.StringFlag{
                        Name: "in",
                    },
                },
                Action: func(c *cli.Context) error {
                    startTime := time.Now()

                    projectName := c.Args().Get(0)
                    taskName := c.Args().Get(1)
                    if projectName == "" || taskName == "" {
                        cli.ShowSubcommandHelp(c)
                        return nil
                    }

                    state := getState()
                    if state != nil && state.running {
                        fmt.Println("task already running")
                        return nil
                    }

                    project := getProjectByName(projectName)
                    if project == nil {
                        color.Printf("project <magenta>%s</> doesn't exists\n", projectName)
                        return nil
                    }

                    task := project.getTask(taskName)
                    if task == nil {
                        color.Printf("adding task <blue>%s</>\n", taskName)
                        task = project.addTask(taskName)
                    }

                    if ago, err := time.ParseDuration(c.String("ago")); err == nil {
                        startTime = startTime.Add(0 - ago)
                    }

                    if in, err := time.ParseDuration(c.String("in")); err == nil {
                        startTime = startTime.Add(in)
                    }

                    if c.String("in") == "" {
                        color.Printf("started tracking time at <green>%s</>\n", startTime.Format("15:04"))
                    } else {
                        fmt.Printf("will start tracking time at %s\n", startTime.Format("15:04"))
                    }

                    db.Exec(
                        "insert into frame (task_id, start_time) values ($1, $2)",
                        task.id,
                        startTime.Format(time.RFC3339),
                    )
                    return nil
                },
            },

            {
                Name: "cancel",
                Aliases: []string{"c"},
                Usage: "cancel a running task",
                Action: func(c *cli.Context) error {
                    res, err := db.Exec("delete from frame where end_time is null")
                    db.Exec("delete from task where not exists (select 1 from frame where frame.task_id = task.id)")
                    if err != nil {
                        log.Fatal(err)
                    }
                    n, err := res.RowsAffected()
                    if n == 0 {
                        fmt.Println("no task started")
                    } else {
                        fmt.Println("task cancelled")
                    }
                    return nil
                },
            },

            {
                Name: "stop",
                Aliases: []string{"st"},
                Usage: "stop a running task",
                Flags: []cli.Flag{
                    &cli.StringFlag{
                        Name: "ago",
                    },
                    &cli.StringFlag{
                        Name: "in",
                    },
                },
                Action: func(c *cli.Context) error {
                    state := getState()
                    endTime := time.Now()

                    if ago, err := time.ParseDuration(c.String("ago")); err == nil {
                        endTime = endTime.Add(0 - ago)
                    }

                    if in, err := time.ParseDuration(c.String("in")); err == nil {
                        endTime = endTime.Add(in)
                    }

                    res, err := db.Exec(
                        "update frame set end_time = $1 where end_time is null",
                        endTime.Format(time.RFC3339),
                    )
                    if err != nil {
                        log.Fatal(err)
                    }

                    if n, _ := res.RowsAffected(); n == 0 {
                        fmt.Printf("no task started")
                    } else {
                        fmt.Printf("stopped tracking time at %s\n", endTime.Format("15:04"))
                        dur := endTime.Sub(state.startTime)
                        fmt.Printf("duration: %s\n", dur.Round(time.Second))
                    }

                    return nil
                },
            },

            {
                Name: "log",
                Usage: "display time spent on tasks",
                Flags: []cli.Flag{
                    &cli.StringFlag{
                        Name: "from",
                        Aliases: []string{"f"},
                    },
                    &cli.StringFlag{
                        Name: "to",
                        Aliases: []string{"t"},
                    },
                    &cli.StringFlag{
                        Name: "project",
                        Aliases: []string{"p"},
                    },
                    &cli.StringFlag{
                        Name: "task",
                        Aliases: []string{"T"},
                    },
                },
                Action: func(c *cli.Context) error {
                    from := time.Time{}
                    to := time.Now()
                    if v := c.String("from"); v != "" {
                        from = timeFromShorthand(v)
                    }
                    if v := c.String("to"); v != "" {
                        to = timeFromShorthand(v)
                    }

                    query := `
                        select
                            p.name,
                            t.name,
                            sum(strftime("%s", end_time) - strftime("%s", start_time)) as total,
                            min(start_time) as start_date,
                            max(end_time) as end_date,
                            (
                                select
                                    sum(strftime("%s", end_time) - strftime("%s", start_time)) as total
                                from frame f2
                                left join task t on t.id = task_id
                                left join project p on p.id = t.project_id
                                where
                                    f2.task_id = task_id
                                group by t.project_id
                            ) as project_total
                        from frame f
                        left join task t on t.id = task_id
                        left join project p on p.id = t.project_id
                    `
                    query += `
                        group by
                            task_id
                    `


                    var params []interface{}
                    var whereConds []string

                    whereConds = append(whereConds, "start_date > ?")
                    params = append(params, from.Format(time.RFC3339))

                    whereConds = append(whereConds, "end_date < ?")
                    params = append(params, to.Format(time.RFC3339))

                    if p := c.String("project"); p != "" {
                        whereConds = append(whereConds, "p.name like ?")
                        params = append(params, "%" + p + "%")
                    }

                    if t := c.String("task"); t != "" {
                        whereConds = append(whereConds, "t.name like ?")
                        params = append(params, "%" + t + "%")
                    }

                    // Add where conditions to query
                    if len(whereConds) != 0 {
                        query += "having\n" + strings.Join(whereConds, "\nand\n")
                    }

                    rows, err := db.Query(query, params...)
                    if err != nil {
                        log.Fatal(err)
                    }
                    var prevProject string

                    type row struct {
                        projectName string
                        taskName string
                        totalDuration time.Duration
                        startDate Time
                        endDate Time
                        projectDuration time.Duration
                    }

                    for rows.Next() {
                        r := row{}
                        rows.Scan(
                            &r.projectName,
                            &r.taskName,
                            &r.totalDuration,
                            (*Time)(&r.startDate),
                            (*Time)(&r.endDate),
                            &r.projectDuration,
                        )
                        r.totalDuration *= time.Second
                        r.projectDuration *= time.Second
                        if r.projectName != prevProject {
                            hours := r.projectDuration.Hours()
                            s := ""
                            if hours != 1 {
                                s = "s"
                            }
                            color.Printf("Project: <magenta>%s</> (%.2f hour%s)\n", r.projectName, hours, s)
                            prevProject = r.projectName
                        }
                        hours := r.totalDuration.Hours()
                        s := ""
                        if hours != 1 {
                            s = "s"
                        }
                        color.Printf("  <blue>%s</> (%.2f hour%s)\n", r.taskName, hours, s)
                    }
                    return nil
                },
            },

            {
                Name: "status",
                Usage: "display status of running task",
                Action: func(c *cli.Context) error {
                    state := getState()
                    color.Printf("Running: <magenta>%s</> ", state.task.project.name)
                    hours := state.timeElapsed.Hours()
                    s := ""
                    if hours != 1 {
                        s = "s"
                    }
                    color.Printf("<blue>%s</> (%.2f hour%s)\n", state.task.name, hours, s)
                    color.Printf("Started at <green>%s</> (%s ago)\n", state.startTime.Format("15:04"), state.timeElapsed.Round(time.Second))
                    return nil
                },
            },

            {
                Name: "projects",
                Usage: "list projects",
                Action: func(c *cli.Context) error {
                    for _, project := range getProjects() {
                        color.Magenta.Println(project.name)
                    }
                    return nil
                },
            },

            {
                Name: "project",
                Subcommands: []*cli.Command{
                    {
                        Name: "add",
                        Usage: "add a new project",
                        ArgsUsage: "name",
                        Action: func (c *cli.Context) error {
                            name := c.Args().Get(0)
                            if name == "" {
                                cli.ShowSubcommandHelp(c)
                                return nil
                            }
                            if getProjectByName(name) != nil {
                                color.Printf("project <magenta>%s</> already exists\n", name)
                                return nil
                            }
                            db.Exec("insert into project (name) values ($1)", name)
                            color.Printf("added project <magenta>%s</>\n", name)
                            return nil
                        },
                    },
                    {
                        Name: "rename",
                        Usage: "rename a new project",
                        ArgsUsage: "old_name new_name",
                        Action: func (c *cli.Context) error {
                            oldName := c.Args().Get(0)
                            newName := c.Args().Get(1)
                            if oldName == "" || newName == "" {
                                cli.ShowSubcommandHelp(c)
                                return nil
                            }
                            if getProjectByName(oldName) == nil {
                                fmt.Printf("project \"%s\" doesn't exists\n", oldName)
                                return nil
                            }
                            db.Exec("update project set name = $1 where name = $2", newName, oldName)
                            fmt.Printf("renamed project \"%s\" to \"%s\"\n", oldName, newName)
                            return nil
                        },
                    },
                    {
                        Name: "remove",
                        Aliases: []string{"rm"},
                        Usage: "delete a project",
                        ArgsUsage: "name",
                        Action: func (c *cli.Context) error {
                            name := c.Args().Get(0)
                            if name == "" {
                                cli.ShowSubcommandHelp(c)
                                return nil
                            }
                            if getProjectByName(name) == nil {
                                fmt.Printf("project \"%s\" doesn't exists\n", name)
                                return nil
                            }
                            db.Exec("delete from project where name = $1", name)
                            fmt.Printf("deleted project \"%s\"\n", name)
                            return nil
                        },
                    },
                },
            },
        },
    }

    err := app.Run(os.Args)
    if err != nil {
        log.Fatal(err)
    }
}
