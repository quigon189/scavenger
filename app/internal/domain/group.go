package domain

import "time"

type Group struct {
	ID             int64
	Number         int64
	StartYear      int64
	EndYear        int64
	Specialty      string
	ShortSpecialty string
	CreatedAt      time.Time
}
