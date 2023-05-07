package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/gookit/color"
	"github.com/jasonwoodland/track/pkg/completion"
	"github.com/jasonwoodland/track/pkg/db"
	"github.com/jasonwoodland/track/pkg/mytime"

	// "github.com/jasonwoodland/track/pkg/mytime"
	"github.com/jasonwoodland/track/pkg/util"
	"github.com/jasonwoodland/track/pkg/view"
	"github.com/urfave/cli/v2"
)

var Daily = &cli.Command{
	Name:         "daily",
	Usage:        "Display daily report for time spent on projects and tasks",
	ArgsUsage:    "[project] [task]",
	BashComplete: completion.ProjectTaskCompletion,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "from",
			Aliases:  []string{"f"},
			Usage:    "Start date",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "to",
			Aliases: []string{"t"},
			Usage:   "End date",
		},
	},
	Action: func(c *cli.Context) error {
		// showFrames := c.Bool("frames")

		from := time.Time{}
		to := time.Now()
		if v := c.String("from"); v != "" {
			from = util.TimeFromShorthand(v).Add(-24 * time.Hour)
		} else {
			log.Fatalln("from flag required")
		}
		if v := c.String("to"); v != "" {
			to = util.TimeFromShorthand(v)
		} else {
			to = time.Now()
		}

		query := `
			with recursive dates(date) as (
				values(?)
				union all
				select strftime("%Y-%m-%dT%H:%M:%SZ", date(date, '+1 day'))
				from dates
				where date < ?
			)
			select
				dates.date,
		        	p.name,
		        	t.name,
		        	(
					select
						sum(strftime("%s", f2.end_time) - strftime("%s", f2.start_time))
					from frame f2
					where
						strftime('%Y-%m-%d', f2.end_time) = strftime('%Y-%m-%d', dates.date)
				) as total,
		                (
		                	select
		                		sum(strftime("%s", f2.end_time) - strftime("%s", f2.start_time)) as total
		                	from frame f2
		                	where
		                		f2.task_id = t.id
		                	and
						strftime('%Y-%m-%d', f2.end_time) = strftime('%Y-%m-%d', dates.date)
		                	and
		                		f2.end_time not like '0001-%'
		                ) as task_total,
		                (
		                	select
		                		sum(strftime("%s", f2.end_time) - strftime("%s", f2.start_time)) as total
		                	from frame f2
		                	left join task t2 on t2.id = f2.task_id
		                	where
		                		t2.project_id = p.id
		                	and
						strftime('%Y-%m-%d', f2.end_time) = strftime('%Y-%m-%d', dates.date)
		                	and
		                		f2.end_time not like '0001-%'
		                ) as project_total
			from dates
			left join frame f on strftime('%Y-%m-%d', f.end_time) = strftime('%Y-%m-%d', dates.date)
			left join task t on t.id = f.task_id
			left join project p on p.id = t.project_id
			group by t.id, dates.date
			order by dates.date
		`

		params := []interface{}{
			from.Format("2006-01-02"),
			to.Format("2006-01-02"),
		}

		rows, err := db.Db.Query(query, params...)
		if err != nil {
			log.Fatal(err)
		}
		var prevDate time.Time
		var prevProj string

		type row struct {
			date            time.Time
			projectName     string
			taskName        string
			totalDuration   time.Duration
			taskDuration    time.Duration
			projectDuration time.Duration
		}

		var totalDuration time.Duration

		for rows.Next() {
			r := row{}
			rows.Scan(
				(*mytime.Time)(&r.date),
				&r.projectName,
				&r.taskName,
				&r.totalDuration,
				&r.taskDuration,
				&r.projectDuration,
			)

			r.totalDuration *= time.Second
			r.taskDuration *= time.Second
			r.projectDuration *= time.Second
			dateFmt := "Mon Jan 02"

			if prevDate != r.date {
				if prevDate != (time.Time{}) {
					color.Println()
				}
				prevDate = r.date
				prevProj = ""
				color.Printf(
					view.DailyDateHours,
					r.date.Format(dateFmt),
					util.GetHours(r.totalDuration),
				)
			}

			if r.projectName != "" && prevProj != r.projectName {
				// if prevProj != "" && prevProj != r.projectName {
				// 	color.Println()
				// }
				prevProj = r.projectName
				color.Printf(
					view.DailyHoursProject,
					util.GetHours(r.projectDuration),
					r.projectName,
				)
			}
			if r.taskName != "" {
				color.Printf(
					view.DailyHoursTask,
					util.GetHours(r.taskDuration),
					50,
					r.taskName,
				)
			}
			totalDuration += r.totalDuration
		}
		fmt.Println()
		fmt.Printf(view.TotalHours, totalDuration.Hours())
		return nil
	},
}
