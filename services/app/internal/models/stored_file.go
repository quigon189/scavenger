package models

type StoredFile struct {
	ID       int
	Filename string
	Path     string
	URL      string
	Size     int64
}
