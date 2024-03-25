package repository

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"gitlet/blob"
	"gitlet/commits"
	"gitlet/constcoe"
	"gitlet/stage"
	"gitlet/utils"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

/*
*   .gitlet
*      |--objects
*      |     |--commit and blob
*      |--refs
*      |    |--heads
*      |         |--master
*      |--HEAD
*      |--stage
 */

var CurrCommit *commits.Commits
var CurrBranch string

func mkdirAll(dir string, mode fs.FileMode) {
	err := os.MkdirAll(dir, mode)
	if err != nil {
		fmt.Printf("mkdir: %s: Failed to create.", dir)
		runtime.Goexit()
	}
}

func Init() {
	err := os.Mkdir(constcoe.GitletDir, 0777)
	if err != nil {
		fmt.Println("A Gitlet version-control system already exists in the current directory.")
		runtime.Goexit()
	}
	mkdirAll(constcoe.ObjectDir, 0777)
	mkdirAll(constcoe.RefsDir, 0777)
	mkdirAll(constcoe.HeadsDir, 0777)

	InitCommit()
	InitHead()
	InitHeads()
}

func InitHead() {
	err := utils.WriteContents(constcoe.HeadsFile, "master")
	utils.Handle(err)
}

func InitHeads() {
	err := utils.WriteContents(constcoe.HeadsDir+"master", CurrCommit.ID)
	utils.Handle(err)
}

func InitCommit() {
	initCommit := commits.CreateCommits("initial commit", make(map[string][]byte), nil)
	CurrCommit = initCommit
	initCommit.Save()
}

func CheckIfInit() {
	if !utils.FileExists(constcoe.GitletDir) && !utils.FileExists(constcoe.HeadsDir) {
		log.Fatal("you didn't initializ the gitlet in this directory")
		runtime.Goexit()
	}
}

func Add(fileName string) {
	filePath, err := utils.GetFileFromCWD(fileName)
	utils.Handle(err)
	if !utils.FileExists(filePath) {
		log.Fatal("File does not exist.")
		runtime.Goexit()
	}
	blob := blob.CreateBlob(filePath)
	storeBlob(blob)
}

func storeBlob(blob *blob.Blob) {
	currCommit := ReadCurrCommit()
	addStage := readStage(constcoe.AddStage)
	removeStage := readStage(constcoe.RemoveStage)
	if !utils.ContainsValue(currCommit.PathToBlobID, blob.ID) || !removeStage.IsNewBlob(blob) {
		if addStage.IsNewBlob(blob) {
			//如果并不在准备删除的缓存区中
			if removeStage.IsNewBlob(blob) {
				blob.Save()
				if addStage.IsFilePathExists(blob.FilePath) {
					addStage.DeleteByPath(blob.FilePath)
				}
				//将该blob添加到缓存区
				addStage.Add(blob)
				addStage.SaveStage(constcoe.AddStage)
			} else {
				removeStage.DeleteByPath(blob.FilePath)
				removeStage.SaveStage(constcoe.RemoveStage)
			}
		}
	}
}

func readStage(pattern string) *stage.Stage {

	if pattern == constcoe.AddStage {
		if !utils.FileExists(constcoe.AddStageFile) {
			return &stage.Stage{}
		}
		stage := &stage.Stage{}
		utils.ReadStruct(constcoe.AddStageFile, stage)
		return stage
	} else if pattern == constcoe.RemoveStage {
		if !utils.FileExists(constcoe.RemoveStageFile) {
			return &stage.Stage{}
		}
		stage := &stage.Stage{}
		utils.ReadStruct(constcoe.RemoveStageFile, stage)
		return stage
	} else {
		utils.Handle(errors.New("Stage place wrong!"))
	}
	return nil
}

// func readAddStage() *stage.Stage {
// 	if !utils.FileExists(constcoe.AddStageFile) {
// 		return &stage.Stage{}
// 	}
// 	stage := &stage.Stage{}
// 	utils.ReadStruct(constcoe.AddStageFile, stage)
// 	return stage
// }

