package domain

import (
	"time"
)

type Discipline struct {
	ID          int64
	Title       string
	Description string
	TeacherID   int64
	GroupID     int64
	Archived    bool
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type MaterialKind string

const (
	MaterialTheoryType MaterialKind = "theory"
	MaterialPractical  MaterialKind = "practical"
	MaterialNote       MaterialKind = "note"
)

type Material interface {
	GetID() string
	GetDisciplineID() int64
	GetVisible() bool
	GetDisplayOrder() int
}

type baseMaterial struct {
	ID           string
	DisciplineID int64
	Title        string
	DisplayOrder int
	Visible      bool
	CreatedAt    time.Time
}

func (m baseMaterial) GetID() string          { return m.ID }
func (m baseMaterial) GetDisciplineID() int64 { return m.DisciplineID }
func (m baseMaterial) GetVisible() bool       { return m.Visible }
func (m baseMaterial) GetDisplayOrder() int   { return m.DisplayOrder }

type baseAssignableMaterial struct {
	baseMaterial
	Deadline time.Time
	Grade    int
}

type TheoryMaterial struct {
	baseMaterial
	Content string
}

type NoteMaterial struct {
	baseMaterial
	Description string
}

type PracticalMaterial struct {
	baseAssignableMaterial
	Content string
	Number  int
}
