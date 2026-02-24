package postgres

import (
	"context"
	"errors"

	"scavenger/core/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type authRepo struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *authRepo {
	return &authRepo{db: db}
}

func (r *authRepo) CreateCode(ctx context.Context, code *models.RegistrationCode) error {
	query := `
	INSERT INTO auth.registration_codes (code, name, email, role, group_id, expires_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(ctx, query,
		code.Code,
		code.Name,
		code.Email,
		code.Role,
		code.GroupID,
		code.ExpiresAt,
	)
	return err
}

func (r *authRepo) GetCode(ctx context.Context, codeStr string) (*models.RegistrationCode, error) {
	query := `
	SELECT code, name, email, role, group_id, expires_at
	FROM auth.registration_codes
	WHERE code = $1
	`

	code := models.RegistrationCode{}
	err := r.db.QueryRow(ctx, query, codeStr).Scan(
		&code.Code,
		&code.Name,
		&code.Email,
		&code.Role,
		&code.GroupID,
		&code.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}
	return &code, nil
}

func (r *authRepo) GetAllCodes(ctx context.Context) ([]models.RegistrationCode, error) {
	query := `
	SELECT code, name, email, role, group_id, expires_at
	FROM auth.registration_codes
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return nil, nil
		}
		return nil, err
	}

	codes := []models.RegistrationCode{}
	for rows.Next() {
		code := models.RegistrationCode{}
		err := rows.Scan(
			&code.Code,
			&code.Name,
			&code.Email,
			&code.Role,
			&code.GroupID,
			&code.ExpiresAt,
		)
		if err != nil {
			return nil, err
		}
		codes = append(codes, code)
	}
	return codes, nil
}

func (r *authRepo) DeleteCode(ctx context.Context, code string) error {
	query := `
	DELETE FROM auth.registration_codes
	WHERE code = $1
	`

	_, err := r.db.Exec(ctx, query, code)
	return err
}

func (r *authRepo) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO auth.users (username, email, name, password_hash, role, group_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(ctx, query,
		user.Username,
		user.Email,
		user.Name,
		user.PasswordHash,
		user.Role,
		user.GroupID,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *authRepo) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, username, email, name, password_hash, role, group_id, created_at, updated_at
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

func (r *authRepo) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	query := `
		SELECT id, username, email, name, password_hash, role, group_id, created_at, updated_at
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

func (r *authRepo) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE auth.users
		SET username = $1, email = $2, name = $3, password_hash = $4, group_id = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
		RETURNING updated_at
	`

	return r.db.QueryRow(ctx, query,
		user.Username,
		user.Email,
		user.Name,
		user.PasswordHash,
		user.GroupID,
		user.ID,
	).Scan(&user.UpdatedAt)
}

func (r *authRepo) UpdatePassword(ctx context.Context, userID int, passwordHash string) error {
	query := `
		UPDATE auth.users
		SET password_hash = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	_, err := r.db.Exec(ctx, query, passwordHash, userID)
	return err
}

func (r *authRepo) DeleteUser(ctx context.Context, userID int) error {
	query := `DELETE FROM auth.users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

func (r *authRepo) HealthCheack(ctx context.Context) error {
	return r.db.Ping(ctx)
}
