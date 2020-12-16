package services

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/urfave/cli/v2"
)

type AwsSdkWrapper interface {
	DescribeInstances(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error)
}

type EC2Service struct {
	wrap AwsSdkWrapper
}

type AwsClient struct{}

func getEc2Client() (*ec2.EC2, error) {
	region := "ap-northeast-1"
	s, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	return ec2.New(s, &aws.Config{Region: &region}), nil
}

func (ac *AwsClient) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	svc, err := getEc2Client()
	if err != nil {
		return nil, err
	}
	res, err := svc.DescribeInstances(input)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (es *EC2Service) GetEC2List(c *cli.Context) (*[]BasicEC2Info, error) {
	var ec2List []BasicEC2Info

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
	res, err := es.wrap.DescribeInstances(input)
	if err != nil {
		return nil, err
	}

	for _, r := range res.Reservations {
		for _, instance := range r.Instances {
			nameTag := ""
			publicIPAddress := ""
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					nameTag = *tag.Value
				}
			}
			if instance.PublicIpAddress != nil {
				publicIPAddress = *instance.PublicIpAddress
			}
			ec2List = append(ec2List, BasicEC2Info{
				Name:             nameTag,
				InstanceID:       *instance.InstanceId,
				PrivateIPAddress: *instance.PrivateIpAddress,
				PublicIPAddress:  publicIPAddress,
				InstanceType:     *instance.InstanceType,
				InstanceState:    *instance.State.Name,
			})
		}
	}
	return &ec2List, nil
}

func (es *EC2Service) GetEC2(name string) (*BasicEC2Info, error) {
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
	res, err := es.wrap.DescribeInstances(input)
	if err != nil {
		return nil, err
	}

	nameTag := ""
	publicIPAddress := ""
	instance := res.Reservations[0].Instances[0]
	for _, tag := range instance.Tags {
		if *tag.Key == "Name" {
			nameTag = *tag.Value
		}
	}
	if instance.PublicIpAddress != nil {
		publicIPAddress = *instance.PublicIpAddress
	}
	ec2Info := &BasicEC2Info{
		Name:             nameTag,
		InstanceID:       *instance.InstanceId,
		PrivateIPAddress: *instance.PrivateIpAddress,
		PublicIPAddress:  publicIPAddress,
		InstanceType:     *instance.InstanceType,
		InstanceState:    *instance.State.Name,
	}
	return ec2Info, nil
}
