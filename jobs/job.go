package jobs

import (
	"cloud.google.com/go/storage"
	"context"
	"ha-backuper/config"
	"ha-backuper/persistence"
)

type JobWorker struct {
	Context          context.Context
	StorageClient    *storage.Client
	FirestoreService *persistence.FirestoreService
	Config           *config.BackuperConfig
}
