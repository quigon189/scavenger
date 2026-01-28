package models

import "time"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

type Group struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`

	Students    []Student    `json:"students"`
	Disciplines []Discipline `json:"disciplines"`
}

type Student struct {
	ID        int       `json:"id"`
	GroupID   int       `json:"group_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User  User  `json:"user"`
	Group Group `json:"group"`
}

type Teacher struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User User `json:"user"`

	Disciplines []Discipline `json:"disciplines"`
}