// func readRemoveStage() *stage.Stage {
// 	if !utils.FileExists(constcoe.RemoveStageFile) {
// 		return &stage.Stage{}
// 	}
// 	stage := &stage.Stage{}
// 	utils.ReadStruct(constcoe.RemoveStageFile, stage)
// 	return stage
// }

func ReadCurrCommit() *commits.Commits {
	currBranch := readCurrBranch()
	headsFile := constcoe.HeadsDir + currBranch
	currCommitID, err := utils.ReadContents(headsFile)
	utils.Handle(err)

	currCommitFile := constcoe.ObjectDir + string(hex.EncodeToString(currCommitID)) + ".txt"
	commit := &commits.Commits{}
	utils.ReadStruct(currCommitFile, commit)
	log.Println("currCommit: ", commit)
	return commit
}

func readCurrBranch() string {
	branch, err := utils.ReadContentsAsString(constcoe.HeadsFile)
	utils.Handle(err)
	return branch
}

func Commit(message string) {
	CurrCommit = ReadCurrCommit()
	log.Println("currCommit: ", CurrCommit)
	//找到缓存区的blob，新建一个commit，parent是currCommit
	addBlobMap := findBlobMap("add")
	removeBlobMap := findBlobMap("remove")
	if addBlobMap == nil && removeBlobMap == nil {
		utils.Handle(errors.New("No changes added to the commit."))
	}

	blobMap := CurrCommit.PathToBlobID
	if blobMap == nil {
		blobMap = make(map[string][]byte)
	}
	for path, addBlob := range addBlobMap {
		blobMap[path] = addBlob
	}
	for path, _ := range removeBlobMap {
		delete(blobMap, path)
	}

	parentID := [][]byte{}
	parentID = append(parentID, CurrCommit.ID)
	for _, v := range parentID {
		log.Println("parentID: ", string(v))
	}

	newCommit := commits.CreateCommits(message, blobMap, parentID)
	log.Println("newCommit: ", newCommit)
	newCommit.Save()

	addStage := readStage(constcoe.AddStage)
	addStage.ClearStage()
	addStage.SaveStage(constcoe.AddStage)

	removeStage := readStage(constcoe.RemoveStage)
	removeStage.ClearStage()
	removeStage.SaveStage(constcoe.RemoveStage)

	saveHeads(newCommit)

}

func findBlobMap(pattern string) map[string][]byte {
	stage := readStage(pattern)
	blobList := stage.GetBlobList()
	blobMap := make(map[string][]byte)

	for _, blob := range blobList {
		blobMap[blob.FilePath] = blob.ID
	}
	return blobMap
}

func saveHeads(commit *commits.Commits) {
	currBranch := readCurrBranch()
	headsFile := constcoe.HeadsDir + currBranch
	utils.WriteContents(headsFile, commit.ID)
}

func Remove(fileName string) {
	filePath, _ := utils.GetFileFromCWD(fileName)
	currCommit := ReadCurrCommit()
	addStage := readStage(constcoe.AddStage)

	_, isStaged := addStage.PathToBlobID[filePath]
	_, isTracked := currCommit.PathToBlobID[filePath]

	if isStaged {
		delete(addStage.PathToBlobID, filePath)
		addStage.SaveStage(constcoe.AddStage)
	}
	if isTracked {
		removeStage := readStage(constcoe.RemoveStage)
		removeBlob := blob.GetBlobById(currCommit.PathToBlobID[filePath])
		removeStage.Add(removeBlob)
		removeStage.SaveStage(constcoe.RemoveStage)
		utils.DeleteFile(filePath)
	} else {
		utils.Handle(errors.New("No reason to remove the file."))
	}
}

