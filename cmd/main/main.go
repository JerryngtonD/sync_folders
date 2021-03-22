package main

import (
	"context"
	"fmt"
	"github.com/fsnotify/fsnotify"
	configs "gopackages/configs"
	internal "gopackages/internal"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
)

func main() {
	config, err := configs.ReadWatcherConfiguration()
	if err != nil {
		fmt.Println("some error was thrown while reading config")
		return
	}

	// initialize infinite waiting till signal of exit
	var wg sync.WaitGroup
	wg.Add(1)

	go watchDir(config, &wg)

	wg.Wait()
}

func watchDir(
	config configs.WatcherConfig,
	wg *sync.WaitGroup,
) {
	// initialize channel for get exit signal
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	// background context for cancellation
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()

	// add folders to watcher observation
	watcher.Add(config.FirstDir)
	watcher.Add(config.SecondDir)

	// initialize logger for writing to text.log
	logFile, err := os.OpenFile("text.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer logFile.Close()
	logger := log.New(logFile, "", log.LstdFlags)

	// infinite process cycle
	for {
		select {
		// watch for events
		case event := <-watcher.Events:
			go handleEvent(event, config, logger, cancel)
		// watch for errors
		case err := <-watcher.Errors:
			fmt.Println("ERROR", err)
		// watch for exit
		case sig := <-c:
			fmt.Printf("Got %s signal. Aborting...\n", sig)
			cancel()
			wg.Done()
		}
	}
}

func handleEvent(
	event fsnotify.Event,
	config configs.WatcherConfig,
	logger *log.Logger,
	cancel context.CancelFunc) {
	fileName, dir := internal.GetFileNameAndDirFromPath(event.Name)

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

	err := syncState(event, pathFrom, pathTo, logger)
	if err != nil {
		fmt.Println("got an error with attempt to sync state")
		cancel()
	}
}

func syncState(event fsnotify.Event, pathFrom string, pathTo string, logger *log.Logger) error {
	if event.Op == fsnotify.Create {
		err := internal.CopyFile(pathFrom, pathTo, logger)
		if err != nil {
			return fmt.Errorf("got error with copying file")
		}
	} else if event.Op == fsnotify.Remove || event.Op == fsnotify.Rename {
		err := internal.DeleteFile(pathTo, logger)
		if err != nil {
			return fmt.Errorf("got error with deleting file")
		}
	}
	return nil
}
