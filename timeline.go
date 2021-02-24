package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/urfave/cli/v2"
)

var Timeline = &cli.Command{
	Name:         "timeline",
	Usage:        "Display a timeline showing time spent on tasks for a given date range",
	ArgsUsage:    "[project] [task]",
	BashComplete: ProjectTaskCompletion,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "from",
			Aliases:  []string{"f"},
			Usage:    "Start date for the timeline",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "to",
			Usage:   "End date for the timeline",
			Aliases: []string{"t"},
		},
	},
	Action: func(c *cli.Context) error {
		from := time.Time{}
		if v := c.String("from"); v != "" {
			from = TimeFromShorthand(v)
		}

		to := time.Now()
		if v := c.String("to"); v != "" {
			to = TimeFromShorthand(v)
		}

		to = to.Add(-24 * time.Hour)

		query := `
			with recursive date(d) as (
				select datetime(?)
				union all
				select datetime(d, '+1 day') from date where d < ?
			)
			select
				d,
				p.id,
				p.name,
				t.id,
				t.name
			from date
			left join frame f on f.start_time > date.d and f.end_time < datetime(d, '+1 day')
			left join task t on t.id = f.task_id
			left join project p on p.id = t.project_id
		`

		var params []interface{}
		var whereConds []string

		params = append(params, from.Format(time.RFC3339))
		params = append(params, to.Format(time.RFC3339))

		if p := c.Args().Get(0); p != "" {
			whereConds = append(whereConds, "(p.name like ? or p.name is null)")
			params = append(params, "%"+p+"%")
		}

		if t := c.Args().Get(1); t != "" {
			whereConds = append(whereConds, "(t.name like ? or t.name is null)")
			params = append(params, "%"+t+"%")
		}

		// Add where conditions to query
		if len(whereConds) != 0 {
			query += "where\n" + strings.Join(whereConds, "\nand\n")
		}

		rows, err := Db.Query(query, params...)
		if err != nil {
			log.Fatal(err)
		}
		var chart = make(map[string]map[int]bool)
		var tasks = make(map[int]string)
		var projects = make(map[int]string)

		for rows.Next() {
			var date string
			var projectId int
			var projectName string
			var taskId int
			var taskName string

			rows.Scan(
				&date,
				&projectId,
				&projectName,
				&taskId,
				&taskName,
			)

			if chart[date] == nil {
				chart[date] = make(map[int]bool)
			}
			if projectName != "" {
				projects[projectId] = projectName
			}
			if taskName != "" {
				tasks[taskId] = taskName
			}
			chart[date][taskId] = true
		}

		longest := 0
		for _, t := range tasks {
			if len(t) > longest {
				longest = len(t)
			}
		}

		longestProject := 0
		for _, p := range projects {
			if len(p) > longestProject {
				longestProject = len(p)
			}
		}

		dates := make([]string, 0, len(chart))
		for d := range chart {
			dates = append(dates, d)
		}
		sort.Strings(dates)

		fmt.Printf(strings.Repeat(" ", longest+longestProject+2))
		for _, date := range dates {
			d, _ := time.Parse("2006-01-02 00:00:00", date)
			color.Printf("<gray>%3v</>", d.Day())

		}
		fmt.Printf("\n")

		taskIds := make([]int, 0, len(tasks))
		for k := range tasks {
			taskIds = append(taskIds, k)
		}
		sort.Ints(taskIds)

		for _, taskId := range taskIds {
			color.Printf("<magenta>%-"+strconv.Itoa(longestProject)+"v</> ", GetTaskById(int64(taskId)).project.name)
			color.Printf("<blue>%-"+strconv.Itoa(longest)+"v</> <gray>┃</>", tasks[taskId])
			for di, date := range dates {
				if chart[date][taskId] {
					var prev, next bool
					if di > 0 {
						prev = chart[dates[di-1]][taskId]
					}
					if di < len(dates)-1 {
						next = chart[dates[di+1]][taskId]
					}
					if prev && next {
						color.Printf("<green>━●━</>")
					} else if prev {
						color.Printf("<green>━● </>")
					} else if next {
						color.Printf("<green> ●━</>")
					} else {
						color.Printf("<green> ● </>")
					}
				} else {
					fmt.Printf("   ")
				}
			}
			color.Printf("<gray>┃</>\n")
		}
		return nil
	},
}
