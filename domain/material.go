package domain

import (
	"fmt"
	"strings"
)

type MaterialKind string

const (
	CourseMaterial  MaterialKind = "course"
	LightSimulation MaterialKind = "light"
	RouteMap        MaterialKind = "route"
	SafetyVideo     MaterialKind = "safety"
)

type MaterialStatus string

const (
	StatusDraft     MaterialStatus = "draft"
	StatusReview    MaterialStatus = "review"
	StatusPublished MaterialStatus = "published"
	StatusArchived  MaterialStatus = "archived"
)

type Material struct {
	ID              string         `json:"id"`
	Title           string         `json:"title"`
	Kind            MaterialKind   `json:"kind"`
	Subject         string         `json:"subject"`
	Campus          string         `json:"campus"`
	VersionDate     string         `json:"version_date"`
	URI             string         `json:"uri"`
	Description     string         `json:"description"`
	Status          MaterialStatus `json:"status"`
	Tags            []string       `json:"tags"`
	DurationSeconds int            `json:"duration_seconds"`
	RouteNumber     string         `json:"route_number"`
	CreatedBy       string         `json:"created_by"`
	UpdatedBy       string         `json:"updated_by"`
}

func (m Material) Validate() error {
	if strings.TrimSpace(m.ID) == "" {
		return fmt.Errorf("material id is required")
	}
	if strings.TrimSpace(m.Title) == "" {
		return fmt.Errorf("material title is required")
	}
	if m.Kind != CourseMaterial && m.Kind != LightSimulation && m.Kind != RouteMap && m.Kind != SafetyVideo {
		return fmt.Errorf("unsupported material kind %q", m.Kind)
	}
	if strings.TrimSpace(m.Subject) == "" {
		return fmt.Errorf("subject is required")
	}
	if strings.TrimSpace(m.Campus) == "" {
		return fmt.Errorf("campus is required")
	}
	if strings.TrimSpace(m.VersionDate) == "" {
		return fmt.Errorf("version date is required")
	}
	if strings.TrimSpace(m.URI) == "" {
		return fmt.Errorf("uri is required")
	}
	if m.DurationSeconds < 0 {
		return fmt.Errorf("duration must be non-negative")
	}
	if m.Status == "" {
		return fmt.Errorf("status is required")
	}
	return nil
}

func (m Material) IsVisible() bool { return m.Status == StatusPublished || m.Status == StatusReview }

func (m Material) MatchesSubject(subject string) bool {
	return strings.EqualFold(m.Subject, strings.TrimSpace(subject))
}

func (m Material) MatchesCampus(campus string) bool {
	return strings.EqualFold(m.Campus, strings.TrimSpace(campus))
}

func (m Material) HasTag(tag string) bool {
	for _, candidate := range m.Tags {
		if strings.EqualFold(candidate, strings.TrimSpace(tag)) {
			return true
		}
	}
	return false
}

func (m Material) SearchText() string {
	return strings.ToLower(strings.Join([]string{m.Title, string(m.Kind), m.Subject, m.Campus, m.Description, strings.Join(m.Tags, " "), m.RouteNumber}, " "))
}

func NewMaterial(id string, title string, kind MaterialKind, subject string, campus string, versionDate string, uri string, creator string) Material {
	return Material{ID: id, Title: title, Kind: kind, Subject: subject, Campus: campus, VersionDate: versionDate, URI: uri, Status: StatusDraft, CreatedBy: creator, UpdatedBy: creator, Tags: []string{}}
}
