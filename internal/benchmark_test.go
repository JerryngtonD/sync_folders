package internal

import (
	"os"
	"runtime"
	"strings"
	"testing"
)

const filePath = "mock_data.txt"
const firstDirName = "first"

func BenchmarkCopyFile(b *testing.B) {
	os.MkdirAll(firstDirName, os.ModePerm)

	for i := 0; i < b.N; i++ {
		if runtime.GOOS == "windows" {
			CopyFile(filePath, strings.Join([]string{firstDirName, filePath}, "\\"), nil)
		} else {
			CopyFile(filePath, strings.Join([]string{firstDirName, filePath}, "/"), nil)
		}

	}

	defer DeleteDirectory(firstDirName)
}
