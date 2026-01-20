package services

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"scavenger/internal/database"
	"scavenger/internal/models"

	"github.com/robfig/cron/v3"
)

type BackupService struct {
	config        *models.Config
	db            *database.Database
	fileStorage   string
	cron          cron.Cron
	cronID        cron.EntryID
	stop          chan struct{}
	backupRunning bool
}

func NewBackupService(cfg *models.Config, db *database.Database) *BackupService {
	return &BackupService{
		config:      cfg,
		db:          db,
		fileStorage: cfg.FS.BasePath,
		cron:        *cron.New(),
		stop:        make(chan struct{}),
	}
}

func (s *BackupService) Start() error {
	if s.config.Backup.Enabled && s.config.Backup.Schedule != "" {
		if err := s.scheduleBackup(); err != nil {
			return fmt.Errorf("failed to schedule buckup: %v", err)
		}
	}

	log.Println("Backup service started")
	return nil
}

func (s *BackupService) Stop() {
	if s.cronID != 0 {
		s.cronID = 0
		s.cron.Stop()
	}
	close(s.stop)
	log.Println("Backup service stoped")
}

func (s *BackupService) UpdateConfig(cfg *models.BackupConfig) error {
	if s.cronID != 0 {
		s.Stop()
	}

	s.config.Backup = *cfg

	if s.config.Backup.Enabled && s.config.Backup.Schedule != "" {
		return s.scheduleBackup()
	}

	// Сохранить изменения в конфиге
	// if err := config.SaveConfig(*s); err != nil {
	// 	return fmt.Errorf("failed to save config: %v", err)
	// }
	return nil
}

func (s *BackupService) RunBackup(bkType string) error {
	s.backupRunning = true
	defer func() { s.backupRunning = false }()

	startTime := time.Now()
	backupLog := &models.BackupLog{
		StartedAt: startTime,
		Type:      bkType,
		Status:    "running",
	}

	tempDir, err := os.MkdirTemp("", "backup-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	if err := s.backupDB(tempDir); err != nil {
		backupLog.Status = "error"
		backupLog.ErrorMsg = fmt.Sprintf("Database backup failed: %v", err)
		return err
	}

	if err := s.backupFiles(tempDir); err != nil {
		backupLog.Status = "error"
		backupLog.ErrorMsg = fmt.Sprintf("Database backup failed: %v", err)
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	finalBakupName := fmt.Sprintf("backup_%s.tar.gz", timestamp)
	finalBakupPath := filepath.Join(s.config.Backup.BackupDir, finalBakupName)

	if err := s.createArchive(finalBakupPath, tempDir); err != nil {
		backupLog.Status = "error"
		backupLog.ErrorMsg = fmt.Sprintf("Archive creation failed: %v", err)
		return err
	}

	info, err := os.Stat(finalBakupPath)
	if err != nil {
		backupLog.Status = "error"
		backupLog.ErrorMsg = fmt.Sprintf("Failed to get backup file info: %v", err)
		return err
	}

	backupLog.Status = "success"
	backupLog.BackupPath = finalBakupPath
	backupLog.Size = info.Size()
	backupLog.FinishedAt = time.Now()
	backupLog.Duration = time.Since(startTime).Seconds()

	if err := s.cleanupOldBackups(); err != nil {
		log.Printf("Failed to cleanup old backups: %v", err)
	}

	log.Printf("Backup completed successfully: %s (%.2f MB)", finalBakupPath, float64(info.Size())/1024/1024)
	return nil
}


func (s *BackupService) scheduleBackup() error {
	id, err := s.cron.AddFunc(s.config.Backup.Schedule, func() {
		if s.backupRunning {
			log.Println("Backup already running, skipping...")
			return
		}

		log.Println("Starting scheduled backup...")
		if err := s.RunBackup("schedule"); err != nil {
			log.Printf("Failed to schedule backup: %v", err)
		}
	})
	if err != nil {
		return fmt.Errorf("invalide schedule format: %v", err)
	}

	s.cronID = id
	s.cron.Start()
	return nil
}

func (s *BackupService) backupDB(tempDir string) error {
	sourceDB := s.config.DB.DataSource
	dbBackupPath := filepath.Join(tempDir, sourceDB)
	err := copyFile(sourceDB, dbBackupPath)
	if err != nil {
		return fmt.Errorf("failed to copy database file: %v", err)
	}

	return nil
}

func (s *BackupService) backupFiles(tempDir string) error {
	err := filepath.Walk(s.fileStorage, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(s.fileStorage, path)
		if err != nil {
			return nil
		}

		if !info.IsDir() {
			dbBackupPath := filepath.Join(tempDir, relPath)
			err := copyFile(path, dbBackupPath)
			if err != nil {
				return fmt.Errorf("failed to copy file: %v", err)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk files storage: %v", err)
	}

	return nil
}

func (s *BackupService) createArchive(outputPath, baseDir string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzWriter := gzip.NewWriter(file)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	err = filepath.Walk(baseDir, func(path string, info fs.FileInfo, err error) error {
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(baseDir, path)
        if err != nil {
            return err
        }
        header.Name = relPath
        
        if err := tarWriter.WriteHeader(header); err != nil {
            return err
        }
        
        if !info.IsDir() {
            data, err := os.Open(path)
            if err != nil {
                return err
            }
            defer data.Close()
            
            if _, err := io.Copy(tarWriter, data); err != nil {
                return err
            }
        }

		return nil
	})	
	return err
}

func (s *BackupService) cleanupOldBackups() error {
	entries, err := os.ReadDir(s.config.Backup.BackupDir)
	if err != nil {
		return err
	}

	var backups []os.FileInfo
	for _, entry := range entries {
		if info, err := entry.Info(); err == nil {
			if strings.HasPrefix(info.Name(), "backup_") && strings.HasSuffix(info.Name(), ".tar.gz") {
				backups = append(backups, info)
			}
		}
	}

	sort.Slice(backups, func(i, j int) bool {
		return backups[i].ModTime().After(backups[j].ModTime())
	})

	for i, backup := range backups {
		if i >= s.config.Backup.MaxBackups {
			backupPath := filepath.Join(s.config.Backup.BackupDir, backup.Name())
            if err := os.Remove(backupPath); err != nil {
                log.Printf("Failed to delete old backup %s: %v", backupPath, err)
            } else {
                log.Printf("Deleted old backup: %s", backupPath)
            }
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destanation, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destanation.Close()

	_, err = io.Copy(destanation, source)
	return err
}
