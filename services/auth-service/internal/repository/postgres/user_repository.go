package postgres

import (
	"context"
	"errors"

	"auth-service/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO auth.users (username, email, name, password_hash, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(ctx, query,
		user.Username,
		user.Email,
		user.Name,
		user.PasswordHash,
		user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, username, email, name, password_hash, role, created_at, updated_at
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

func (r *UserRepository) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	query := `
		SELECT id, username, email, name, password_hash, role, created_at, updated_at
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

func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE auth.users
		SET username = $1, email = $2, name = $3, password_hash = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
		RETURNING updated_at
	`

	return r.db.QueryRow(ctx, query,
		user.Username,
		user.Email,
		user.Name,
		user.PasswordHash,
		user.ID,
	).Scan(&user.UpdatedAt)
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userID int, passwordHash string) error {
	query := `
		UPDATE auth.users
		SET password_hash = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	_, err := r.db.Exec(ctx, query, passwordHash, userID)
	return err
}

func (r *UserRepository) DeleteUser(ctx context.Context, userID int) error {
	query := `DELETE FROM auth.users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

func (r *UserRepository) HealthCheack(ctx context.Context) error {
	return r.db.Ping(ctx)
}
