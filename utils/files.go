package utils

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

// 读取文件并返回一个字符串切片
func ReadFileAsLine(path string) (error, []string) {
	lineSlice := make([]string, 0)

	if !IsFile(path) {
		return os.ErrNotExist, nil
	}
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		return err, nil
	}

	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		lineSlice = append(lineSlice, line)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err, nil
			}
		}
	}

	return nil, lineSlice
}

func ReadDir(path string) []string {
	var files []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if IsFile(path) {
			files = append(files, info.Name())
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}
