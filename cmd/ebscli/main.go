package main

import (
	"os"

	"github.com/dnabic/ebscli"
)

func main() {
	exitCode := ebscli.Main(os.Args)
	os.Exit(exitCode)
}
