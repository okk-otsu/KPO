package grading

import (
	"slices"
	"testing"
)

func TestRubricDuplicateNames(t *testing.T) {
	first, _ := newCriterion("a", "Correctness", 10, 50)
	duplicate, _ := newCriterion("b", " correctness ", 10, 50)
	criteria := []Criterion{first, duplicate}
	if _, err := makeTestRubric(criteria); err == nil {
		t.Fatal("expected duplicate name to be rejected")
	}
}

func TestAddCriterion(t *testing.T) {
	first, _ := newCriterion("a", "Correctness", 10, 50)
	unchangedCriteria := []Criterion{first}
	rubric, err := makeTestRubric(unchangedCriteria)
	if err != nil {
		t.Fatal(err)
	}

	duplicateName, _ := newCriterion("b", " correctness ", 20, 25)
	duplicates := []Criterion{duplicateName}
	for _, criterion := range duplicates {
		if err := rubric.AddCriterion(criterion.ID(), criterion.Name(), criterion.MaxScore(), criterion.Weight()); err == nil {
			t.Fatal("expected duplicate to be rejected")
		}

		criteria := rubric.Criteria()
		if !slices.Equal(criteria, unchangedCriteria) {
			t.Fatal("failed addition changed rubric")
		}
	}

	added, _ := newCriterion("b", "Style", 20, 25)
	if err := rubric.AddCriterion(added.ID(), added.Name(), added.MaxScore(), added.Weight()); err != nil {
		t.Fatal(err)
	}

	criteria := rubric.Criteria()
	want := []Criterion{first, added}
	if !slices.Equal(criteria, want) {
		t.Fatal("added criterion was not preserved")
	}
}

func TestRubricRejectsDuplicateID(t *testing.T) {
	rubric := NewRubric()
	if err := rubric.AddCriterion("a", "Correctness", 10, 50); err != nil {
		t.Fatal(err)
	}
	if err := rubric.AddCriterion("a", "Style", 10, 50); err == nil {
		t.Fatal("expected duplicate ID to be rejected")
	}
}

func makeTestRubric(criteria []Criterion) (Rubric, error) {
	rubric := NewRubric()
	for _, criterion := range criteria {
		if err := rubric.AddCriterion(criterion.ID(), criterion.Name(), criterion.MaxScore(), criterion.Weight()); err != nil {
			return Rubric{}, err
		}
	}
	return rubric, nil
}
