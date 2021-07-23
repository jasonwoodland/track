package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
)

var Cleanup = func() {}

func main() {
	openDb()
	defer Db.Close()

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
		Cleanup()
		os.Exit(0)
	}()

	app := &cli.App{
		Name:                 "track",
		Usage:                "Track time for projects and tasks",
		EnableBashCompletion: true,

		Commands: cli.Commands{
			Start,
			Cancel,
			Stop,
			Timeline,
			Log,
			Status,
			ProjectCmds,
			TaskCmds,
			Projects,
			FrameCmds,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
