package util

import (
	"fmt"
	"os"
	"path/filepath"
)

func IsFileExist(filePath string) bool {
	fi, err := os.Stat(filePath)
	if err != nil {
		return os.IsExist(err)
	}
	return !fi.IsDir()
}

func IsFileNotExist(filePath string) bool {
	fi, err := os.Stat(filePath)
	if err != nil {
		return os.IsNotExist(err)
	}
	return fi.IsDir()
}

func IsDirNotExist(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return os.IsNotExist(err)
	}
	return !fi.IsDir()
}

func IsDirExist(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return fi.IsDir()
}

func IsExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return os.IsExist(err)
	}
	return true
}

func IsNotExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return os.IsNotExist(err)
	}
	return false
}

// CreateFile force create file, if file exists , it will be rewritten
func CreateFile(filePath string, value []byte) error {
	dir := filepath.Dir(filePath)
	if err := CreateDirIfNotExist(dir); err != nil {
		return fmt.Errorf("create dir %s failed: %s", dir, err)
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0644))
	if err != nil {
		return fmt.Errorf("open file %s failed: %s", filePath, err)
	}
	defer file.Close()

	if _, err := file.Write(value); err != nil {
		return err
	}

	return nil
}

// CreateDirIfNotExist create directory if given dirPath not exist.
func CreateDirIfNotExist(dirPath string) error {
	if IsDirExist(dirPath) {
		return nil
	}
	return os.MkdirAll(dirPath, os.FileMode(0755))
}
