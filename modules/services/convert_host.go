package services

import (
	"fmt"

	"github.com/urfave/cli"
)

func ConvertHost(c *cli.Context) error {
	if c.String("name") == "" {
		fmt.Errorf("-n --name is required")
	}
	instance, err := getEC2(c.String("name"))
	if err != nil {
		fmt.Println(err)
	}

	if c.Bool("id") {
		fmt.Println(instance.InstanceId)
	} else if c.Bool("pip") {
		fmt.Println(instance.PrivateIpAddress)
	} else if c.Bool("gip") {
		fmt.Println(instance.PublicIpAddress)
	} else {
		fmt.Println(instance.InstanceId)
	}

	return nil
}
