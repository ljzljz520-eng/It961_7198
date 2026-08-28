package workflow

import (
	"sort"

	"drivingmaterials/domain"
	"drivingmaterials/persistence"
)

type Coverage struct {
	Campus    string
	Subjects  []string
	Kinds     []domain.MaterialKind
	Published int
	Draft     int
}

func BuildCoverage(store *persistence.Store, campus string) (Coverage, error) {
	materials, err := store.ListMaterials()
	if err != nil {
		return Coverage{}, err
	}
	c := Coverage{Campus: campus, Subjects: []string{}, Kinds: []domain.MaterialKind{}}
	subjects := map[string]bool{}
	kinds := map[domain.MaterialKind]bool{}
	for _, m := range materials {
		if campus != "" && m.Campus != campus {
			continue
		}
		subjects[m.Subject] = true
		kinds[m.Kind] = true
		if m.Status == domain.StatusPublished {
			c.Published++
		}
		if m.Status == domain.StatusDraft {
			c.Draft++
		}
	}
	for subject := range subjects {
		c.Subjects = append(c.Subjects, subject)
	}
	for kind := range kinds {
		c.Kinds = append(c.Kinds, kind)
	}
	sort.Strings(c.Subjects)
	sort.Slice(c.Kinds, func(i, j int) bool { return c.Kinds[i] < c.Kinds[j] })
	return c, nil
}
