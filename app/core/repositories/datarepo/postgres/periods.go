package postgres

import (
	"context"
	"fmt"

	"scavenger/core/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PeriodRepository struct {
	db *pgxpool.Pool
}

func NewPeriodRepository(db *pgxpool.Pool) *PeriodRepository {
	return &PeriodRepository{db: db}
}

func (r *PeriodRepository) Create(ctx context.Context, period *models.Period) error {
	query := `
INSERT INTO data.periods (name, half_year, start_date, end_date)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at
	`

	return r.db.QueryRow(ctx, query,
		period.Name,
		period.HalfYear,
		period.StartDate,
		period.EndDate,
	).Scan(
		&period.ID,
		&period.CreatedAt,
	)
}

func (r *PeriodRepository) GetByID(ctx context.Context, id int) (*models.Period, error) {
	query := `
SELECT id, name, half_year, start_date, end_date, created_at
FROM data.periods
WHERE id = $1
	`

	period := &models.Period{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&period.ID,
		&period.Name,
		&period.HalfYear,
		&period.StartDate,
		&period.EndDate,
		&period.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get period: %v", err)
	}

	return period, nil
}

func (r *PeriodRepository) GetActive(ctx context.Context) (*models.Period, error) {
	query := `
SELECT id, name, half_year, start_date, end_date, created_at
FROM data.periods
WHERE start_date <= CURRENT_DATE
AND end_date >= CURRENT_DATE
ORDER BY start_date DESC
LIMIT 1
	`

	period := &models.Period{}

	err := r.db.QueryRow(ctx, query).Scan(
		&period.ID,
		&period.Name,
		&period.HalfYear,
		&period.StartDate,
		&period.EndDate,
		&period.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get active period: %v", err)
	}

	return period, nil
}

func (r *PeriodRepository) Update(ctx context.Context, period *models.Period) error {
	query := `
UPDATE data.periods
SET name = $1, half_year = $2, start_date = $3, 
    end_date = $4 
WHERE id = $5
	`

	_, err := r.db.Exec(ctx, query,
		period.Name,
		period.HalfYear,
		period.StartDate,
		period.EndDate,
		period.ID,
	)
	return err
}

func (r *PeriodRepository) Delete(ctx context.Context, id int) error {
	query := `
DELETE FROM data.periods
WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *PeriodRepository) GetAll(ctx context.Context) ([]models.Period, error) {
	query := `
SELECT id, name, half_year, start_date, end_date, created_at
FROM data.periods
ORDER BY start_date DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query periods: %v", err)
	}
	defer rows.Close()

	var periods []models.Period

	for rows.Next() {
		var period models.Period

		err := rows.Scan(
			&period.ID,
			&period.Name,
			&period.HalfYear,
			&period.StartDate,
			&period.EndDate,
			&period.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan period: %v", err)
		}

		periods = append(periods, period)
	}

	return periods, nil
}