func Log() {
	currCommit := ReadCurrCommit()
	log.Println(currCommit)

	for currCommit.ParentID != nil {
		var isMerged bool
		if len(currCommit.ParentID) == 2 {
			isMerged = true
		}
		printCommit(isMerged, currCommit)
		parentCommitID := currCommit.ParentID
		currCommit = readCommitByID(parentCommitID[0])
	}
	printCommit(false, currCommit)
}

func GlobalLog() {
	commitFiles := utils.GetCommitFile(constcoe.ObjectDir)
	commit := &commits.Commits{}
	for _, commitFile := range commitFiles {
		utils.ReadStruct(commitFile, commit)
		if len(commit.ParentID) == 2 {
			printCommit(true, commit)
		} else {
			printCommit(false, commit)
		}
	}
}

func readCommitByID(commitID []byte) *commits.Commits {
	commit := &commits.Commits{}
	filePath := constcoe.ObjectDir + string(hex.EncodeToString(commitID)) + ".txt"
	_ = utils.ReadStruct(filePath, commit)
	return commit
}

func printCommit(isMerged bool, currCommit *commits.Commits) {

	if isMerged {
		firstParent := currCommit.ParentID[1][:7]
		secParent := currCommit.ParentID[2][:7]
		fmt.Println("===")
		fmt.Println("commit: ", string(hex.EncodeToString(currCommit.ID)))
		fmt.Println("Merge: ", string(firstParent), " ", string(secParent))
		fmt.Println("Date: ", utils.FormatUnixTimestamp(currCommit.TimeStamp))
	} else {
		fmt.Println("===")
		fmt.Println("commit: ", string(hex.EncodeToString(currCommit.ID)))
		fmt.Println("Date: ", utils.FormatUnixTimestamp(currCommit.TimeStamp))
		fmt.Println(currCommit.Message)
	}
}

func FindCommitByMessage(message string) {
	commitList := []*commits.Commits{}
	commitFiles := utils.GetCommitFile(constcoe.ObjectDir)
	for _, commitFile := range commitFiles {
		commit := &commits.Commits{}
		utils.ReadStruct(commitFile, commit)
		if commit.Message == message {
			commitList = append(commitList, commit)
		}
	}
	if len(commitList) == 0 {
		utils.Handle(errors.New("Found no commit with that message."))
	}
	for _, commit := range commitList {
		fmt.Println(string(hex.EncodeToString(commit.ID)))
	}
}

func Status() {
	//branches
	currBranch := readCurrBranch()
	branchList := utils.GetBranchList(constcoe.HeadsDir)
	printBranchesStatus(currBranch, branchList)

	//stage
	addStage := readStage(constcoe.AddStage)
	fmt.Println("=== Staged Files ===")
	for _, stagedBlob := range addStage.GetBlobList() {
		fmt.Println(stagedBlob.FileName)
	}
	fmt.Println("")

	//removed file
	removed := readStage(constcoe.RemoveStage)
	fmt.Println("=== Removed Files ===")
	for _, removedStage := range removed.GetBlobList() {
		fmt.Println(removedStage.FileName)
	}
	fmt.Println("")

	//=== Modifications Not Staged For Commit ===
	fmt.Println("=== Modifications Not Staged For Commit ===")
	printModifiedNotStaged()
	fmt.Println("")

	//=== Untracked Files ===
	fmt.Println("=== Untracked Files ===")
}

func printBranchesStatus(currBranch string, branchList []string) {
	fmt.Println("=== Branches ===")
	fmt.Println("*", currBranch)
	for _, branch := range branchList {
		if branch != currBranch {
			fmt.Println(branch)
		}
	}
	fmt.Println("")
}

