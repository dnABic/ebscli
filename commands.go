package ebscli

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	_ "github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type listArgs struct {
	name         string
	awsRegion    string
	ebsFilterTag string
	ebsFilterId  string
	ec2Id        string
	attachedOnly bool
}

type attachArgs struct {
	name         string
	awsRegion    string
	ebsFilterTag string
	ebsFilterId  string
	ec2Id        string
}

//func getInstanceId(region string) {
//	sess, err := session.NewSession()
//	if err != nil {
//		log.Fatalf("failed to create session %v\n", err)
//	}
//
//	endpoint_url := os.Getenv("EBSCLI_ENDPOINT_URL")
//	svc := ec2metadata.New(sess, &aws.Config{Region: aws.String(region), Endpoint: aws.String(endpoint_url)})
//	fmt.Println(svc.GetMetadata("InstanceID"))
//}

func listEbs(args listArgs) {
	//getInstanceId(args.awsRegion)
	sess, err := session.NewSession()
	if err != nil {
		log.Fatalf("failed to create session %v\n", err)
	}

	endpoint_url := os.Getenv("EBSCLI_ENDPOINT_URL")
	ec2conn := ec2.New(sess, &aws.Config{Region: aws.String(args.awsRegion), Endpoint: aws.String(endpoint_url)})

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

func attachEbs(args attachArgs) {
	sess, err := session.NewSession()
	if err != nil {
		log.Fatalf("failed to create session %v\n", err)
	}

	ec2conn := ec2.New(sess, &aws.Config{Region: aws.String(args.awsRegion)})

	var paramsEbs *ec2.DescribeVolumesInput
	paramsEbs = nil
	var filterByTag []*ec2.Filter
	var volumeIds []*string

	if len(args.ebsFilterTag) > 0 {
		filterByTag = getFilterTag(args.ebsFilterTag)
	}

	if len(args.ebsFilterId) > 0 {
		volumeIds = getVolumeIds(args.ebsFilterId)
	}

	paramsEbs = &ec2.DescribeVolumesInput{
		DryRun:    aws.Bool(false),
		Filters:   filterByTag,
		VolumeIds: volumeIds,
	}

	respEbs, err := ec2conn.DescribeVolumes(paramsEbs)

	if err != nil {
		//fmt.Println(err.Error())
		//return nil
		log.Fatalf(err.Error())
	}

	for _, volume := range respEbs.Volumes {
		//fmt.Println(volume)
		fmt.Println(*volume.VolumeId)
		if len(args.ec2Id) > 0 {
			fmt.Println(args.ec2Id)
		}
	}
	//	params := &ec2.AttachVolumeInput{
	//		Device:     aws.String("String"), // Required
	//		InstanceId: aws.String("String"), // Required
	//		VolumeId:   aws.String("String"), // Required
	//	}
	//	resp, err := ec2conn.AttachVolume(params)
	//	if err != nil {
	//		log.Fatalf(err.Error())
	//	}

	fmt.Println(respEbs)
	//return nil
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
