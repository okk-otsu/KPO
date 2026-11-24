package jsonrepo

import (
	"os"
	"path/filepath"
	"pos/internal/application"
	"pos/internal/domain/grading"
	"pos/internal/domain/student"
	"reflect"
	"testing"
)

func TestRepositorySaveAndLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	repository := NewRepository(path)
	learner, err := student.NewStudent("s1", "Student A", "group-1")
	if err != nil {
		t.Fatal(err)
	}
	rubric := grading.NewRubric()
	if err := rubric.AddCriterion("c1", "Homework", 20, 30); err != nil {
		t.Fatal(err)
	}
	grade, err := grading.NewGrade(learner.ID(), rubric, "c1", 15)
	if err != nil {
		t.Fatal(err)
	}
	students := []student.Student{learner}
	grades := []grading.Grade{grade}
	want := application.State{Students: students, Rubric: rubric, Grades: grades}
	if err := repository.Save(want); err != nil {
		t.Fatal(err)
	}
	got, err := NewRepository(path).Load()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(toDTO(got), toDTO(want)) {
		t.Fatalf("loaded state = %#v, want %#v", got, want)
	}
	empty := application.State{}
	if err := repository.Save(empty); err != nil {
		t.Fatal(err)
	}
	got, err = repository.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(toDTO(got), toDTO(empty)) {
		t.Fatalf("expected overwritten empty state, got %#v", got)
	}
}

func TestRepositoryMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	got, err := NewRepository(path).Load()
	if err != nil {
		t.Fatal(err)
	}
	empty := application.State{}
	if !reflect.DeepEqual(toDTO(got), toDTO(empty)) {
		t.Fatalf("expected empty state, got %#v", got)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("Load created a file or stat failed: %v", err)
	}
}

func TestRepositoryRejectsInvalidFile(t *testing.T) {
	inputs := []string{
		"",
		"{broken",
		`{} {}`,
		`{"students":[{"id":"s1","name":"","group":"group-1"}]}`,
		`{"grades":[{"student_id":"missing","criterion_id":"c1","score":1}]}`,
	}
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "state.json")
			if err := os.WriteFile(path, []byte(input), 0600); err != nil {
				t.Fatal(err)
			}
			if _, err := NewRepository(path).Load(); err == nil {
				t.Fatal("expected error for invalid file")
			}
		})
	}
}

func TestRepositoryFileErrors(t *testing.T) {
	dir := t.TempDir()
	if _, err := NewRepository(dir).Load(); err == nil {
		t.Fatal("expected error reading a directory")
	}
	state := application.State{}
	if err := NewRepository(filepath.Join(dir, "missing", "state.json")).Save(state); err == nil {
		t.Fatal("expected error for missing parent directory")
	}
	if err := NewRepository(dir).Save(state); err == nil {
		t.Fatal("expected error replacing a directory")
	}
	entries, err := os.ReadDir(filepath.Dir(dir))
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if matched, _ := filepath.Match(".pos-state-*.json", entry.Name()); matched {
			t.Fatalf("temporary file was not removed: %s", entry.Name())
		}
	}
}
