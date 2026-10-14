package viewmodels

import "scavenger/internal/domain"

// DisciplineGroupView — секция аккордеона с несколькими дисциплинами.
type DisciplineGroupView struct {
	ID       string // уникальный якорь, например "group-702" или "teacher-5"
	Title    string // заголовок секции
	Subtitle string // опционально: "3 дисциплины"
	Items    []DisciplineCardView
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

type MaterialPageView struct {
	Discipline domain.Discipline
	Group      *domain.Group
	Teacher    *domain.User
	User       *domain.User // текущий пользователь

	Prev domain.Material
	Next domain.Material

	// Тело. Заполнен ровно один из указателей.
	Theory    *TheoryPageBody
	Practical *PracticalPageBody
}

type TheoryPageBody struct {
	Material domain.TheoryMaterial
	HTML     string // уже отрендеренный markdown
}

type PracticalPageBody struct {
	Material domain.PracticalMaterial
	HTML     string

	// Студенту — оценка (nil = не оценено).
	Grade *int

	// Преподавателю — статистика сдач.
	Progress *Progress
}
