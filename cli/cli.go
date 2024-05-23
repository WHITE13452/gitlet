package cli

import (
	"fmt"
	"gitlet/repository"
	"os"
	"runtime"
)

type Commandline struct{}

func (cli *Commandline) validateArgs(num int) {
	if len(os.Args) != num {
		fmt.Println("args' num no good")
		runtime.Goexit()
	}
}

func (cli *Commandline) Run() {

	// initGitCommand := flag.NewFlagSet("init", flag.ExitOnError)
	// addCommand := flag.NewFlagSet("add", flag.ExitOnError)
	// commitCommand := flag.NewFlagSet("commit", flag.ExitOnError)
	// rmCommand := flag.NewFlagSet("rm", flag.ExitOnError)
	// logCommand := flag.NewFlagSet("log", flag.ExitOnError)
	// globallogCommand := flag.NewFlagSet("global-log", flag.ExitOnError)
	// findCommand := flag.NewFlagSet("find", flag.ExitOnError)
	// statusCommand := flag.NewFlagSet("status", flag.ExitOnError)
	// checkoutCommand := flag.NewFlagSet("checkout", flag.ExitOnError)
	// branchCommand := flag.NewFlagSet("branch", flag.ExitOnError)
	// rmBranchCommand := flag.NewFlagSet("rm-branch", flag.ExitOnError)
	// restCommand := flag.NewFlagSet("reset", flag.ExitOnError)
	// mergeCommand := flag.NewFlagSet("merge", flag.ExitOnError)

	// addFileRefFileName := addCommand.String("filename", "", "The name of the file you want to add")
	// commitFileWithMsg := commitCommand.String("message", "", "Enter the commit message")
	// rmCommitWithFileName := rmCommand.String("fileName", "", "Enter the file name")
	// findCommitByMessage := findCommand.String("message", "", "Enter the commit message")
	// checkoutByFileName := checkoutCommand.String("fileName", "", "Which file do you want to check")
	// checkoutByCommitIDAndFileName := checkoutCommand.String("commitID", "", "You need to give commitID as second arg")
	// checkoutByBranchName := checkoutCommand.String("branchName", "", "Enter Branch Name here")
	// createBranchByName := branchCommand.String("branchName", "", "Enter Branch Name here")
	// rmBranchByName := rmBranchCommand.String("branchName", "", "Enter Branch Name here")
	// resetCommitByCommitID := restCommand.String("commitID", "", "Enter Commit ID here")
	// mergeBranchName := mergeCommand.String("branchName", "", "Enter a branch Name here")

	switch os.Args[1] {
	case "init":
		cli.validateArgs(2)
		repository.Init()
	case "add":
		cli.validateArgs(3)
		repository.CheckIfInit()
		repository.Add(os.Args[2])

	case "commit":
		cli.validateArgs(3)
		repository.CheckIfInit()
		repository.Commit(os.Args[2])

	case "rm":
		cli.validateArgs(3)
		repository.CheckIfInit()
		repository.Remove(os.Args[2])

	case "log":
		cli.validateArgs(2)
		repository.CheckIfInit()
		repository.Log()

	case "global-log":
		cli.validateArgs(2)
		repository.CheckIfInit()
		repository.GlobalLog()

	case "find":
		cli.validateArgs(3)
		repository.CheckIfInit()
		repository.FindCommitByMessage(os.Args[2])
		
	case "status":
		cli.validateArgs(2)
		repository.CheckIfInit()
		repository.Status()

	case "checkout":
		if len(os.Args) < 3 {
			fmt.Println("Please provide at least one argument.")
			runtime.Goexit()
		}
		args := os.Args[2:]
		if len(args) == 0 {
			fmt.Println("Please provide fileName and commitID or branchName")
			runtime.Goexit()
		}
		repository.CheckIfInit()
		if len(args) == 2  && args[0] == "--"{ // "-- filename"
			repository.CheckOutByCommitIdAndFileName("", args[1])
		} else if len(args) == 3 && args[1] == "--"{ // "commitID -- filename"
			repository.CheckOutByCommitIdAndFileName(args[0], args[2])
		} else if len(args) == 1 {
			repository.CheckOutByBranchName(args[0])
		} else
		{
			fmt.Println("Please provide fileName and commitID or branchName")
			runtime.Goexit()
		}

	case "branch":
		cli.validateArgs(3)
		repository.CheckIfInit()
		repository.CreateBranch(os.Args[2])

	case "rm-branch":
		cli.validateArgs(3)
		repository.CheckIfInit()
		repository.RemoveBranch(os.Args[2])

	case "reset":
		cli.validateArgs(3)
		repository.CheckIfInit()
		repository.ResetByCommitId(os.Args[2])

	case "merge":
		cli.validateArgs(3)
		repository.CheckIfInit()
		repository.Merge(os.Args[2])

	default:
		fmt.Println("no commands like this")
		runtime.Goexit()
	}

	// if initGitCommand.Parsed() {
	// 	repository.Init()
	// }
	// if addCommand.Parsed() {
	// 	repository.CheckIfInit()
	// 	repository.Add(*addFileRefFileName)
	// }
	// if commitCommand.Parsed() {
	// 	repository.CheckIfInit()
	// 	repository.Commit(*commitFileWithMsg)
	// }
	// if rmCommand.Parsed() {
	// 	repository.CheckIfInit()
	// 	repository.Remove(*rmCommitWithFileName)
	// }
	// if logCommand.Parsed() {
	// 	repository.CheckIfInit()
	// 	repository.Log()
	// }
	// if globallogCommand.Parsed() {
	// 	repository.CheckIfInit()
	// 	repository.GlobalLog()
	// }
	// if findCommand.Parsed() {
	// 	repository.CheckIfInit()
	// 	repository.FindCommitByMessage(*findCommitByMessage)
	// }
	// if statusCommand.Parsed() {
	// 	repository.CheckIfInit()
	// 	repository.Status()
	// }
	// if checkoutCommand.Parsed() {
	// 	repository.CheckIfInit()
	// 	if *checkoutByCommitIDAndFileName == "" {
	// 		if *checkoutByFileName == "" {
	// 			checkoutCommand.Usage()
	// 			runtime.Goexit()	
	// 		} else {
	// 			repository.CheckOutByCommitIdAndFileName("", *checkoutByFileName)
	// 		}
	// 	} else {
	// 		if *checkoutByFileName == "" { 
	// 			checkoutCommand.Usage()
	// 			runtime.Goexit()
	// 		} else {
	// 			repository.CheckOutByCommitIdAndFileName(*checkoutByCommitIDAndFileName, *checkoutByFileName)
	// 		}
	// 	}
	// 	if *checkoutByCommitIDAndFileName == "" && *checkoutByFileName == "" && *checkoutByBranchName != "" {
	// 		repository.CheckOutByBranchName(*checkoutByBranchName)
	// 	}
	// }
	// if branchCommand.Parsed() {
	// 	repository.CheckIfInit()
	// 	repository.CreateBranch(*createBranchByName)
	// }
	// if rmBranchCommand.Parsed() {
	// 	repository.CheckIfInit()
	// 	repository.RemoveBranch(*rmBranchByName)
	// }
	// if restCommand.Parsed() {
	// 	repository.CheckIfInit()
	// 	repository.ResetByCommitId(*resetCommitByCommitID)	
	// }
	// if mergeCommand.Parsed() {
	// 	repository.CheckIfInit()
	// 	repository.Merge(*mergeBranchName)
		
	// }
}
