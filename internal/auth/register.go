package auth

import (
	"scavenger/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(username, name, password, role, group string) (*models.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &models.User{
		Username: username,
		Name: name,
		Role: role,
		Group: group,
		PasswordHash: string(passwordHash),
	}, nil
}
