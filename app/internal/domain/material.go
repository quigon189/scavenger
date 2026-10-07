package domain

import "time"

type MaterialKind string

const (
	MaterialTheoryType    MaterialKind = "theory"
	MaterialPracticalType MaterialKind = "practical"
	MaterialNoteType      MaterialKind = "note"
)

type Material interface {
	GetID() string
	GetDisciplineID() int64
	GetTitle() string
	GetVisible() bool
	GetDisplayOrder() int
	Kind() MaterialKind
}

type BaseMaterial struct {
	ID           string
	DisciplineID int64
	Title        string
	DisplayOrder int
	Visible      bool
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

func (m BaseMaterial) GetID() string          { return m.ID }
func (m BaseMaterial) GetDisciplineID() int64 { return m.DisciplineID }
func (m BaseMaterial) GetTitle() string       { return m.Title }
func (m BaseMaterial) GetVisible() bool       { return m.Visible }
func (m BaseMaterial) GetDisplayOrder() int   { return m.DisplayOrder }

type TheoryMaterial struct {
	BaseMaterial
	Content string
}

func (m TheoryMaterial) Kind() MaterialKind { return MaterialTheoryType }

type NoteMaterial struct {
	BaseMaterial
	Description string
}

func (m NoteMaterial) Kind() MaterialKind { return MaterialNoteType }

type PracticalMaterial struct {
	BaseMaterial
	Deadline time.Time
	Content  string
	Number   int
}

func (m PracticalMaterial) Kind() MaterialKind { return MaterialPracticalType }
