package database

import (
	"database/sql"
	"scavenger/internal/models"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

type Database struct {
	db *sql.DB
	cfg models.DatebaseConfig
}

func NewDB(cfg *models.Config, dataSource string) (*Database, error) {
	db, err := sql.Open("sqlite3", dataSource)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Database{db: db, cfg: cfg.DB}, nil
}

func (d *Database) Migrate() error {
	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	if err := goose.Up(d.db, d.cfg.MigrationsPath); err != nil {
		return err
	}

	return nil
}
