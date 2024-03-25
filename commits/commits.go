package commits

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gitlet/blob"
	"gitlet/constcoe"
	"gitlet/utils"
	"time"
)

type Commits struct {
	ID        			[]byte //当前commit的hash
	TimeStamp 			int64
	Message   			string //提交时的信息
	ParentID  			[][]byte
	PathToBlobID   		map[string][]byte //path, blobid
}


func CreateCommits(msg string, blobs map[string][]byte, parentID [][]byte) *Commits {
	commits := &Commits{
		ID:        []byte{},
		TimeStamp: time.Now().Unix(),
		Message:   msg,
		ParentID:  parentID,
		PathToBlobID:     blobs,
	}
	commits.ID = utils.GenerateHash(*commits)

	return commits
}

func (c *Commits) Save() {
	fmt.Println("len id: ", len(string(hex.EncodeToString(c.ID))))
	filePath := constcoe.ObjectDir + string(hex.EncodeToString(c.ID)) + ".txt"
	// var content bytes.Buffer
	// encoder := gob.NewEncoder(&content)
	// err := encoder.Encode(c)
	// utils.Handle(err)
	// err = os.WriteFile(filePath, content.Bytes(), 0777)
	// utils.Handle(err)
	utils.SaveStruct(filePath, c)
}

func (c *Commits) generateHash() []byte {
	blobsByte, err := json.Marshal(c.PathToBlobID)
	utils.Handle(err)
	information := bytes.Join([][]byte{utils.ToHexInt(c.TimeStamp), []byte(c.Message), append(c.ParentID[0], c.ParentID[1]...), blobsByte}, []byte{})
	hash := sha1.Sum(information)
	return hash[:]
}

// func (c *Commits) ContainsBlobValue(content []byte) bool {
// 	for _, v := range c.Blobs {
// 		if  bytes.Equal(v, content){
// 			return true
// 		}
// 	}
// 	return false
// }

func (c *Commits) GetFileName() []string {
	fileList := []string{}
	for _, blobID := range c.PathToBlobID {
		blob := blob.GetBlobById(blobID)
		fileList = append(fileList, blob.FileName)
	}
	return fileList
}

func  (c *Commits) GetBlobByFileName(fileName string) *blob.Blob {
	filePath, err := utils.GetFileFromCWD(fileName)
	utils.Handle(err)
	blobID := c.PathToBlobID[filePath]
	return blob.GetBlobById(blobID)
}

func (c *Commits) GetBlobList()  {
	
}