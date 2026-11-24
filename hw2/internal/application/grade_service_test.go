package application

import (
	"errors"
	"slices"
	"testing"

	"pos/internal/domain/grading"
	"pos/internal/domain/student"
)

func gradeTestState(t *testing.T) State {
	t.Helper()
	first, err := student.NewStudent("s1", "Alex", "G1")
	if err != nil {
		t.Fatal(err)
	}
	second, err := student.NewStudent("s2", "Sam", "G1")
	if err != nil {
		t.Fatal(err)
	}
	students := []student.Student{first, second}
	rubric := grading.NewRubric()
	if err := rubric.AddCriterion("c1", "Correctness", 10, 50); err != nil {
		t.Fatal(err)
	}
	if err := rubric.AddCriterion("c2", "Style", 10, 50); err != nil {
		t.Fatal(err)
	}
	state := State{Students: students, Rubric: rubric}
	return state
}

func TestServiceSetGradeAddsAndReplaces(t *testing.T) {
	state := gradeTestState(t)
	repository := testRepository{state: state}
	service := NewService(&repository)
	cases := []struct {
		studentID   student.StudentID
		criterionID grading.CriterionID
		score       int
	}{
		{"s1", "c1", 3}, {"s2", "c1", 7}, {"s1", "c2", 8}, {"s1", "c1", 10},
	}
	for _, tc := range cases {
		if err := service.SetGrade(tc.studentID, tc.criterionID, tc.score); err != nil {
			t.Fatal(err)
		}
	}
	grades := repository.state.Grades
	if len(grades) != 3 || repository.saveCalls != 4 {
		t.Fatal("expected three grades and four saves")
	}
	for _, tc := range cases[1:] {
		found := false
		for _, grade := range grades {
			if grade.StudentID() == tc.studentID && grade.CriterionID() == tc.criterionID {
				found = true
				if grade.Score() != tc.score {
					t.Fatal("unexpected grade score")
				}
			}
		}
		if !found {
			t.Fatal("grade was not saved")
		}
	}
}

func TestServiceSetGradeRejectsInvalidInput(t *testing.T) {
	cases := []struct {
		name        string
		studentID   student.StudentID
		criterionID grading.CriterionID
		score       int
	}{
		{"unknown student", "missing", "c1", 5},
		{"unknown criterion", "s1", "missing", 5},
		{"negative score", "s1", "c1", -1},
		{"above maximum", "s1", "c1", 11},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			state := gradeTestState(t)
			repository := testRepository{state: state}
			service := NewService(&repository)
			if err := service.SetGrade(tc.studentID, tc.criterionID, tc.score); err == nil {
				t.Fatal("expected error")
			}
			if repository.saveCalls != 0 || len(repository.state.Grades) != 0 {
				t.Fatal("invalid grade must not be saved")
			}
		})
	}
}

func TestServiceSetGradeRepositoryErrors(t *testing.T) {
	failure := errors.New("storage unavailable")
	state := gradeTestState(t)
	repository := testRepository{state: state, loadErr: failure}
	service := NewService(&repository)
	if err := service.SetGrade("s1", "c1", 0); !errors.Is(err, failure) {
		t.Fatalf("expected load error, got %v", err)
	}
	if repository.saveCalls != 0 {
		t.Fatal("must not save after load failure")
	}
	repository.loadErr = nil
	if err := service.SetGrade("s1", "c1", 0); err != nil {
		t.Fatal(err)
	}
	original := slices.Clone(repository.state.Grades)
	repository.saveErr = failure
	if err := service.SetGrade("s1", "c1", 10); !errors.Is(err, failure) {
		t.Fatalf("expected save error, got %v", err)
	}
	if !slices.Equal(repository.state.Grades, original) {
		t.Fatal("failed save changed stored grade")
	}
	if repository.saveCalls != 2 {
		t.Fatal("expected two save attempts")
	}
}

