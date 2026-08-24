package catalog

import (
	"fmt"
	"strings"

	"drivingmaterials/domain"
	"drivingmaterials/persistence"
)

type BatchResult struct {
	Accepted []domain.Material
	Issues   []BatchIssue
}

func (s *Service) RegisterBatch(materials []domain.Material) BatchResult {
	result := BatchResult{Accepted: make([]domain.Material, 0), Issues: ValidateBatch(materials)}
	issueIDs := map[string]bool{}
	for _, issue := range result.Issues {
		issueIDs[issue.ID] = true
	}
	for _, material := range materials {
		if issueIDs[material.ID] {
			continue
		}
		material = catalogMaterial(material)
		if err := s.RegisterMaterial(material); err != nil {
			result.Issues = append(result.Issues, BatchIssue{ID: material.ID, Message: err.Error()})
			continue
		}
		result.Accepted = append(result.Accepted, material)
	}
	return result
}

func catalogMaterial(m domain.Material) domain.Material {
	m.Tags = NormalizeTags(m.Tags)
	if m.Status == "" {
		m.Status = domain.StatusDraft
	}
	return ApplyTagPolicy(m)
}

func (s *Service) RegisterCourse(course domain.Course) error {
	if err := EnsureCourse(course); err != nil {
		return err
	}
	return s.store.PutCourse(course)
}

func (s *Service) CourseCatalog(campus string) ([]domain.Course, error) {
	courses, err := s.store.ListCourses()
	if err != nil {
		return nil, err
	}
	filtered := make([]domain.Course, 0)
	for _, course := range courses {
		if campus == "" || strings.EqualFold(course.Campus, campus) {
			filtered = append(filtered, course)
		}
	}
	return filtered, nil
}

func (s *Service) MaterialSummary(id string) (string, error) {
	m, err := s.store.GetMaterial(id)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s | %s | %s | %s | %s", m.ID, m.Title, m.Kind, m.Campus, m.VersionDate), nil
}

func (s *Service) CloneForCampus(id, campus, actor string) (domain.Material, error) {
	m, err := s.store.GetMaterial(id)
	if err != nil {
		return m, err
	}
	if strings.TrimSpace(campus) == "" {
		return m, fmt.Errorf("campus is required")
	}
	m.ID = m.ID + "-" + strings.ToLower(strings.ReplaceAll(campus, " ", "-"))
	m.Campus = campus
	m.CreatedBy = actor
	m.UpdatedBy = actor
	if err := s.RegisterMaterial(m); err != nil {
		return m, err
	}
	return m, nil
}

func (s *Service) ArchiveStale(campus, version, actor string) (int, error) {
	materials, err := s.store.ListMaterials()
	if err != nil {
		return 0, err
	}
	count := 0
	for _, m := range materials {
		if m.Campus != campus || m.VersionDate >= version || m.Status == domain.StatusArchived {
			continue
		}
		if _, err := s.SetStatus(m.ID, domain.StatusArchived, actor); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

func (s *Service) Store() *persistence.Store { return s.store }
