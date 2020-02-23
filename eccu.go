package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/yukkyun/eccu/modules/commands"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			commands.ListCommand(),
			commands.FinderSearchCommand(),
			commands.ConvertHostCommand(),
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug, d",
				Value: false,
				Usage: "output debug level logs if you set this flag",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		os.Exit(1)
	}
}
