package main

import (
	"context"
	"log/slog"
	"scavenger/internal/service"
	"time"
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

	groupInput = service.GroupInput{
		Number: 502,
		StartYear: 2024,
		DurationOfStudy: 4,
		Specialty: "Системное и сетевое администрирование",
		ShortSpecialty: "СА",
	}

	g, err = svc.Group.Create(ctx, groupInput)
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

	discInput = service.DisciplineInput{
		Title: "Test 2",
		GroupID: g.ID,
		TeacherID: t.ID,
		Archived: false,
	}

	d, err = svc.Discipline.Create(ctx, discInput)
	if err != nil {
		return err
	}

	discInput = service.DisciplineInput{
		Title: "Test 3",
		GroupID: 1,
		TeacherID: t.ID,
		Archived: false,
	}

	d, err = svc.Discipline.Create(ctx, discInput)
	if err != nil {
		return err
	}


	discs, err := svc.Discipline.List(ctx)
	if err != nil {
		return err
	}
	slog.Debug("list disciplines", "discs", discs)

	note, err := svc.Material.CreateNote(ctx, service.MaterialNoteInput{
		DisciplineID: 1,
		Title: "Test Note",
		Description: "This is test note",
		Visible: false,
	})
	if err != nil {
		return err
	}

	slog.Debug("create note", "material", note)

	theory, err := svc.Material.CreateTheory(ctx, service.MaterialTheoryInput{
		DisciplineID: 1,
		Title: "Test Theory",
		Visible: true,
		Content: `# Заголовок первого уровня (H1)

Это стандартный абзац текста. В Goldmark вы можете использовать **жирный текст**, *курсив* и ~~зачёркнутый текст~~.

## Заголовок второго уровня (H2)

### Цитата

> «Goldmark — это быстрый и расширяемый парсер Markdown на языке Go».
> — Документация Hugo

### Списки

Маркированный список:
* Первый элемент
* Второй элемент
  * Вложенный элемент
  * Еще один вложенный элемент
* Третий элемент

Нумерованный список:
1. Шаг первый
2. Шаг второй
3. Шаг третий

Чек-лист (задачи):
- [x] Изучить синтаксис CommonMark
- [x] Настроить парсер Goldmark
- [ ] Написать документацию

### Ссылки и изображения

Пример ссылки: посетите [официальный сайт Go](https://go.dev).

Пример изображения:
![Логотип Markdown](https://example.com)

### Таблица

| Язык | Назначение | Парсер по умолчанию |
| :--- | :--- | :--- |
| **Go** | Веб и утилиты | Goldmark |
| **Ruby** | Блоги | Kramdown |
| **JavaScript** | Сборка и сайты | Remark |

### Формулы

$$I_{1}=U_{1}/R_{2}$$

Еще какая то формула: $i=\delta q/\delta t$

### Блок кода (подсветка синтаксиса)

` + "```go\n" + `package main

import "fmt"

func main() {
    p := "Goldmark"
    fmt.Printf("Hello, %s!\\n", p)
}
` + "```\n" + `
---

*Горизонтальная линия выше разделяет секции документа.*`,
	})

	slog.Debug("create theory", "material", theory)

	pract, err := svc.Material.CreatePractical(ctx, service.MaterialPracticalInput{
		DisciplineID: 1,
		Title: "Test Practical",
		Visible: true,
		Number: 1,
		Content: "This in md Practical material",
		Deadline: time.Now().Add(30 * time.Hour),	
	})

	slog.Debug("create practical", "material", pract)

	materials, err := svc.Material.ListByDisciplineID(ctx, d.ID, true)
	if err != nil {
		return err
	}

	slog.Debug("List materials for discipline", "discipline", d, "materials", materials)

	
	return nil
}
