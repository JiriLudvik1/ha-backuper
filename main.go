package main

import (
	"context"
	"ha-backuper/config"
	"ha-backuper/jobs"
	"path"
	"time"
)

const tempFolder = "temp"

func main() {
	ctx := context.Background()

	configuration, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	backupPath := createBackupPath()
	err = CompressFolder(configuration.HomeAssistantPath, backupPath)
	if err != nil {
		panic(err)
	}

	worker := jobs.NewJobWorker(ctx, configuration)
	err = worker.Backup(backupPath)
	if err != nil {
		panic(err)
	}
}

func createBackupPath() string {
	filename := time.Now().Format("2006-01-02T15-04-05") + ".zip"
	return path.Join(tempFolder, filename)
}
