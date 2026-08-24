package persistence

import (
	"fmt"
	"os"

	"go.etcd.io/bbolt"
)

func (s *Store) Export(path string) error {
	if path == "" {
		return fmt.Errorf("export path is required")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return fmt.Errorf("store is closed")
	}
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create export: %w", err)
	}
	defer file.Close()
	if err := s.db.View(func(tx *bbolt.Tx) error { _, err := tx.WriteTo(file); return err }); err != nil {
		return fmt.Errorf("export database: %w", err)
	}
	return file.Sync()
}
