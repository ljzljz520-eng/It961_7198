package persistence

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"drivingmaterials/domain"
	"go.etcd.io/bbolt"
)

var (
	errNotFound      = errors.New("record not found")
	materialBucket   = []byte("materials")
	courseBucket     = []byte("courses")
	reviewBucket     = []byte("reviews")
	auditBucket      = []byte("audits")
	planBucket       = []byte("lesson_plans")
	completionBucket = []byte("lesson_completions")
)

type AuditEvent struct {
	ID         string `json:"id"`
	Action     string `json:"action"`
	MaterialID string `json:"material_id"`
	Actor      string `json:"actor"`
	Detail     string `json:"detail"`
	CreatedAt  string `json:"created_at"`
}

type Store struct {
	db   *bbolt.DB
	mu   sync.RWMutex
	path string
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, fmt.Errorf("database path is required")
	}
	db, err := bbolt.Open(filepath.Clean(path), 0600, &bbolt.Options{Timeout: time.Second})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	s := &Store{db: db, path: filepath.Clean(path)}
	if err := s.initialize(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) initialize() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, name := range [][]byte{materialBucket, courseBucket, reviewBucket, auditBucket, planBucket, completionBucket} {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) PutLessonPlan(plan domain.LessonPlan) error {
	if err := plan.Validate(); err != nil {
		return err
	}
	value, err := encode(plan)
	if err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store is closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(planBucket).Put([]byte(plan.ID), value) })
}

func (s *Store) GetLessonPlan(id string) (domain.LessonPlan, error) {
	var plan domain.LessonPlan
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return plan, errors.New("store is closed")
	}
	err := s.db.View(func(tx *bbolt.Tx) error {
		raw := tx.Bucket(planBucket).Get([]byte(id))
		if raw == nil {
			return errNotFound
		}
		return decode(raw, &plan)
	})
	return plan, err
}

func (s *Store) ListLessonPlans(campus, subject string) ([]domain.LessonPlan, error) {
	plans := make([]domain.LessonPlan, 0)
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return nil, errors.New("store is closed")
	}
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(planBucket).ForEach(func(_, raw []byte) error {
			var p domain.LessonPlan
			if err := decode(raw, &p); err != nil {
				return err
			}
			if (campus == "" || p.Campus == campus) && (subject == "" || p.Subject == subject) {
				plans = append(plans, p)
			}
			return nil
		})
	})
	return plans, err
}

func (s *Store) PutCompletion(completion domain.LessonCompletion) error {
	if completion.PlanID == "" || completion.Coach == "" {
		return errors.New("completion plan and coach are required")
	}
	value, err := encode(completion)
	if err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store is closed")
	}
	key := completion.PlanID + ":" + completion.Coach
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(completionBucket).Put([]byte(key), value) })
}

func (s *Store) ListCompletions(planID string) ([]domain.LessonCompletion, error) {
	values := make([]domain.LessonCompletion, 0)
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return nil, errors.New("store is closed")
	}
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(completionBucket).ForEach(func(_, raw []byte) error {
			var c domain.LessonCompletion
			if err := decode(raw, &c); err != nil {
				return err
			}
			if planID == "" || c.PlanID == planID {
				values = append(values, c)
			}
			return nil
		})
	})
	return values, err
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}

func (s *Store) Path() string { s.mu.RLock(); defer s.mu.RUnlock(); return s.path }

func (s *Store) PutMaterial(m domain.Material) error {
	if err := m.Validate(); err != nil {
		return err
	}
	value, err := encode(m)
	if err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store is closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(materialBucket).Put([]byte(m.ID), value) })
}

func (s *Store) GetMaterial(id string) (domain.Material, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var m domain.Material
	if s.db == nil {
		return m, errors.New("store is closed")
	}
	err := s.db.View(func(tx *bbolt.Tx) error {
		raw := tx.Bucket(materialBucket).Get([]byte(id))
		if raw == nil {
			return errNotFound
		}
		return decode(raw, &m)
	})
	return m, err
}

func (s *Store) ListMaterials() ([]domain.Material, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]domain.Material, 0)
	if s.db == nil {
		return nil, errors.New("store is closed")
	}
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(materialBucket).ForEach(func(_, raw []byte) error {
			var m domain.Material
			if err := decode(raw, &m); err != nil {
				return err
			}
			result = append(result, m)
			return nil
		})
	})
	return result, err
}

func (s *Store) PutCourse(c domain.Course) error {
	if err := c.Validate(); err != nil {
		return err
	}
	value, err := encode(c)
	if err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store is closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(courseBucket).Put([]byte(c.Key()), value) })
}

func (s *Store) ListCourses() ([]domain.Course, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]domain.Course, 0)
	if s.db == nil {
		return nil, errors.New("store is closed")
	}
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(courseBucket).ForEach(func(_, raw []byte) error {
			var c domain.Course
			if err := decode(raw, &c); err != nil {
				return err
			}
			result = append(result, c)
			return nil
		})
	})
	return result, err
}

func (s *Store) PutReview(r domain.Review) error {
	if err := r.Validate(); err != nil {
		return err
	}
	value, err := encode(r)
	if err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store is closed")
	}
	key := fmt.Sprintf("%s:%d", r.MaterialID, r.Version)
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(reviewBucket).Put([]byte(key), value) })
}

func (s *Store) ListReviews(materialID string) ([]domain.Review, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]domain.Review, 0)
	if s.db == nil {
		return nil, errors.New("store is closed")
	}
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(reviewBucket).ForEach(func(_, raw []byte) error {
			var r domain.Review
			if err := decode(raw, &r); err != nil {
				return err
			}
			if materialID == "" || r.MaterialID == materialID {
				result = append(result, r)
			}
			return nil
		})
	})
	return result, err
}

func (s *Store) AppendAudit(e AuditEvent) error {
	if e.ID == "" || e.Action == "" || e.Actor == "" {
		return errors.New("audit event requires id, action, and actor")
	}
	value, err := encode(e)
	if err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store is closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(auditBucket).Put([]byte(e.ID), value) })
}

func (s *Store) ListAudits() ([]AuditEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]AuditEvent, 0)
	if s.db == nil {
		return nil, errors.New("store is closed")
	}
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(auditBucket).ForEach(func(_, raw []byte) error {
			var e AuditEvent
			if err := decode(raw, &e); err != nil {
				return err
			}
			result = append(result, e)
			return nil
		})
	})
	return result, err
}
