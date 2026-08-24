package query

import (
	"sort"
	"strings"

	"drivingmaterials/domain"
	"drivingmaterials/persistence"
)

type SearchResult struct {
	Material domain.Material `json:"material"`
	Labels   []string        `json:"labels"`
	Score    int             `json:"score"`
}

type Searcher struct {
	store *persistence.Store
}

func NewSearcher(store *persistence.Store) *Searcher {
	return &Searcher{store: store}
}

func (s *Searcher) Search(filter domain.MaterialFilter) ([]SearchResult, error) {
	filter = filter.Normalize()
	materials, err := s.store.ListMaterials()
	if err != nil {
		return nil, err
	}
	results := make([]SearchResult, 0)
	for _, m := range materials {
		if !filter.Accepts(m) {
			continue
		}
		results = append(results, s.decorate(m, filter))
	}
	sort.SliceStable(results, func(i, j int) bool {
		if results[i].Score == results[j].Score {
			return results[i].Material.Title < results[j].Material.Title
		}
		return results[i].Score > results[j].Score
	})
	return results, nil
}

func (s *Searcher) decorate(m domain.Material, filter domain.MaterialFilter) SearchResult {
	// Build an independent label slice per result. A shared scratch buffer would
	// alias every result's Labels to the same backing array, so a later Search
	// would silently overwrite labels produced by an earlier Search.
	labels := make([]string, 0, 8)
	labels = append(labels, string(m.Kind), m.Subject, m.Campus, m.VersionDate)
	if filter.Query != "" {
		labels = append(labels, "text-match")
	}
	if filter.Tag != "" {
		labels = append(labels, "tag:"+filter.Tag)
	}
	return SearchResult{Material: m, Labels: labels, Score: score(m, filter)}
}

func score(m domain.Material, filter domain.MaterialFilter) int {
	score := 0
	if filter.Subject != "" && m.MatchesSubject(filter.Subject) {
		score += 5
	}
	if filter.Campus != "" && m.MatchesCampus(filter.Campus) {
		score += 4
	}
	if filter.VersionDate != "" && m.VersionDate == filter.VersionDate {
		score += 3
	}
	if filter.Kind != "" && m.Kind == filter.Kind {
		score += 2
	}
	if filter.Query != "" && strings.Contains(m.SearchText(), filter.Query) {
		score += 1
	}
	return score
}

func GroupByKind(results []SearchResult) map[domain.MaterialKind][]SearchResult {
	grouped := make(map[domain.MaterialKind][]SearchResult)
	for _, result := range results {
		grouped[result.Material.Kind] = append(grouped[result.Material.Kind], result)
	}
	return grouped
}

func Summarize(results []SearchResult) string {
	parts := make([]string, 0, len(results))
	for _, result := range results {
		parts = append(parts, result.Material.ID+":"+result.Material.Title)
	}
	return strings.Join(parts, " | ")
}
