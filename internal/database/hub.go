package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func OpenHub(dataDir, databaseURL string, dev bool) (*gorm.DB, error) {
	cfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}
	if dev {
		cfg.Logger = logger.Default.LogMode(logger.Info)
	}

	if databaseURL != "" && strings.HasPrefix(databaseURL, "postgres") {
		db, err := gorm.Open(postgres.Open(databaseURL), cfg)
		if err != nil {
			return nil, fmt.Errorf("open postgres: %w", err)
		}
		log.Println("Connected to PostgreSQL")
		return db, nil
	}

	// SQLite fallback
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}

	dbPath := filepath.Join(dataDir, "hub.db")
	dsn := dbPath + "?_pragma=journal_mode(wal)&_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)"
	db, err := gorm.Open(sqlite.Open(dsn), cfg)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	log.Println("Connected to SQLite:", dbPath)
	return db, nil
}
