package query

import (
	"sort"
	"strings"

	"drivingmaterials/domain"
)

type Facet struct {
	Value string
	Count int
}

type Facets struct {
	Subjects []Facet
	Campuses []Facet
	Versions []Facet
	Kinds    []Facet
	Statuses []Facet
}

func BuildFacets(materials []domain.Material, filter domain.MaterialFilter) Facets {
	filter = filter.Normalize()
	subjects := map[string]int{}
	campuses := map[string]int{}
	versions := map[string]int{}
	kinds := map[string]int{}
	statuses := map[string]int{}
	for _, m := range materials {
		if !filter.Accepts(m) {
			continue
		}
		subjects[m.Subject]++
		campuses[m.Campus]++
		versions[m.VersionDate]++
		kinds[string(m.Kind)]++
		statuses[string(m.Status)]++
	}
	return Facets{Subjects: toFacets(subjects), Campuses: toFacets(campuses), Versions: toFacets(versions), Kinds: toFacets(kinds), Statuses: toFacets(statuses)}
}

func toFacets(values map[string]int) []Facet {
	result := make([]Facet, 0, len(values))
	for value, count := range values {
		result = append(result, Facet{Value: value, Count: count})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Count == result[j].Count {
			return result[i].Value < result[j].Value
		}
		return result[i].Count > result[j].Count
	})
	return result
}

func (f Facets) Total() int {
	total := 0
	for _, facet := range f.Subjects {
		total += facet.Count
	}
	return total
}

func (f Facets) SubjectValues() []string {
	values := make([]string, 0, len(f.Subjects))
	for _, facet := range f.Subjects {
		values = append(values, strings.TrimSpace(facet.Value))
	}
	return values
}

func (f Facets) ContainsKind(kind domain.MaterialKind) bool {
	for _, facet := range f.Kinds {
		if facet.Value == string(kind) {
			return true
		}
	}
	return false
}

func Narrow(materials []domain.Material, campus, subject string) []domain.Material {
	filter := domain.MaterialFilter{Campus: campus, Subject: subject}
	return FilterMaterials(materials, filter)
}

func SearchByTag(materials []domain.Material, tag string) []domain.Material {
	filter := domain.MaterialFilter{Tag: tag}
	return FilterMaterials(materials, filter)
}
