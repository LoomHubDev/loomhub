package store

import (
	"database/sql"
	"testing"
	"time"

	"github.com/LoomHubDev/loomhub/internal/database"
	"github.com/LoomHubDev/loomhub/internal/models"
	_ "modernc.org/sqlite"
)

func testDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:?_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatal(err)
	}
	if err := database.Migrate(db); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestUserCRUD(t *testing.T) {
	db := testDB(t)
	owners := NewOwnerStore(db)
	users := NewUserStore(db)

	now := time.Now().UTC().Format(time.RFC3339)
	id := "test-user-id"

	// Create owner first
	err := owners.Create(&models.Owner{
		ID: id, Name: "testuser", Type: "user", CreatedAt: now,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create user
	err = users.Create(&models.User{
		ID: id, Username: "testuser", Email: "test@example.com",
		DisplayName: "Test User", PasswordHash: "hash123",
		CreatedAt: now, UpdatedAt: now,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Get by ID
	u, err := users.GetByID(id)
	if err != nil {
		t.Fatal(err)
	}
	if u.Username != "testuser" {
		t.Errorf("Username = %q, want %q", u.Username, "testuser")
	}

	// Get by username
	u, err = users.GetByUsername("testuser")
	if err != nil {
		t.Fatal(err)
	}
	if u.Email != "test@example.com" {
		t.Errorf("Email = %q, want %q", u.Email, "test@example.com")
	}

	// Get by email
	u, err = users.GetByEmail("test@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if u.ID != id {
		t.Errorf("ID = %q, want %q", u.ID, id)
	}

	// Count
	count, err := users.Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("Count = %d, want 1", count)
	}

	// Update
	u.DisplayName = "Updated Name"
	u.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	if err := users.Update(u); err != nil {
		t.Fatal(err)
	}

	u2, _ := users.GetByID(id)
	if u2.DisplayName != "Updated Name" {
		t.Errorf("DisplayName = %q, want %q", u2.DisplayName, "Updated Name")
	}
}

func TestOwnerNameUniqueness(t *testing.T) {
	db := testDB(t)
	owners := NewOwnerStore(db)

	now := time.Now().UTC().Format(time.RFC3339)

	err := owners.Create(&models.Owner{ID: "1", Name: "taken", Type: "user", CreatedAt: now})
	if err != nil {
		t.Fatal(err)
	}

	// Same name should fail
	err = owners.Create(&models.Owner{ID: "2", Name: "taken", Type: "org", CreatedAt: now})
	if err == nil {
		t.Error("expected error for duplicate owner name")
	}

	exists, _ := owners.NameExists("taken")
	if !exists {
		t.Error("NameExists should return true")
	}

	exists, _ = owners.NameExists("available")
	if exists {
		t.Error("NameExists should return false for non-existent name")
	}
}
