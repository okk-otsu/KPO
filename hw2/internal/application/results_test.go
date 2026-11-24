package application

import (
	"errors"
	"testing"
)

func TestListResultsErrorsAndEmpty(t *testing.T) {
	repository := &testRepository{}
	service := NewService(repository)
	results, err := service.ListResults()
	if err != nil || len(results) != 0 {
		t.Fatalf("expected empty results, got %v, %v", results, err)
	}
	wantErr := errors.New("load failed")
	repository.loadErr = wantErr
	if _, err := service.ListResults(); !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
	repository.loadErr = nil
	if err := service.AddStudent("Student A", "group-1"); err != nil {
		t.Fatal(err)
	}
	if _, err := service.ListResults(); err == nil {
		t.Fatal("expected error for empty rubric")
	}
}
