package components

import "scavenger/internal/domain"

type DisciplineCardView struct {
	Discipline domain.Discipline
	Group      *domain.Group // может быть nil, если группу не удалось/не нужно подгружать
	Teacher    *domain.User
}

// DisciplineGroupView — секция аккордеона с несколькими дисциплинами.
type DisciplineGroupView struct {
	ID       string // уникальный якорь, например "group-702" или "teacher-5"
	Title    string // заголовок секции
	Subtitle string // опционально: "3 дисциплины"
	Items    []DisciplineCardView
}

type NoteCardView struct {
	Material domain.NoteMaterial
}

type TheoryCardView struct {
	Material domain.TheoryMaterial
}

type PracticalCardView struct {
	Material domain.PracticalMaterial
	Grade    *int
	Progress  *Progress
}

type Progress struct {
	Done  int
	Total int
}

func (p *Progress) Percent() int {
	if p == nil || p.Total == 0 {
		return 0
	}
	return p.Done * 100 / p.Total
}

type MaterialStats struct {
	Progress *Progress
	Grade   *int
}

type MaterialCardView struct {
	Note      *NoteCardView
	Theory    *TheoryCardView
	Practical *PracticalCardView
}

func NewMaterialCardView(m domain.Material, stats MaterialStats) MaterialCardView {
	switch v := m.(type) {
	case domain.NoteMaterial:
		return MaterialCardView{Note: &NoteCardView{Material: v}}
	case *domain.NoteMaterial:
		return MaterialCardView{Note: &NoteCardView{Material: *v}}
	case domain.TheoryMaterial:
		return MaterialCardView{Theory: &TheoryCardView{Material: v}}
	case *domain.TheoryMaterial:
		return MaterialCardView{Theory: &TheoryCardView{Material: *v}}
	case domain.PracticalMaterial:
		return MaterialCardView{Practical: &PracticalCardView{
			Material: v,
			Progress: stats.Progress,
			Grade: stats.Grade,
		}}
	case *domain.PracticalMaterial:
		return MaterialCardView{Practical: &PracticalCardView{
			Material: *v,
			Progress: stats.Progress,
			Grade: stats.Grade,
		}}
	default:
		return MaterialCardView{}
	}
}
