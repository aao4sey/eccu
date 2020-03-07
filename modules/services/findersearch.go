package services

import (
	"encoding/json"
	"fmt"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/urfave/cli/v2"
	common "github.com/yukkyun/eccu/modules/services/common"
)

func FinderSearch(c *cli.Context) error {
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

	idx, err := fuzzyfinder.FindMulti(
		*ec2List,
		func(i int) string {
			return (*ec2List)[i].Name
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintf("Name: %s\nInstanceId: %s\nPrivateIP: %s\nPubricIP: %s\nInstance State: %s\nInstance Type: %s",
				(*ec2List)[i].Name,
				(*ec2List)[i].InstanceID,
				(*ec2List)[i].PrivateIPAddress,
				(*ec2List)[i].PublicIPAddress,
				(*ec2List)[i].InstanceState,
				(*ec2List)[i].InstanceType,
			)
		}))
	if err != nil {
		return err
	}
	fmt.Printf("%s\t%s\t%s\t%s\t%s\t%s\n",
		(*ec2List)[idx[0]].Name,
		(*ec2List)[idx[0]].InstanceID,
		(*ec2List)[idx[0]].PrivateIPAddress,
		(*ec2List)[idx[0]].PublicIPAddress,
		(*ec2List)[idx[0]].InstanceState,
		(*ec2List)[idx[0]].InstanceType,
	)
	return nil
}
