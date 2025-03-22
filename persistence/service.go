package persistence

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"ha-backuper/config"
	"log"
	"time"
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

func (s *FirestoreService) InsertBackupEntity(result *BackupEntity) error {
	err := s.writeDocument(s.config.FirestoreCollection, nil, result)
	if err != nil {
		return err
	}

	log.Printf("BackupEntity created: %v", result)
	return nil
}

func (s *FirestoreService) GetDeletableBackups() (map[string]*BackupEntity, error) {
	oneWeekAgo := time.Now().AddDate(0, 0, -7)
	query := s.client.Collection(s.config.FirestoreCollection).
		Where("location", "==", s.config.LocationIdentifier).
		Where("isDeleted", "==", false).
		Where("uploadedAt", "<=", oneWeekAgo)

	docs, err := query.Documents(s.ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to execute Firestore query: %w", err)
	}

	result := make(map[string]*BackupEntity)

	for _, doc := range docs {
		var backupEntity BackupEntity
		if err := doc.DataTo(&backupEntity); err != nil {
			return nil, fmt.Errorf("failed to parse document data: %w", err)
		}

		result[doc.Ref.ID] = &backupEntity
	}

	return result, nil
}

func (s *FirestoreService) SetBackupsAsDeleted(backupIds []string) error {
	for _, backupId := range backupIds {
		docRef := s.client.Collection(s.config.FirestoreCollection).Doc(backupId)
		_, err := docRef.Update(s.ctx, []firestore.Update{
			{Path: "isDeleted", Value: true},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *FirestoreService) GetLatestBackup() (*BackupEntity, error) {
	query := s.client.Collection(s.config.FirestoreCollection).
		Where("location", "==", s.config.LocationIdentifier).
		OrderBy("uploadedAt", firestore.Desc).
		Limit(1)
	docs, err := query.Documents(s.ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to execute Firestore query: %w", err)
	}

	if len(docs) == 0 {
		return nil, fmt.Errorf("no backups found")
	}

	var backupEntity BackupEntity
	if err := docs[0].DataTo(&backupEntity); err != nil {
		return nil, fmt.Errorf("failed to parse document data: %w", err)
	}

	return &backupEntity, nil
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
