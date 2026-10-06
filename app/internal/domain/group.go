package domain

import (
	"fmt"
	"time"
)

type Group struct {
	ID             int64
	Number         int64
	StartYear      int64
	EndYear        int64
	Specialty      string
	ShortSpecialty string
	CreatedAt      time.Time
	UpdatedAt      *time.Time
}

func (g *Group) IsActive() bool {
	currentYear := int64(time.Now().Year())
	currentMonth := int64(time.Now().Month())

	if currentYear > g.EndYear || currentYear < g.StartYear {
		return false
	}

	if currentYear == g.StartYear && currentMonth < 9 {
		return false
	}

	if currentYear == g.EndYear && currentMonth > 6 {
		return false
	}

	return true
}

func (g *Group) Name() string {
	return fmt.Sprintf("%s-%d", g.ShortSpecialty, g.Number)
}
