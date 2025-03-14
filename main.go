package main

import (
	"context"
	"fmt"
	"ha-backuper/config"
	"ha-backuper/jobs"
	"os"
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

	backupPath, err := createBackupPath(configuration.LocationIdentifier)
	if err != nil {
		panic(err)
	}

	err = CompressFolder(configuration.HomeAssistantPath, backupPath)
	if err != nil {
		panic(err)
	}

	worker := jobs.NewJobWorker(ctx, configuration)
	result, err := worker.Backup(backupPath)
	if err != nil {
		panic(err)
	}
	err = worker.FirestoreService.BackupCreatedInsert(result)
	if err != nil {
		panic(err)
	}
}

func createBackupPath(location string) (string, error) {
	const tempFolderName = "temp"
	formattedNow := time.Now().Format("2006-01-02T15-04-05")
	filename := formattedNow + "-" + location + ".zip"
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	tempFolderPath := path.Join(currentDir, tempFolderName)
	err = os.MkdirAll(tempFolderPath, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create temp folder: %w", err)
	}

	backupFilePath := path.Join(tempFolderPath, filename)
	return backupFilePath, nil
}
