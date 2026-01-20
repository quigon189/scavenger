package models

import "strconv"

type Discipline struct {
	ID      int
	Name    string
	GroupID *int

	Group Group
	Labs  []Lab
}

func (d *Discipline) IDtoStr() string {
	return strconv.Itoa(d.ID)
}

func (d *Discipline) AvgMark() string {
	avg := 0
	counter := 0
	for _, lab := range d.Labs {
		for _, report := range lab.Reports {
			if report.Grade != 0 {
				avg = avg + report.Grade
				counter ++
			}
		}
	}

	if avg == 0 {
		return "-"
	}

	return strconv.FormatFloat(float64(avg) / float64(counter), 'f', 2, 64)
}

func (d *Discipline) SubmittedWorks() int {
	counter := 0
	for _, lab := range d.Labs {
		for _, report := range lab.Reports {
			if report.Status == "submitted" {
				counter ++
			}
		}
	}

	return counter
}
