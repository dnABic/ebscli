package main

import (
	"os"

	"github.com/dnabic/ebscli"
	"github.com/dnabic/ebscli/version"
)

func main() {
	exitCode := ebscli.Main(os.Args, version.GetVersion())
	os.Exit(exitCode)
}
