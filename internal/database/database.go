package database

import (
	"context"
	"database/sql"
	"scavenger/internal/models"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3"
)

type Database struct {
	db  *sql.DB
	cfg models.DatebaseConfig
}

func NewDB(cfg *models.Config) (*Database, error) {
	driver := ydb.MustOpen(context.TODO(), cfg.DB.DataSource) //TODO: разобраться что за контекст вообще
	connector := ydb.MustConnector(driver, ydb.WithPositionalArgs())
	db := sql.OpenDB(connector)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Database{db: db, cfg: cfg.DB}, nil
}

func (d *Database) Migrate() error {
	gooseDB, err := sql.Open("ydb", d.cfg.DataSource+
		"?go_query_mode=scripting&"+
		"go_fake_tx=scripting&"+
		"go_query_bind=declare,numeric",
	)
	if err != nil {
		return err
	}
	defer gooseDB.Close()

	if err := goose.SetDialect("ydb"); err != nil {
		return err
	}

	if err := goose.Up(gooseDB, d.cfg.MigrationsPath); err != nil {
		return err
	}

	return nil
}
