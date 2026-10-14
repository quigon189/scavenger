package components

import (
	"fmt"
	"scavenger/internal/domain"
	"scavenger/web/templates/viewmodels"
	"strconv"
)

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
	case *domain.NoteMaterial:
		return *v, true
	}
	return domain.NoteMaterial{}, false
}

func asTheory(m domain.Material) (domain.TheoryMaterial, bool) {
	switch v := m.(type) {
	case *domain.TheoryMaterial:
		return *v, true
	}
	return domain.TheoryMaterial{}, false
}

func asPractical(m domain.Material) (domain.PracticalMaterial, bool) {
	switch v := m.(type) {
	case *domain.PracticalMaterial:
		return *v, true
	}
	return domain.PracticalMaterial{}, false
}

func materialID(v viewmodels.MaterialCardView) string {
	switch {
	case v.Note != nil:
		return v.Note.Material.ID
	case v.Theory != nil:
		return v.Theory.Material.ID
	case v.Practical != nil:
		return v.Practical.Material.ID
	}
	return ""
}

func MaterialTitle(v viewmodels.MaterialPageView) string {
	switch {
	case v.Theory != nil:
		return v.Theory.Material.Title
	case v.Practical != nil:
		return fmt.Sprintf("Работа №%d: %s", v.Practical.Material.Number, v.Practical.Material.Title)
	}
	return ""
}

func materialIDFromPage(v viewmodels.MaterialPageView) string {
	switch {
	case v.Theory != nil:
		return v.Theory.Material.ID
	case v.Practical != nil:
		return v.Practical.Material.ID
	}
	return ""
}

func materialVisible(v viewmodels.MaterialCardView) bool {
	switch {
	case v.Note != nil:
		return v.Note.Material.Visible
	case v.Theory != nil:
		return v.Theory.Material.Visible
	case v.Practical != nil:
		return v.Practical.Material.Visible
	}
	return true
}

func itoa64(v int64) string { return strconv.FormatInt(v, 10) }

func humanSize(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for m := n / unit; m >= unit; m /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(n)/float64(div), "KMGTPE"[exp])
}