func printModifiedNotStaged() {
	//filesWeNeed := []string{}
	currCommit := ReadCurrCommit()
	addStage := readStage(constcoe.AddStage)
	removeStage := readStage(constcoe.RemoveStage)

	blobFiles := utils.GetBlobFile(constcoe.ObjectDir)
	blob := &blob.Blob{}
	for _, blobFilePath := range blobFiles {
		utils.ReadStruct(blobFilePath, blob)

		filePath := blob.FilePath

		_, isTracked := currCommit.PathToBlobID[filePath]

		_, isStagedForAdd := addStage.PathToBlobID[filePath]
		_, isStagedForRemoval := removeStage.PathToBlobID[filePath]
		file, err := os.Open(filePath)
		defer file.Close()

		isDeleted := false
		if err != nil {
			isDeleted = true
		}

		isChanged := false
		fileContent, _ := os.ReadFile(filePath)
		if !bytes.Equal(fileContent, blob.Content) {
			isChanged = true
		}

		//is tracked in the current commit,changed in the working directory
		if isTracked && isChanged {
			//but not staged
			if !isStagedForAdd {
				fmt.Println(filepath.Base(filePath), "(modified)")
			}
		}

		if isStagedForAdd {
			//Staged for addition
			if !isChanged {
				//but with different contents than in the working directory
				//filesWeNeed = append(filesWeNeed, filePath)
				fmt.Println(filepath.Base(filePath), "(modified)")
			}
			if isDeleted {
				//filesWeNeed = append(filesWeNeed, filePath)
				fmt.Println(filepath.Base(filePath), "(deleted)")
			}
		}

		if !isStagedForRemoval && isTracked && isDeleted {
			//filesWeNeed = append(filesWeNeed, filePath)
			fmt.Println(filepath.Base(filePath), "(deleted)")
		}
	}

	// for _, path := range filesWeNeed {
	// 	fmt.Println(filepath.Base(path))
	// }
}

func CheckOutByCommitIdAndFileName(commitID, fileName string) {
	// var currCommit *commits.Commits
	if commitID != "" {
		decodeCommitID := getByteCommitID(commitID)
		CurrCommit = readCommitByID(decodeCommitID)
		if CurrCommit == nil {
			utils.Handle(errors.New("No commit with that id exists."))
		}
	} else {
		CurrCommit = ReadCurrCommit()
	}

	fileNameList := CurrCommit.GetFileName()
	isExist := false
	for _, file := range fileNameList {
		if file == fileName {
			blob := CurrCommit.GetBlobByFileName(fileName)
			WriteBlobToCWDFile(blob)
			isExist = true
			break
		}
	}
	if !isExist {
		utils.Handle(errors.New("File does not exist in that commit."))
	}
}

// CheckOutByBranchName checks out the specified branch by name.
// It compares the files between the current branch and the checkout branch,
// and performs the necessary file operations to update the working directory.
// If the specified branch does not exist or is the current branch, an error is returned.
func CheckOutByBranchName(branchName string) {
	currBranch := readCurrBranch()
	if currBranch == branchName {
		utils.Handle(errors.New("No need to checkout the current branch."))
	}

	allBranch := utils.ReadAllBranch()

	if !branchExists(branchName, allBranch) {
		utils.Handle(errors.New("No such branch exists."))
	}
	//全局变量
	CurrCommit = ReadCurrCommit()
	checkedBranchCommit := ReadCommitByBranch(branchName)

	changeCommitTo(checkedBranchCommit)

	//将branch改为checked branch
	utils.WriteContents(constcoe.HeadsFile, branchName)
}

func branchExists(branchName string, allBranches []string) bool {
	for _, branch := range allBranches {
		if branchName == branch {
			return true
		}
	}
	return false
}

func WriteBlobToCWDFile(blob *blob.Blob) {
	filePath, err := utils.GetFileFromCWD(blob.FileName)
	utils.Handle(err)
	err = utils.WriteContents(filePath, blob.Content)
	utils.Handle(err)
}

func ReadCommitByBranch(branchName string) *commits.Commits {
	branchFileName := constcoe.HeadsDir + branchName
	commitID, err := utils.ReadContents(branchFileName)
	utils.Handle(err)
	// decodeID, err := hex.DecodeString(newCommitID)
	// log.Println("newCommitID : ",newCommitID)
	//utils.Handle(err)
	return readCommitByID(commitID)
}

