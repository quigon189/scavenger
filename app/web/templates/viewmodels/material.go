package viewmodels

import "scavenger/internal/domain"

type MaterialFormData struct {
	CSRF string

	ID           string
	DisciplineID int64
	Type         string // "note" | "theory" | "practical"
	IsEdit       bool

	Title   string
	Visible bool

	// Заметка
	Description string

	Content string //общий для теории и практики

	// Практика
	Number   int
	Deadline string // значение для <input type="datetime-local">, формат "2006-01-02T15:04"

	Attachments []domain.File

	ErrMsg string
}

func (f *MaterialFormData) FormAction() string {
	if f.IsEdit {
		return "/material/edit"
	}
	return "/material/new/" + f.Type
}

type AttachmentItemView struct {
	File domain.File
}
