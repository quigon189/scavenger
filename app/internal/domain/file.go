package domain

import "time"

type File struct {
	ID           string
	StorageKey   string
	OriginalName string
	ContentType  string
	SHA256       string
	Size         int64
	OwnerID      int64
	CreatedAt    time.Time
}
