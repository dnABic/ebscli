package ebscli

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/urfave/cli"
	"log"
	"os"
	"strings"
)

// Main entry point for ebscli application
func Main(args []string, version string) int {
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
				},
				cli.StringFlag{
					Name:  "tag, t",
					Value: "",
					Usage: "Volume filter by tags, eg. \"tag-key=tag-value,another-tag-key=another-tag-value\"",
				},
				cli.BoolFlag{
					Name:  "attached, a",
					Usage: "If set to true, lists only ebs volumes which are attached to the host from where ebscli is executed",
				},
				cli.StringFlag{
					Name:  "id, i",
					Value: "",
					Usage: "Volume filter by ids, eg. \"id1,id2,id3\"",
				},
				cli.StringFlag{
					Name:  "ec2-id, e",
					Value: "",
					Usage: "Filter by volumes attached to specified ec2 ID",
				},
			),
			Action: func(c *cli.Context) error {
				args := listArgs{
					name:         c.Args().First(),
					awsRegion:    c.String("region"),
					ebsFilterTag: c.String("tag"),
					ebsFilterId:  c.String("id"),
					ec2Id:        c.String("ec2-id"),
					attachedOnly: c.Bool("attached"),
				}
				listEbs(args)
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
				},
				cli.StringFlag{
					Name:  "tag, t",
					Value: "",
					Usage: "Volume filter by tags, eg. \"tag-key=tag-value,another-tag-key=another-tag-value\"",
				},
				cli.StringFlag{
					Name:  "id, i",
					Value: "",
					Usage: "Volume filter by ids, eg. \"id1,id2,id3\"",
				},
				cli.StringFlag{
					Name:  "ec2-id, e",
					Value: "",
					Usage: "ec2 instance ID to which volume should be attached",
				},
			),
			Action: func(c *cli.Context) error {
				args := attachArgs{
					name:         c.Args().First(),
					awsRegion:    c.String("region"),
					ebsFilterTag: c.String("tag"),
					ebsFilterId:  c.String("id"),
					ec2Id:        c.String("ec2-id"),
				}
				attachEbs(args)
				return nil
			},
		},
		{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "prints the ebscli version",

			Flags: append(commonFlags),

			Action: func(c *cli.Context) error {
				log.Printf("ebscli version: %s", version)
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
	return exitCode
	//_ = exitCode
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
