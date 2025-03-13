package main

import (
	"ha-backuper/config"
	"time"
)

const tempFolder = "/temp"

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	now := time.Now()
	archiveName := now.Format("2006-01-02T15-04-05") + ".zip"
	backupPath := tempFolder + archiveName

	err = CompressFolder(config.HomeAssistantPath, backupPath)
	if err != nil {
		panic(err)
	}
}
