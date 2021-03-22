package main

import (
	"github.com/stretchr/testify/require"
	"gopackages/configs"
	internal "gopackages/internal"
	"os"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

const fileName = "some.yml"

func TestCoreFunction(t *testing.T) {
	req := require.New(t)

	var config configs.WatcherConfig

	if runtime.GOOS == "windows" {
		config = configs.WatcherConfig{
			FirstDir:  "first\\",
			SecondDir: "second\\",
		}
	} else {
		config = configs.WatcherConfig{
			FirstDir:  "first/",
			SecondDir: "second/",
		}
	}

	_ = os.MkdirAll(config.FirstDir, os.ModePerm)

	_ = os.MkdirAll(config.SecondDir, os.ModePerm)

	var wg sync.WaitGroup
	wg.Add(1)

	var isExist = false

	go func() {
		time.Sleep(1 * time.Second)
		// Create the file.
		if runtime.GOOS == "windows" {
			file, err := os.Create(strings.Join([]string{config.FirstDir, fileName}, ""))
			if err != nil {
				panic("Unable to create tag file!")
			}
			defer file.Close()
		} else {
			file, err := os.Create(strings.Join([]string{config.FirstDir, fileName}, ""))
			if err != nil {
				panic("Unable to create tag file!")
			}
			defer file.Close()
		}

	}()

	go func() {
		time.Sleep(2 * time.Second)
		isExist = internal.Exists(strings.Join([]string{config.SecondDir, fileName}, ""))

		wg.Done()
	}()

	go watchDir(config, &wg)
	wg.Wait()

	req.Equal(true, isExist)
}
