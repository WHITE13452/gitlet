package cli

import (
	"flag"
	"fmt"
	"gitlet/repository"
	"gitlet/utils"
	"os"
	"runtime"
)

type commandline struct{}

func (cli *commandline) validateArgs(num int) {
	if len(os.Args) != num {
		fmt.Println("args' num no good")
		runtime.Goexit()
	}
}

func (cli *commandline) Run() {

	initGitCommand := flag.NewFlagSet("init", flag.ExitOnError)
	addCommand := flag.NewFlagSet("add", flag.ExitOnError)
	commitCommand := flag.NewFlagSet("commit", flag.ExitOnError)
	rmCommand := flag.NewFlagSet("rm", flag.ExitOnError)
	logCommand := flag.NewFlagSet("log", flag.ExitOnError)
	globallogCommand := flag.NewFlagSet("global-log", flag.ExitOnError)
	findCommand := flag.NewFlagSet("find", flag.ExitOnError)
	statusCommand := flag.NewFlagSet("status", flag.ExitOnError)
	checkoutCommand := flag.NewFlagSet("checkout", flag.ExitOnError)
	branchCommand := flag.NewFlagSet("branch", flag.ExitOnError)
	rmBranchCommand := flag.NewFlagSet("rm-branch", flag.ExitOnError)
	restCommand := flag.NewFlagSet("reset", flag.ExitOnError)
	mergeCommand := flag.NewFlagSet("merge", flag.ExitOnError)

	addFileRefFileName := addCommand.String("filename", "", "The name of the file you want to add")
	commitFileWithMsg := commitCommand.String("message", "", "Enter the commit message")
	rmCommitWithFileName := rmCommand.String("fileName", "", "Enter the file name")
	findCommitByMessage := findCommand.String("message", "", "Enter the commit message")
	checkoutByFileName := checkoutCommand.String("fileName", "", "Which file do you want to check")
	checkoutByCommitIDAndFileName := checkoutCommand.String("commitID", "", "You need to give commitID as second arg")
	checkoutByBranchName := checkoutCommand.String("branchName", "", "Enter Branch Name here")
	createBranchByName := branchCommand.String("branchName", "", "Enter Branch Name here")
	rmBranchByName := rmBranchCommand.String("branchName", "", "Enter Branch Name here")
	resetCommitByCommitID := restCommand.String("commitID", "", "Enter Commit ID here")
	mergeBranchName := mergeCommand.String("branchName", "", "Enter a branch Name here")

	switch os.Args[1] {
	case "init":
		cli.validateArgs(1)
		err := initGitCommand.Parse(os.Args[2:])
		utils.Handle(err)
	case "add":
		cli.validateArgs(2)
		err := addCommand.Parse(os.Args[2:])
		utils.Handle(err)
	case "commit":
		cli.validateArgs(2)
		err := commitCommand.Parse(os.Args[2:])
		utils.Handle(err)
	case "rm":
		cli.validateArgs(2)
		err := rmCommand.Parse(os.Args[2:])
		utils.Handle(err)
	case "log":
		cli.validateArgs(1)
		err := logCommand.Parse(os.Args[2:])
		utils.Handle(err)
	case "global-log":
		cli.validateArgs(1)
		err := globallogCommand.Parse(os.Args[2:])
		utils.Handle(err)
	case "find":
		cli.validateArgs(2)
		err := findCommand.Parse(os.Args[2:])
		utils.Handle(err)
	case "status":
		cli.validateArgs(1)
		err := statusCommand.Parse(os.Args[2:])
		utils.Handle(err)
	case "checkout":
		
		err := checkoutCommand.Parse(os.Args[2:])
		utils.Handle(err)
	case "branch":
		cli.validateArgs(2)
		err := branchCommand.Parse(os.Args[2:])
		utils.Handle(err)
	case "rm-branch":
		cli.validateArgs(2)
		err := rmBranchCommand.Parse(os.Args[2:])
		utils.Handle(err)
	case "reset":
		cli.validateArgs(2)
		err := restCommand.Parse(os.Args[2:])
		utils.Handle(err)
	case "merge":
		cli.validateArgs(2)
		err := mergeCommand.Parse(os.Args[2:])
		utils.Handle(err)
	default:
		fmt.Println("no commands like this")
		runtime.Goexit()
	}

	if initGitCommand.Parsed() {
		repository.Init()
	}
	if addCommand.Parsed() {
		repository.CheckIfInit()
		repository.Add(*addFileRefFileName)
	}
	if commitCommand.Parsed() {
		repository.CheckIfInit()
		repository.Commit(*commitFileWithMsg)
	}
	if rmCommand.Parsed() {
		repository.CheckIfInit()
		repository.Remove(*rmCommitWithFileName)
	}
	if logCommand.Parsed() {
		repository.CheckIfInit()
		repository.Log()
	}
	if globallogCommand.Parsed() {
		repository.CheckIfInit()
		repository.GlobalLog()
	}
	if findCommand.Parsed() {
		repository.CheckIfInit()
		repository.FindCommitByMessage(*findCommitByMessage)
	}
	if statusCommand.Parsed() {
		repository.CheckIfInit()
		repository.Status()
	}
	if checkoutCommand.Parsed() {
		repository.CheckIfInit()
		if *checkoutByCommitIDAndFileName == "" {
			if *checkoutByFileName == "" {
				checkoutCommand.Usage()
				runtime.Goexit()	
			} else {
				repository.CheckOutByCommitIdAndFileName("", *checkoutByFileName)
			}
		} else {
			if *checkoutByFileName == "" { 
				checkoutCommand.Usage()
				runtime.Goexit()
			} else {
				repository.CheckOutByCommitIdAndFileName(*checkoutByCommitIDAndFileName, *checkoutByFileName)
			}
		}
		if *checkoutByCommitIDAndFileName == "" && *checkoutByFileName == "" && *checkoutByBranchName != "" {
			repository.CheckOutByBranchName(*checkoutByBranchName)
		}
	}
	if branchCommand.Parsed() {
		repository.CheckIfInit()
		repository.CreateBranch(*createBranchByName)
	}
	if rmBranchCommand.Parsed() {
		repository.CheckIfInit()
		repository.RemoveBranch(*rmBranchByName)
	}
	if restCommand.Parsed() {
		repository.CheckIfInit()
		repository.ResetByCommitId(*resetCommitByCommitID)	
	}
	if mergeCommand.Parsed() {
		repository.CheckIfInit()
		repository.Merge(*mergeBranchName)
		
	}
}
