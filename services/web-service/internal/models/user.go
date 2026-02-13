package models

import "time"

type User struct {
	ID        int
	Username  string
	Name      string
	Role      string
	Email     string
	Group     string
	Status    string
	GroupID   *int
	CreatedAt time.Time
}

type Group struct {
	ID        int
	Name      string
	CreatedAt time.Time
}