func OverWriteFileFromCommit(filePath map[string]bool, commit *commits.Commits) {
	if filePath == nil {
		return
	}
	for file := range filePath {
		blob := commit.GetBlobByFileName(file)
		WriteBlobToCWDFile(blob)
	}
}

func ClearAllStage() {
	addStage := readStage(constcoe.AddStage)
	removeStage := readStage(constcoe.RemoveStage)

	addStage.ClearStage()
	addStage.SaveStage(constcoe.AddStage)

	removeStage.ClearStage()
	removeStage.SaveStage(constcoe.RemoveStage)
}

func CreateBranch(branchName string) {
	allBranch := utils.ReadAllBranch()
	if branchExists(branchName, allBranch) {
		utils.Handle(errors.New("A branch with that name already exists."))
	}
	//将当前commit的ID写入新建的branch文件中
	currCommit := ReadCurrCommit()
	utils.WriteContents(constcoe.HeadsDir+branchName, currCommit.ID)
}

func RemoveBranch(branchName string) {
	currBranch := readCurrBranch()
	if currBranch == branchName {
		utils.Handle(errors.New("Cannot remove the current branch."))
	}
	if !branchExists(branchName, utils.ReadAllBranch()) {
		utils.Handle(errors.New("A branch with that name does not exist."))
	}
	utils.DeleteFile(constcoe.HeadsDir + branchName)
}

func ResetByCommitId(commitID string) {
	decodeCommitID := getByteCommitID(commitID)
	newCommit := readCommitByID(decodeCommitID)
	if newCommit == nil {
		utils.Handle(errors.New("No commit with that id exists."))
	}
	//currCommit 作为全局变量，在changeCommitTo中会使用
	CurrCommit = ReadCurrCommit()
	changeCommitTo(newCommit)

	currBranch := readCurrBranch()
	changeBranchHeadTo(decodeCommitID, currBranch)

}

// changeCommitTo updates the current branch commit to match the given commit.
// It performs the following steps:
// 1. Finds the files that exist only in the current branch commit but not in the checkout branch commit.
// 2. Deletes the files that exist only in the current branch commit.
// 3. Finds the files that exist in both the checkout branch and current branch commits.
// 4. Overwrites the files in the current branch commit with the versions from the checkout branch commit.
// 5. Finds the files that exist only in the checkout branch commit.
// 6. Writes the files that exist only in the checkout branch commit to the directory.
// 7. Checks if there are any untracked files in the way and throws an error if so.
// 8. Overwrites the files in the current branch commit with the versions from the checkout branch commit.
// 9. Clears the staging area.
func changeCommitTo(commit *commits.Commits) {

	//找到只在current branch commit中存在但是check out branch commit中不存在的文件
	filesOnlyTrackedInCurrCommit := make(map[string]bool)
	checkedBranchCommitFiles := commit.GetFileName()

	filesTrackedInCurrCommit := CurrCommit.GetFileName()
	filesTrackedInCurrCommitMap := make(map[string]bool)

	for _, fileName := range filesTrackedInCurrCommit {
		filesOnlyTrackedInCurrCommit[fileName] = true
		filesTrackedInCurrCommitMap[fileName] = true
	}

	for _, fileName := range checkedBranchCommitFiles {
		delete(filesOnlyTrackedInCurrCommit, fileName)
	}

	//删除filesOnlyTrackedInCurrCommit
	for fileName, _ := range filesOnlyTrackedInCurrCommit {
		filePath, err := utils.GetFileFromCWD(fileName)
		utils.Handle(err)
		utils.DeleteFile(filePath)
	}

	//找到在checkout branch和currbranch中都存在的文件
	//用checkout branch中的文件替代currbranch中的文件版本（相同文件名但blobID不同）
	bothCommitTrackedFiles := make(map[string]bool)
	for _, fileName := range checkedBranchCommitFiles {
		if _, ok := filesTrackedInCurrCommitMap[fileName]; ok {
			bothCommitTrackedFiles[fileName] = true
		}
	}
	OverWriteFileFromCommit(bothCommitTrackedFiles, commit)

	//只在checkout branch中存在的文件，就直接写进目录
	filesOnlyTrackedIncheckedBranchCommit := make(map[string]bool)

	for _, file := range checkedBranchCommitFiles {
		filesOnlyTrackedIncheckedBranchCommit[file] = true
	}
	for _, filePath := range filesTrackedInCurrCommit {
		delete(filesOnlyTrackedIncheckedBranchCommit, filePath)
	}

	for fileName := range filesOnlyTrackedIncheckedBranchCommit {
		filePath, err := utils.GetFileFromCWD(fileName)
		utils.Handle(err)
		//如果这个时候在check Branch中的文件已经存在目录里了
		//说明目录中执行checkout前增加了新的同文件名的文件还没有commit，那么抛出错误
		if utils.FileExists(filePath) {
			utils.Handle(errors.New("There is an untracked file in the way; delete it, or add and commit it first."))
		}
	}

	OverWriteFileFromCommit(filesOnlyTrackedIncheckedBranchCommit, commit)
	ClearAllStage()
}

