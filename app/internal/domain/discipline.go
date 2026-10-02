package domain

import "time"

type Discipline struct {
	ID          int64
	Title       string
	Description string
	TeacherID   int64
	GroupID     int64
	Archived    bool
	CreatedAt   time.Time
}

type MaterialKind string

const (
	MaterialTheory MaterialKind = "theory"
)

type Material interface {
	ID() int64
	Title() string
	Number() int
	DisciplineID() int64
	StartDate() time.Time
	EndDate() time.Time
	Visible() bool
	Kind() MaterialKind
	CreatedAt() time.Time
}

type baseMaterial struct {
	id           int64
	disciplineID int64
	title        string
	number       int
	startDate    time.Time
	endDate      time.Time
	visible      bool
	createdAt    time.Time
}

func (m baseMaterial) ID() int64            { return m.id }
func (m baseMaterial) Title() string        { return m.title }
func (m baseMaterial) Number() int          { return m.number }
func (m baseMaterial) DisciplineID() int64  { return m.disciplineID }
func (m baseMaterial) StartDate() time.Time { return m.startDate }
func (m baseMaterial) EndDate() time.Time   { return m.endDate }
func (m baseMaterial) Visible() bool        { return m.visible }
func (m baseMaterial) CreatedAt() time.Time { return m.createdAt }

type TheoryMaterial struct {
	baseMaterial
	Content string
}

func (m TheoryMaterial) Kind() MaterialKind { return MaterialTheory }
