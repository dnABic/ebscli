package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/urfave/cli"
	"log"
	"os"
)

var version = "0.0.2"

func main() {
	app := cli.NewApp()
	app.Name = "ebscli"
	app.Usage = "manage ebs volumes"
	app.Version = version
	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list ebs volumes",
			Action: func(c *cli.Context) error {
				sess, err := session.NewSession()
				if err != nil {
					log.Fatalf("failed to create session %v\n", err)
				}

				awsRegion := "us-east-1"
				ec2conn := ec2.New(sess, &aws.Config{Region: aws.String(awsRegion)})

				resp, err := ec2conn.DescribeInstances(nil)
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
						}
					}
				}

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
