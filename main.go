package main

import (
	"gitlet/cli"
	"os"
)

func main()  {
	defer os.Exit(0)
	cmd := cli.Commandline{}
	cmd.Run()
}