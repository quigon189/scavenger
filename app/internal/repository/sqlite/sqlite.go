package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"scavenger/internal/repository"
	"time"

	_ "modernc.org/sqlite"
)

type queryer interface {
	ExecContext(ctx context.Context, q string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, q string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, q string, args ...any) *sql.Row
}

type DB struct {
	DB *sql.DB
	q  queryer
}

func Open(ctx context.Context, dsn string) (*DB, error) {
	if err := os.MkdirAll(filepath.Dir(dsn), 0o755); err != nil {
		return nil, fmt.Errorf("mkdir: %w", err)
	}

	dsn = dsn + "?_pragma=journal_mode(WAL)&_pragma=foreign_keys(ON)&_pragma=busy_timeout(5000)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	// TODO: добавить миграции
	return &DB{DB: db, q: db}, nil
}

func (d *DB) Close() error {
	return d.DB.Close()
}

func (d *DB) WithTx(ctx context.Context, fn func(repository.Repository) error) error {
	tx, err := d.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	txDB := &DB{DB: d.DB, q: tx}
	if err := fn(txDB); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (d *DB) Users() repository.UserRepo       { return &userRepo{q: d.q} }
func (d *DB) Sessions() repository.SessionRepo { return &sessionRepo{q: d.q} }
