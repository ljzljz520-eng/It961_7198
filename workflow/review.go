package workflow

import (
	"fmt"

	"drivingmaterials/catalog"
	"drivingmaterials/domain"
	"drivingmaterials/persistence"
)

type ReviewService struct {
	catalog *catalog.Service
	store   *persistence.Store
}

func NewReviewService(catalogService *catalog.Service, store *persistence.Store) *ReviewService {
	return &ReviewService{catalog: catalogService, store: store}
}

func (s *ReviewService) Submit(review domain.Review) (domain.Material, error) {
	if err := review.Validate(); err != nil {
		return domain.Material{}, err
	}
	m, err := s.store.GetMaterial(review.MaterialID)
	if err != nil {
		return m, err
	}
	if m.Status != domain.StatusReview {
		return m, fmt.Errorf("material is not awaiting review")
	}
	if err := s.store.PutReview(review); err != nil {
		return m, err
	}
	if review.Decision == "approve" {
		return s.catalog.SetStatus(m.ID, domain.StatusPublished, review.Reviewer)
	}
	updated, err := s.catalog.SetStatus(m.ID, domain.StatusDraft, review.Reviewer)
	return updated, err
}

func (s *ReviewService) Pending() ([]domain.Material, error) { return s.catalog.ListAll() }

func (s *ReviewService) ReviewsFor(id string) ([]domain.Review, error) {
	return s.store.ListReviews(id)
}

func (s *ReviewService) Prepare(review domain.Review) error {
	if review.Decision == "reject" && review.Notes == "" {
		return fmt.Errorf("review notes are required")
	}
	return review.Validate()
}
