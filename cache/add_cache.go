package cache

import (
	"encoding/hex"
	"gitlet/constcoe"
)

type BlobCache struct {
	ID           []byte //hash
	Content      []byte //文件内容
	FileName     string //对应的文件
	FilePath     string //对应的文件路径
	BlobFileName string //blob文件名

}


func (bc *BlobCache) GetBlobKey() string {
	return constcoe.CACHE_BLOB + hex.EncodeToString(bc.ID)
}


