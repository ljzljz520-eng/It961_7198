package persistence

import (
	"errors"

	"go.etcd.io/bbolt"
)

type Counts struct {
	Materials   int
	Courses     int
	Reviews     int
	Audits      int
	Plans       int
	Completions int
}

func (s *Store) Counts() (Counts, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return Counts{}, errors.New("store is closed")
	}
	counts := Counts{}
	err := s.db.View(func(tx *bbolt.Tx) error {
		counts.Materials = tx.Bucket(materialBucket).Stats().KeyN
		counts.Courses = tx.Bucket(courseBucket).Stats().KeyN
		counts.Reviews = tx.Bucket(reviewBucket).Stats().KeyN
		counts.Audits = tx.Bucket(auditBucket).Stats().KeyN
		counts.Plans = tx.Bucket(planBucket).Stats().KeyN
		counts.Completions = tx.Bucket(completionBucket).Stats().KeyN
		return nil
	})
	return counts, err
}
