package application

import (
	"errors"
	"testing"

	"pos/internal/domain/grading"
	"pos/internal/domain/student"
)

func TestServiceCalculateFinalScore(t *testing.T) {
	state := gradeTestState(t)
	repository := testRepository{state: state}
	service := NewService(&repository)
	if err := service.SetGrade("s1", "c1", 5); err != nil {
		t.Fatal(err)
	}
	if err := service.SetGrade("s1", "c2", 10); err != nil {
		t.Fatal(err)
	}
	if err := service.SetGrade("s2", "c1", 0); err != nil {
		t.Fatal(err)
	}
	if err := service.SetGrade("s2", "c2", 0); err != nil {
		t.Fatal(err)
	}
	repository.saveCalls = 0
	score, err := service.CalculateFinalScore("s1")
	if err != nil {
		t.Fatal(err)
	}
	if score != 8 {
		t.Fatalf("expected rounded score 8, got %d", score)
	}
	if repository.saveCalls != 0 {
		t.Fatal("calculation must not save state")
	}
}

func TestServiceCalculateFinalScoreMissingGrades(t *testing.T) {
	state := gradeTestState(t)
	repository := testRepository{state: state}
	service := NewService(&repository)
	score, err := service.CalculateFinalScore("s1")
	if err != nil || score != 0 {
		t.Fatalf("expected zero without error, got %d, %v", score, err)
	}
	if err := service.SetGrade("s1", "c1", 10); err != nil {
		t.Fatal(err)
	}
	score, err = service.CalculateFinalScore("s1")
	if err != nil || score != 5 {
		t.Fatalf("expected missing criterion to count as zero, got %d, %v", score, err)
	}
}

func TestServiceCalculateFinalScoreErrors(t *testing.T) {
	state := gradeTestState(t)
	repository := testRepository{state: state}
	service := NewService(&repository)
	if _, err := service.CalculateFinalScore("missing"); err == nil {
		t.Fatal("expected unknown student error")
	}
	failure := errors.New("storage unavailable")
	repository.loadErr = failure
	if _, err := service.CalculateFinalScore("s1"); !errors.Is(err, failure) {
		t.Fatalf("expected load error, got %v", err)
	}
	repository.loadErr = nil
	repository.state.Rubric = grading.NewRubric()
	if _, err := service.CalculateFinalScore("s1"); err == nil {
		t.Fatal("expected empty rubric error")
	}
	repository.state = state
	grade, err := grading.NewGrade("s1", state.Rubric, "c1", 5)
	if err != nil {
		t.Fatal(err)
	}
	duplicates := []grading.Grade{grade, grade}
	repository.state.Grades = duplicates
	if _, err := service.CalculateFinalScore(student.StudentID("s1")); err == nil {
		t.Fatal("expected domain validation error for duplicate grades")
	}
	if repository.saveCalls != 0 {
		t.Fatal("calculation must not save state")
	}
}

func TestServiceCurrentScoreAndExport(t *testing.T) {
	state := gradeTestState(t)
	repository := testRepository{state: state}
	service := NewService(&repository)
	if _, err := service.CalculateCurrentScore("s1"); err == nil {
		t.Fatal("expected no grades error")
	}
	if _, err := service.CalculateCurrentScore("missing"); err == nil {
		t.Fatal("expected unknown student error")
	}
	if err := service.SetGrade("s1", "c1", 8); err != nil {
		t.Fatal(err)
	}
	repository.saveCalls = 0
	current, err := service.CalculateCurrentScore("s1")
	if err != nil || current != 8 {
		t.Fatalf("current = %d, %v", current, err)
	}
	results, err := service.ListResults()
	if err != nil || len(results) != 2 || results[0].Score != 4 {
		t.Fatalf("export results = %v, %v", results, err)
	}
	if repository.saveCalls != 0 {
		t.Fatal("calculation must not save state")
	}
	failure := errors.New("storage unavailable")
	repository.loadErr = failure
	if _, err := service.CalculateCurrentScore("s1"); !errors.Is(err, failure) {
		t.Fatalf("expected load error, got %v", err)
	}
}
