package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	svc := ec2.New(sess, &aws.Config{Region: aws.String("us-east-1")})

	resp, err := svc.DescribeInstances(nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("> Number of reservation sets: ", len(resp.Reservations))
	for idx, res := range resp.Reservations {
		fmt.Println("  > Number of instances: ", len(res.Instances))
		for _, inst := range resp.Reservations[idx].Instances {
			fmt.Println("    - Instance ID: ", *inst.InstanceId)
			fmt.Println("Block devices count", len(inst.BlockDeviceMappings))
			for cnt, ebs := range inst.BlockDeviceMappings {
				fmt.Println(cnt)
				fmt.Println("Block device name", *ebs.DeviceName)
				volume := ebs.Ebs
				fmt.Println(*volume.VolumeId)
			}
		}
	}
}
