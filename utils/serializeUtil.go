package utils

import (
	"bytes"
	"encoding/gob"
	"errors"
	"os"
)

//写文件
func WriteContents(filePath string, contents ...interface{}) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, obj := range contents {
		switch v := obj.(type) {
		case []byte:
			_, err = f.Write(v)
		case string:
			_, err = f.WriteString(v)
		default:
			return errors.New("unsupported type")
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func ReadContents(filePath string) ([]byte, error) {
	if !FileExists(filePath) {
		Handle(errors.New("File does not exists !"))
	}
	return os.ReadFile(filePath)
}

func ReadContentsAsString(filePath string) (string, error) {
	content, err := ReadContents(filePath)
	return string(content), err
}

func ReadStruct(filePath string, expectedType interface{}) error {
	fileContent, err := ReadContents(filePath)
	Handle(err)
	decoder := gob.NewDecoder(bytes.NewBuffer(fileContent))
	err = decoder.Decode(expectedType)
	Handle(err)
	return nil
}

func SaveStruct(filePath string, savedStruct interface{}) {
	var content bytes.Buffer
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(savedStruct)
	Handle(err)
	err = os.WriteFile(filePath, content.Bytes(), 0777)
	Handle(err)
}

func DeleteFile(filePath string) error {
	fileInfo, err := os.Stat(filePath)
	Handle(err)
	if fileInfo.IsDir() {
		Handle(errors.New("can not delete a dir"))
	}

	err = os.Remove(filePath)
	Handle(err)
	return nil
}

