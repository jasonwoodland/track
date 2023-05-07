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
	"github.com/jasonwoodland/track/pkg/view"
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
			Usage:   "Start date from which to include frames",
		},
		&cli.StringFlag{
			Name:    "to",
			Usage:   "End date from which to include frames",
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
				(
					select
						sum(strftime("%s", f2.end_time) - strftime("%s", f2.start_time)) as total
					from frame f2
					where
						f2.task_id = t.id
					and
						f2.end_time >= ?
					and
						f2.end_time <= ?
					and
						f2.end_time not like '0001-%'
				) as task_total,
				(
					select
						sum(strftime("%s", f2.end_time) - strftime("%s", f2.start_time)) as total
					from frame f2
					left join task t2 on t2.id = task_id
					where
						t2.project_id = p.id
					and
						f2.end_time >= ?
					and
						f2.end_time <= ?
					and
						f2.end_time not like '0001-%'
				) as project_total
			from frame f
			left join task t on t.id = task_id
			left join project p on p.id = t.project_id
			where
				f.end_time >= ?
			and
				f.end_time <= ?
			group by task_id
			`

		params := []interface{}{
			from.Format("2006-01-02"),
			to.Format("2006-01-02"),
			from.Format("2006-01-02"),
			to.Format("2006-01-02"),
			from.Format("2006-01-02"),
			to.Format("2006-01-02"),
		}

		var whereConds []string

		whereConds = append(whereConds, "f.end_time >= ?")
		params = append(params, from.Format("2006-01-02"))

		whereConds = append(whereConds, "f.end_time <= ?")
		params = append(params, to.Format("2006-01-02"))

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
				color.Printf(view.ProjectHours, r.projectName, hours)
				prevProject = r.projectName
			}

			color.Printf(
				view.FrameTimesDurationTask,
				r.startDate.Format("Mon Jan 02"),
				r.endDate.Format("Mon Jan 02 2006"),
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
						view.FrameTimesDuration,
						i,
						frame.StartTime.Format("Mon Jan 02 15:04"),
						frame.EndTime.Format("15:04"),
						util.GetHours(frame.EndTime.Sub(frame.StartTime)),
					)
				}
				fmt.Println()
			}
		}
		fmt.Println()
		fmt.Printf(view.TotalHours, totalDuration.Hours())
		return nil
	},
}
