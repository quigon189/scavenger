package components

import "scavenger/internal/domain"

func kindLabel(k domain.MaterialKind) string {
	switch k {
	case domain.MaterialNoteType:
		return "Заметка"
	case domain.MaterialTheoryType:
		return "Теория"
	case domain.MaterialPracticalType:
		return "Практика"
	}
	return "Материал"
}

func kindBadgeClass(k domain.MaterialKind) string {
	switch k {
	case domain.MaterialNoteType:
		return "text-bg-warning"
	case domain.MaterialTheoryType:
		return "text-bg-primary"
	case domain.MaterialPracticalType:
		return "text-bg-success"
	}
	return "text-bg-secondary"
}

// asNote / asTheory / asPractical приводят domain.Material к конкретному типу,
// принимая как значение, так и указатель.
func asNote(m domain.Material) (domain.NoteMaterial, bool) {
	switch v := m.(type) {
	case domain.NoteMaterial:
		return v, true
	case *domain.NoteMaterial:
		return *v, true
	}
	return domain.NoteMaterial{}, false
}

func asTheory(m domain.Material) (domain.TheoryMaterial, bool) {
	switch v := m.(type) {
	case domain.TheoryMaterial:
		return v, true
	case *domain.TheoryMaterial:
		return *v, true
	}
	return domain.TheoryMaterial{}, false
}

func asPractical(m domain.Material) (domain.PracticalMaterial, bool) {
	switch v := m.(type) {
	case domain.PracticalMaterial:
		return v, true
	case *domain.PracticalMaterial:
		return *v, true
	}
	return domain.PracticalMaterial{}, false
}
