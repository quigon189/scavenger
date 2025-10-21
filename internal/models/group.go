package models

import "strconv"

type Group struct {
	ID          int
	Name        string
	Disciplines []Discipline
	Students    []User
}

func (g *Group) StudentCount() int {
	return len(g.Students)
}

func (g *Group) IDtoStr() string {
	return strconv.Itoa(g.ID)
}
