package main

import (
	"context"
	"log/slog"
	"scavenger/internal/service"
)

func Seed(ctx context.Context, svc *service.Services) error {
	users := []service.UserInput{
		{
			Email:    "t@edu",
			Password: "123456",
			Role:     "teacher",
			FullName: "Teacher Teacherov",
		},
		{
			Email:    "s@edu",
			Password: "123456",
			Role:     "student",
			FullName: "Student Studentov",
		},
	}

	for _, ui := range users {
		if _, err := svc.User.ByEmail(ctx, ui.Email); err == nil {
			continue
		}

		if u, err := svc.User.Create(ctx, ui); err != nil {
			return err
		} else {
			slog.Debug("create user", "user", u)
		}
	}

	groupInput := service.GroupInput{
		Number:          702,
		StartYear:       2026,
		DurationOfStudy: 3,
		Specialty:       "Веб-разработка",
		ShortSpecialty:  "ВР",
	}

	g, err := svc.Group.Create(ctx, groupInput)
	if err != nil {
		return err
	}

	slog.Debug("create group", "group", g)

	stud, err := svc.User.ByEmail(ctx, "s@edu")
	if err != nil {
		return err
	}

	if err := svc.Group.AddStudent(ctx, g.ID, stud.ID); err != nil {
		return err
	}

	groups, err := svc.Group.List(ctx)
	if err != nil {
		return err
	}

	for _, g := range groups {
		studs, err := svc.Group.Students(ctx, g.ID)
		if err != nil {
			return err
		}
		slog.Debug("group info", "group_number", g.Number, "students", studs)
	}

	return nil
}
