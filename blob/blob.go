// 保存文件内容（不包括文件名等）。
// gitlet会保存一个文件的多个版本，所以一个文件会对应多个blobs
// 每个都在不同的commit中被追踪
package blob

import (
	"encoding/hex"
	"gitlet/constcoe"
	"gitlet/utils"
	"log"
	"path"
)

type Blob struct {
	ID           []byte //hash
	Content      []byte //文件内容
	FileName     string //对应的文件
	FilePath     string //对应的文件路径
	BlobFileName string //blob文件名
}

func CreateBlob(filePath string) *Blob {
	content, err := utils.ReadContents(filePath)
	utils.Handle(err)
	blob := &Blob{
		ID:           []byte{},
		Content:      content,
		FileName:     path.Base(filePath),
		FilePath:     filePath,
		BlobFileName: "",
	}
	blob.ID = utils.GenerateHash(*blob)
	blob.BlobFileName = constcoe.ObjectDir + hex.EncodeToString(blob.ID)
	return blob
}

func (b *Blob) Save() {
	utils.SaveStruct(b.BlobFileName, b)
}

func GetBlobById(id []byte) *Blob {
	blob := &Blob{}
	filePath := constcoe.ObjectDir + string(hex.EncodeToString(id))
	error := utils.ReadStruct(filePath, blob)
	log.Println("Blob: ", blob)
	utils.Handle(error)
	return blob
}

