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

const tempFolderName = "temp"

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
	fmt.Printf("Backup compressed at %s\n", backupPath)

	worker := jobs.NewJobWorker(ctx, configuration)
	result, err := worker.Backup(backupPath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Backup created cloud storage at %s\n", result.StoragePath)

	err = worker.FirestoreService.BackupCreatedInsert(result)
	if err != nil {
		panic(err)
	}
	fmt.Println("Backup created in firestore")

	err = os.Remove(backupPath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Backup deleted at %s\n", backupPath)

	cleanupResult, err := worker.Cleanup()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Cleanup result: %v\n", cleanupResult)

	fmt.Println("Done")
}

func createBackupPath(location string) (string, error) {
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
