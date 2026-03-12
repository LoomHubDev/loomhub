package sync

import (
	"testing"
)

func TestObjectStoreWriteAndRead(t *testing.T) {
	dir := t.TempDir()
	store := NewObjectStore(dir)

	content := []byte("hello world")
	hash := HashContent(content)

	if err := store.Write(hash, content); err != nil {
		t.Fatal(err)
	}

	if !store.Has(hash) {
		t.Error("Has should return true after Write")
	}

	got, err := store.Read(hash)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hello world" {
		t.Errorf("Read = %q, want %q", got, "hello world")
	}

	size, err := store.Size(hash)
	if err != nil {
		t.Fatal(err)
	}
	if size != int64(len(content)) {
		t.Errorf("Size = %d, want %d", size, len(content))
	}
}

func TestObjectStoreDedup(t *testing.T) {
	dir := t.TempDir()
	store := NewObjectStore(dir)

	content := []byte("dedup test")
	hash := HashContent(content)

	// Write twice — should not error
	store.Write(hash, content)
	if err := store.Write(hash, content); err != nil {
		t.Fatal(err)
	}
}

func TestObjectStoreHashMismatch(t *testing.T) {
	dir := t.TempDir()
	store := NewObjectStore(dir)

	err := store.Write("badhash", []byte("content"))
	if err == nil {
		t.Error("expected error for hash mismatch")
	}
}

func TestObjectStoreDelete(t *testing.T) {
	dir := t.TempDir()
	store := NewObjectStore(dir)

	content := []byte("delete me")
	hash := HashContent(content)
	store.Write(hash, content)

	if err := store.Delete(hash); err != nil {
		t.Fatal(err)
	}
	if store.Has(hash) {
		t.Error("Has should return false after Delete")
	}
}

