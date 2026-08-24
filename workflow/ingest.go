package workflow

import (
	"fmt"

	"drivingmaterials/catalog"
	"drivingmaterials/domain"
	"drivingmaterials/persistence"
)

type IngestService struct {
	catalog *catalog.Service
	store   *persistence.Store
}

func NewIngestService(catalogService *catalog.Service, store *persistence.Store) *IngestService {
	return &IngestService{catalog: catalogService, store: store}
}

type IngestRequest struct {
	ID              string
	Title           string
	Kind            domain.MaterialKind
	Subject         string
	Campus          string
	VersionDate     string
	URI             string
	Description     string
	Tags            []string
	DurationSeconds int
	RouteNumber     string
	Creator         string
}

func (s *IngestService) Ingest(req IngestRequest) (domain.Material, error) {
	m := domain.NewMaterial(req.ID, req.Title, req.Kind, req.Subject, req.Campus, req.VersionDate, req.URI, req.Creator)
	m.Description = req.Description
	m.Tags = catalog.NormalizeTags(req.Tags)
	m.DurationSeconds = req.DurationSeconds
	m.RouteNumber = req.RouteNumber
	m = catalog.ApplyTagPolicy(m)
	if err := s.catalog.RegisterMaterial(m); err != nil {
		return m, err
	}
	return m, nil
}

func (s *IngestService) Publish(id, actor string) (domain.Material, error) {
	m, err := s.catalog.SetStatus(id, domain.StatusReview, actor)
	if err != nil {
		return m, err
	}
	return s.catalog.SetStatus(id, domain.StatusPublished, actor)
}

func (s *IngestService) Archive(id, actor string) error {
	_, err := s.catalog.SetStatus(id, domain.StatusArchived, actor)
	return err
}

func (s *IngestService) ValidateRequests(requests []IngestRequest) error {
	materials := make([]domain.Material, 0, len(requests))
	for _, req := range requests {
		materials = append(materials, domain.NewMaterial(req.ID, req.Title, req.Kind, req.Subject, req.Campus, req.VersionDate, req.URI, req.Creator))
	}
	if issues := catalog.ValidateBatch(materials); len(issues) > 0 {
		return fmt.Errorf("batch contains %d issues", len(issues))
	}
	return nil
}
