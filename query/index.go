package query

import (
	"sort"
	"strings"

	"drivingmaterials/domain"
)

type Index struct {
	byToken  map[string]map[string]bool
	byCampus map[string]map[string]bool
}

func NewIndex() *Index {
	return &Index{byToken: map[string]map[string]bool{}, byCampus: map[string]map[string]bool{}}
}

func (i *Index) Add(m domain.Material) {
	for _, token := range strings.Fields(m.SearchText()) {
		if i.byToken[token] == nil {
			i.byToken[token] = map[string]bool{}
		}
		i.byToken[token][m.ID] = true
	}
	campus := strings.ToLower(m.Campus)
	if i.byCampus[campus] == nil {
		i.byCampus[campus] = map[string]bool{}
	}
	i.byCampus[campus][m.ID] = true
}

func (i *Index) Remove(m domain.Material) {
	for _, token := range strings.Fields(m.SearchText()) {
		if ids := i.byToken[token]; ids != nil {
			delete(ids, m.ID)
			if len(ids) == 0 {
				delete(i.byToken, token)
			}
		}
	}
	campus := strings.ToLower(m.Campus)
	if ids := i.byCampus[campus]; ids != nil {
		delete(ids, m.ID)
	}
}

func (i *Index) IDsForQuery(query string) []string {
	ids := map[string]bool{}
	for _, token := range strings.Fields(strings.ToLower(query)) {
		for id := range i.byToken[token] {
			ids[id] = true
		}
	}
	result := make([]string, 0, len(ids))
	for id := range ids {
		result = append(result, id)
	}
	sort.Strings(result)
	return result
}

func (i *Index) IDsForCampus(campus string) []string {
	result := make([]string, 0)
	for id := range i.byCampus[strings.ToLower(strings.TrimSpace(campus))] {
		result = append(result, id)
	}
	sort.Strings(result)
	return result
}
