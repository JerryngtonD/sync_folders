package internal

import (
	"github.com/stretchr/testify/require"
	"os"
	"runtime"
	"testing"
)

const dirName = "someDir"

const entirePathWindows = "\\evgeny\\Desktop\\some.yml"
const entirePathUnix = "/evgeny/Desktop/some.yml"

const expPathWithoutFileWindows = "\\evgeny\\Desktop\\"
const expPathWithoutFileUnix = "/evgeny/Desktop/"

const expFileName = "some.yml"

func TestFolderExists(t *testing.T) {
	req := require.New(t)

	_ = os.MkdirAll(dirName, os.ModePerm)
	defer DeleteDirectory(dirName)

	isExist := Exists(dirName)

	req.Equal(true, isExist)
}

func TestFolderNotExists(t *testing.T) {
	req := require.New(t)

	isExist := Exists(dirName)

	req.Equal(false, isExist)
}

func TestDeleteDirectory(t *testing.T) {
	req := require.New(t)

	_ = os.MkdirAll(dirName, os.ModePerm)

	isExist := Exists(dirName)
	req.Equal(true, isExist)

	DeleteDirectory(dirName)

	isExist = Exists(dirName)
	req.Equal(false, isExist)
}

func TestGetFileNameAndDirFromPath(t *testing.T) {
	req := require.New(t)

	if runtime.GOOS == "windows" {
		fileName, dir := GetFileNameAndDirFromPath(entirePathWindows)
		req.Equal(expFileName, fileName)
		req.Equal(expPathWithoutFileWindows, dir)
	} else {
		fileName, dir := GetFileNameAndDirFromPath(entirePathUnix)
		req.Equal(expFileName, fileName)
		req.Equal(expPathWithoutFileUnix, dir)
	}

}
