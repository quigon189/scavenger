package models

import "time"

type BackupConfig struct {
	Enabled       bool   `yaml:"enabled"`
	BackupDir     string `yaml:"backup_dir"`
	Schedule      string `yaml:"schedule"`
	MaxBackups    int    `yaml:"max_backups"`
	UpdatedAt     time.Time
}

type BackupLog struct {
	StartedAt  time.Time
	FinishedAt time.Time
	Type       string
	Status     string
	ErrorMsg   string
	BackupPath string
	Size       int64
	Duration   float64
}
