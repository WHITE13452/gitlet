package stage

import (
	"gitlet/blob"
	"gitlet/constcoe"
	"gitlet/utils"
)

//暂存区域
type Stage struct {
	PathToBlobID	map[string][]byte	//[path,hash]
}

//在缓存区中是否存在该blob
func (s *Stage) IsNewBlob(blob *blob.Blob) bool {
	if !utils.ContainsValue(s.PathToBlobID, blob.ID) {
		return true
	}
	return false
}

func (s *Stage) IsFilePathExists(path string) bool {
	if  _, ok := s.PathToBlobID[path]; ok {
		return true
	}
	return false
}

func (s *Stage) DeleteByPath(path string)  {
	delete(s.PathToBlobID, path)
}

func (s *Stage) Add(b *blob.Blob)  {
	if s.PathToBlobID == nil {
		s.PathToBlobID = make(map[string][]byte)
	}
	s.PathToBlobID[b.FilePath] = b.ID
}

// func (s *Stage) SaveAddStage()  {
// 	utils.SaveStruct(constcoe.AddStageFile, s)
// }

// func (s *Stage) SaveRemoveStage()  {
// 	utils.SaveStruct(constcoe.RemoveStageFile, s)
// }

func (s *Stage) SaveStage(pattern string)  {
	if pattern == constcoe.AddStage {
		utils.SaveStruct(constcoe.AddStageFile, s)
	} else if pattern == constcoe.RemoveStage {
		utils.SaveStruct(constcoe.RemoveStageFile, s)
	}
}

func (s *Stage) ClearStage()  {
	s.PathToBlobID = make(map[string][]byte)
}

func (s *Stage) GetBlobList() []*blob.Blob {
	blobList := []*blob.Blob{}
	for _, blobID := range s.PathToBlobID {
		blobList = append(blobList, blob.GetBlobById(blobID))
	}
	return blobList
}
