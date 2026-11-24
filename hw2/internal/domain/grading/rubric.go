package grading

import (
	"fmt"
	"slices"
	"strings"
)

type Rubric struct {
	criteria []Criterion
}

func NewRubric() Rubric {
	return Rubric{}
}

func (r Rubric) Criteria() []Criterion {
	return slices.Clone(r.criteria)
}

func (r Rubric) FindCriterion(id CriterionID) (Criterion, error) {
	for _, criterion := range r.criteria {
		if criterion.ID() == id {
			return criterion, nil
		}
	}
	return Criterion{}, fmt.Errorf("criterion not found in rubric: %s", id)
}

func (r *Rubric) AddCriterion(id CriterionID, name string, maxScore, weight int) error {
	criterion, err := newCriterion(id, name, maxScore, weight)
	if err != nil {
		return err
	}
	if err := r.checkConflict(criterion); err != nil {
		return err
	}
	r.criteria = append(r.criteria, criterion)
	return nil
}

func (r Rubric) checkConflict(criterion Criterion) error {
	for _, existing := range r.criteria {
		if existing.ID() == criterion.ID() {
			return fmt.Errorf("duplicate criterion id: %s", criterion.ID())
		}
		if strings.EqualFold(existing.Name(), criterion.Name()) {
			return fmt.Errorf("duplicate criterion name: %s", criterion.Name())
		}
	}
	return nil
}
