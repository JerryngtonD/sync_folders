package internal

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
)

func GetFileNameAndDirFromPath(path string) (string, string) {
	const PATH_WINDOWS_SEPARATOR = "\\"
	const PATH_UNIX_SEPARATOR = "/"

	var pathChunks []string
	if runtime.GOOS == "windows" {
		pathChunks = strings.Split(path, PATH_WINDOWS_SEPARATOR)
	} else {
		pathChunks = strings.Split(path, PATH_UNIX_SEPARATOR)
	}

	fileName := pathChunks[len(pathChunks)-1]
	dir := path[:len(path)-len(fileName)]

	return fileName, dir
}

func CopyFile(fromPath string, toPath string, logger *log.Logger) error {
	if !Exists(toPath) {
		from, err := os.OpenFile(fromPath, os.O_RDONLY, 0666)

		if err != nil {
			log.Fatal(err)
		}
		defer from.Close()

		to, err := os.OpenFile(toPath, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer to.Close()

		_, err = io.Copy(to, from)
		if err != nil {
			return fmt.Errorf("unexpected error with copying")
		}
		if logger != nil {
			logger.Println("Successfully sync copying file from: " + fromPath + " to path: " + toPath)
		}
	}
	return nil
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func DeleteFile(filePath string, logger *log.Logger) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	logger.Println("Successfully sync removing file: " + filePath)
	return nil
}

func DeleteDirectory(dirPath string) error {
	err := os.RemoveAll(dirPath)
	if err != nil {
		fmt.Println("got error with deleting presented directory: " + dirPath)
		return err
	}
	return nil
}
