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

// MaterialFormData хранит введённые значения для повторного показа
// формы при ошибке валидации.
type MaterialFormData struct {
	Type    string // "note" | "theory" | "practical"
	Title   string
	Visible bool

	// Заметка
	Description string

	// Теория
	TheoryContent string

	// Практика
	Number          int
	PracticalContent string
	Deadline        string // значение для <input type="datetime-local">, формат "2006-01-02T15:04"

	// Общее
	DisplayOrder int
}
