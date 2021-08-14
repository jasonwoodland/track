package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/jasonwoodland/track/pkg/completion"
	"github.com/jasonwoodland/track/pkg/db"
	"github.com/jasonwoodland/track/pkg/model"
	"github.com/jasonwoodland/track/pkg/mytime"
	"github.com/jasonwoodland/track/pkg/util"
	"github.com/urfave/cli/v2"
)

var Log = &cli.Command{
	Name:         "log",
	Usage:        "Display time spent on projects and tasks",
	ArgsUsage:    "[project] [task]",
	BashComplete: completion.ProjectTaskCompletion,
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
			from = util.TimeFromShorthand(v)
		}
		if v := c.String("to"); v != "" {
			to = util.TimeFromShorthand(v)
		}

		query := `
			select
				p.name,
				t.name,
				sum(strftime("%s", end_time) - strftime("%s", start_time)) as total,

				-- We're only filtering by these max start_date/min end_date,
				-- So if any frames are within the --from/--to flags, we inlcude the
				-- in the results
				min(start_time) as start_date,
				max(end_time) as end_date,
				sum(strftime("%s", end_time) - strftime("%s", start_time)) task_total,
				(
					select
						sum(strftime("%s", end_time) - strftime("%s", start_time)) as total
					from frame f2
					left join task t2 on t2.id = task_id
					where
						t2.project_id = p.id
					and
						start_time > ?
					and
						end_time < ?
					and
						end_time not like '0001-%'
				) as project_total
			from frame f
			left join task t on t.id = task_id
			left join project p on p.id = t.project_id
			group by task_id
			`

		params := []interface{}{
			from.Format(time.RFC3339),
			to.Format(time.RFC3339),
		}

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
			order by p.name, start_date
		`

		rows, err := db.Db.Query(query, params...)
		if err != nil {
			log.Fatal(err)
		}
		var prevProject string

		type row struct {
			projectName     string
			taskName        string
			totalDuration   time.Duration
			startDate       time.Time
			endDate         time.Time
			taskDuration    time.Duration
			projectDuration time.Duration
		}

		var totalDuration time.Duration

		for rows.Next() {
			r := row{}
			rows.Scan(
				&r.projectName,
				&r.taskName,
				&r.totalDuration,
				(*mytime.Time)(&r.startDate),
				(*mytime.Time)(&r.endDate),
				&r.taskDuration,
				&r.projectDuration,
			)

			r.totalDuration *= time.Second
			r.taskDuration *= time.Second
			r.projectDuration *= time.Second

			totalDuration += r.taskDuration

			if r.projectName != prevProject {
				hours := r.projectDuration.Hours()
				if prevProject != "" {
					fmt.Println()
				}
				color.Printf("Project: <magenta>%s</> %.2fh\n", r.projectName, hours)
				prevProject = r.projectName
			}

			color.Printf(
				"  <green>%s - %s</> %6s <blue>%-*s</>\n",
				r.startDate.Format("Mon Jan 02"),
				r.endDate.Format("Mon Jan 02"),
				util.GetHours(r.taskDuration),
				50,
				r.taskName,
			)

			if showFrames {
				frames := model.GetProjectByName(r.projectName).GetTask(r.taskName).GetFrames()

				for i, frame := range frames {
					// Don't print frames that fall outside of the --from/--to flags
					if frame.StartTime.Before(from) || frame.EndTime.After(to) {
						continue
					}

					color.Printf(
						"    <gray>[%v]</> <green>%s - %s</> %s\n",
						i,
						frame.StartTime.Format("Mon Jan 02 15:04"),
						frame.EndTime.Format("15:04"),
						util.GetHours(frame.EndTime.Sub(frame.StartTime)),
					)
				}
			}
		}
		fmt.Println()
		fmt.Printf("%s total\n", util.GetHours(totalDuration))
		return nil
	},
}
