package helper

import (
	"errors"
	"os"
)

func PathIsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func MkdirAll(path string) error {
	if !PathIsExist(path) {
		return os.MkdirAll(path, os.ModePerm)
	} else {
		fileInfo, _ := os.Stat(path)
		if !fileInfo.IsDir() {
			return errors.New(path + " 不是一个目录")
		}
	}
	return nil
}
