package viewmodels

import "scavenger/internal/domain"

type DisciplineCardView struct {
	Discipline domain.Discipline
	Group      *domain.Group // может быть nil, если группу не удалось/не нужно подгружать
	Teacher    *domain.User
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

type MaterialCardView struct {
	Note      *NoteCardView
	Theory    *TheoryCardView
	Practical *PracticalCardView
}
