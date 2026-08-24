package query

import (
	"sort"

	"drivingmaterials/domain"
)

func FilterMaterials(materials []domain.Material, filter domain.MaterialFilter) []domain.Material {
	filter = filter.Normalize()
	result := make([]domain.Material, 0)
	for _, m := range materials {
		if filter.Accepts(m) {
			result = append(result, m)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Title < result[j].Title })
	return result
}

func LatestByCampus(materials []domain.Material) map[string]domain.Material {
	latest := map[string]domain.Material{}
	for _, m := range materials {
		old, ok := latest[m.Campus]
		if !ok || m.VersionDate > old.VersionDate {
			latest[m.Campus] = m
		}
	}
	return latest
}

func KindsPresent(materials []domain.Material) []domain.MaterialKind {
	found := map[domain.MaterialKind]bool{}
	for _, m := range materials {
		found[m.Kind] = true
	}
	result := make([]domain.MaterialKind, 0, len(found))
	for k := range found {
		result = append(result, k)
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result
}
