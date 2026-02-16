package authrepo

import (
	"context"
	"errors"

	"scavenger/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgUserRepo struct {
	db *pgxpool.Pool
}

func NewPgUserRepository(db *pgxpool.Pool) *PgUserRepo {
	return &PgUserRepo{db: db}
}

func (r *PgUserRepo) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO auth.users (username, email, name, password_hash, role, group_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, status, created_at, updated_at
	`

	return r.db.QueryRow(ctx, query,
		user.Username,
		user.Email,
		user.Name,
		user.PasswordHash,
		user.Role,
		user.GroupID,
	).Scan(&user.ID, &user.Status, &user.CreatedAt, &user.UpdatedAt)
}

func (r *PgUserRepo) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, username, email, name, password_hash, role, status, group_id, created_at, updated_at
		FROM auth.users
		WHERE username = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.Role,
		&user.Status,
		&user.GroupID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *PgUserRepo) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	query := `
		SELECT id, username, email, name, password_hash, role, status, group_id, created_at, updated_at
		FROM auth.users
		WHERE id = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.Role,
		&user.Status,
		&user.GroupID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil

}

func (r *PgUserRepo) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE auth.users
		SET username = $1, email = $2, name = $3, password_hash = $4, status = $5, group_id = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
		RETURNING updated_at
	`

	return r.db.QueryRow(ctx, query,
		user.Username,
		user.Email,
		user.Name,
		user.PasswordHash,
		user.Status,
		user.GroupID,
		user.ID,
	).Scan(&user.UpdatedAt)
}

func (r *PgUserRepo) UpdatePassword(ctx context.Context, userID int, passwordHash string) error {
	query := `
		UPDATE auth.users
		SET password_hash = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	_, err := r.db.Exec(ctx, query, passwordHash, userID)
	return err
}

func (r *PgUserRepo) DeleteUser(ctx context.Context, userID int) error {
	query := `DELETE FROM auth.users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

func (r *PgUserRepo) HealthCheack(ctx context.Context) error {
	return r.db.Ping(ctx)
}
