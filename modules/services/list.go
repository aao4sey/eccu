package services

import (
	"encoding/json"

	"github.com/urfave/cli/v2"
	common "github.com/yukkyun/eccu/modules/services/common"
)

func GetEc2List(c *cli.Context) error {
	common.SetLogFilter(c.Bool("debug"))

	// Get cache
	var ec2List *[]BasicEC2Info
	cache, err := getCache()
	if err == nil {
		json.Unmarshal([]byte(*cache), &ec2List)
	} else {
		ec2svc := &EC2Service{
			&AwsClient{},
		}
		ec2List, err = ec2svc.GetEC2List(c)
		if err != nil {
			return err
		}

		// Write cache
		result, _ := json.Marshal(&ec2List)
		err = putCache(string(result))
		if err != nil {
			return err
		}
	}

	for _, b := range *ec2List {
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
