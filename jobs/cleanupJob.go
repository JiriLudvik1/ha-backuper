package jobs

import (
	"cloud.google.com/go/storage"
	"errors"
	"ha-backuper/persistence"
	"sync"
)

type CleanupResult struct {
	SuccessfulDeletedIds []string
	FailedDeletedIds     []string
}

func (j *JobWorker) Cleanup() (CleanupResult, error) {
	result := CleanupResult{}
	backupsToDelete, err := j.FirestoreService.GetDeletableBackups()
	if err != nil {
		return result, err
	}
	if len(backupsToDelete) == 0 {
		return result, nil
	}

	successfulDeletedIds := make([]string, 0)
	failedDeletedIds := make([]string, 0)
	wg := &sync.WaitGroup{}
	for id, backup := range backupsToDelete {
		wg.Add(1)
		go func(backup *persistence.BackupEntity) {
			defer wg.Done()
			deleteError := j.deleteBackup(backup)
			if deleteError != nil {
				failedDeletedIds = append(failedDeletedIds, id)
				return
			}
			successfulDeletedIds = append(successfulDeletedIds, id)
		}(backup)
	}
	wg.Wait()

	result.SuccessfulDeletedIds = successfulDeletedIds
	result.FailedDeletedIds = failedDeletedIds

	err = j.FirestoreService.SetBackupsAsDeleted(successfulDeletedIds)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (j *JobWorker) deleteBackup(backup *persistence.BackupEntity) error {
	err := j.StorageClient.Bucket(j.Config.BucketName).
		Object(backup.StoragePath).
		Delete(j.Context)

	if err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return nil
		}
		// Return any other errors
		return err
	}

	// Success case: no errors
	return nil
}
