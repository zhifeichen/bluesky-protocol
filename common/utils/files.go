package common

import (
	"io/ioutil"
	"os"
	"strings"
)

/**
判断文件是否存在
*/
func IsExist(fp string) bool {
	_, err := os.Stat(fp)
	return err == nil || os.IsExist(err)
}

func ToBytes(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}

func ReadToString(filePath string) (string, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func ReadToTrimString(filePath string) (string, error) {
	str, err := ReadToString(filePath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(str), nil
}
