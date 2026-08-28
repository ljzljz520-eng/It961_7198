package catalog

import (
	"fmt"
	"sort"

	"drivingmaterials/domain"
	"drivingmaterials/persistence"
)

type PlanService struct{ store *persistence.Store }

func NewPlanService(store *persistence.Store) *PlanService { return &PlanService{store: store} }

func (p *PlanService) Create(id, name, subject, campus, version string) (domain.LessonPlan, error) {
	plan := domain.NewLessonPlan(id, name, subject, campus, version)
	if err := plan.Validate(); err != nil && id != "" {
		return plan, err
	}
	return plan, nil
}

func (p *PlanService) AddMaterialStep(planID, title, materialID string, minutes int, required bool) (domain.LessonPlan, error) {
	plan, err := p.store.GetLessonPlan(planID)
	if err != nil {
		return plan, err
	}
	if _, err := p.store.GetMaterial(materialID); err != nil {
		return plan, fmt.Errorf("step material: %w", err)
	}
	if err := plan.AddStep(title, materialID, minutes, required); err != nil {
		return plan, err
	}
	if err := p.store.PutLessonPlan(plan); err != nil {
		return plan, err
	}
	return plan, nil
}

func (p *PlanService) Publish(planID string) (domain.LessonPlan, error) {
	plan, err := p.store.GetLessonPlan(planID)
	if err != nil {
		return plan, err
	}
	if err := plan.Publish(); err != nil {
		return plan, err
	}
	if err := p.store.PutLessonPlan(plan); err != nil {
		return plan, err
	}
	return plan, nil
}

func (p *PlanService) Find(campus, subject string) ([]domain.LessonPlan, error) {
	plans, err := p.store.ListLessonPlans(campus, subject)
	if err != nil {
		return nil, err
	}
	sort.Slice(plans, func(i, j int) bool {
		if plans[i].VersionDate == plans[j].VersionDate {
			return plans[i].Name < plans[j].Name
		}
		return plans[i].VersionDate > plans[j].VersionDate
	})
	return plans, nil
}

func (p *PlanService) RecordCompletion(completion domain.LessonCompletion) error {
	if completion.PlanID == "" || completion.Coach == "" {
		return fmt.Errorf("completion plan and coach are required")
	}
	plan, err := p.store.GetLessonPlan(completion.PlanID)
	if err != nil {
		return err
	}
	if err := completion.Validate(plan); err != nil {
		return err
	}
	return p.store.PutCompletion(completion)
}

func (p *PlanService) Completion(planID string) ([]domain.LessonCompletion, error) {
	return p.store.ListCompletions(planID)
}

func (p *PlanService) CompletionStatus(planID string) (int, error) {
	plan, err := p.store.GetLessonPlan(planID)
	if err != nil {
		return 0, err
	}
	values, err := p.store.ListCompletions(planID)
	if err != nil {
		return 0, err
	}
	best := 0
	for _, completion := range values {
		percent := completion.CompletionPercent(plan)
		if percent > best {
			best = percent
		}
	}
	return best, nil
}
