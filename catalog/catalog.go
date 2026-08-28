package catalog

import (
	"fmt"
	"sort"

	"drivingmaterials/domain"
	"drivingmaterials/persistence"
)

type Service struct {
	store *persistence.Store
	clock LabelClock
}

type LabelClock interface{ Today() string }

type StaticClock struct{ Date string }

func (c StaticClock) Today() string { return c.Date }

func NewService(store *persistence.Store, clock LabelClock) *Service {
	if clock == nil {
		clock = StaticClock{Date: "2026-01-01"}
	}
	return &Service{store: store, clock: clock}
}

func (s *Service) RegisterMaterial(m domain.Material) error {
	if m.Status == "" {
		m.Status = domain.StatusDraft
	}
	if err := m.Validate(); err != nil {
		return err
	}
	if _, err := s.store.GetMaterial(m.ID); err == nil {
		return fmt.Errorf("material %s already exists", m.ID)
	}
	if err := s.store.PutMaterial(m); err != nil {
		return err
	}
	return s.store.AppendAudit(persistence.AuditEvent{ID: "register-" + m.ID, Action: "register", MaterialID: m.ID, Actor: m.CreatedBy, Detail: m.Title, CreatedAt: s.clock.Today()})
}

func (s *Service) UpdateMaterial(m domain.Material) error {
	if err := m.Validate(); err != nil {
		return err
	}
	old, err := s.store.GetMaterial(m.ID)
	if err != nil {
		return err
	}
	if old.Status == domain.StatusArchived {
		return fmt.Errorf("archived material cannot be updated")
	}
	if m.UpdatedBy == "" {
		m.UpdatedBy = old.UpdatedBy
	}
	if err := s.store.PutMaterial(m); err != nil {
		return err
	}
	return s.store.AppendAudit(persistence.AuditEvent{ID: "update-" + m.ID + "-" + m.VersionDate, Action: "update", MaterialID: m.ID, Actor: m.UpdatedBy, Detail: m.VersionDate, CreatedAt: s.clock.Today()})
}

func (s *Service) ListAll() ([]domain.Material, error) {
	values, err := s.store.ListMaterials()
	if err != nil {
		return nil, err
	}
	sort.Slice(values, func(i, j int) bool {
		if values[i].VersionDate == values[j].VersionDate {
			return values[i].Title < values[j].Title
		}
		return values[i].VersionDate > values[j].VersionDate
	})
	return values, nil
}

func (s *Service) SetStatus(id string, status domain.MaterialStatus, actor string) (domain.Material, error) {
	m, err := s.store.GetMaterial(id)
	if err != nil {
		return m, err
	}
	if !validTransition(m.Status, status) {
		return m, fmt.Errorf("invalid status transition %s -> %s", m.Status, status)
	}
	m.Status = status
	m.UpdatedBy = actor
	if err := s.store.PutMaterial(m); err != nil {
		return m, err
	}
	if err := s.store.AppendAudit(persistence.AuditEvent{ID: "status-" + id + "-" + string(status), Action: "status", MaterialID: id, Actor: actor, Detail: string(status), CreatedAt: s.clock.Today()}); err != nil {
		return m, err
	}
	return m, nil
}

func validTransition(from, to domain.MaterialStatus) bool {
	if from == domain.StatusDraft && (to == domain.StatusReview || to == domain.StatusArchived) {
		return true
	}
	if from == domain.StatusReview && (to == domain.StatusPublished || to == domain.StatusDraft || to == domain.StatusArchived) {
		return true
	}
	if from == domain.StatusPublished && to == domain.StatusArchived {
		return true
	}
	return false
}
