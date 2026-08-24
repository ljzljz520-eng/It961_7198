package domain

import "strings"

type MaterialFilter struct {
	Subject         string
	Campus          string
	VersionDate     string
	Kind            MaterialKind
	Status          MaterialStatus
	Tag             string
	Query           string
	IncludeArchived bool
}

func (f MaterialFilter) Normalize() MaterialFilter {
	f.Subject = strings.TrimSpace(f.Subject)
	f.Campus = strings.TrimSpace(f.Campus)
	f.VersionDate = strings.TrimSpace(f.VersionDate)
	f.Tag = strings.TrimSpace(f.Tag)
	f.Query = strings.TrimSpace(strings.ToLower(f.Query))
	return f
}

func (f MaterialFilter) Empty() bool {
	return f.Subject == "" && f.Campus == "" && f.VersionDate == "" && f.Kind == "" && f.Status == "" && f.Tag == "" && f.Query == ""
}

func (f MaterialFilter) Accepts(m Material) bool {
	n := f.Normalize()
	if n.Subject != "" && !m.MatchesSubject(n.Subject) {
		return false
	}
	if n.Campus != "" && !m.MatchesCampus(n.Campus) {
		return false
	}
	if n.VersionDate != "" && m.VersionDate != n.VersionDate {
		return false
	}
	if n.Kind != "" && m.Kind != n.Kind {
		return false
	}
	if n.Status != "" && m.Status != n.Status {
		return false
	}
	if n.Tag != "" && !m.HasTag(n.Tag) {
		return false
	}
	if n.Query != "" && !strings.Contains(m.SearchText(), n.Query) {
		return false
	}
	if !n.IncludeArchived && m.Status == StatusArchived {
		return false
	}
	return true
}
