package sync

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "github.com/glebarez/go-sqlite" // registers "sqlite" driver for sql.Open
)

type LoomDBManager struct {
	dataDir string
	mu      sync.Mutex
	dbs     map[string]*managedDB
	locks   map[string]*sync.Mutex
}

type managedDB struct {
	db       *sql.DB
	lastUsed time.Time
}

func NewLoomDBManager(dataDir string) *LoomDBManager {
	return &LoomDBManager{
		dataDir: dataDir,
		dbs:     make(map[string]*managedDB),
		locks:   make(map[string]*sync.Mutex),
	}
}

func (m *LoomDBManager) Open(diskPath string) (*sql.DB, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if mdb, ok := m.dbs[diskPath]; ok {
		mdb.lastUsed = time.Now()
		return mdb.db, nil
	}

	dbPath := filepath.Join(m.dataDir, diskPath, "loom.db")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("loom database not found: %s", dbPath)
	}

	db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(wal)&_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, fmt.Errorf("open loom db: %w", err)
	}

	db.SetMaxOpenConns(4)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(time.Hour)

	m.dbs[diskPath] = &managedDB{db: db, lastUsed: time.Now()}
	return db, nil
}

func (m *LoomDBManager) Lock(loomID string) {
	m.mu.Lock()
	lk, ok := m.locks[loomID]
	if !ok {
		lk = &sync.Mutex{}
		m.locks[loomID] = lk
	}
	m.mu.Unlock()
	lk.Lock()
}

func (m *LoomDBManager) Unlock(loomID string) {
	m.mu.Lock()
	lk := m.locks[loomID]
	m.mu.Unlock()
	if lk != nil {
		lk.Unlock()
	}
}

func (m *LoomDBManager) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, mdb := range m.dbs {
		mdb.db.Close()
	}
	m.dbs = make(map[string]*managedDB)
}
