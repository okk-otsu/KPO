package grading

import "testing"

func TestNewGrade(t *testing.T) {
	criterion, err := newCriterion("a", "Correctness", 20, 100)
	if err != nil {
		t.Fatal(err)
	}
	criteria := []Criterion{criterion}
	rubric, err := makeTestRubric(criteria)
	if err != nil {
		t.Fatal(err)
	}

	validScores := []int{0, 12, 20}
	for _, score := range validScores {
		grade, err := NewGrade("student-1", rubric, "a", score)
		if err != nil {
			t.Fatalf("score %d: %v", score, err)
		}
		if grade.StudentID() != "student-1" || grade.CriterionID() != "a" || grade.Score() != score {
			t.Fatalf("grade does not preserve input: %+v", grade)
		}
	}

	invalidScores := []int{-1, 21}
	for _, score := range invalidScores {
		if _, err := NewGrade("student-1", rubric, "a", score); err == nil {
			t.Errorf("expected score %d to be rejected", score)
		}
	}
	if _, err := NewGrade("student-1", rubric, "outside", 10); err == nil {
		t.Error("expected unknown criterion to be rejected")
	}
	if _, err := NewGrade("   ", rubric, "a", 10); err == nil {
		t.Error("expected blank student ID to be rejected")
	}
}

func TestFindCriterion(t *testing.T) {
	criterion, err := newCriterion("a", "Correctness", 10, 100)
	if err != nil {
		t.Fatal(err)
	}
	criteria := []Criterion{criterion}
	rubric, err := makeTestRubric(criteria)
	if err != nil {
		t.Fatal(err)
	}
	found, err := rubric.FindCriterion("a")
	if err != nil || found != criterion {
		t.Fatalf("expected stored criterion, got %+v, %v", found, err)
	}
	if _, err := rubric.FindCriterion("missing"); err == nil {
		t.Fatal("expected missing criterion to be rejected")
	}
}
