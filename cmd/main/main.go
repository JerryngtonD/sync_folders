package main

import (
	"context"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"sync_folders/configs"
)

//func checkExit(c chan os.Signal) {
//	select {
//	case sig := <-c:
//
//	}
//}

func main() {
	config, err := configs.ReadWatcherConfiguration()
	if err != nil {
		fmt.Errorf("some error was thrown while reading config")
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup
	wg.Add(1)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()

	watcher.Add(config.FirstDir)
	watcher.Add(config.SecondDir)

	go watchDir(watcher, config, c, cancel, &wg)

	wg.Wait()
}

func handleEvent(event fsnotify.Event, config configs.WatcherConfig) {
	fileName, dir := getFileNameAndDirFromPath(event.Name)

	master := config.FirstDir
	slave := config.SecondDir

	masterFilePath := strings.Join([]string{master, fileName}, "")
	slaveFilePath := strings.Join([]string{slave, fileName}, "")

	var pathFrom string
	var pathTo string

	if dir == master {
		pathFrom = masterFilePath
		pathTo = slaveFilePath
	} else if dir == slave {
		pathFrom = slaveFilePath
		pathTo = masterFilePath
	}

	syncState(event, pathFrom, pathTo)
}

func getFileNameAndDirFromPath(path string) (string, string) {
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

func copyFile(fromPath string, toPath string) error {
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
			return fmt.Errorf("unexpected error with copying: ", err)
		}
		return nil
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

func deleteFile(filePath string) error {
	err := os.Remove(filePath)
	return err
}

func syncState(event fsnotify.Event, pathFrom string, pathTo string) {
	if event.Op == fsnotify.Create {
		go copyFile(pathFrom, pathTo)
	} else if event.Op == fsnotify.Remove || event.Op == fsnotify.Rename {
		go deleteFile(pathTo)
	}
}

func watchDir(
	watcher *fsnotify.Watcher,
	config configs.WatcherConfig,
	c chan os.Signal,
	cancel context.CancelFunc,
	wg *sync.WaitGroup) {
	for {
		select {
		// watch for events
		case event := <-watcher.Events:
			handleEvent(event, config)
		// watch for errors
		case err := <-watcher.Errors:
			fmt.Println("ERROR", err)
		// watch for exit
		case sig := <-c:
			fmt.Printf("Got %s signal. Aborting...\n", sig)
			log.Print("a")
			cancel()
			log.Print("b")
			wg.Done()
		}
	}
}
