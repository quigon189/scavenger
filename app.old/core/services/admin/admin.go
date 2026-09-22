package admin

import (
	"context"
	"fmt"
	"math/rand/v2"
	"net/mail"
	"regexp"
	"slices"
	"strings"
	"time"

	"scavenger/core/models"
	"scavenger/core/repositories"
	"scavenger/core/services/base"
	"scavenger/core/storage"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type AdminService struct {
	base.BaseService
	dataRepo     repositories.DataRepo
	userRepo     repositories.UserRepo
	codeRepo     repositories.CodeRepo
	groupRepo    repositories.GroupRepo
	storage      storage.Storage
	bucketName   string
	signedURLTTL time.Duration
}

func NewAdminService(
	dataRepo repositories.DataRepo,
	userRepo repositories.UserRepo,
	codeRepo repositories.CodeRepo,
) *AdminService {
	return &AdminService{
		dataRepo: dataRepo,
		userRepo: userRepo,
		codeRepo: codeRepo,
	}
}

func (s *AdminService) adminRequire(ctx context.Context) error {
	return s.RequireRole(ctx, "admin")
}

func (s *AdminService) GetAdminDashboard(ctx context.Context) (*models.AdminDashboard, error) {
	if err := s.adminRequire(ctx); err != nil {
		return nil, err
	}

	return s.dataRepo.Admin.GetAdminDashboard(ctx)
}

func (s *AdminService) GenerateCode(ctx context.Context, code *models.RegistrationCode) error {
	if err := s.adminRequire(ctx); err != nil {
		return err
	}

	reName := regexp.MustCompile(`^[А-ЯЁ][а-яё]+ [А-ЯЁ][а-яё]+( [А-ЯЁ][а-яё]+)?$`)

	errs := []string{}
	if !reName.MatchString(code.Name) {
		errs = append(errs, "неверно указано имя")
	}
	if _, err := mail.ParseAddress(code.Email); err != nil {
		errs = append(errs, "неверно указана электронная почта")
	}
	if !slices.Contains([]string{"student", "admin", "teacher"}, code.Role) {
		errs = append(errs, "неверно указана роль")
	}
	if code.ExpiresAt.Before(time.Now()) {
		errs = append(errs, "неверно указан срок годности")
	}

	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, "; "))
	}

	for range 10 {
		code.Code = generateCode()
		err := s.codeRepo.CreateCode(ctx, code)
		if err != nil {
			continue
		}
		break
	}
	return nil
}

func (s *AdminService) GetAllCodes(ctx context.Context) ([]models.RegistrationCode, error) {
	if err := s.adminRequire(ctx); err != nil {
		return nil, err
	}

	return s.codeRepo.GetAllCodes(ctx)
}

func (s *AdminService) RevokeCode(ctx context.Context, code string) error {
	if err := s.adminRequire(ctx); err != nil {
		return err
	}
	return s.codeRepo.DeleteCode(ctx, code)
}

func (s *AdminService) GetAllGroups(ctx context.Context) ([]models.Group, error) {
	if err := s.adminRequire(ctx); err != nil {
		return nil, err
	}

	return s.groupRepo.GetAll(ctx)
}

func (s *AdminService) CreateGroup(ctx context.Context, group *models.Group) error {
	if err := s.adminRequire(ctx); err != nil {
		return err
	}
	return s.groupRepo.Create(ctx, group)
}

func generateCode() string {
	code := make([]byte, 8)
	l := len(charset)
	for i := range code {
		code[i] = charset[rand.IntN(l)]
	}
	return string(code)
}
