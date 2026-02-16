package models

import "context"

type UserRole string

const (
	StudentRole UserRole = "student"
	AdminRole   UserRole = "admin"
	TeacherRole UserRole = "teacher"
)

type User struct {
	ID           int
	Username     string `yaml:"username"`
	Name         string `yaml:"name"`
	PasswordHash string `yaml:"password"`
	Theme        string

	RoleID   int
	RoleName string

	GroupID   int
	GroupName string
}

type AdminStats struct {
	TotalReports     int
	PendingReports   int
	TotalDisciplines int
	TotalGroups      int
}

func GetUserFromContext(ctx context.Context) *User {
	if user, ok := ctx.Value("user").(User); ok {
		return &user
	}
	return &User{}

}
