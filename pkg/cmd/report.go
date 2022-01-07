package cmd

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gookit/color"
	"github.com/jasonwoodland/track/pkg/db"
	"github.com/jasonwoodland/track/pkg/mytime"
	"github.com/jasonwoodland/track/pkg/util"
	"github.com/jasonwoodland/track/pkg/view"
	"github.com/urfave/cli/v2"
)

var Report = &cli.Command{
	Name:      "report",
	Usage:     "Display monthly report for time spent on projects and tasks",
	ArgsUsage: "month",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "csv",
			Usage: "Output CSV format",
		},
		&cli.BoolFlag{
			Name:    "monthly",
			Aliases: []string{"m"},
			Usage:   "Output monthly tracked hours",
		},
	},
	Action: func(c *cli.Context) error {
		var (
			fromDate = util.MonthFromShorthand(c.Args().Get(0))
			toDate   = time.Date(fromDate.Year(), fromDate.Month()+1, 1, 0, 0, 0, 0, time.UTC)
			query    string
			params   []interface{}
			monthly  = c.Bool("monthly")
		)

		query = `
			select
				p.name,
				t.name,
				iif(
					t.monthly or ?,
					(select min(start_time) from frame where task_id = t.id and end_time > ?),
					min(f.start_time)
				) start_time,
				iif(
					t.monthly or ?,
				    (select max(end_time) from frame where task_id = t.id and end_time < ?),
					max(f.end_time)
				) end_time,
				iif(
					t.monthly or ?,
					sum(case when start_time > ? and end_time < ? then strftime("%s", end_time) - strftime("%s", start_time) else 0 end),
					sum(strftime("%s", end_time) - strftime("%s", start_time))
				) total,
				(t.monthly or ?) monthly
			from task t
			left join frame f on f.task_id = t.id
			left join project p on p.id = t.project_id
			group by t.id
			having
				(end_time > ? and end_time < ?) or ((monthly or ?) = true and total > 0)
			order by p.name, start_time;
		`

		params = []interface{}{
			monthly,
			fromDate.Format(time.RFC3339),
			monthly,
			toDate.Format(time.RFC3339),
			monthly,
			fromDate.Format(time.RFC3339),
			toDate.Format(time.RFC3339),
			monthly,
			fromDate.Format(time.RFC3339),
			toDate.Format(time.RFC3339),
			monthly,
		}

		rows, err := db.Db.Query(query, params...)
		if err != nil {
			log.Fatal(err)
		}

		type row struct {
			projectName  string
			taskName     string
			startDate    time.Time
			endDate      time.Time
			taskDuration time.Duration
			monthly      bool
		}

		if c.Bool("csv") {
			w := csv.NewWriter(os.Stdout)

			w.Write([]string{
				"Project",
				"Task",
				"Start",
				"End",
				"Total",
			})

			numRows := 0

			for rows.Next() {
				numRows++
				r := row{}
				rows.Scan(
					&r.projectName,
					&r.taskName,
					(*mytime.Time)(&r.startDate),
					(*mytime.Time)(&r.endDate),
					&r.taskDuration,
					&r.monthly,
				)
				r.taskDuration *= time.Second

				marker := ""
				if r.monthly {
					marker = "*"
				}

				if err := w.Write([]string{
					r.projectName,
					r.taskName + marker,
					r.startDate.Format("Mon Jan 02 2006"),
					r.endDate.Format("Mon Jan 02 2006"),
					fmt.Sprintf("%.2f", r.taskDuration.Hours()),
				}); err != nil {
					log.Fatalln("error outputting csv:", err)
				}
			}

			w.Write([]string{
				"Total",
				"",
				"",
				"",
				fmt.Sprintf("=SUM(E2:E%d)", numRows+1),
			})

			w.Flush()
		} else {
			var lastProjectName string

			for rows.Next() {
				r := row{}
				rows.Scan(
					&r.projectName,
					&r.taskName,
					(*mytime.Time)(&r.startDate),
					(*mytime.Time)(&r.endDate),
					&r.taskDuration,
					&r.monthly,
				)
				r.taskDuration *= time.Second

				if lastProjectName != r.projectName {
					if lastProjectName != "" {
						color.Println()
					}
					color.Printf(view.Project, r.projectName)
				}

				marker := ""
				if r.monthly {
					marker = "*"
				}

				color.Printf(
					view.FrameTimesDurationTask,
					r.startDate.Format("Mon Jan 02"),
					r.endDate.Format("Mon Jan 02"),
					util.GetHours(r.taskDuration),
					50,
					r.taskName+marker,
				)

				lastProjectName = r.projectName
			}

			color.Println()
		}

		return nil
	},
}
