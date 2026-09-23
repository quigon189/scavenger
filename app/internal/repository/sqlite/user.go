package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"scavenger/internal/domain"
	"scavenger/internal/repository"
	"time"
)

type userRepo struct{ q queryer }

func (r *userRepo) Create(ctx context.Context, u *domain.User) error {
	res, err := r.q.ExecContext(ctx,
		`INSERT INTO users(email, password_hash, full_name, role, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		u.Email, u.PasswordHash, u.FullName, string(u.Role), time.Now().UTC())
	if err != nil {
		return err
	}
	u.ID, _ = res.LastInsertId()
	return nil
}

func (r *userRepo) ByID(ctx context.Context, id int64) (*domain.User, error) {
	return r.scanOne(r.q.QueryRowContext(ctx,
		`SELECT id, email, password_hash, full_name, role, created_at
			FROM users WHERE id = ?`, id))
}

func (r *userRepo) ByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.scanOne(r.q.QueryRowContext(ctx,
		`SELECT id, email, password_hash, full_name, role, created_at
			FROM users WHERE email = ?`, email))
}

func (r *userRepo) List(ctx context.Context) ([]domain.User, error) {
	rows, err := r.q.QueryContext(ctx,
		`SELECT id, email, password_hash, full_name, role, created_at
		FROM users ORDER BY id`)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		var role string
		err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &role, &u.CreatedAt)
		if err != nil {
			return nil, err
		}
		u.Role = domain.Role(role)
		users = append(users, u)
	}

	return users, nil
}

func (r *userRepo) scanOne(row *sql.Row) (*domain.User, error) {
	var u domain.User
	var role string
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &role, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	u.Role = domain.Role(role)
	return &u, nil
}
