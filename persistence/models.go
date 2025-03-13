package persistence

import "time"

type BackupEntity struct {
	StoragePath string    `firestore:"storagePath"`
	BucketName  string    `firestore:"bucketName"`
	UploadedAt  time.Time `firestore:"uploadedAt"`
	Location    string    `firestore:"location"`
	IsDeleted   bool      `firestore:"isDeleted"`
}
