package jobs

import (
	"fmt"
	"ha-backuper/compression"
	"ha-backuper/persistence"
	"io"
	"os"
	"path/filepath"
	"time"
)

const restoreFolderName = "restore"

func (j *JobWorker) RestoreLatestBackup() error {
	latestBackup, err := j.FirestoreService.GetLatestBackup()
	if err != nil {
		return err
	}

	if latestBackup == nil || latestBackup.StoragePath == "" || latestBackup.BucketName == "" {
		return fmt.Errorf("invalid backup entity: missing StoragePath or BucketName")
	}

	localBackupFilePath, err := createRestorePath(latestBackup)
	if err != nil || localBackupFilePath == "" {
		return fmt.Errorf("failed to create restore path: %w", err)
	}

	err = j.getRestoreFile(localBackupFilePath, latestBackup)
	if err != nil {
		return fmt.Errorf("failed to get the backup file: %w", err)
	}

	fmt.Println("Backup successfully obtained")

	err = renameFolderToOldSuffix(j.Config.HomeAssistantPath)
	if err != nil {
		return fmt.Errorf("failed to rename the backup folder: %w", err)
	}

	fmt.Printf("Original folder successfully renamed to %s_old\n", j.Config.HomeAssistantPath)

	err = compression.DecompressFolder(localBackupFilePath, j.Config.HomeAssistantPath)
	if err != nil {
		return fmt.Errorf("failed to decompress the backup: %w", err)
	}

	fmt.Printf("Backup successfully decompressed to %s\n", j.Config.HomeAssistantPath)
	return nil
}

func (j *JobWorker) downloadFileFromBucket(bucketName, objectName, localFilePath string) error {
	bucket := j.StorageClient.Bucket(bucketName)
	object := bucket.Object(objectName)

	reader, err := object.NewReader(j.Context)
	if err != nil {
		return fmt.Errorf("failed to create reader for object %s: %w", objectName, err)
	}
	defer reader.Close()

	localFile, err := os.Create(localFilePath)
	if err != nil {
		return fmt.Errorf("failed to create local file %s: %w", localFilePath, err)
	}
	defer localFile.Close()

	if _, err := io.Copy(localFile, reader); err != nil {
		return fmt.Errorf("failed to write object contents to local file: %w", err)
	}

	return nil
}

func (j *JobWorker) getRestoreFile(localBackupFilePath string, latestBackup *persistence.BackupEntity) error {
	if _, err := os.Stat(localBackupFilePath); os.IsNotExist(err) {
		err = j.downloadFileFromBucket(latestBackup.BucketName, latestBackup.StoragePath, localBackupFilePath)
		if err != nil {
			return fmt.Errorf("failed to download the backup file: %w", err)
		}
		fmt.Printf("Backup successfully downloaded from cloud storage to %s\n", localBackupFilePath)
		return nil
	} else if err != nil {
		return fmt.Errorf("error checking file existence at %s: %w", localBackupFilePath, err)
	}

	fmt.Printf("File already exists at %s. Skipping download.\n", localBackupFilePath)
	return nil
}

func createRestorePath(latestBackup *persistence.BackupEntity) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	restoreFolderPath := filepath.Join(currentDir, restoreFolderName)
	err = os.MkdirAll(restoreFolderPath, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create restore folder: %w", err)
	}

	fileName := "restored" + filepath.Base(latestBackup.StoragePath)
	return filepath.Join(restoreFolderPath, fileName), nil
}

func renameFolderToOldSuffix(folderPath string) error {
	timestamp := time.Now().Format("2006_01_02__15_04_05")
	newFolderName := folderPath + "_old_" + timestamp
	newFolderPath := newFolderName

	err := os.Rename(folderPath, newFolderPath)
	if err != nil {
		return err
	}
	return nil

}
