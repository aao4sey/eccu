package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/urfave/cli/v2"
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

func getNameTag(tags []*ec2.Tag) string {
	name := ""
	for _, tag := range tags {
		if *tag.Key == "Name" {
			name = *tag.Value
		}
	}
	return name
}

func getTags(tags []*ec2.Tag) []EC2Tag {
	var returnTags []EC2Tag
	for _, tag := range tags {
		returnTags = append(returnTags, EC2Tag{
			Key:   *tag.Key,
			Value: *tag.Value,
		})
	}
	return returnTags
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
			publicIpAddress := ""
			privateIpAddress := ""

			if instance.PublicIpAddress != nil {
				publicIpAddress = *instance.PublicIpAddress
			}
			if instance.PrivateIpAddress != nil {
				privateIpAddress = *instance.PrivateIpAddress
			}

			ec2List = append(ec2List, BasicEC2Info{
				Name:             getNameTag(instance.Tags),
				InstanceId:       *instance.InstanceId,
				PrivateIpAddress: privateIpAddress,
				PublicIpAddress:  publicIpAddress,
				InstanceType:     *instance.InstanceType,
				InstanceState:    *instance.State.Name,
				Tags:             getTags(instance.Tags),
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

	publicIpAddress := ""
	instance := res.Reservations[0].Instances[0]
	if instance.PublicIpAddress != nil {
		publicIpAddress = *instance.PublicIpAddress
	}
	ec2Info := BasicEC2Info{
		Name:             getNameTag(instance.Tags),
		InstanceId:       *instance.InstanceId,
		PrivateIpAddress: *instance.PrivateIpAddress,
		PublicIpAddress:  publicIpAddress,
		InstanceType:     *instance.InstanceType,
		InstanceState:    *instance.State.Name,
		Tags:             getTags(instance.Tags),
	}
	return ec2Info, nil
}
