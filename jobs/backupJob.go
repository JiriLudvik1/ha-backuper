package jobs

import (
	"fmt"
	"ha-backuper/persistence"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func (j *JobWorker) Backup(backupPath string) (*persistence.BackupEntity, error) {
	file, err := os.Open(backupPath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
		return nil, err
	}
	defer file.Close()

	fileName := filepath.Base(backupPath)
	folderName := j.Config.LocationIdentifier
	objectName := fmt.Sprintf("%s/%s", folderName, fileName)
	bucketName := j.Config.BucketName

	bucket := j.StorageClient.Bucket(bucketName)
	object := bucket.Object(objectName)

	writer := object.NewWriter(j.Context)
	defer writer.Close()

	if _, err := file.Seek(0, 0); err != nil {
		log.Fatalf("Failed to reset file pointer: %v", err)
	}

	if _, err = io.Copy(writer, file); err != nil {
		log.Fatalf("Failed to write file to bucket: %v", err)
	}

	fmt.Printf("Successfully backed up file %s to bucket %s in folder %s\n", backupPath, bucketName, folderName)
	result := &persistence.BackupEntity{
		BucketName:  bucketName,
		StoragePath: objectName,
		UploadedAt:  time.Now(),
		Location:    j.Config.LocationIdentifier,
		IsDeleted:   false,
	}
	return result, nil
}
