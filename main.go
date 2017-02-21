package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/urfave/cli"
	"log"
	"os"
	"strings"
)

var version = "0.0.2"

func main() {
	var awsRegion string
	var ebsFilterTag string
	var ebsFilterId string
	var attachedOnly bool
	var ec2Id string

	commonFlags := []cli.Flag{}

	app := cli.NewApp()

	app.Name = "ebscli"
	app.Usage = "manage ebs volumes"
	app.Version = version

	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list ebs volumes",

			Flags: append(commonFlags,
				cli.StringFlag{
					Name:   "region, r",
					Value:  "us-east-1",
					Usage:  "AWS Region",
					EnvVar: "AWS_DEFAULT_REGION",

					Destination: &awsRegion,
				},
				cli.StringFlag{
					Name:  "tag, t",
					Value: "",
					Usage: "Volume filter by tags, eg. \"tag-key=tag-value,another-tag-key=another-tag-value\"",

					Destination: &ebsFilterTag,
				},
				cli.BoolFlag{
					Name:  "attached, a",
					Usage: "If set to true, lists only ebs volumes which are attached to the host from where ebscli is executed",

					Destination: &attachedOnly,
				},
				cli.StringFlag{
					Name:  "id, i",
					Value: "",
					Usage: "Volume filter by ids, eg. \"id1,id2,id3\"",

					Destination: &ebsFilterId,
				},
				cli.StringFlag{
					Name:  "ec2-id, e",
					Value: "",
					Usage: "Filter by volumes attached to specified ec2 ID",

					Destination: &ec2Id,
				},
			),
			Action: func(c *cli.Context) error {
				sess, err := session.NewSession()
				if err != nil {
					log.Fatalf("failed to create session %v\n", err)
				}

				ec2conn := ec2.New(sess, &aws.Config{Region: aws.String(awsRegion)})

				var paramsEbs *ec2.DescribeVolumesInput
				paramsEbs = nil
				var filterByTag []*ec2.Filter
				var volumeIds []*string

				var paramsInstance *ec2.DescribeInstancesInput
				paramsInstance = nil

				if len(ebsFilterTag) > 0 {
					filterByTag = getFilterTag(ebsFilterTag)
				}

				if len(ebsFilterId) > 0 {
					volumeIds = getVolumeIds(ebsFilterId)
				}

				if attachedOnly || len(ec2Id) > 0 {
					var instanceIds []*string

					if len(ec2Id) > 0 {
						instanceIds = append(instanceIds, &ec2Id)
					}

					paramsInstance = &ec2.DescribeInstancesInput{
						DryRun:      aws.Bool(false),
						InstanceIds: instanceIds,
					}

					resp, err := ec2conn.DescribeInstances(paramsInstance)
					if err != nil {
						log.Fatalf("There was an error listing instances in %s: %s", awsRegion, err.Error())
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
					return nil
				}

				fmt.Println(respEbs)
				return nil
			},
		},
		{
			Name:    "attach",
			Aliases: []string{"a"},
			Usage:   "attach ebs volume",

			Flags: append(commonFlags,
				cli.StringFlag{
					Name:   "region, r",
					Value:  "us-east-1",
					Usage:  "AWS Region",
					EnvVar: "AWS_DEFAULT_REGION",

					Destination: &awsRegion,
				},
				cli.StringFlag{
					Name:  "tag, t",
					Value: "",
					Usage: "Volume filter by tags, eg. \"tag-key=tag-value,another-tag-key=another-tag-value\"",

					Destination: &ebsFilterTag,
				},
				cli.StringFlag{
					Name:  "id, i",
					Value: "",
					Usage: "Volume filter by ids, eg. \"id1,id2,id3\"",

					Destination: &ebsFilterId,
				},
			),

			Action: func(c *cli.Context) error {
				sess, err := session.NewSession()
				if err != nil {
					log.Fatalf("failed to create session %v\n", err)
				}

				ec2conn := ec2.New(sess, &aws.Config{Region: aws.String(awsRegion)})

				var paramsEbs *ec2.DescribeVolumesInput
				paramsEbs = nil
				var filterByTag []*ec2.Filter
				var volumeIds []*string

				if len(ebsFilterTag) > 0 {
					filterByTag = getFilterTag(ebsFilterTag)
				}

				if len(ebsFilterId) > 0 {
					volumeIds = getVolumeIds(ebsFilterId)
				}

				paramsEbs = &ec2.DescribeVolumesInput{
					DryRun:    aws.Bool(false),
					Filters:   filterByTag,
					VolumeIds: volumeIds,
				}

				respEbs, err := ec2conn.DescribeVolumes(paramsEbs)

				if err != nil {
					fmt.Println(err.Error())
					return nil
				}

				fmt.Println(respEbs)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	exitCode := 0
	if err != nil {
		if _, ok := err.(*cli.ExitError); !ok {
			fmt.Println(err)
		}
		exitCode = 1
	}
	//return exitCode
	_ = exitCode
}

func getFilterTag(ebsFilterTag string) []*ec2.Filter {
	var filterByTag []*ec2.Filter
	tagList := strings.Split(ebsFilterTag, ",")
	for _, tag := range tagList {
		tagParams := strings.Split(tag, "=")
		if len(tagParams) != 2 {
			log.Fatalf("Invalid parameter value %s\n", tag)
		}
		tagName := "tag:" + tagParams[0]
		tagValue := tagParams[1]
		filterElement := ec2.Filter{
			Name: aws.String(tagName),
			Values: []*string{
				aws.String(tagValue), // Required
			},
		}
		filterByTag = append(filterByTag, &filterElement)
	}
	return filterByTag
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
