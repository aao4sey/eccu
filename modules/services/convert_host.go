package services

import (
	"fmt"
	"log"

	"github.com/urfave/cli/v2"
	common "github.com/yukkyun/eccu/modules/services/common"
)

func ConvertHost(c *cli.Context) error {
	common.SetLogFilter(c.Bool("debug"))
	if c.String("name") == "" {
		log.Print("[ERROR] -n --name is required")
	}
	ec2svc := &EC2Service{
		&AwsClient{},
	}
	ec2Info, err := ec2svc.GetEC2(c.String("name"))
	if err != nil {
		return err
	}

	if c.Bool("id") {
		fmt.Println(ec2Info.InstanceID)
	} else if c.Bool("pip") {
		fmt.Println(ec2Info.PrivateIPAddress)
	} else if c.Bool("gip") {
		fmt.Println(ec2Info.PublicIPAddress)
	} else {
		fmt.Println(ec2Info.InstanceID)
	}

	return nil
}
