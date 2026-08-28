package presenter

import (
	"fmt"
	"io"

	"drivingmaterials/domain"
)

func WriteMaterial(w io.Writer, m domain.Material) error {
	_, err := fmt.Fprintf(w, "%s: %s [%s] %s/%s %s\n", m.ID, m.Title, m.Status, m.Campus, m.Subject, m.VersionDate)
	return err
}

func StatusLabel(status domain.MaterialStatus) string {
	switch status {
	case domain.StatusDraft:
		return "DRAFT"
	case domain.StatusReview:
		return "IN REVIEW"
	case domain.StatusPublished:
		return "PUBLISHED"
	case domain.StatusArchived:
		return "ARCHIVED"
	default:
		return "UNKNOWN"
	}
}
