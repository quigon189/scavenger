package models

import "time"

type Discipline struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	TeacherID   int       `json:"teacher_id"`
	GroupID     int       `json:"group_id"`
	PeriodID    int       `json:"period_id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Teacher Teacher `json:"teacher"`
	Group   Group   `json:"group"`
	Period  Period  `json:"period"`

	Labs []Lab `json:"labs"`
}
