package utils

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"gitlet/constcoe"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

func Handle(err error)  {
	if err != nil {
		log.Panic(err)
	}
}

//检查文件是否存在
func FileExists(fileAddr string) bool {
	if _, err := os.Stat(fileAddr); os.IsNotExist(err) {
		return	false
	}
	return true
}

// FormatUnixTimestamp 格式化 Unix 时间戳为字符串
func FormatUnixTimestamp(unixTimestamp int64) string {
	// 将 Unix 时间戳转换为 time.Time 类型
	t := time.Unix(unixTimestamp, 0)

	// 使用 time.Format 函数将时间格式化为字符串
	return t.Format("Mon Jan 02 15:04:05 2006 -0700")
}
//将int64转换为字节串
func ToHexInt(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func GenerateHash(obj interface{}) []byte {
	var hashChunks [][]byte

	val := reflect.ValueOf(obj)
	for i := 0; i < val.NumField(); i++ {
		fileld := val.Field(i)
		//如果字段是string，转成byte
		switch fileld.Kind() {
		case reflect.String:
			hashChunks = append(hashChunks, []byte(fileld.String()))
		case reflect.Slice:
			// 如果字段是切片类型
			if fileld.Type().Elem().Kind() == reflect.Uint8 {
				// 如果切片元素类型是字节，则直接拼接
				hashChunks = append(hashChunks, fileld.Bytes())
			} else {
				// 如果切片元素类型不是字节，可以按需处理
				// 这里给出一个示例，将切片中的每个元素转为字符串再拼接
				for j := 0; j < fileld.Len(); j++ {
					element := fileld.Index(j)
					hashChunks = append(hashChunks, []byte(fmt.Sprintf("%v", element.Interface())))
				}
			}
		case reflect.Int64:
			hashChunks = append(hashChunks, ToHexInt(fileld.Int()))
		case reflect.Map:
			tmp, _ := json.Marshal(fileld)
			hashChunks = append(hashChunks, tmp)
		}
	}

	hashBytes := bytes.Join(hashChunks, []byte{})
	
	hash := sha1.Sum(hashBytes)

	return hash[:]
	
}

func ContainsValue(byteMap map[string][]byte, value []byte) bool {
	for _, v := range byteMap {
		if bytes.Equal(v, value) {
			return true
		}
	}
	return false
}

//返回文件的绝对路径
func GetFileFromCWD(fileName string) (string, error) {
	//如果是绝对路径
	if filepath.IsAbs(fileName) {
		return fileName, nil
	}
	//如果是相对路径
	cwd, err := os.Getwd()
	Handle(err)
	return filepath.Join(cwd, fileName) ,nil
}

func GetCommitFile(folderPath string) []string {
	var txtFilePaths []string
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		Handle(err)

		// 判断是否为文件并且后缀为txt
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".txt") {
			txtFilePaths = append(txtFilePaths, path)
		}

		return nil
	})
	Handle(err)

	return txtFilePaths
}

func GetBranchList(folderPath string) []string {
	var branchFiles []string
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		Handle(err)
		branchFiles = append(branchFiles, info.Name())
		return nil
	})
	Handle(err)
	return branchFiles
}

func GetBlobFile(folderPath string) []string {
	var branchFiles []string
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		Handle(err)

		// 判断是否为文件并且后缀为txt
		if !info.IsDir() && !strings.HasSuffix(info.Name(), ".txt") {
			branchFiles = append(branchFiles, path)
		}

		return nil
	})
	Handle(err)

	return branchFiles
}

// 给出commitID的前n个字符用来匹配commit文件
func MatchFileByPrefix(folderPath, prefix string) string {
	var matchingFile string
	err := filepath.Walk(folderPath, func(path string, info fs.FileInfo, err error) error {
		Handle(err)
		if !info.IsDir() {
			fileName := info.Name()
			if strings.HasPrefix(fileName, prefix) {
				ext := filepath.Ext(fileName)
				matchingFile = fileName[: len(fileName) - len(ext)]
			}
		}
		return nil
	})
	Handle(err)
	return matchingFile
}

func ReadAllBranch() []string {
	var branches []string
	err := filepath.Walk(constcoe.HeadsDir, func(path string, info fs.FileInfo, err error) error {
		Handle(err)
		if !info.IsDir() {
			fileName := info.Name()
			branches = append(branches, fileName)
		}
		return nil
	})
	Handle(err)
	return branches
}

func GetAllWorkingFile() []string {
	files, err := filepath.Glob("*")
	Handle(err)

	var workingFiles []string
	for _, file := range files {
		fileInfo, err := os.Stat(file)
		Handle(err)

		if !fileInfo.IsDir() {
			workingFiles = append(workingFiles, file)
		}
	}

	return workingFiles
}