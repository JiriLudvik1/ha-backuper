package jobs

import (
	"cloud.google.com/go/storage"
	"context"
	"google.golang.org/api/option"
	"ha-backuper/config"
	"ha-backuper/persistence"
	"log"
)

type JobWorker struct {
	Context          context.Context
	StorageClient    *storage.Client
	FirestoreService *persistence.FirestoreService
	Config           *config.BackuperConfig
}

func NewJobWorker(ctx context.Context, config *config.BackuperConfig) *JobWorker {
	firestoreService, err := persistence.NewFirestoreService(ctx, config)
	if err != nil {
		panic(err)
	}

	storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile(config.ServiceAccountPath))
	if err != nil {
		panic(err)
	}

	worker := &JobWorker{
		Context:          ctx,
		StorageClient:    storageClient,
		FirestoreService: firestoreService,
		Config:           config,
	}
	go worker.closeOnContext()
	return worker
}

func (j *JobWorker) closeOnContext() {
	<-j.Context.Done()
	err := j.StorageClient.Close()
	if err != nil {
		log.Printf("Error closing storage client: %v", err)
		return
	}
	log.Println("Storage client closed")
}
