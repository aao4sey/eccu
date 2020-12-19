package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
	command "github.com/yukkyun/eccu/modules/commands"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			command.ListCommand(),
			command.FuzzySearchCommand(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
