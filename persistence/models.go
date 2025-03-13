package persistence

import "time"

type Backup struct {
	StoragePath string    `firestore:"storagePath"`
	UploadedAt  time.Time `firestore:"uploadedAt"`
	Location    string    `firestore:"location"`
	IsDeleted   bool      `firestore:"isDeleted"`
}
