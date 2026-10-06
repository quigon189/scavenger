package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"scavenger/internal/domain"
	"scavenger/internal/repository"
	"time"
)

type groupRepo struct{ q queryer }

func (r *groupRepo) Save(ctx context.Context, g *domain.Group) error {
	if g.ID == 0 {
		createdAt := time.Now()
		res, err := r.q.ExecContext(ctx,
			`INSERT INTO groups(number, start_year, end_year, specialty, short_specialty, created_at)
			VALUES (?, ?, ?, ?, ?, ?)`,
			g.Number, g.StartYear, g.EndYear, g.Specialty, g.ShortSpecialty, createdAt.UTC())
		if err != nil {
			return err
		}
		g.ID, err = res.LastInsertId()
		if err != nil {
			return err
		}
		g.CreatedAt = createdAt
		return nil
	}

	updatedAt := time.Now()
	_, err := r.q.ExecContext(ctx,
		`UPDATE groups SET
			number          = ?,
			start_year      = ?,
			end_year        = ?,
			specialty       = ?,
			short_specialty = ?,
			updated_at      = ?
		WHERE id = ?`,
		g.Number, g.StartYear, g.EndYear, g.Specialty, g.ShortSpecialty, updatedAt.UTC(), g.ID)
	if err != nil {
		return err
	}
	g.UpdatedAt = &updatedAt
	return nil
}

func (r *groupRepo) ByID(ctx context.Context, gid int64) (*domain.Group, error) {
	var g domain.Group
	row := r.q.QueryRowContext(ctx, `
		SELECT id, number, start_year, end_year, specialty, short_specialty, created_at, updated_at
		FROM groups
		WHERE id = ?
		`, gid)
	err := row.Scan(&g.ID, &g.Number, &g.StartYear, &g.EndYear,
		&g.Specialty, &g.ShortSpecialty, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return &g, nil
}

func (r *groupRepo) ByStudentID(ctx context.Context, sid int64) (*domain.Group, error) {
	var gid int64
	row := r.q.QueryRowContext(ctx, `
		SELECT group_id
		FROM group_students
		WHERE student_id = ?
		`, sid)
	if err := row.Scan(&gid); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return r.ByID(ctx, gid)
}

func (r *groupRepo) List(ctx context.Context) ([]domain.Group, error) {
	rows, err := r.q.QueryContext(ctx, `
	SELECT id, number, start_year, end_year, specialty, short_specialty, created_at, updated_at
	FROM groups
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []domain.Group
	for rows.Next() {
		var group domain.Group
		err := rows.Scan(
			&group.ID,
			&group.Number,
			&group.StartYear,
			&group.EndYear,
			&group.Specialty,
			&group.ShortSpecialty,
			&group.CreatedAt,
			&group.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func (r *groupRepo) StudentIDs(ctx context.Context, gid int64) ([]int64, error) {
	rows, err := r.q.QueryContext(ctx, `
		SELECT student_id
		FROM group_students
		WHERE group_id = ?
		`, gid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var studentIDs []int64
	for rows.Next() {
		var sid int64
		if err := rows.Scan(&sid); err != nil {
			return nil, err
		}
		studentIDs = append(studentIDs, sid)
	}
	return studentIDs, nil
}

func (r *groupRepo) AddStudent(ctx context.Context, gid int64, sid int64) error {
	_, err := r.q.ExecContext(ctx, `
		INSERT INTO group_students(group_id, student_id)
		VALUES (?, ?)
		`, gid, sid)
	return err
}

func (r *groupRepo) RemoveStudent(ctx context.Context, sid int64) error {
	_, err := r.q.ExecContext(ctx, `
		DELETE FROM group_students
		WHERE student_id = ?
		`, sid)
	return err
}
