package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
	command "github.com/yukkyun/eccu/modules/commands"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			command.ListCommand(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
