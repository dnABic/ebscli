package ebscli

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"log"
	"strings"
)

type listArgs struct {
	name         string
	awsRegion    string
	ebsFilterTag string
	ebsFilterId  string
	ec2Id        string
	attachedOnly bool
}

func listEbs(args listArgs) {
	sess, err := session.NewSession()
	if err != nil {
		log.Fatalf("failed to create session %v\n", err)
	}

	ec2conn := ec2.New(sess, &aws.Config{Region: aws.String(args.awsRegion)})

	var paramsEbs *ec2.DescribeVolumesInput
	paramsEbs = nil
	var filterByTag []*ec2.Filter
	var volumeIds []*string

	var paramsInstance *ec2.DescribeInstancesInput
	paramsInstance = nil

	if len(args.ebsFilterTag) > 0 {
		filterByTag = getFilterTag(args.ebsFilterTag)
	}

	if len(args.ebsFilterId) > 0 {
		volumeIds = getVolumeIds(args.ebsFilterId)
	}

	if args.attachedOnly || len(args.ec2Id) > 0 {
		var instanceIds []*string

		if len(args.ec2Id) > 0 {
			instanceIds = append(instanceIds, &args.ec2Id)
		}

		paramsInstance = &ec2.DescribeInstancesInput{
			DryRun:      aws.Bool(false),
			InstanceIds: instanceIds,
		}

		resp, err := ec2conn.DescribeInstances(paramsInstance)
		if err != nil {
			log.Fatalf("There was an error listing instances in %s: %s", args.awsRegion, err.Error())
		}

		for idx, res := range resp.Reservations {
			fmt.Println("Number of instances: ", len(res.Instances))
			for _, inst := range resp.Reservations[idx].Instances {
				fmt.Println("    - Instance ID: ", *inst.InstanceId)
				fmt.Println("Block devices count", len(inst.BlockDeviceMappings))
				for cnt, ebs := range inst.BlockDeviceMappings {
					fmt.Println(cnt)
					fmt.Println("Block device name", *ebs.DeviceName)
					volume := ebs.Ebs
					fmt.Println(*volume.VolumeId)
					volumeIds = append(volumeIds, volume.VolumeId)
				}
			}
		}

	}

	paramsEbs = &ec2.DescribeVolumesInput{
		DryRun:    aws.Bool(false),
		Filters:   filterByTag,
		VolumeIds: volumeIds,
	}

	respEbs, err := ec2conn.DescribeVolumes(paramsEbs)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(respEbs)
	}
}

func getVolumeIds(ebsFilterId string) []*string {
	var volumeIds []*string
	volumeIdList := strings.Split(ebsFilterId, ",")
	for _, id := range volumeIdList {
		volumeId := id
		volumeIds = append(volumeIds, &volumeId)
	}

	return volumeIds
}
