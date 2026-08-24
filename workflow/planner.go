package workflow

import (
	"fmt"
	"sort"

	"drivingmaterials/catalog"
	"drivingmaterials/domain"
	"drivingmaterials/persistence"
	"drivingmaterials/query"
)

type Planner struct {
	plans *catalog.PlanService
	store *persistence.Store
}

func NewPlanner(store *persistence.Store) *Planner {
	return &Planner{plans: catalog.NewPlanService(store), store: store}
}

func (p *Planner) BuildForCampus(id, name, subject, campus, version string) (domain.LessonPlan, error) {
	materials, err := p.store.MaterialsForSubject(subject)
	if err != nil {
		return domain.LessonPlan{}, err
	}
	candidates := query.AvailableForCoach(materials, campus)
	if len(candidates) == 0 {
		return domain.LessonPlan{}, fmt.Errorf("no visible materials for campus %s", campus)
	}
	plan := domain.NewLessonPlan(id, name, subject, campus, version)
	for _, m := range candidates {
		required := m.Kind == domain.CourseMaterial || m.Kind == domain.LightSimulation
		minutes := 5
		if m.DurationSeconds > 0 {
			minutes = m.DurationSeconds / 60
			if minutes < 1 {
				minutes = 1
			}
		}
		if err := plan.AddStep(m.Title, m.ID, minutes, required); err != nil {
			return plan, err
		}
	}
	if err := plan.Publish(); err != nil {
		return plan, err
	}
	if err := p.store.PutLessonPlan(plan); err != nil {
		return plan, err
	}
	return plan, nil
}

func (p *Planner) Refresh(planID string) (domain.LessonPlan, error) {
	plan, err := p.store.GetLessonPlan(planID)
	if err != nil {
		return plan, err
	}
	materials, err := p.store.MaterialsForSubject(plan.Subject)
	if err != nil {
		return plan, err
	}
	visible := query.AvailableForCoach(materials, plan.Campus)
	sort.Slice(visible, func(i, j int) bool { return visible[i].Kind < visible[j].Kind })
	updated := domain.NewLessonPlan(plan.ID, plan.Name, plan.Subject, plan.Campus, plan.VersionDate)
	for _, m := range visible {
		if err := updated.AddStep(m.Title, m.ID, 5, m.Kind == domain.CourseMaterial); err != nil {
			return plan, err
		}
	}
	updated.Published = plan.Published
	if err := p.store.PutLessonPlan(updated); err != nil {
		return plan, err
	}
	return updated, nil
}

func (p *Planner) Complete(planID, coach string, steps []int, notes string) (domain.LessonCompletion, error) {
	plan, err := p.store.GetLessonPlan(planID)
	if err != nil {
		return domain.LessonCompletion{}, err
	}
	completion := domain.LessonCompletion{PlanID: planID, Coach: coach, Completed: steps, Notes: notes}
	if err := completion.Validate(plan); err != nil {
		return completion, err
	}
	if err := p.store.PutCompletion(completion); err != nil {
		return completion, err
	}
	return completion, nil
}

func (p *Planner) Progress(planID string) (int, error) {
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
		if percent := completion.CompletionPercent(plan); percent > best {
			best = percent
		}
	}
	return best, nil
}
