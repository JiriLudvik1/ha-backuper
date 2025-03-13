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

func (s *FirestoreService) closeOnContext() {
	<-s.ctx.Done()

	err := s.client.Close()
	if err != nil {
		log.Printf("Error closing firestore client: %v", err)
		return
	}

	log.Println("Firestore client closed")
}
