package models

import "time"

type File struct {
	ID          int       `json:"id"`
	UUID        string    `json:"uuid"`
	Filename    string    `json:"filename"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	Bucket      string    `json:"bucket"`
	Path        string    `json:"path"`
	CreatedAt   time.Time `json:"created_at"`
}

type LabFile struct {
	ID     int `json:"id"`
	LabID  int `json:"lab_id"`
	FileID int `json:"file_id"`

	Lab  Lab  `json:"lab"`
	File File `json:"file"`
}

type ReportFile struct {
	ID       int `json:"id"`
	ReportID int `json:"report_id"`
	FileID   int `json:"file_id"`

	Report LabReport `json:"report"`
	File   File      `json:"file"`
}
