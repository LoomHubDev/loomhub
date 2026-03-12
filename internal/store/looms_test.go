package store

import (
	"testing"
	"time"

	"github.com/LoomHubDev/loomhub/internal/models"
)

func TestLoomCRUD(t *testing.T) {
	db := testDB(t)
	owners := NewOwnerStore(db)
	users := NewUserStore(db)
	looms := NewLoomStore(db)

	now := time.Now().UTC().Format(time.RFC3339)

	// Setup owner + user
	owners.Create(&models.Owner{ID: "u1", Name: "flakerimi", Type: "user", CreatedAt: now})
	users.Create(&models.User{
		ID: "u1", Username: "flakerimi", Email: "f@example.com",
		PasswordHash: "hash", CreatedAt: now, UpdatedAt: now,
	})

	// Create loom
	err := looms.Create(&models.Loom{
		ID: "l1", OwnerID: "u1", Name: "my-app", Description: "Test loom",
		Visibility: "public", DefaultStream: "main", DiskPath: "looms/flakerimi/my-app",
		CreatedAt: now, UpdatedAt: now,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Get by owner+name
	l, err := looms.GetByOwnerAndName("flakerimi", "my-app")
	if err != nil {
		t.Fatal(err)
	}
	if l.FullName != "flakerimi/my-app" {
		t.Errorf("FullName = %q, want %q", l.FullName, "flakerimi/my-app")
	}
	if l.Description != "Test loom" {
		t.Errorf("Description = %q, want %q", l.Description, "Test loom")
	}

	// List by owner
	items, total, err := looms.ListByOwner("u1", 50, 0)
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 {
		t.Errorf("total = %d, want 1", total)
	}
	if len(items) != 1 {
		t.Errorf("items = %d, want 1", len(items))
	}

	// Update
	l.Description = "Updated description"
	l.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	if err := looms.Update(l); err != nil {
		t.Fatal(err)
	}

	l2, _ := looms.GetByOwnerAndName("flakerimi", "my-app")
	if l2.Description != "Updated description" {
		t.Errorf("Description = %q, want %q", l2.Description, "Updated description")
	}

	// Delete
	if err := looms.Delete("l1"); err != nil {
		t.Fatal(err)
	}
	_, err = looms.GetByOwnerAndName("flakerimi", "my-app")
	if err == nil {
		t.Error("expected error after delete")
	}
}

func TestLoomNameUniquenessPerOwner(t *testing.T) {
	db := testDB(t)
	owners := NewOwnerStore(db)
	users := NewUserStore(db)
	looms := NewLoomStore(db)

	now := time.Now().UTC().Format(time.RFC3339)

	owners.Create(&models.Owner{ID: "u1", Name: "user1", Type: "user", CreatedAt: now})
	users.Create(&models.User{
		ID: "u1", Username: "user1", Email: "u1@example.com",
		PasswordHash: "hash", CreatedAt: now, UpdatedAt: now,
	})

	looms.Create(&models.Loom{
		ID: "l1", OwnerID: "u1", Name: "app", Visibility: "public",
		DefaultStream: "main", DiskPath: "looms/user1/app", CreatedAt: now, UpdatedAt: now,
	})

	// Same owner + name should fail
	err := looms.Create(&models.Loom{
		ID: "l2", OwnerID: "u1", Name: "app", Visibility: "public",
		DefaultStream: "main", DiskPath: "looms/user1/app2", CreatedAt: now, UpdatedAt: now,
	})
	if err == nil {
		t.Error("expected error for duplicate loom name per owner")
	}
}
