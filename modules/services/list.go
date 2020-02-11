package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/urfave/cli"
)

type BasicEC2Info struct {
	Name             string
	InstanceId       string
	PrivateIpAddress string
	PublicIpAddress  string
	InstanceType     string
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
	region := "ap-northeast-1"
	s, err := session.NewSession()
	if err != nil {
		fmt.Println(err)
		return err
	}
	svc := ec2.New(s, &aws.Config{Region: &region})

	var input *ec2.DescribeInstancesInput
	if c.String("status") != "all" {
		filter := []*ec2.Filter{
			{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String(c.String("status")),
				},
			},
		}
		input = &ec2.DescribeInstancesInput{
			Filters: filter,
		}
	} else {
		input = nil
	}

	res, err := svc.DescribeInstances(input)

	if err != nil {
		fmt.Println(err)
		return err
	}

	var ec2List []BasicEC2Info
	for _, r := range res.Reservations {
		for _, instance := range r.Instances {
			nameTag := ""
			publicIpAddress := ""
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					nameTag = *tag.Value
				}
			}
			if instance.PublicIpAddress != nil {
				publicIpAddress = *instance.PublicIpAddress
			}
			ec2List = append(ec2List, BasicEC2Info{
				Name:             nameTag,
				InstanceId:       *instance.InstanceId,
				PrivateIpAddress: *instance.PrivateIpAddress,
				PublicIpAddress:  publicIpAddress,
				InstanceType:     *instance.InstanceType,
			})
		}
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
