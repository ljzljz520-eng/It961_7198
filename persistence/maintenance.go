package persistence

import (
	"errors"
	"fmt"
	"sort"

	"drivingmaterials/domain"
	"go.etcd.io/bbolt"
)

type IntegrityIssue struct {
	Bucket  string
	Key     string
	Message string
}

func (s *Store) ValidateIntegrity() ([]IntegrityIssue, error) {
	issues := make([]IntegrityIssue, 0)
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return nil, errors.New("store is closed")
	}
	err := s.db.View(func(tx *bbolt.Tx) error {
		for _, name := range [][]byte{materialBucket, courseBucket, reviewBucket, auditBucket, planBucket, completionBucket} {
			bucket := tx.Bucket(name)
			if bucket == nil {
				issues = append(issues, IntegrityIssue{Bucket: string(name), Message: "bucket missing"})
				continue
			}
			err := bucket.ForEach(func(key, raw []byte) error {
				if len(key) == 0 {
					issues = append(issues, IntegrityIssue{Bucket: string(name), Message: "empty key"})
				}
				if len(raw) == 0 {
					issues = append(issues, IntegrityIssue{Bucket: string(name), Key: string(key), Message: "empty value"})
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return issues, err
}

func (s *Store) RemoveMaterial(id string) error {
	if id == "" {
		return errors.New("material id is required")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store is closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		if tx.Bucket(materialBucket).Get([]byte(id)) == nil {
			return errNotFound
		}
		return tx.Bucket(materialBucket).Delete([]byte(id))
	})
}

func (s *Store) MaterialIDs() ([]string, error) {
	ids := make([]string, 0)
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return nil, errors.New("store is closed")
	}
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(materialBucket).ForEach(func(key, _ []byte) error { ids = append(ids, string(key)); return nil })
	})
	sort.Strings(ids)
	return ids, err
}

func (s *Store) ReplaceMaterial(m domain.Material) error {
	if err := m.Validate(); err != nil {
		return err
	}
	current, err := s.GetMaterial(m.ID)
	if err != nil {
		return err
	}
	if current.VersionDate > m.VersionDate {
		return fmt.Errorf("cannot replace newer material %s", m.ID)
	}
	return s.PutMaterial(m)
}

func (s *Store) CountByStatus(status domain.MaterialStatus) (int, error) {
	materials, err := s.ListMaterials()
	if err != nil {
		return 0, err
	}
	count := 0
	for _, m := range materials {
		if m.Status == status {
			count++
		}
	}
	return count, nil
}

func (s *Store) CountByCampus(campus string) (int, error) {
	materials, err := s.ListMaterials()
	if err != nil {
		return 0, err
	}
	count := 0
	for _, m := range materials {
		if m.Campus == campus {
			count++
		}
	}
	return count, nil
}

func (s *Store) MaterialsForSubject(subject string) ([]domain.Material, error) {
	materials, err := s.ListMaterials()
	if err != nil {
		return nil, err
	}
	filtered := make([]domain.Material, 0)
	for _, m := range materials {
		if m.MatchesSubject(subject) {
			filtered = append(filtered, m)
		}
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].VersionDate > filtered[j].VersionDate })
	return filtered, nil
}

func (s *Store) DeleteAudit(id string) error {
	if id == "" {
		return errors.New("audit id is required")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store is closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(auditBucket).Delete([]byte(id)) })
}

func (s *Store) AuditCountForMaterial(id string) (int, error) {
	audits, err := s.ListAudits()
	if err != nil {
		return 0, err
	}
	count := 0
	for _, event := range audits {
		if event.MaterialID == id {
			count++
		}
	}
	return count, nil
}

func (s *Store) Reopen() (*Store, error) {
	if err := s.Close(); err != nil {
		return nil, err
	}
	return Open(s.path)
}
