package database

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestMigrate(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:?_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := Migrate(db); err != nil {
		t.Fatalf("Migrate failed: %v", err)
	}

	// Verify tables exist
	tables := []string{
		"owners", "users", "access_tokens", "organizations", "org_members",
		"looms", "loom_collaborators", "pins", "weave_requests", "wr_comments",
		"reviews", "labels", "wr_labels", "webhooks", "webhook_deliveries",
		"activities", "object_refs",
	}

	for _, table := range tables {
		var name string
		err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
		if err != nil {
			t.Errorf("table %q not found: %v", table, err)
		}
	}
}

func TestMigrateIdempotent(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:?_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Run twice — should not error
	if err := Migrate(db); err != nil {
		t.Fatal(err)
	}
	if err := Migrate(db); err != nil {
		t.Fatalf("second Migrate failed: %v", err)
	}
}
