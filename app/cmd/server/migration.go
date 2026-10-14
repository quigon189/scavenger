package main

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func Migrate(ctx context.Context, dialect string, db *sql.DB) error {
	if err := goose.SetDialect(dialect); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}
