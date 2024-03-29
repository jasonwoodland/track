package completion

import (
	"fmt"
	"strings"

	"github.com/jasonwoodland/track/pkg/model"
	"github.com/urfave/cli/v2"
)

func ShowFlagCompletion(c *cli.Context) bool {
	if c.Args().Len() > 0 && c.Args().Get(c.Args().Len() - 1)[0] == '-' {
		for _, f := range c.Command.Flags {
			desc := f.String()[strings.Index(f.String(), "\t")+1:]
			for _, n := range f.Names() {
				if len(n) == 1 {
					fmt.Printf("-%s", n)
				} else {
					fmt.Printf("--%s", n)
				}
				fmt.Printf(":%s\n", desc)
			}
		}
		return true
	}
	return false
}

func ProjectCompletion(c *cli.Context) {
	if ShowFlagCompletion(c) {
		return
	}

	if c.NArg() == 0 {
		for _, p := range model.GetProjects() {
			fmt.Println(p.Name)
		}
		return
	}
}

func ProjectTaskCompletion(c *cli.Context) {
	if ShowFlagCompletion(c) {
		return
	}

	if c.NArg() == 0 {
		for _, p := range model.GetProjects() {
			fmt.Println(p.Name)
		}
		return
	}

	p := model.GetProjectByName(c.Args().Get(0))

	if c.NArg() == 1 {
		for _, t := range p.GetTasks() {
			fmt.Println(t.Name)
		}
	}
}

func ProjectTaskProjectTaskCompletion(c *cli.Context) {
	if ShowFlagCompletion(c) {
		return
	}

	if c.NArg() == 0 {
		for _, p := range model.GetProjects() {
			fmt.Println(p.Name)
		}
		return
	}

	p := model.GetProjectByName(c.Args().Get(0))

	if c.NArg() == 1 {
		for _, t := range p.GetTasks() {
			fmt.Println(t.Name)
		}
	}

	if c.NArg() == 2 {
		for _, p := range model.GetProjects() {
			fmt.Println(p.Name)
		}
		return
	}

	p = model.GetProjectByName(c.Args().Get(0))

	if c.NArg() == 3 {
		for _, t := range p.GetTasks() {
			fmt.Println(t.Name)
		}
	}
}

func ProjectTaskFrameCompletion(c *cli.Context) {
	if ShowFlagCompletion(c) {
		return
	}

	if c.NArg() == 0 {
		for _, p := range model.GetProjects() {
			fmt.Println(p.Name)
		}
		return
	}

	p := model.GetProjectByName(c.Args().Get(0))

	if c.NArg() == 1 {
		for _, t := range p.GetTasks() {
			fmt.Println(t.Name)
		}
		return
	}

	t := p.GetTask(c.Args().Get(1))

	if c.NArg() == 2 {
		frames := t.GetFrames()
		for i, f := range frames {
			fmt.Printf(
				"%v:%s - %s\n",
				i,
				f.StartTime.Format("Mon Jan 02 15:04"),
				f.EndTime.Format("15:04"),
			)
		}
	}
}
func ProjectTaskFrameProjectTaskCompletion(c *cli.Context) {
	if ShowFlagCompletion(c) {
		return
	}

	if c.NArg() < 3 {
		if c.NArg() == 0 {
			for _, p := range model.GetProjects() {
				fmt.Println(p.Name)
			}
			return
		}

		p := model.GetProjectByName(c.Args().Get(0))

		if c.NArg() == 1 {
			for _, t := range p.GetTasks() {
				fmt.Println(t.Name)
			}
			return
		}

		t := p.GetTask(c.Args().Get(1))

		if c.NArg() == 2 {
			frames := t.GetFrames()
			for i, f := range frames {
				fmt.Printf(
					"%v:%s - %s\n",
					i,
					f.StartTime.Format("Mon Jan 02 15:04"),
					f.EndTime.Format("15:04"),
				)
			}
		}
	}

	if c.NArg() == 3 {
		for _, p := range model.GetProjects() {
			fmt.Println(p.Name)
		}
		return
	}

	p := model.GetProjectByName(c.Args().Get(0))

	if c.NArg() == 4 {
		for _, t := range p.GetTasks() {
			fmt.Println(t.Name)
		}
		return
	}
}
