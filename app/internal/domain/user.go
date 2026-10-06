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
	UpdatedAt    *time.Time
}

func (r Role) Valid() bool {
	switch r {
	case RoleStudent, RoleTeacher:
		return true
	}
	return false
}

func (u *User) ViewRole() string {
	switch u.Role {
	case RoleStudent:
		return "Студент"
	case RoleTeacher:
		return "Преподаватель"
	default:
		return "Без роли"
	}
}