func changeBranchHeadTo(commitID []byte, branchName string) {
	utils.WriteContents(constcoe.HeadsDir+branchName, commitID)
}

func getByteCommitID(commitID string) []byte {
	if len(commitID) != 40 {
		commitID = utils.MatchFileByPrefix(constcoe.ObjectDir, commitID)
	}
	decodeCommitID, err := hex.DecodeString(commitID)
	utils.Handle(err)
	return decodeCommitID
}

//记得不同的版本的文件名可能是一样的，但是文件内容不同时，add之后blobID就会不同，所以检查不同的blobID即可
func Merge(branchName string)  {
	CurrBranch = readCurrBranch()
	//检查缓存区是否非空
	addStage := readStage(constcoe.AddStage)
	removeStage := readStage(constcoe.RemoveStage)
	if len(addStage.PathToBlobID) != 0 || len(removeStage.PathToBlobID) != 0 {
		utils.Handle(errors.New("You have uncommitted changes."))
	}

	if !branchExists(branchName, utils.ReadAllBranch()) {
		utils.Handle(errors.New("A branch with that name does not exist."))
	}
	//不能和自己合并
	checkIfMergeWithSelf(branchName)

	CurrCommit = ReadCurrCommit()
	mergeCommit := ReadCommitByBranch(branchName)
	commonAncestor, err := findCommonAncestor(mergeCommit)
	utils.Handle(err)
	checkIfAncestorIsSameCommAsGivenBranch(commonAncestor, mergeCommit)
	checkIfAncestorIsCurrBranch(commonAncestor, branchName)
	
	workingFiles := utils.GetAllWorkingFile()
	for _, file := range workingFiles {
		if CurrCommit.PathToBlobID[file] == nil &&  mergeCommit.PathToBlobID[file] != nil {
			utils.Handle(errors.New("There is an untracked file in the way; delete it, or add and commit it first."))
		}
	}

	for filePath, currCommitblobID := range CurrCommit.PathToBlobID {
		mergeCommitBlobID := mergeCommit.PathToBlobID[filePath]
		splitCommitBlobID := commonAncestor.PathToBlobID[filePath]
		//如果当前commit和merge commit的blobID相同，那么保持原状
		if mergeCommitBlobID != nil && bytes.Equal(splitCommitBlobID, mergeCommitBlobID) {
			continue
		}
		if splitCommitBlobID != nil {
			//如果当前branch中没有修改
			if bytes.Equal(splitCommitBlobID, currCommitblobID) {
				//但是merge branch中有已经删除了的文件，那么删除
				if mergeCommitBlobID == nil {
					removeStage := readStage(constcoe.RemoveStage)
					removeStage.Add(blob.GetBlobById(currCommitblobID))
					removeStage.SaveStage(constcoe.RemoveStage)
					if utils.FileExists(filePath) {
						utils.DeleteFile(filePath)
					}
				} else { //如果merge branch中有修改，那么直接写入
					blob := blob.GetBlobById(mergeCommitBlobID)
					WriteBlobToCWDFile(blob)
					addStage := readStage(constcoe.AddStage)
					addStage.Add(blob)
					addStage.SaveStage(constcoe.AddStage)
				}
			} else { //在当前branch中有修改
				conflictContent := conflictFileContents(currCommitblobID, mergeCommitBlobID)
				utils.WriteContents(filePath, conflictContent)
				Add(filePath)
				fmt.Println("Encountered a merge conflict.")
			}
		} else if mergeCommitBlobID != nil && !bytes.Equal(mergeCommitBlobID, currCommitblobID) {
			conflictContent := conflictFileContents(currCommitblobID, mergeCommitBlobID)
			utils.WriteContents(filePath, conflictContent)
			Add(filePath)
			fmt.Println("Encountered a merge conflict.")
		}
	}
	for filePath, mergeCommitBlobID := range mergeCommit.PathToBlobID {
		currCommitBlobID := CurrCommit.PathToBlobID[filePath]
		splitCommitBlobID := commonAncestor.PathToBlobID[filePath]
		//分离点和当前branch中都不存在，在merge branch中存在
		if splitCommitBlobID == nil && currCommitBlobID == nil {
			blob := blob.GetBlobById(mergeCommitBlobID)
			WriteBlobToCWDFile(blob)
			addStage := readStage(constcoe.AddStage)
			addStage.Add(blob)
			addStage.SaveStage(constcoe.AddStage)
		}
	}

}