func TestServiceListGrades(t *testing.T) {
	state := gradeTestState(t)
	first, err := grading.NewGrade("s1", state.Rubric, "c1", 5)
	if err != nil {
		t.Fatal(err)
	}
	other, err := grading.NewGrade("s2", state.Rubric, "c1", 7)
	if err != nil {
		t.Fatal(err)
	}
	second, err := grading.NewGrade("s1", state.Rubric, "c2", 8)
	if err != nil {
		t.Fatal(err)
	}
	grades := []grading.Grade{first, other, second}
	state.Grades = grades
	repository := testRepository{state: state}
	service := NewService(&repository)
	got, err := service.ListGrades("s1")
	if err != nil {
		t.Fatal(err)
	}
	want := []grading.Grade{first, second}
	if !slices.Equal(got, want) {
		t.Fatalf("unexpected grades: %v", got)
	}
	got[0] = second
	if !slices.Equal(repository.state.Grades, grades) {
		t.Fatal("returned slice shares stored data")
	}
	if repository.saveCalls != 0 {
		t.Fatal("listing must not save state")
	}
}

func TestServiceListGradesEmptyAndErrors(t *testing.T) {
	state := gradeTestState(t)
	repository := testRepository{state: state}
	service := NewService(&repository)
	grades, err := service.ListGrades("s1")
	if err != nil || len(grades) != 0 {
		t.Fatalf("expected empty list without error, got %v, %v", grades, err)
	}
	if _, err := service.ListGrades("missing"); err == nil {
		t.Fatal("expected unknown student error")
	}
	failure := errors.New("storage unavailable")
	repository.loadErr = failure
	if _, err := service.ListGrades("s1"); !errors.Is(err, failure) {
		t.Fatalf("expected load error, got %v", err)
	}
	if repository.saveCalls != 0 {
		t.Fatal("listing must not save state")
	}
}

func TestServiceRemoveGrade(t *testing.T) {
	state := gradeTestState(t)
	repository := testRepository{state: state}
	service := NewService(&repository)
	if err := service.SetGrade("s1", "c1", 5); err != nil {
		t.Fatal(err)
	}
	if err := service.SetGrade("s2", "c1", 7); err != nil {
		t.Fatal(err)
	}
	if err := service.SetGrade("s1", "c2", 8); err != nil {
		t.Fatal(err)
	}
	original := slices.Clone(repository.state.Grades)
	want := []grading.Grade{original[1], original[2]}
	if err := service.RemoveGrade("s1", "c1"); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(repository.state.Grades, want) {
		t.Fatal("wrong grades after removal")
	}
	if repository.saveCalls != 4 {
		t.Fatal("removal was not saved")
	}
	if err := service.RemoveGrade("s1", "c1"); err == nil {
		t.Fatal("expected missing grade error")
	}
	if repository.saveCalls != 4 {
		t.Fatal("missing grade must not trigger save")
	}
}

func TestServiceRemoveGradeMissingAndStorageErrors(t *testing.T) {
	state := gradeTestState(t)
	repository := testRepository{state: state}
	service := NewService(&repository)
	if err := service.RemoveGrade("missing", "c1"); err == nil {
		t.Fatal("expected unknown student error")
	}
	if err := service.RemoveGrade("s1", "c1"); err == nil {
		t.Fatal("expected missing grade error")
	}
	if repository.saveCalls != 0 {
		t.Fatal("invalid removal must not save")
	}
	failure := errors.New("storage unavailable")
	repository.loadErr = failure
	if err := service.RemoveGrade("s1", "c1"); !errors.Is(err, failure) {
		t.Fatalf("expected load error, got %v", err)
	}
	if repository.saveCalls != 0 {
		t.Fatal("load failure must not trigger save")
	}
	repository.loadErr = nil
	if err := service.SetGrade("s1", "c1", 5); err != nil {
		t.Fatal(err)
	}
	original := slices.Clone(repository.state.Grades)
	repository.saveErr = failure
	if err := service.RemoveGrade("s1", "c1"); !errors.Is(err, failure) {
		t.Fatalf("expected save error, got %v", err)
	}
	if !slices.Equal(repository.state.Grades, original) {
		t.Fatal("failed save changed stored grades")
	}
	if repository.saveCalls != 2 {
		t.Fatal("expected removal save attempt")
	}
	repository.saveErr = nil
	if err := service.RemoveGrade("s1", "c1"); err != nil {
		t.Fatal(err)
	}
	if len(repository.state.Grades) != 0 {
		t.Fatal("last grade was not removed")
	}
}
