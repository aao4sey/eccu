package services

import (
	"fmt"
	"log"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/urfave/cli/v2"
)

func formatTags(tags []EC2Tag) string {
	s := ""
	for _, tag := range tags {
		s += tag.Key + ": " + tag.Value + "\n"
	}
	return s
}

func formatEC2Info(ec2info BasicEC2Info) string {
	result := "[[ Basic Information ]]\n"
	result += fmt.Sprintf("%-30s: %s\n", "Name", ec2info.Name)
	result += fmt.Sprintf("%-30s: %s\n", "InstanceId", ec2info.InstanceId)
	result += fmt.Sprintf("%-30s: %s\n", "PrivateIP", ec2info.PrivateIpAddress)
	result += fmt.Sprintf("%-30s: %s\n", "PublicIP", ec2info.PublicIpAddress)
	result += fmt.Sprintf("%-30s: %s\n", "InstanceState", ec2info.InstanceState)
	result += fmt.Sprintf("%-30s: %s\n", "InstanceType", ec2info.InstanceType)
	result += "\n[[ Tag Information ]]\n"
	for _, tag := range ec2info.Tags {
		result += fmt.Sprintf("%-30s: %s\n", tag.Key, tag.Value)
	}
	return result
}

func IsValid(status string) bool {
	var preDefinedEc2Status = [...]string{
		"pending",
		"running",
		"stopping",
		"stopped",
		"shutting-down",
		"terminated",
	}
	fmt.Println(status)
	for _, s := range preDefinedEc2Status {
		if s == status {
			return true
		}
	}
	return false
}

func EC2FuzzySearch(c *cli.Context) error {
	if !IsValid(c.String("status")) && c.String("status") != "" {
		var message interface{}
		message = "Status is unexpected."
		return cli.Exit(message, 1)
	}
	ec2List, err := getEc2List(c)
	if err != nil {
		fmt.Println(err)
	}
	idx, err := fuzzyfinder.FindMulti(
		ec2List,
		func(i int) string {
			return ec2List[i].Name
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return formatEC2Info(ec2List[i])
		}))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(formatEC2Info(ec2List[idx[0]]))
	return nil
}
