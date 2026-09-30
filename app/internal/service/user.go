package service

import (
	"context"
	"errors"
	"net/mail"
	"regexp"
	"strings"

	"scavenger/internal/auth"
	"scavenger/internal/domain"
	"scavenger/internal/repository"
)

var (
	ErrInvalidEmail = errors.New("invalid email")
	ErrInvalidFullName = errors.New("invalid full name")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidRole = errors.New("invalid role")
)

type UserService struct {
	repo repository.Repository
}

func NewUser(repo repository.Repository) *UserService {
	return &UserService{repo: repo}
}

type UserInput struct {
	Email    string
	Password string
	Role     string
	FullName string
}

func (s *UserService) Create(ctx context.Context, in UserInput) (*domain.User, error) {
	var user domain.User

	if !isValidEmail(in.Email) {
		return nil, errors.New("invalid email")
	}

	if !isValidFullName(in.FullName) {
		return nil, errors.New("invalid full name")
	}

	if !isValidPassword(in.Password) {
		return nil, errors.New("invalid password")
	}

	passwordHash, err := auth.HashPassword(in.Password)
	if err != nil {
		return nil, err
	}

	role := domain.Role(in.Role)

	if !role.Valid() {
		return nil, errors.New("invalid role")
	}

	user.Email = in.Email
	user.FullName = in.FullName
	user.PasswordHash = passwordHash
	user.Role = role

	if err := s.repo.Users().Create(ctx, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) List(ctx context.Context) ([]domain.User, error) {
	return s.repo.Users().List(ctx)
}

func (s *UserService) Deactivate(ctx context.Context, id int64) error {
	return s.repo.WithTx(ctx, func(r repository.Repository) error {
		user, err := r.Users().ByID(ctx, id)
		if err != nil {
			return err
		}

		user.IsActive = false

		return r.Users().Update(ctx, user)
	})
}

func (s *UserService) Activate(ctx context.Context, id int64) error {
	return s.repo.WithTx(ctx, func(r repository.Repository) error {
		user, err := r.Users().ByID(ctx, id)
		if err != nil {
			return err
		}

		user.IsActive = true

		return r.Users().Update(ctx, user)
	})
}

func (s *UserService) UpdatePassword(ctx context.Context, id int64, password string) error {
	return s.repo.WithTx(ctx, func(r repository.Repository) error {
		user, err := r.Users().ByID(ctx, id)
		if err != nil {
			return err
		}

		if isValidPassword(password) {
			passwordHash, err := auth.HashPassword(password)
			if err != nil {
				return err
			}
			user.PasswordHash = passwordHash

			return r.Users().Update(ctx, user)
		}

		return ErrInvalidPassword
	})
}

func (s *UserService) UpdateFullName(ctx context.Context, id int64, fullName string) error {
	return s.repo.WithTx(ctx, func(r repository.Repository) error {
		user, err := r.Users().ByID(ctx, id)
		if err != nil {
			return err
		}

		if !isValidFullName(fullName) {
			return ErrInvalidFullName
		}

		user.FullName = fullName

		return r.Users().Update(ctx, user)
	})
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isValidFullName(name string) bool {
	name = strings.TrimSpace(name)
	parts := strings.Fields(name)
	if len(parts) < 2 {
		return false
	}

	validNameRegex := regexp.MustCompile(`^[\p{L}\s-]+$`)
	return validNameRegex.MatchString(name)
}

func isValidPassword(password string) bool {
	return len(password) >= 6
}
