package service

import (
	"context"
	"errors"
	"scavenger/internal/domain"
	"scavenger/internal/repository"
)

var (
	ErrInvalidInput   = errors.New("invalid group input")
	ErrUnactiveGroup  = errors.New("unactive group")
	ErrNotStudentRole = errors.New("user don't have student role")
)

type GroupService struct {
	repo repository.Repository
}

func NewGroup(repo repository.Repository) *GroupService {
	return &GroupService{repo: repo}
}

type GroupInput struct {
	Number          int
	StartYear       int
	DurationOfStudy int
	Specialty       string
	ShortSpecialty  string
}

type GroupFilter struct {
	Number    *int64
	ShortSpecialty *string
	IsActive  *bool
	StartYear *int64
}

func (gf *GroupFilter) Apply(groups []domain.Group) []domain.Group {
	var filteredGroups []domain.Group
	if gf == nil {
		return groups
	}
	
	for _, group := range groups {
		if gf.Number != nil && *gf.Number != group.Number {
			continue
		}
		if gf.IsActive != nil && *gf.IsActive != group.IsActive() {
			continue
		}
		if gf.ShortSpecialty != nil && *gf.ShortSpecialty != group.ShortSpecialty {
			continue
		}
		if gf.ShortSpecialty != nil && *gf.StartYear != group.StartYear {
			continue
		}
		filteredGroups = append(filteredGroups, group)
	}

	return filteredGroups
}

func (s *GroupService) Create(ctx context.Context, in GroupInput) (*domain.Group, error) {
	if in.Number == 0 {
		return nil, ErrInvalidInput
	}
	if in.StartYear < 2000 || in.StartYear > 3000 {
		return nil, ErrInvalidInput
	}
	if in.DurationOfStudy == 0 {
		return nil, ErrInvalidInput
	}
	if in.Specialty == "" || in.ShortSpecialty == "" {
		return nil, ErrInvalidInput
	}

	group := domain.Group{
		Number:         int64(in.Number),
		StartYear:      int64(in.StartYear),
		EndYear:        int64(in.StartYear + in.DurationOfStudy),
		Specialty:      in.Specialty,
		ShortSpecialty: in.ShortSpecialty,
	}

	if err := s.repo.Groups().Save(ctx, &group); err != nil {
		return nil, err
	}
	return &group, nil
}

func (s *GroupService) ByID(ctx context.Context, gid int64) (*domain.Group, error) {
	return s.repo.Groups().ByID(ctx, gid)
}

func (s *GroupService) ByStudentID(ctx context.Context, sid int64) (*domain.Group, error) {
	user, err := s.repo.Users().ByID(ctx, sid)
	if err != nil {
		return nil, err
	}
	if user.Role != domain.RoleStudent {
		return nil, ErrNotStudentRole
	}
	return s.repo.Groups().ByStudentID(ctx, user.ID)
}

func (s *GroupService) List(ctx context.Context, filter *GroupFilter) ([]domain.Group, error) {
	groups, err := s.repo.Groups().List(ctx)
	if err != nil {
		return nil, err
	}

	return filter.Apply(groups), nil
}

func (s *GroupService) Students(ctx context.Context, gid int64) ([]domain.User, error) {
	var students []domain.User

	studentIDs, err := s.repo.Groups().StudentIDs(ctx, gid)
	if err != nil {
		return nil, err
	}

	for _, id := range studentIDs {
		student, err := s.repo.Users().ByID(ctx, id)
		if err != nil {
			return nil, err
		}

		students = append(students, *student)
	}

	return students, nil
}

func (s *GroupService) AddStudent(ctx context.Context, gid int64, sid int64) error {
	group, err := s.repo.Groups().ByID(ctx, gid)
	if err != nil {
		return err
	}
	if !group.IsActive() {
		return ErrUnactiveGroup
	}
	user, err := s.repo.Users().ByID(ctx, sid)
	if err != nil {
		return err
	}
	if !user.IsActive {
		return ErrUnactiveUser
	}
	if user.Role != domain.RoleStudent {
		return ErrNotStudentRole
	}
	return s.repo.Groups().AddStudent(ctx, group.ID, user.ID)
}

func (s *GroupService) RemoveStudent(ctx context.Context, sid int64) error {
	user, err := s.repo.Users().ByID(ctx, sid)
	if err != nil {
		return err
	}

	group, err := s.repo.Groups().ByStudentID(ctx, user.ID)
	if err != nil {
		return err
	}

	if !group.IsActive() {
		return ErrUnactiveGroup
	}

	return s.repo.Groups().RemoveStudent(ctx, user.ID)
}
