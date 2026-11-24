package grading

import (
	"fmt"
	"math"
)

// ScoreStrategy defines a replaceable rule for calculating a weighted score.
type ScoreStrategy interface {
	Calculate(rubric Rubric, grades []Grade) (int, error)
}

type FinalScoreStrategy struct{}
type CurrentScoreStrategy struct{}

func CalculateFinalScore(rubric Rubric, grades []Grade) (int, error) {
	return FinalScoreStrategy{}.Calculate(rubric, grades)
}

func (FinalScoreStrategy) Calculate(rubric Rubric, grades []Grade) (int, error) {
	return calculateWeightedScore(rubric, grades, false)
}

func (CurrentScoreStrategy) Calculate(rubric Rubric, grades []Grade) (int, error) {
	return calculateWeightedScore(rubric, grades, true)
}

func calculateWeightedScore(rubric Rubric, grades []Grade, onlyGraded bool) (int, error) {
	criteria := rubric.Criteria()
	if len(criteria) == 0 {
		return 0, fmt.Errorf("rubric must not be empty")
	}

	scores := make(map[CriterionID]int, len(grades))
	for _, grade := range grades {

		_, err := rubric.FindCriterion(grade.CriterionID())
		if err != nil {
			return 0, err
		}

		if _, exists := scores[grade.CriterionID()]; exists {
			return 0, fmt.Errorf("duplicate grade for criterion: %s", grade.CriterionID())
		}

		scores[grade.CriterionID()] = grade.Score()
	}

	var weightedSum, totalWeight float64
	for _, criterion := range criteria {
		if criterion.MaxScore() <= 0 || criterion.Weight() <= 0 {
			return 0, fmt.Errorf("criterion maximum and weight must be positive: %s", criterion.ID())
		}
		score, graded := scores[criterion.ID()]
		if onlyGraded && !graded {
			continue
		}
		ratio := float64(score) / float64(criterion.MaxScore())
		weightedSum += ratio * float64(criterion.Weight())
		totalWeight += float64(criterion.Weight())
	}
	if totalWeight == 0 {
		return 0, fmt.Errorf("нет выставленных оценок")
	}
	return int(math.Round(10 * weightedSum / totalWeight)), nil
}
