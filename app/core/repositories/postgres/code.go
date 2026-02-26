package postgres

import (
	"context"
	"errors"

	"scavenger/core/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type codeRepo struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *codeRepo {
	return &codeRepo{db: db}
}

func (r *codeRepo) CreateCode(ctx context.Context, code *models.RegistrationCode) error {
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

func (r *codeRepo) GetCode(ctx context.Context, codeStr string) (*models.RegistrationCode, error) {
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

func (r *codeRepo) GetAllCodes(ctx context.Context) ([]models.RegistrationCode, error) {
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

func (r *codeRepo) DeleteCode(ctx context.Context, code string) error {
	query := `
	DELETE FROM auth.registration_codes
	WHERE code = $1
	`

	_, err := r.db.Exec(ctx, query, code)
	return err
}

func (r *codeRepo) HealthCheack(ctx context.Context) error {
	return r.db.Ping(ctx)
}
