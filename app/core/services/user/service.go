package user

import (
	"context"
	"fmt"
	"net/mail"
	"regexp"
	"scavenger/core/models"
	"scavenger/core/repositories"
	"scavenger/core/services/base"
	"scavenger/core/utils"
	"slices"
	"strings"
)

type UserService struct {
	base.BaseService
	userRepo repositories.UserRepo
}

func New(userRepo repositories.UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) RegisterUser(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
	existingUser, err := s.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %v", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("user with username %s already exists", req.Username)
	}

	reName := regexp.MustCompile(`^[А-ЯЁ][а-яё]+ [А-ЯЁ][а-яё]+( [А-ЯЁ][а-яё]+)?$`)

	errs := []string{}
	if !reName.MatchString(req.Name) {
		errs = append(errs, "неверно указано имя")
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		errs = append(errs, "неверно указана электронная почта")
	}
	if !slices.Contains([]string{"student", "admin", "teacher"}, req.Role) {
		errs = append(errs, "неверно указана роль")
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf(strings.Join(errs, "; "))
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	user := &models.User{
		Username:     req.Username,
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         req.Role,
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return user, nil

}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return s.userRepo.GetUserByUsername(ctx, username)
}

func (s *UserService) ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	if !utils.CheckPasswordHash(oldPassword, user.PasswordHash) {
		return fmt.Errorf("invalid old password")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %v", err)
	}

	if err := s.userRepo.UpdatePassword(ctx, user.ID, hashedPassword); err != nil {
		return fmt.Errorf("failed to update password: %v", err)
	}

	return nil
}
