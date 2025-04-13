package main

import (
	"context"
	"fmt"
	"ha-backuper/compression"
	"ha-backuper/config"
	"ha-backuper/jobs"
	"os"
	"path"
	"time"
)

const TempFolderName = "temp"

func main() {
	ctx := context.Background()
	configuration, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	worker := jobs.NewJobWorker(ctx, configuration)

	args := os.Args
	if len(args) == 1 {
		fmt.Println("No arguments provided - running backup.")
		performBackup(worker)
		return
	}

	if len(args) == 2 && args[1] == "restore-latest" {
		fmt.Println("Restoring latest backup.")
		err = worker.RestoreLatestBackup()
		if err != nil {
			panic(err)
		}
		return
	}

	fmt.Println("Unknown arguments provided - exiting.")
}

func performBackup(worker *jobs.JobWorker) {
	backupPath, err := createBackupPath(worker.Config.LocationIdentifier)
	if err != nil {
		panic(err)
	}

	err = compression.CompressFolder(worker.Config.HomeAssistantPath, backupPath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Backup compressed at %s\n", backupPath)

	result, err := worker.Backup(backupPath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Backup created cloud storage at %s\n", result.StoragePath)

	err = worker.FirestoreService.InsertBackupEntity(result)
	if err != nil {
		panic(err)
	}
	fmt.Println("Backup created in firestore")

	err = os.Remove(backupPath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Backup deleted at %s\n", backupPath)

	cleanupResult, err := worker.CleanupStorageBackups()
	if err != nil {
		panic(err)
	}
	fmt.Printf("CleanupStorageBackups result: %v\n", cleanupResult)

	deletedDocumentIds, err := worker.FirestoreService.DeleteOldBackupRecords()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deleted old backup records: %v\n", deletedDocumentIds)

	fmt.Println("Backup process done")
}

func createBackupPath(location string) (string, error) {
	formattedNow := time.Now().Format("2006-01-02T15-04-05")
	filename := formattedNow + "-" + location + ".zip"
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	tempFolderPath := path.Join(currentDir, TempFolderName)
	err = os.MkdirAll(tempFolderPath, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create temp folder: %w", err)
	}

	backupFilePath := path.Join(tempFolderPath, filename)
	return backupFilePath, nil
}
