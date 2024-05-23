package cache

import (
	"encoding/hex"
	"gitlet/constcoe"
)

type CommitCache struct {
	ID        			[]byte //当前commit的hash
	TimeStamp 			int64
	Message   			string //提交时的信息
	ParentID  			[][]byte
	PathToBlobID   		map[string][]byte //path, blobid
}

func (cc *CommitCache) GetCommitKey() string {
	return constcoe.CACHE_COMMIT + hex.EncodeToString(cc.ID)
}