package services

import (
	"fmt"
	"log"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/urfave/cli/v2"
)

func FinderSearch(c *cli.Context) error {
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
			return fmt.Sprintf("Name: %s\nInstanceId: %s\nPrivateIP: %s\nPubricIP: %s\nInstance State: %s\nInstance Type: %s",
				ec2List[i].Name,
				ec2List[i].InstanceId,
				ec2List[i].PrivateIpAddress,
				ec2List[i].PublicIpAddress,
				ec2List[i].InstanceState,
				ec2List[i].InstanceType,
			)
		}))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\t%s\t%s\t%s\t%s\t%s\n",
		ec2List[idx[0]].Name,
		ec2List[idx[0]].InstanceId,
		ec2List[idx[0]].PrivateIpAddress,
		ec2List[idx[0]].PublicIpAddress,
		ec2List[idx[0]].InstanceState,
		ec2List[idx[0]].InstanceType,
	)
	return nil
}
