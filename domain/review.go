package domain

import "fmt"

type Review struct {
	MaterialID string `json:"material_id"`
	Reviewer   string `json:"reviewer"`
	Decision   string `json:"decision"`
	Notes      string `json:"notes"`
	Version    int    `json:"version"`
}

func (r Review) Validate() error {
	if r.MaterialID == "" {
		return fmt.Errorf("material id is required")
	}
	if r.Reviewer == "" {
		return fmt.Errorf("reviewer is required")
	}
	if r.Decision != "approve" && r.Decision != "reject" {
		return fmt.Errorf("decision must be approve or reject")
	}
	if r.Decision == "reject" && r.Notes == "" {
		return fmt.Errorf("rejection needs notes")
	}
	if r.Version < 1 {
		return fmt.Errorf("review version must be positive")
	}
	return nil
}
