package test

import (
	// "gitlet/repository"
	// "gitlet/stage"
	// "gitlet/utils"
	// "log"

	"gitlet/gredis"
	"gitlet/repository"
	"testing"
)

func Test_init(t *testing.T) {
	gredis.Setup()
	// repository.Init()
	// repository.Add("1.txt")
	// repository.Add("2.txt")
	// repository.Add("3.txt")
	// repository.Add("4.txt")
	// repository.Add("5.txt")
	repository.Commit("test commit.txt")
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
	// log.Println(repository.ReadCurrCommit())
	// repository.ResetByCommitId("a6aeb76e1e5f")
	// log.Println(repository.ReadCurrCommit())

	// ancestor := &commits.Commits{
	// 	PathToBlobID: map[string][]byte{
	// 		"file1.txt": []byte("blob1"),
	// 		"file2.txt": []byte("blob2"),
	// 	},
	// }
	// fmt.Printf("ancestor: %#v\n", ancestor)
	
	// newCommit := &commits.Commits{
	// 	PathToBlobID: map[string][]byte{
	// 		"file2.txt": []byte("blob2"),
	// 		"file3.txt": []byte("blob3"),
	// 	},
	// }
	
	// mergeCommit := &commits.Commits{
	// 	PathToBlobID: map[string][]byte{
	// 		"file1.txt": []byte("blob1"),
	// 		"file3.txt": []byte("blob3"),
	// 	},
	// }
	
	// result := repository.FindAllBlobID(ancestor, newCommit, mergeCommit)
	// for _, content := range result {
	// 	fmt.Println(string(content))
	// }
}
