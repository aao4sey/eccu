package utilcommand

import (
	"github.com/urfave/cli/v2"
	services "github.com/yukkyun/eccu/modules/services"
)

func ListCommand() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "out",
				Aliases: []string{"o"},
				Usage:   "choose output type (default: tsv)",
			},
			&cli.StringFlag{
				Name:        "status",
				Aliases:     []string{"s"},
				Usage:       "choose ec2 status (default: -)",
				DefaultText: "all",
			},
		},
		Action: services.GetEc2List,
	}
}

func FinderSearchCommand() *cli.Command {
	return &cli.Command{
		Name:   "fs",
		Usage:  "",
		Action: services.FinderSearch,
	}
}
