package domain

import "fmt"

type Course struct {
	Code       string `json:"code"`
	Name       string `json:"name"`
	Subject    string `json:"subject"`
	Campus     string `json:"campus"`
	Instructor string `json:"instructor"`
	Active     bool   `json:"active"`
}

func (c Course) Validate() error {
	if c.Code == "" {
		return fmt.Errorf("course code is required")
	}
	if c.Name == "" {
		return fmt.Errorf("course name is required")
	}
	if c.Subject == "" {
		return fmt.Errorf("course subject is required")
	}
	if c.Campus == "" {
		return fmt.Errorf("course campus is required")
	}
	return nil
}

func (c Course) Key() string { return c.Campus + ":" + c.Code }

func (c Course) IsAvailable() bool { return c.Active && c.Instructor != "" }
