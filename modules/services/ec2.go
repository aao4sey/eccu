package services

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli/v2"
)

const (
	cacheExpireTime = 15
)

func getCacheDirPath(dir *string) error {
	homedir, err := homedir.Dir()
	if err != nil {
		return err
	}
	*dir = homedir + "/.config/eccu"
	return nil
}

func getEc2Client() (*ec2.EC2, error) {
	region := "ap-northeast-1"
	s, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	return ec2.New(s, &aws.Config{Region: &region}), nil
}

func getCache(data *string) error {
	var cacheDirPath string
	err := getCacheDirPath(&cacheDirPath)
	if err != nil {
		return err
	}

	cacheFilePath := cacheDirPath + "/.cache"
	result, err := ioutil.ReadFile(cacheFilePath)
	if err != nil {
		return err
	}

	var s syscall.Stat_t
	syscall.Stat(cacheFilePath, &s)
	now := time.Now()

	sec, _ := s.Mtim.Unix()
	lastFileModifiedTime := time.Unix(sec, 0)

	duration := now.Sub(lastFileModifiedTime)
	log.Printf("[DEBUG] %s", duration)
	if duration.Minutes() > cacheExpireTime {
		return errors.New("cache expired")
	}
	*data = string(result)
	return nil
}

func putCache(data string) error {
	var cacheDirPath string
	err := getCacheDirPath(&cacheDirPath)
	if err != nil {
		return err
	}

	if f, err := os.Stat(cacheDirPath); os.IsNotExist(err) || f.IsDir() {
		err := os.MkdirAll(cacheDirPath, 0777)
		if err != nil {
			return err
		}
	}
	file, err := os.Create(cacheDirPath + "/.cache")
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(([]byte)(data))
	return nil
}

func getEc2List(c *cli.Context, ec2List *[]BasicEC2Info) error {
	var cache string
	err := getCache(&cache)
	if err == nil {
		json.Unmarshal([]byte(cache), ec2List)

		return nil
	}
	svc, err := getEc2Client()

	if err != nil {
		return err
	}

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
		return err
	}

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
			*ec2List = append(*ec2List, BasicEC2Info{
				Name:             nameTag,
				InstanceId:       *instance.InstanceId,
				PrivateIpAddress: *instance.PrivateIpAddress,
				PublicIpAddress:  publicIpAddress,
				InstanceType:     *instance.InstanceType,
				InstanceState:    *instance.State.Name,
			})
		}
	}

	result, _ := json.Marshal(&ec2List)
	err = putCache(string(result))
	if err != nil {
		return err
	}
	return nil
}

func getEC2(name string, result *BasicEC2Info) error {
	svc, err := getEc2Client()
	if err != nil {
		return err
	}
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
		return err
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
	*result = BasicEC2Info{
		Name:             nameTag,
		InstanceId:       *instance.InstanceId,
		PrivateIpAddress: *instance.PrivateIpAddress,
		PublicIpAddress:  publicIpAddress,
		InstanceType:     *instance.InstanceType,
		InstanceState:    *instance.State.Name,
	}
	return nil
}
