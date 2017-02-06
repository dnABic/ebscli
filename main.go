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
			),
			Action: func(c *cli.Context) error {
				sess, err := session.NewSession()
				if err != nil {
					log.Fatalf("failed to create session %v\n", err)
				}

				//fmt.Println("WE have CONFIG", c.String("region"))
				ec2conn := ec2.New(sess, &aws.Config{Region: aws.String(awsRegion)})
				var paramsEbs *ec2.DescribeVolumesInput
				paramsEbs = nil
				var tagFilter []*ec2.Filter
				var volumeIds []*string

				if len(ebsFilterTag) > 0 {
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
						tagFilter = append(tagFilter, &filterElement)
					}

				}

				if len(ebsFilterId) > 0 {
					volumeIds = append(volumeIds, &ebsFilterId)
				}

				paramsEbs = &ec2.DescribeVolumesInput{
					DryRun:    aws.Bool(false),
					Filters:   tagFilter,
					VolumeIds: volumeIds,
				}

				respEbs, err := ec2conn.DescribeVolumes(paramsEbs)

				if err != nil {
					fmt.Println(err.Error())
					return nil
				}

				fmt.Println(respEbs)
				return nil

				var params *ec2.DescribeInstancesInput

				if len(ebsFilterTag) > 0 {
					tagList := strings.Split(ebsFilterTag, ",")
					tagParams := strings.Split(tagList[0], "=")
					//tagName := strings.Join("tag:", tagParams[0])
					tagName := "tag:" + tagParams[0]
					tagValue := tagParams[1]

					params = &ec2.DescribeInstancesInput{
						DryRun: aws.Bool(false),
						Filters: []*ec2.Filter{
							{
								Name: aws.String(tagName),
								Values: []*string{
									aws.String(tagValue),
								},
							},
							//More values...
						},
					}
				} else {
					params = nil
				}

				resp, err := ec2conn.DescribeInstances(params)
				if err != nil {
					fmt.Println("there was an error listing instances in", awsRegion, err.Error())
					log.Fatal(err.Error())
				}

				fmt.Println("Number of reservation sets: ", len(resp.Reservations))
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

							if len(ebsFilterTag) > 0 {
								tagList := strings.Split(ebsFilterTag, ",")
								tagParams := strings.Split(tagList[0], "=")
								tagName := "tag:" + tagParams[0]
								tagValue := tagParams[1]

								paramsEBS := &ec2.DescribeVolumesInput{
									DryRun: aws.Bool(false),
									Filters: []*ec2.Filter{
										{ // Required
											Name: aws.String(tagName),
											Values: []*string{
												aws.String(tagValue), // Required
												// More values...
											},
										},
										// More values...
									},
									//MaxResults: aws.Int64(1),
									//	NextToken:  aws.String("String"),
									//		VolumeIds: []*string{
									//			aws.String("String"), // Required
									//			// More values...
									//		},
								}
								respEBS, err := ec2conn.DescribeVolumes(paramsEBS)

								if err != nil {
									// Print the error, cast err to awserr.Error to get the Code and
									// Message from an error.
									fmt.Println(err.Error())
									return nil
								}

								// Pretty-print the response data.
								fmt.Println(respEBS)
							}
						}
					}
				}

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
