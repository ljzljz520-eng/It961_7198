package catalog

import (
	"sort"
	"strings"

	"drivingmaterials/domain"
)

type Snapshot struct {
	Campus      string
	VersionDate string
	Counts      map[domain.MaterialKind]int
	Titles      []string
}

func BuildSnapshot(materials []domain.Material, campus string, version string) Snapshot {
	s := Snapshot{Campus: campus, VersionDate: version, Counts: map[domain.MaterialKind]int{}, Titles: []string{}}
	for _, m := range materials {
		if campus != "" && !strings.EqualFold(m.Campus, campus) {
			continue
		}
		if version != "" && m.VersionDate != version {
			continue
		}
		s.Counts[m.Kind]++
		s.Titles = append(s.Titles, m.Title)
	}
	sort.Strings(s.Titles)
	return s
}

func (s Snapshot) Total() int {
	total := 0
	for _, n := range s.Counts {
		total += n
	}
	return total
}

func (s Snapshot) HasKind(kind domain.MaterialKind) bool { return s.Counts[kind] > 0 }
