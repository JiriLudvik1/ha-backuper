package jobs

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func (j *JobWorker) Backup(backupPath string) error {
	file, err := os.Open(backupPath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
		return err
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
	return nil
}
