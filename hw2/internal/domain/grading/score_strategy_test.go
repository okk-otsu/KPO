package grading

import "testing"

func TestScoreStrategies(t *testing.T) {
	rubric := NewRubric()
	if err := rubric.AddCriterion("a", "Practice", 20, 40); err != nil {
		t.Fatal(err)
	}
	if err := rubric.AddCriterion("b", "Exam", 10, 60); err != nil {
		t.Fatal(err)
	}
	practice, err := NewGrade("s", rubric, "a", 15)
	if err != nil {
		t.Fatal(err)
	}
	zero, err := NewGrade("s", rubric, "b", 0)
	if err != nil {
		t.Fatal(err)
	}
	exam, err := NewGrade("s", rubric, "b", 8)
	if err != nil {
		t.Fatal(err)
	}
	partial := []Grade{practice}
	withZero := []Grade{practice, zero}
	complete := []Grade{practice, exam}
	duplicate := []Grade{practice, practice}
	unknown := Grade{studentID: "s", criterionID: "missing", score: 1}
	unknownGrades := []Grade{unknown}
	emptyRubric := NewRubric()
	cases := []struct {
		name                     string
		rubric                   Rubric
		grades                   []Grade
		final, current           int
		finalError, currentError bool
	}{
		{"partial and rounding", rubric, partial, 3, 8, false, false},
		{"explicit zero counts", rubric, withZero, 3, 3, false, false},
		{"complete", rubric, complete, 8, 8, false, false},
		{"no grades", rubric, nil, 0, 0, false, true},
		{"empty rubric", emptyRubric, nil, 0, 0, true, true},
		{"duplicate", rubric, duplicate, 0, 0, true, true},
		{"unknown criterion", rubric, unknownGrades, 0, 0, true, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			strategies := []struct {
				strategy  ScoreStrategy
				want      int
				wantError bool
			}{
				{FinalScoreStrategy{}, tc.final, tc.finalError},
				{CurrentScoreStrategy{}, tc.current, tc.currentError},
			}
			for _, strategy := range strategies {
				got, err := strategy.strategy.Calculate(tc.rubric, tc.grades)
				if (err != nil) != strategy.wantError || got != strategy.want {
					t.Fatalf("%T: got %d, %v; want %d, error %v", strategy.strategy, got, err, strategy.want, strategy.wantError)
				}
			}
		})
	}
}
