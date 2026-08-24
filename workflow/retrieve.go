package workflow

import (
	"drivingmaterials/catalog"
	"drivingmaterials/domain"
	"drivingmaterials/persistence"
	"drivingmaterials/query"
)

type RetrieveService struct {
	catalog  *catalog.Service
	searcher *query.Searcher
	store    *persistence.Store
}

func NewRetrieveService(catalogService *catalog.Service, searcher *query.Searcher, store *persistence.Store) *RetrieveService {
	return &RetrieveService{catalog: catalogService, searcher: searcher, store: store}
}

type RetrievalReport struct {
	Filter   domain.MaterialFilter
	Results  []query.SearchResult
	Snapshot catalog.Snapshot
	Total    int
}

func (s *RetrieveService) Retrieve(filter domain.MaterialFilter) (RetrievalReport, error) {
	results, err := s.searcher.Search(filter)
	if err != nil {
		return RetrievalReport{}, err
	}
	materials, err := s.catalog.ListAll()
	if err != nil {
		return RetrievalReport{}, err
	}
	snapshot := catalog.BuildSnapshot(materials, filter.Campus, filter.VersionDate)
	return RetrievalReport{Filter: filter, Results: results, Snapshot: snapshot, Total: len(results)}, nil
}

func (s *RetrieveService) Material(id string) (domain.Material, error) {
	return s.store.GetMaterial(id)
}

func (s *RetrieveService) History(id string) ([]persistence.AuditEvent, error) {
	audits, err := s.store.ListAudits()
	if err != nil {
		return nil, err
	}
	result := make([]persistence.AuditEvent, 0)
	for _, event := range audits {
		if event.MaterialID == id {
			result = append(result, event)
		}
	}
	return result, nil
}

func (s *RetrieveService) CompareVersions(id string, otherVersion string) (bool, error) {
	m, err := s.store.GetMaterial(id)
	if err != nil {
		return false, err
	}
	return m.VersionDate == otherVersion, nil
}
