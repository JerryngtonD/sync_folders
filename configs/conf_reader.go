package configs

import (
	"encoding/json"
	"errors"
	"os"
)

const ConfigPath = "configs/watcher.json"

type WatcherConfig struct {
	FirstDir  string
	SecondDir string
}

func ReadWatcherConfiguration() (WatcherConfig, error) {
	file, _ := os.Open(ConfigPath)
	decoder := json.NewDecoder(file)
	config := new(WatcherConfig)
	err := decoder.Decode(&config)
	if err != nil {
		return WatcherConfig{}, errors.New("can't decode config")
	}
	return *config, nil
}
