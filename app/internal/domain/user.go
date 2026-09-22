package domain

import "time"

type Role string

const (
	RoleStudent Role = "student"
	RoleTeacher Role = "teacher"
)

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	FullName     string
	Role         Role
	IsActive     bool
	CreatedAt    time.Time
}

func (r Role) Valid() bool {
	switch r {
	case RoleStudent, RoleTeacher:
		return true
	}
	return false
}
