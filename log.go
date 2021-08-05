package main

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/urfave/cli/v2"
	"log"
	"strings"
	"time"
)

var Log = &cli.Command{
	Name:         "log",
	Usage:        "Display time spent on projects and tasks",
	ArgsUsage:    "[project] [task]",
	BashComplete: ProjectTaskCompletion,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "from",
			Aliases: []string{"f"},
			Usage:   "Start date from which to include tasks (TODO)",
		},
		&cli.StringFlag{
			Name:    "to",
			Usage:   "End date from which to include tasks (TODO)",
			Aliases: []string{"t"},
		},
		&cli.BoolFlag{
			Name:    "frames",
			Aliases: []string{"x"},
			Usage:   "Show individual frames for each task",
		},
	},
	Action: func(c *cli.Context) error {
		showFrames := c.Bool("frames")

		from := time.Time{}
		to := time.Now()
		if v := c.String("from"); v != "" {
			from = TimeFromShorthand(v)
		}
		if v := c.String("to"); v != "" {
			to = TimeFromShorthand(v)
		}

		query := `
			select
				p.name,
				t.name,
				sum(strftime("%s", end_time) - strftime("%s", start_time)) as total,

				-- We're only filtering by these max start_date/min end_date,
				-- So if any frames are within the --from/--to flags, we inlcude the
				-- in the results
				max(start_time) as start_date,
				min(end_time) as end_date,
				(
					select
						sum(strftime("%s", end_time) - strftime("%s", start_time)) as total
					from frame f2
					left join task t2 on t2.id = task_id
					where
						t2.project_id = p.id
					and
						end_time not like '0001-%'
				) as project_total
			from frame f
			left join task t on t.id = task_id
			left join project p on p.id = t.project_id
			group by task_id
			`

		var params []interface{}
		var whereConds []string

		whereConds = append(whereConds, "start_date > ?")
		params = append(params, from.Format(time.RFC3339))

		whereConds = append(whereConds, "end_date < ?")
		params = append(params, to.Format(time.RFC3339))

		if p := c.Args().Get(0); p != "" {
			whereConds = append(whereConds, "p.name like ?")
			params = append(params, "%"+p+"%")
		}

		if t := c.Args().Get(1); t != "" {
			whereConds = append(whereConds, "t.name like ?")
			params = append(params, "%"+t+"%")
		}

		// Add where conditions to query
		if len(whereConds) != 0 {
			query += "having\n" + strings.Join(whereConds, "\nand\n")
		}

		query += `
			order by (
				select max(start_time)
				from frame f3
				left join task t3 on t3.id = task_id
				where
					t3.project_id = p.id
			)
		`

		rows, err := Db.Query(query, params...)
		if err != nil {
			log.Fatal(err)
		}
		var prevProject string

		type row struct {
			projectName     string
			taskName        string
			totalDuration   time.Duration
			startDate       Time
			endDate         Time
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
				if prevProject != "" {
					fmt.Println()
				}
				color.Printf("Project: <magenta>%s</> (%.2f hour%s)\n", r.projectName, hours, s)
				prevProject = r.projectName
			}
			color.Printf("  <blue>%s</> (%s)\n", r.taskName, GetHours(r.totalDuration))

			if showFrames {
				frames := GetProjectByName(r.projectName).GetTask(r.taskName).GetFrames()
				for i, frame := range frames {
					// Don't print frames that fall outside of the --from/--to flags
					if frame.startTime.Before(from) {
						continue
					}
					if frame.endTime.After(to) {
						continue
					}
					color.Printf(
						"    <gray>[%v]</> <green>%s - %s</> <default>(%s)</>\n",
						i,
						frame.startTime.Format("Mon Jan 02 15:04"),
						frame.endTime.Format("15:04"),
						GetHours(frame.endTime.Sub(frame.startTime)),
					)
				}
			}
		}
		fmt.Println()
		return nil
	},
}
