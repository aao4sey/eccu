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

func ConvertHostCommand() *cli.Command {
	return &cli.Command{
		Name:  "conv-host",
		Usage: "",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "Host Name(AWS EC2 Name tag value)",
			},
			&cli.BoolFlag{
				Name:    "id",
				Aliases: []string{"i"},
				Usage:   "Convert to instance-id",
			},
			&cli.BoolFlag{
				Name:    "pip",
				Aliases: []string{"p"},
				Usage:   "Convert to private ip",
			},
			&cli.BoolFlag{
				Name:    "gip",
				Aliases: []string{"g"},
				Usage:   "Convert to public(global) ip",
			},
		},
		Action: services.ConvertHost,
	}
}
