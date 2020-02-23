package services

import (
	"fmt"

	"github.com/urfave/cli/v2"
	common "github.com/yukkyun/eccu/modules/services/common"
)

type BasicEC2Info struct {
	Name             string
	InstanceId       string
	PrivateIpAddress string
	PublicIpAddress  string
	InstanceType     string
	InstanceState    string
}

func (b *BasicEC2Info) ShowTsv() {
	fmt.Printf("%s\t%s\t%s\t%s\t%s\n",
		b.Name,
		b.InstanceId,
		b.PrivateIpAddress,
		b.PublicIpAddress,
		b.InstanceType,
	)
}

func (b *BasicEC2Info) ShowCsv() {
	fmt.Printf("%s,%s,%s,%s,%s\n",
		b.Name,
		b.InstanceId,
		b.PrivateIpAddress,
		b.PublicIpAddress,
		b.InstanceType,
	)
}

func GetEc2List(c *cli.Context) error {
	common.SetLogFilter(c.Bool("debug"))
	var ec2List []BasicEC2Info
	err := getEc2List(c, &ec2List)
	if err != nil {
		return err
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
