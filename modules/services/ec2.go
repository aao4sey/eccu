package services

import (
	"fmt"
)

type BasicEC2Info struct {
	Name             string
	InstanceID       string
	PrivateIPAddress string
	PublicIPAddress  string
	InstanceType     string
	InstanceState    string
}

func (b *BasicEC2Info) ShowTsv() {
	fmt.Printf("%s\t%s\t%s\t%s\t%s\n",
		b.Name,
		b.InstanceID,
		b.PrivateIPAddress,
		b.PublicIPAddress,
		b.InstanceType,
	)
}

func (b *BasicEC2Info) ShowCsv() {
	fmt.Printf("%s,%s,%s,%s,%s\n",
		b.Name,
		b.InstanceID,
		b.PrivateIPAddress,
		b.PublicIPAddress,
		b.InstanceType,
	)
}
