package models

import "time"

type Period struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	HalfYear  int       `json:"half_year"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
}
