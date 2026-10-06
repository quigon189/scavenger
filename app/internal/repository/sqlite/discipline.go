package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"scavenger/internal/domain"
	"scavenger/internal/repository"
	"time"
)

type disciplineRepo struct{ q queryer }

func (r *disciplineRepo) Save(ctx context.Context, d *domain.Discipline) error {
	if d.ID == 0 {
		createdAt := time.Now()
		res, err := r.q.ExecContext(ctx, `
			INSERT INTO disciplines (title, description, teacher_id, group_id, archived, created_at)
			VALUES (?, ?, ?, ?, ?, ?)
			`, d.Title, d.Description, d.TeacherID, d.GroupID, d.Archived, createdAt.UTC(), createdAt.UTC())
		if err != nil {
			return err
		}
		d.ID, err = res.LastInsertId()
		if err != nil {
			return err
		}
		d.CreatedAt = createdAt

		return nil
	}

	updatedAt := time.Now()

	_, err := r.q.ExecContext(ctx, `
		UPDATE disciplines SET	
			title       = ?,
			description = ?,
			teacher_id  = ?,
			group_id    = ?,
			archived    = ?,
			updated_at  = ?
		WHERE id = ?
		`, d.Title, d.Description, d.TeacherID, d.GroupID, d.Archived, updatedAt.UTC(), d.ID)
	if err != nil {
		return err
	}

	d.UpdatedAt = &updatedAt

	return nil
}

func (r *disciplineRepo) ByID(ctx context.Context, id int64) (*domain.Discipline, error) {
	return scanDiscipline(r.q.QueryRowContext(ctx, `
		SELECT id, title, description, teacher_id, group_id, archived, created_at, updated_at
		FROM disciplines
		WHERE id = ?
		`, id))
}

func (r *disciplineRepo) ByTeacherID(ctx context.Context, tid int64) ([]domain.Discipline, error) {
	rows, err := r.q.QueryContext(ctx, `
		SELECT id, title, description, teacher_id, group_id, archived, created_at, updated_at
		FROM disciplines
		WHERE teacher_id = ?
		`, tid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanDisciplineRows(rows)
}

func (r *disciplineRepo) ByGroupID(ctx context.Context, gid int64) ([]domain.Discipline, error) {
	rows, err := r.q.QueryContext(ctx, `
		SELECT id, title, description, teacher_id, group_id, archived, created_at, updated_at
		FROM disciplines
		WHERE group_id = ?
		`, gid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanDisciplineRows(rows)
}

func (r *disciplineRepo) List(ctx context.Context) ([]domain.Discipline, error) {
	rows, err := r.q.QueryContext(ctx, `
		SELECT id, title, description, teacher_id, group_id, archived, created_at, updated_at
		FROM disciplines
		`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanDisciplineRows(rows)
}

func (r *disciplineRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.q.ExecContext(ctx, `
		DELETE FROM disciplines
		WHERE id = ?
		`, id)
	return err
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanDiscipline(s rowScanner) (*domain.Discipline, error) {
	var d domain.Discipline

	err := s.Scan(&d.ID, &d.Title, &d.Description, &d.TeacherID, &d.GroupID,
		&d.Archived, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return &d, nil
}

func scanDisciplineRows(rows *sql.Rows) ([]domain.Discipline, error) {
	var disciplines []domain.Discipline
	for rows.Next() {
		d, err := scanDiscipline(rows)
		if err != nil {
			return nil, err
		}
		disciplines = append(disciplines, *d)
	}
	return disciplines, nil
}
