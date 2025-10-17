package auth

import (
	"scavenger/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(username, name, password, role string) (*models.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &models.User{
		Username: username,
		Name: name,
		RoleName: role,
		PasswordHash: string(passwordHash),
	}, nil
}

func GeneratePassHash(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passwordHash), nil
}

func ComparePassword(user models.User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
}
