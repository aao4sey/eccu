package services

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

type BasicEC2Info struct {
	Name             string
	InstanceId       string
	PrivateIpAddress string
	PublicIpAddress  string
	InstanceType     string
	InstanceState    string
	Tags             []EC2Tag
}

type EC2Tag struct {
	Key   string
	Value string
}

func (b *BasicEC2Info) ShowTsv() {
	fmt.Printf("%s\t%s\t%s\t%s\t%s\t%s\n",
		b.Name,
		b.InstanceId,
		b.PrivateIpAddress,
		b.PublicIpAddress,
		b.InstanceType,
		b.InstanceState,
	)
}

func (b *BasicEC2Info) ShowCsv() {
	fmt.Printf("%s,%s,%s,%s,%s,%s\n",
		b.Name,
		b.InstanceId,
		b.PrivateIpAddress,
		b.PublicIpAddress,
		b.InstanceType,
		b.InstanceState,
	)
}

func GetEc2List(c *cli.Context) error {
	ec2List, err := getEc2List(c)
	if err != nil {
		fmt.Println(err)
	}

	for _, b := range ec2List {
		switch c.String("out") {
		case "tsv":
			b.ShowTsv()
		case "csv":
			b.ShowCsv()
		default:
			b.ShowTsv()
		}
	}
	return nil
}
