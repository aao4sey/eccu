package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/urfave/cli"
)

func getEc2Client() *ec2.EC2 {
	region := "ap-northeast-1"
	s, err := session.NewSession()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return ec2.New(s, &aws.Config{Region: &region})
}

func getEc2List(c *cli.Context) ([]BasicEC2Info, error) {
	svc := getEc2Client()

	var input *ec2.DescribeInstancesInput
	if c.String("status") != "all" && c.String("status") != "" {
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
		return nil, err
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
				InstanceState:    *instance.State.Name,
			})
		}
	}
	return ec2List, nil
}

func getEC2(name string) (BasicEC2Info, error) {
	svc := getEc2Client()
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Name"),
				Values: []*string{
					aws.String(name),
				},
			},
		},
	}
	res, err := svc.DescribeInstances(input)
	if err != nil {
		fmt.Println(err)
	}

	nameTag := ""
	publicIpAddress := ""
	instance := res.Reservations[0].Instances[0]
	for _, tag := range instance.Tags {
		if *tag.Key == "Name" {
			nameTag = *tag.Value
		}
	}
	if instance.PublicIpAddress != nil {
		publicIpAddress = *instance.PublicIpAddress
	}
	ec2Info := BasicEC2Info{
		Name:             nameTag,
		InstanceId:       *instance.InstanceId,
		PrivateIpAddress: *instance.PrivateIpAddress,
		PublicIpAddress:  publicIpAddress,
		InstanceType:     *instance.InstanceType,
		InstanceState:    *instance.State.Name,
	}
	return ec2Info, nil
}
