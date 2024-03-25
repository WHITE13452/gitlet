package test

import (
	// "gitlet/repository"
	// "gitlet/stage"
	// "gitlet/utils"
	// "log"

	"gitlet/repository"
	"log"
	"testing"
)

func Test_init(t *testing.T) {
	 //repository.Init()
	//repository.Add("4.txt")
	// repository.Commit("test commit 4.txt")
	// repository.Remove("3.txt")
	// stage := &stage.Stage{}
	// err := utils.ReadStruct("/Users/white/develop/GO_workspace/gitlet/test/.gitlet/remove_stage", stage)
	// utils.Handle(err)
	// log.Println(stage)
	// log.Println(string(hex.EncodeToString(tmp)))
	// repository.Remove("1.txt")
	//repository.Commit("test remove 1.txt")
	// repository.GlobalLog()
	// repository.Log()
	//repository.FindCommitByMessage("test not found")
	//repository.Status()
	// repository.CheckOutByCommitIdAndFileName("a6aeb76e1e5f", "3.txt")
	// repository.CreateBranch("master")
	//repository.CheckOutByBranchName("master")
	//repository.RemoveBranch("test")
	log.Println(repository.ReadCurrCommit())
	repository.ResetByCommitId("a6aeb76e1e5f")
	log.Println(repository.ReadCurrCommit())
}
