package grading

import (
	"fmt"
	"strings"
)

type CriterionID string

type Criterion struct {
	id       CriterionID
	name     string
	maxScore int
	weight   int
}

func newCriterion(id CriterionID, name string, maxScore, weight int) (Criterion, error) {
	if strings.TrimSpace(string(id)) == "" {
		return Criterion{}, fmt.Errorf("id must not be empty")
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return Criterion{}, fmt.Errorf("name must not be empty")
	}

	if maxScore <= 0 {
		return Criterion{}, fmt.Errorf("maxScore must be positive")
	}

	if weight <= 0 || weight > 100 {
		return Criterion{}, fmt.Errorf("weight must be in interval [1; 100]")
	}

	return Criterion{
		id:       id,
		name:     name,
		maxScore: maxScore,
		weight:   weight,
	}, nil
}

func (c Criterion) ID() CriterionID {
	return c.id
}

func (c Criterion) Name() string {
	return c.name
}

func (c Criterion) MaxScore() int {
	return c.maxScore
}

func (c Criterion) Weight() int {
	return c.weight
}
