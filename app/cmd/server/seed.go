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

	isActive := true
	groups, err := svc.Group.List(ctx, &service.GroupFilter{IsActive: &isActive})
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

	t, err := svc.User.ByEmail(ctx, "t@edu")
	if err != nil {
		return err
	}

	discInput := service.DisciplineInput{
		Title: "Test",
		GroupID: g.ID,
		TeacherID: t.ID,
		Archived: false,
	}

	d, err := svc.Discipline.Create(ctx, discInput)
	if err != nil {
		return err
	}

	d.Title = "Test 2"

	svc.Discipline.Update(ctx, d)

	discs, err := svc.Discipline.List(ctx)
	if err != nil {
		return err
	}
	slog.Debug("list disciplines", "discs", discs)

	material, err := svc.Material.CreateNote(ctx, service.MaterialNoteInput{
		DisciplineID: d.ID,
		Title: "Test Note",
		Description: "This is test note",
		Visible: true,
	})
	if err != nil {
		return err
	}

	slog.Debug("create note", "material", material)

	materials, err := svc.Material.ListByDisciplineID(ctx, d.ID)
	if err != nil {
		return err
	}

	slog.Debug("List materials for discipline", "discipline", d, "materials", materials)

	
	return nil
}
