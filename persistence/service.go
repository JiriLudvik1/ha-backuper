package persistence

import (
	"cloud.google.com/go/firestore"
	"context"
	"google.golang.org/api/option"
	"ha-backuper/config"
	"log"
)

type FirestoreService struct {
	ctx    context.Context
	client *firestore.Client
	config *config.BackuperConfig
}

func NewFirestoreService(ctx context.Context, config *config.BackuperConfig) (*FirestoreService, error) {
	client, err := firestore.NewClient(ctx, config.GcloudProject, option.WithCredentialsFile(config.ServiceAccountPath))
	if err != nil {
		return nil, err
	}

	service := &FirestoreService{
		ctx:    ctx,
		client: client,
		config: config,
	}
	go service.closeOnContext()

	return service, nil
}

func (s *FirestoreService) BackupCreatedInsert(result *BackupEntity) error {
	err := s.writeDocument(s.config.FirestoreCollection, nil, result)
	if err != nil {
		return err
	}

	log.Printf("BackupEntity created: %v", result)
	return nil
}

func (s *FirestoreService) writeDocument(collectionName string, documentId *string, data interface{}) error {
	var docRef *firestore.DocumentRef

	if documentId == nil {
		docRef = s.client.Collection(collectionName).NewDoc()
	} else {
		docRef = s.client.Collection(collectionName).Doc(*documentId)
	}

	_, err := docRef.Set(s.ctx, data)
	return err
}

func (s *FirestoreService) closeOnContext() {
	<-s.ctx.Done()

	err := s.client.Close()
	if err != nil {
		log.Printf("Error closing firestore client: %v", err)
		return
	}

	log.Println("Firestore client closed")
}
