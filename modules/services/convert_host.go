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

	var instance BasicEC2Info
	err := getEC2(c.String("name"), &instance)
	if err != nil {
		return err
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
