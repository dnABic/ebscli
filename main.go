package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
)

var version = "0.0.1"

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
				fmt.Println("All done")
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