func conflictFileContents(currCommitBlobID, mergeCommitBlobID []byte) string {
	currContent := getContentFromBlobID(currCommitBlobID)
	mergeContent := getContentFromBlobID(mergeCommitBlobID)

	return fmt.Sprintf("<<<<<<< HEAD\n%s=======\n%s>>>>>>>\n", currContent, mergeContent)
	
}

func  getContentFromBlobID(blobID []byte) string {
	blob :=blob.GetBlobById(blobID)
	return string(blob.Content)
}

func checkIfMergeWithSelf(branchName string)  {
	if CurrBranch == branchName {
		utils.Handle(errors.New("Cannot merge a branch with itself."))
	}
}

func findCommonAncestor(mergeCommit *commits.Commits) (*commits.Commits, error) {
	ancestors := make(map[string]bool)
	for commit1 := CurrCommit; commit1 != nil; commit1 = getNextCommit(commit1) {
		ancestors[string(commit1.ID)] = true
	}
	for commit2 := mergeCommit; commit2 != nil; commit2 = getNextCommit(commit2) {
		if _, ok := ancestors[string(commit2.ID)]; ok {
			return commit2, nil
		}
	}
	return nil, errors.New("No common ancestor found.")
}

func getNextCommit(commit *commits.Commits) *commits.Commits {
	if len(commit.ParentID) > 0 {
		return readCommitByID(commit.ParentID[0])
	}
	return nil
	
}

func checkIfAncestorIsSameCommAsGivenBranch(ancestor, mergeCommit *commits.Commits)  {
	if bytes.Equal(ancestor.ID, mergeCommit.ID) {
		utils.Handle(errors.New("Given branch is an ancestor of the current branch."))
	}
}
func checkIfAncestorIsCurrBranch(ancestor *commits.Commits, mergeBranch string)  {
	if bytes.Equal(ancestor.ID, CurrCommit.ID) {
		CheckOutByBranchName(mergeBranch)
		utils.Handle(errors.New("Current branch fast-forwarded."))
	}
}

