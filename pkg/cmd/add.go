package cmd

import (
	"log"
	"time"

	"github.com/gookit/color"
	"github.com/jasonwoodland/track/pkg/completion"
	"github.com/jasonwoodland/track/pkg/db"
	"github.com/jasonwoodland/track/pkg/model"
	"github.com/jasonwoodland/track/pkg/util"
	"github.com/urfave/cli/v2"
)

var Add = &cli.Command{
	Name:         "add",
	Usage:        "Add a frame to a task",
	ArgsUsage:    "project task duration",
	BashComplete: completion.ProjectTaskFrameCompletion,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "offset",
			Aliases: []string{"o"},
			Usage:   "Duration to offset the frame by (eg. -o -5m will add a frame that finished 5 minutes ago)",
		},
	},
	Action: func(c *cli.Context) error {
		if c.Args().Len() != 3 {
			cli.ShowSubcommandHelp(c)
			return nil
		}

		projectName := c.Args().Get(0)
		taskName := c.Args().Get(1)
		duration, err := time.ParseDuration(c.Args().Get(2))

		if err != nil {
			log.Fatalf("Bad duration: %s", c.Args().Get(2))
		}

		project := model.GetProjectByName(projectName)
		if project == nil {
			color.Printf("Project <magenta>%s</> doesn't exist\n", projectName)
			return nil
		}

		task := project.GetTask(taskName)
		if task == nil {
			color.Printf("Adding task <blue>%s</>\n", taskName)
			task = project.AddTask(taskName)
		}

		startTime := time.Now().Add(0 - duration)
		endTime := time.Now()

		if o, err := time.ParseDuration(c.String("offset")); err == nil {
			startTime = startTime.Add(o)
			endTime = endTime.Add(o)
		}

		db.Db.Exec(
			"insert into frame (task_id, start_time, end_time) values ($1, $2, $3)",
			task.Id,
			startTime.Format(time.RFC3339),
			endTime.Format(time.RFC3339),
		)

		color.Printf(
			"Added: <magenta>%s</> <blue>%s</> (%s, %s total)\033[K\n",
			project.Name,
			task.Name,
			util.GetHours(duration),
			util.GetHours(task.GetTotal()),
		)

		color.Printf(
			"  <gray>[%v]</> <green>%s - %s</> <default>(%s)</>\n",
			task.GetNumFrames()-1,
			startTime.Format("Mon Jan 02 15:04"),
			endTime.Format("15:04"),
			util.GetHours(endTime.Sub(startTime)),
		)
		return nil
	},
}
