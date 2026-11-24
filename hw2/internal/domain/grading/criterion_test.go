package grading

import "testing"

func TestCriterionWeight(t *testing.T) {

	badWeigts := []int{-1, 0, 101}
	for _, weight := range badWeigts {
		if _, err := newCriterion("a", "Correctness", 10, weight); err == nil {
			t.Errorf("expected weight %d to be rejected", weight)
		}
	}

	goodWeights := []int{1, 50, 100}
	for _, weight := range goodWeights {
		if _, err := newCriterion("a", "Correctness", 10, weight); err != nil {
			t.Errorf("weight %d: %v", weight, err)
		}
	}
}
