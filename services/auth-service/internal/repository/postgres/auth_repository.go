package postgres

import (
	"auth-service/internal/models"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO auth.users (username, email, name, password_hash, role, theme)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(ctx, query,
		user.Username,
		user.Email,
		user.Name,
		user.PasswordHash,
		user.Role,
		user.Theme,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *AuthRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, username, email, name, password_hash, role, theme, created_at, updated_at
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
		&user.Theme,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return user, nil
}

func (r *AuthRepository) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	query := `
		SELECT id, username, email, name, password_hash, role, theme, created_at, updated_at
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
		&user.Theme,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return user, nil

}

func (r *AuthRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE auth.users
		SET username = $1, email = $2, name = $3, password_hash = $4, theme = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
	`

	_, err := r.db.Exec(ctx, query,
		user.Username,
		user.Email,
		user.Name,
		user.PasswordHash,
		user.Theme,
		user.ID,
	)

	return err
}

func (r *AuthRepository) UpdatePassword(ctx context.Context, userID int, passwordHash string) error {
	query := `
		UPDATE auth.users
		SET password_hash = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	_, err := r.db.Exec(ctx, query, passwordHash, userID)
	return err
}

func (r *AuthRepository) UpdateTheme(ctx context.Context, userID int, theme string) error {
	query := `
		UPDATE auth.users
		SET theme = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	_, err := r.db.Exec(ctx, query, theme, userID)
	return err
}

func (r *AuthRepository) DeleteUser(ctx context.Context, userID int) error {
	query := `DELETE FROM auth.users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

func (r *AuthRepository) SaveSession(ctx context.Context, userID int, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO auth.sessions (user_id, token, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE
		SET token = $2, expires_at = $3, created_at = CURRENT_TIMESTAMP
	`

	_, err := r.db.Exec(ctx, query, userID, token, expiresAt)
	return err
}

func (r *AuthRepository) GetSession(ctx context.Context, token string) (*models.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.name, u.password_hash, u.role, u.theme, u.created_at, u.updated_at
		FROM auth.sessions s
		JOIN auth.users u ON s.user_id = u.id
		WHERE s.token = $1 AND s.expires_at > CURRENT_TIMESTAMP
	`

	user := &models.User{}
	err := r.db.QueryRow(ctx, query, token).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.Role,
		&user.Theme,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return user, nil
}

func (r *AuthRepository) DeleteSession(ctx context.Context, token string) error {
	query := `DELETE FROM auth.sessions WHERE token = $1`
	_, err := r.db.Exec(ctx, query, token)
	return err
}

func (r *AuthRepository) DeleteExpiriedSessions(ctx context.Context) error {
	query := `DELETE FROM auth.sessions WHERE expires_at <= CURRENT_TIMESTAMP`
	_, err := r.db.Exec(ctx, query)
	return err
}

func (r *AuthRepository) HealthCheack(ctx context.Context) error {
	return r.db.Ping(ctx)
}
