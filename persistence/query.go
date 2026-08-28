package persistence

import (
	"errors"
	"sort"

	"drivingmaterials/domain"
)

func (s *Store) MaterialsByKind(kind domain.MaterialKind) ([]domain.Material, error) {
	materials, err := s.ListMaterials()
	if err != nil {
		return nil, err
	}
	result := make([]domain.Material, 0)
	for _, m := range materials {
		if m.Kind == kind {
			result = append(result, m)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Title < result[j].Title })
	return result, nil
}

func (s *Store) MaterialsByVersion(version string) ([]domain.Material, error) {
	if version == "" {
		return nil, errors.New("version is required")
	}
	materials, err := s.ListMaterials()
	if err != nil {
		return nil, err
	}
	result := make([]domain.Material, 0)
	for _, m := range materials {
		if m.VersionDate == version {
			result = append(result, m)
		}
	}
	return result, nil
}

func (s *Store) UpsertCourse(c domain.Course) error {
	if c.Code == "" {
		return errors.New("course code is required")
	}
	return s.PutCourse(c)
}

func (s *Store) ReviewVersions(id string) ([]int, error) {
	reviews, err := s.ListReviews(id)
	if err != nil {
		return nil, err
	}
	versions := make([]int, 0, len(reviews))
	for _, r := range reviews {
		versions = append(versions, r.Version)
	}
	sort.Ints(versions)
	return versions, nil
}
