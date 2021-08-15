package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/jasonwoodland/track/pkg/cleanup"
	"github.com/jasonwoodland/track/pkg/cmd"
	"github.com/jasonwoodland/track/pkg/db"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
)

func main() {
	db.OpenDb()
	defer db.Db.Close()

	// Fix escaped hashes which zsh completion adds to cli args
	for i, a := range os.Args {
		os.Args[i] = strings.ReplaceAll(a, "\\#", "#")
	}

	// On SIGINT, send the escape seq. to switch back from the alternate screen.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println()
		fmt.Printf("\033[?1049l")
		cleanup.Cleanup()
		os.Exit(0)
	}()

	app := &cli.App{
		Name:                   "track",
		Usage:                  "Track time for projects and tasks",
		EnableBashCompletion:   true,
		UseShortOptionHandling: true,

		Commands: cli.Commands{
			cmd.Start,
			cmd.Cancel,
			cmd.Stop,
			cmd.Add,
			cmd.Timeline,
			cmd.Log,
			cmd.Status,
			cmd.ProjectCmds,
			cmd.TaskCmds,
			cmd.Projects,
			cmd.FrameCmds,
			cmd.Report,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