// func mergeFilesToNewCommit(ancestor, newCommit, mergeCommit *commits.Commits) *commits.Commits {
// 	ancestorBlobs := ancestor.PathToBlobID
// 	currBranchBlobs := newCommit.PathToBlobID
// 	mergeBranchBlobs := mergeCommit.PathToBlobID
// 	newCommitBlobs := make(map[string][]byte)

// 	for fileName, ancestorBlobID := range ancestorBlobs {
// 		currBlobID, isCurrBlobExist := currBranchBlobs[fileName]
// 		mergeBlobID, isMergeBlobExist := mergeBranchBlobs[fileName]

// 		if isCurrBlobExist && isMergeBlobExist {
// 			if bytes.Equal(ancestorBlobID, mergeBlobID) {
// 				if !bytes.Equal(currBlobID, mergeBlobID) {
// 					newCommitBlobs[fileName] = mergeBlobID
// 				}
// 			} else if !bytes.Equal(ancestorBlobID, currBlobID) {
// 				if bytes.Equal(mergeBlobID, currBlobID) {
// 					newCommitBlobs[fileName] = mergeBlobID
// 				} else {
// 					utils.Handle(errors.New("Encountered a merge conflict."))
// 				}
// 			}
// 		} else if isCurrBlobExist {
// 			if !bytes.Equal(ancestorBlobID, currBlobID) {
// 				utils.Handle(errors.New("Encountered a merge conflict."))
// 			}
// 		} else if isMergeBlobExist {
// 			newCommitBlobs[fileName] = mergeBlobID
// 		}
// 	}


// }



func WriteFiles(filePath map[string]bool, commit *commits.Commits)  {
	if filePath == nil {
		return
	}
	for file := range filePath {
		filep, err := utils.GetFileFromCWD(file)
		utils.Handle(err)
		if utils.FileExists(filep) {
			utils.Handle(errors.New("There is an untracked file in the way; delete it, or add and commit it first."))
		}
	}
	OverWriteFileFromCommit(filePath, commit)
}

func findAllBlobID(ancestor, newCommit, mergeCommit *commits.Commits) [][]byte {
	allFiles := [][]byte{}
	allFiles = appendBlobIDs(allFiles, ancestor.PathToBlobID)
	allFiles = appendBlobIDs(allFiles, newCommit.PathToBlobID)
	allFiles = appendBlobIDs(allFiles, mergeCommit.PathToBlobID)
	return removeDuplicates(allFiles)
}

func appendBlobIDs(allFiles [][]byte, blobIDs map[string][]byte) [][]byte {
	for _, blobID := range blobIDs {
		allFiles = append(allFiles, blobID)
	}
	return allFiles
}

//去重
func removeDuplicates(allFiles [][]byte) [][]byte {
	uniqueFiles := make(map[string]bool)
	result := [][]byte{}
	for _, file := range allFiles {
		fileStr := string(file)
		if !uniqueFiles[fileStr] {
			uniqueFiles[fileStr] = true
			result = append(result, file)
		}
	}
	return result
}


func checkIfConflict(ancestor, newCommit, mergeCommit *commits.Commits) {
	allFiles := findAllBlobID(ancestor, newCommit, mergeCommit)

	for _, file := range allFiles {
		filePath := string(file)
		newCommitContent := getContentFromCommit(newCommit, filePath)
		mergeCommitContent := getContentFromCommit(mergeCommit, filePath)

		if newCommitContent != mergeCommitContent {
			// Handle conflict resolution here
			// Replace the contents of the conflicted file with the desired resolution
			// For example:
			// newCommitContent = "Desired resolution"
			// Overwrite the file with newCommitContent
		}
	}
}

func getContentFromCommit(commit *commits.Commits, filePath string) string {
	// Retrieve the content of the file from the commit
	// Implement the logic to fetch the content from the commit
	// For example:
	// content := commit.GetFileContent(filePath)
	// return content
	return ""
}

