package cmd

import (
	"log"
	"time"

	"github.com/gookit/color"
	"github.com/jasonwoodland/track/pkg/db"
	"github.com/jasonwoodland/track/pkg/mytime"
	"github.com/jasonwoodland/track/pkg/util"
	"github.com/urfave/cli/v2"
)

var Report = &cli.Command{
	Name:      "report",
	Usage:     "Display monthly report for time spent on projects and tasks",
	ArgsUsage: "month",
	// BashComplete: ProjectTaskCompletion,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "csv",
			Usage: "Output CSV format",
		},
	},
	Action: func(c *cli.Context) error {
		fromDate := util.MonthFromShorthand(c.Args().Get(0))
		toDate := time.Date(fromDate.Year(), fromDate.Month()+1, 1, 0, 0, 0, 0, time.UTC)

		query := `
			select
				p.name,
				t.name,
				min(f.start_time) start_time,
				max(f.end_time) end_time,
				sum(strftime("%s", end_time) - strftime("%s", start_time)) total
			from task t
			left join frame f on f.task_id = t.id
			left join project p on p.id = t.project_id
			group by t.id
			having
				end_time > ? and end_time < ?
			order by p.name, start_time;
		`

		params := []interface{}{
			fromDate.Format(time.RFC3339),
			toDate.Format(time.RFC3339),
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
		}

		var lastProjectName string

		for rows.Next() {
			r := row{}
			rows.Scan(
				&r.projectName,
				&r.taskName,
				(*mytime.Time)(&r.startDate),
				(*mytime.Time)(&r.endDate),
				&r.taskDuration,
			)
			r.taskDuration *= time.Second

			if lastProjectName != r.projectName {
				if lastProjectName != "" {
					color.Println()
				}
				color.Printf("Project: <magenta>%s</>\n", r.projectName)
			}

			color.Printf(
				"  <green>%s - %s</> %7s <blue>%-*s</>\n",
				r.startDate.Format("Mon Jan 02"),
				r.endDate.Format("Mon Jan 02"),
				util.GetHours(r.taskDuration),
				50,
				r.taskName,
			)

			lastProjectName = r.projectName
		}

		color.Println()

		return nil
	},
}
