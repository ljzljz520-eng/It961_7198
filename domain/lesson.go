package domain

import (
	"fmt"
	"sort"
	"strings"
)

type LessonStep struct {
	Number     int    `json:"number"`
	Title      string `json:"title"`
	MaterialID string `json:"material_id"`
	Minutes    int    `json:"minutes"`
	Required   bool   `json:"required"`
}

type LessonPlan struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Subject     string       `json:"subject"`
	Campus      string       `json:"campus"`
	VersionDate string       `json:"version_date"`
	Steps       []LessonStep `json:"steps"`
	Published   bool         `json:"published"`
}

type LessonCompletion struct {
	PlanID    string `json:"plan_id"`
	Coach     string `json:"coach"`
	Completed []int  `json:"completed"`
	Notes     string `json:"notes"`
}

func (s LessonStep) Validate() error {
	if s.Number < 1 {
		return fmt.Errorf("step number must be positive")
	}
	if strings.TrimSpace(s.Title) == "" {
		return fmt.Errorf("step title is required")
	}
	if strings.TrimSpace(s.MaterialID) == "" {
		return fmt.Errorf("step material is required")
	}
	if s.Minutes <= 0 {
		return fmt.Errorf("step duration must be positive")
	}
	return nil
}

func (p LessonPlan) Validate() error {
	if strings.TrimSpace(p.ID) == "" {
		return fmt.Errorf("plan id is required")
	}
	if strings.TrimSpace(p.Name) == "" {
		return fmt.Errorf("plan name is required")
	}
	if strings.TrimSpace(p.Subject) == "" {
		return fmt.Errorf("plan subject is required")
	}
	if strings.TrimSpace(p.Campus) == "" {
		return fmt.Errorf("plan campus is required")
	}
	if strings.TrimSpace(p.VersionDate) == "" {
		return fmt.Errorf("plan version is required")
	}
	if len(p.Steps) < 1 {
		return fmt.Errorf("plan needs at least one step")
	}
	seen := map[int]bool{}
	for _, step := range p.Steps {
		if err := step.Validate(); err != nil {
			return err
		}
		if seen[step.Number] {
			return fmt.Errorf("duplicate step number %d", step.Number)
		}
		seen[step.Number] = true
	}
	return nil
}

func (p LessonPlan) OrderedSteps() []LessonStep {
	steps := append([]LessonStep(nil), p.Steps...)
	sort.Slice(steps, func(i, j int) bool { return steps[i].Number < steps[j].Number })
	return steps
}

func (p LessonPlan) TotalMinutes() int {
	total := 0
	for _, step := range p.Steps {
		total += step.Minutes
	}
	return total
}

func (p LessonPlan) RequiredCount() int {
	count := 0
	for _, step := range p.Steps {
		if step.Required {
			count++
		}
	}
	return count
}

func (p LessonPlan) Step(number int) (LessonStep, bool) {
	for _, step := range p.Steps {
		if step.Number == number {
			return step, true
		}
	}
	return LessonStep{}, false
}

func (c LessonCompletion) Validate(plan LessonPlan) error {
	if c.PlanID != plan.ID {
		return fmt.Errorf("completion plan mismatch")
	}
	if strings.TrimSpace(c.Coach) == "" {
		return fmt.Errorf("coach is required")
	}
	if len(c.Completed) == 0 {
		return fmt.Errorf("completion needs steps")
	}
	for _, number := range c.Completed {
		if _, ok := plan.Step(number); !ok {
			return fmt.Errorf("unknown completed step %d", number)
		}
	}
	return nil
}

func (c LessonCompletion) IsComplete(plan LessonPlan) bool {
	required := map[int]bool{}
	for _, step := range plan.Steps {
		if step.Required {
			required[step.Number] = true
		}
	}
	for _, number := range c.Completed {
		delete(required, number)
	}
	return len(required) == 0
}

func (c LessonCompletion) CompletionPercent(plan LessonPlan) int {
	if len(plan.Steps) == 0 {
		return 0
	}
	done := map[int]bool{}
	for _, number := range c.Completed {
		done[number] = true
	}
	return len(done) * 100 / len(plan.Steps)
}

func NewLessonPlan(id, name, subject, campus, version string) LessonPlan {
	return LessonPlan{ID: id, Name: name, Subject: subject, Campus: campus, VersionDate: version, Steps: []LessonStep{}}
}

func (p *LessonPlan) AddStep(title, materialID string, minutes int, required bool) error {
	number := len(p.Steps) + 1
	step := LessonStep{Number: number, Title: title, MaterialID: materialID, Minutes: minutes, Required: required}
	if err := step.Validate(); err != nil {
		return err
	}
	p.Steps = append(p.Steps, step)
	return nil
}

func (p *LessonPlan) Publish() error {
	if err := p.Validate(); err != nil {
		return err
	}
	p.Published = true
	return nil
}
