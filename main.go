package main

import (
	"gitlet/cli"
	"gitlet/gredis"
	"os"
)

func main()  {
	defer os.Exit(0)
	gredis.Setup()
	cmd := cli.Commandline{}
	cmd.Run()
}