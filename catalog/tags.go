package catalog

import (
	"sort"
	"strings"

	"drivingmaterials/domain"
)

func NormalizeTags(tags []string) []string {
	seen := map[string]bool{}
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		normalized := strings.ToLower(strings.TrimSpace(tag))
		if normalized == "" || seen[normalized] {
			continue
		}
		seen[normalized] = true
		result = append(result, normalized)
	}
	sort.Strings(result)
	return result
}

func ApplyTagPolicy(m domain.Material) domain.Material {
	m.Tags = NormalizeTags(m.Tags)
	if m.Kind == domain.SafetyVideo && !m.HasTag("safety") {
		m.Tags = append(m.Tags, "safety")
	}
	return m
}
