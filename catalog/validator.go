package catalog

import (
	"fmt"
	"strings"

	"drivingmaterials/domain"
)

type BatchIssue struct {
	Index   int
	ID      string
	Message string
}

func ValidateBatch(materials []domain.Material) []BatchIssue {
	issues := make([]BatchIssue, 0)
	seen := make(map[string]bool)
	for i, m := range materials {
		if err := m.Validate(); err != nil {
			issues = append(issues, BatchIssue{Index: i, ID: m.ID, Message: err.Error()})
			continue
		}
		if seen[m.ID] {
			issues = append(issues, BatchIssue{Index: i, ID: m.ID, Message: "duplicate material id"})
		}
		seen[m.ID] = true
		if strings.Contains(m.Title, "  ") {
			issues = append(issues, BatchIssue{Index: i, ID: m.ID, Message: "title contains repeated spaces"})
		}
	}
	return issues
}

func EnsureCourse(c domain.Course) error {
	if err := c.Validate(); err != nil {
		return err
	}
	if !c.IsAvailable() {
		return fmt.Errorf("course %s is not available", c.Code)
	}
	return nil
}
