package sync

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type ObjectStore struct {
	root string
}

func NewObjectStore(dataDir string) *ObjectStore {
	root := filepath.Join(dataDir, "objects")
	os.MkdirAll(root, 0755)
	return &ObjectStore{root: root}
}

func (s *ObjectStore) objectPath(hash string) string {
	if len(hash) < 2 {
		return filepath.Join(s.root, hash)
	}
	return filepath.Join(s.root, hash[:2], hash)
}

func (s *ObjectStore) Has(hash string) bool {
	_, err := os.Stat(s.objectPath(hash))
	return err == nil
}

func (s *ObjectStore) Read(hash string) ([]byte, error) {
	return os.ReadFile(s.objectPath(hash))
}

// HashContent computes a Loom-compatible content hash.
// Format: sha256("blob:" + len + "\0" + content)
func HashContent(content []byte) string {
	h := sha256.New()
	h.Write([]byte("blob:"))
	h.Write([]byte(strconv.Itoa(len(content))))
	h.Write([]byte{0})
	h.Write(content)
	return hex.EncodeToString(h.Sum(nil))
}

func (s *ObjectStore) Write(hash string, content []byte) error {
	// Verify hash using Loom's content-addressed format
	computedHex := HashContent(content)
	if computedHex != hash {
		return fmt.Errorf("hash mismatch: expected %s, got %s", hash, computedHex)
	}

	path := s.objectPath(hash)

	// Already exists — skip
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	// Create prefix directory
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	// Write atomically via temp file
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, content, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func (s *ObjectStore) Delete(hash string) error {
	return os.Remove(s.objectPath(hash))
}

func (s *ObjectStore) Size(hash string) (int64, error) {
	info, err := os.Stat(s.objectPath(hash))
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}
