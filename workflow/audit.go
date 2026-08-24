package workflow

import (
	"sort"

	"drivingmaterials/persistence"
)

type AuditService struct{ store *persistence.Store }

func NewAuditService(store *persistence.Store) *AuditService { return &AuditService{store: store} }

func (s *AuditService) ForMaterial(id string) ([]persistence.AuditEvent, error) {
	events, err := s.store.ListAudits()
	if err != nil {
		return nil, err
	}
	filtered := make([]persistence.AuditEvent, 0)
	for _, event := range events {
		if event.MaterialID == id {
			filtered = append(filtered, event)
		}
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].ID < filtered[j].ID })
	return filtered, nil
}

func (s *AuditService) Actions(id string) ([]string, error) {
	events, err := s.ForMaterial(id)
	if err != nil {
		return nil, err
	}
	actions := make([]string, 0, len(events))
	for _, event := range events {
		actions = append(actions, event.Action)
	}
	return actions, nil
}

func (s *AuditService) Count(id string) (int, error) {
	events, err := s.ForMaterial(id)
	if err != nil {
		return 0, err
	}
	return len(events), nil
}
