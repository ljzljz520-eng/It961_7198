package query

import (
	"sort"

	"drivingmaterials/domain"
)

type VersionChange struct {
	ID     string
	Title  string
	Older  string
	Newer  string
	Campus string
}

func Compare(materials []domain.Material, campus, subject string) []VersionChange {
	groups := map[string][]domain.Material{}
	for _, m := range materials {
		if (campus == "" || m.Campus == campus) && (subject == "" || m.Subject == subject) {
			groups[m.Title] = append(groups[m.Title], m)
		}
	}
	changes := make([]VersionChange, 0)
	for title, values := range groups {
		if len(values) < 2 {
			continue
		}
		sort.Slice(values, func(i, j int) bool { return values[i].VersionDate < values[j].VersionDate })
		changes = append(changes, VersionChange{ID: values[len(values)-1].ID, Title: title, Older: values[0].VersionDate, Newer: values[len(values)-1].VersionDate, Campus: values[len(values)-1].Campus})
	}
	sort.Slice(changes, func(i, j int) bool { return changes[i].Title < changes[j].Title })
	return changes
}

func Latest(materials []domain.Material) []domain.Material {
	byTitle := map[string]domain.Material{}
	for _, m := range materials {
		old, ok := byTitle[m.Title]
		if !ok || m.VersionDate > old.VersionDate {
			byTitle[m.Title] = m
		}
	}
	result := make([]domain.Material, 0, len(byTitle))
	for _, m := range byTitle {
		result = append(result, m)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Title < result[j].Title })
	return result
}

func AvailableForCoach(materials []domain.Material, campus string) []domain.Material {
	result := make([]domain.Material, 0)
	for _, m := range materials {
		if m.Campus != campus || !m.IsVisible() {
			continue
		}
		result = append(result, m)
	}
	return Latest(result)
}
