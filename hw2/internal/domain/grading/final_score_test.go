package grading

import (
	"testing"

	"pos/internal/domain/student"
)

func TestCalculateFinalScore(t *testing.T) {
	first, err := newCriterion("a", "Correctness", 10, 3)
	if err != nil {
		t.Fatal(err)
	}
	second, err := newCriterion("b", "Style", 20, 2)
	if err != nil {
		t.Fatal(err)
	}
	criteria := []Criterion{first, second}
	rubric, err := makeTestRubric(criteria)
	if err != nil {
		t.Fatal(err)
	}
	makeGrade := func(id student.StudentID, criterionID CriterionID, score int) Grade {
		t.Helper()
		grade, err := NewGrade(id, rubric, criterionID, score)
		if err != nil {
			t.Fatal(err)
		}
		return grade
	}
	fullFirst := makeGrade("s", "a", 10)
	fullSecond := makeGrade("s", "b", 20)
	partialFirst := makeGrade("s", "a", 8)
	partialSecond := makeGrade("s", "b", 15)
	fullGrades := []Grade{fullFirst, fullSecond}
	partialGrades := []Grade{partialFirst, partialSecond}
	missingGrades := []Grade{fullFirst}
	duplicateGrades := []Grade{fullFirst, fullFirst}
	unknown := Grade{studentID: "s", criterionID: "unknown", score: 1}
	unknownGrades := []Grade{unknown}
	emptyRubric := Rubric{}
	cases := []struct {
		name      string
		rubric    Rubric
		grades    []Grade
		want      int
		wantError bool
	}{
		{"maximum with arbitrary weights", rubric, fullGrades, 10, false},
		{"different maxima", rubric, partialGrades, 8, false},
		{"missing grade retains weight", rubric, missingGrades, 6, false},
		{"no grades", rubric, nil, 0, false},
		{"empty rubric", emptyRubric, nil, 0, true},
		{"duplicate", rubric, duplicateGrades, 0, true},
		{"unknown criterion", rubric, unknownGrades, 0, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := CalculateFinalScore(tc.rubric, tc.grades)
			if (err != nil) != tc.wantError {
				t.Fatalf("error = %v, wantError = %v", err, tc.wantError)
			}
			if err == nil && got != tc.want {
				t.Fatalf("score = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestCalculateFinalScoreRounding(t *testing.T) {
	criterion, err := newCriterion("a", "Correctness", 100, 1)
	if err != nil {
		t.Fatal(err)
	}
	criteria := []Criterion{criterion}
	rubric, err := makeTestRubric(criteria)
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct{ score, want int }{{74, 7}, {75, 8}, {76, 8}}
	for _, tc := range cases {
		grade, err := NewGrade("s", rubric, "a", tc.score)
		if err != nil {
			t.Fatal(err)
		}
		grades := []Grade{grade}
		got, err := CalculateFinalScore(rubric, grades)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.want {
			t.Errorf("score %d: got %d, want %d", tc.score, got, tc.want)
		}
	}
}
