package domain

import "time"

type Report struct {
	ID                string
	MaterialID        string
	StudentID         int64
	Text              string
	AttachmentFileIDs []string
	CreatedAt         time.Time
	UpdatedAt         *time.Time
}
