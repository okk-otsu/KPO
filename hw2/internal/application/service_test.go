package application

import (
	"errors"
	"slices"
	"testing"

	"pos/internal/domain/student"

	"github.com/google/uuid"
)

type testRepository struct {
	state     State
	loadErr   error
	saveErr   error
	saveCalls int
}

func (r *testRepository) Load() (State, error) {
	return r.state, r.loadErr
}

func (r *testRepository) Save(state State) error {
	r.saveCalls++
	if r.saveErr != nil {
		return r.saveErr
	}
	r.state = state
	return nil
}

func TestServiceAddStudent(t *testing.T) {
	repository := testRepository{}
	service := NewService(&repository)
	for range 2 {
		if err := service.AddStudent("Alex", "G1"); err != nil {
			t.Fatal(err)
		}
	}
	if repository.saveCalls != 2 || len(repository.state.Students) != 2 {
		t.Fatal("students were not saved")
	}
	students := repository.state.Students
	for _, added := range students {
		if _, err := uuid.Parse(string(added.ID())); err != nil {
			t.Fatalf("invalid student UUID: %v", err)
		}
		if added.Name() != "Alex" || added.Group() != "G1" {
			t.Fatal("student data was not preserved")
		}
	}
	if students[0].ID() == students[1].ID() {
		t.Fatal("students have the same ID")
	}
}

func TestServiceListStudents(t *testing.T) {
	existing, err := student.NewStudent("student-1", "Alex", "G1")
	if err != nil {
		t.Fatal(err)
	}
	want := []student.Student{existing}
	state := State{Students: want}
	repository := testRepository{state: state}
	service := NewService(&repository)
	got, err := service.ListStudents()
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(got, want) {
		t.Fatal("unexpected student list")
	}
	if repository.saveCalls != 0 {
		t.Fatal("listing students must not save state")
	}
}

func TestServiceListStudentsEmpty(t *testing.T) {
	repository := testRepository{}
	service := NewService(&repository)
	students, err := service.ListStudents()
	if err != nil || len(students) != 0 {
		t.Fatalf("expected empty list without error, got %v, %v", students, err)
	}
}

func TestServiceAddCriterion(t *testing.T) {
	repository := testRepository{}
	service := NewService(&repository)
	if err := service.AddCriterion("Correctness", 20, 70); err != nil {
		t.Fatal(err)
	}
	if err := service.AddCriterion("Style", 10, 30); err != nil {
		t.Fatal(err)
	}
	criteria := repository.state.Rubric.Criteria()
	if repository.saveCalls != 2 || len(criteria) != 2 {
		t.Fatal("criteria were not saved")
	}
	if criteria[0].Name() != "Correctness" || criteria[0].MaxScore() != 20 || criteria[0].Weight() != 70 {
		t.Fatal("criterion data was not preserved")
	}
	for _, criterion := range criteria {
		if _, err := uuid.Parse(string(criterion.ID())); err != nil {
			t.Fatalf("invalid criterion UUID: %v", err)
		}
	}
	if criteria[0].ID() == criteria[1].ID() {
		t.Fatal("criteria have the same ID")
	}
	if err := service.AddCriterion("correctness", 10, 50); err == nil {
		t.Fatal("expected duplicate name error")
	}
	if repository.saveCalls != 2 || !slices.Equal(criteria, repository.state.Rubric.Criteria()) {
		t.Fatal("rejected criterion changed saved state")
	}
}

func TestServiceRejectsInvalidInputWithoutSaving(t *testing.T) {
	cases := []struct {
		name string
		run  func(*Service) error
	}{
		{name: "student name", run: func(s *Service) error { return s.AddStudent("", "G1") }},
		{name: "student group", run: func(s *Service) error { return s.AddStudent("Alex", "") }},
		{name: "criterion name", run: func(s *Service) error { return s.AddCriterion("", 10, 50) }},
		{name: "criterion maximum", run: func(s *Service) error { return s.AddCriterion("Style", 0, 50) }},
		{name: "criterion weight", run: func(s *Service) error { return s.AddCriterion("Style", 10, 0) }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repository := testRepository{}
			service := NewService(&repository)
			if err := tc.run(service); err == nil {
				t.Fatal("expected validation error")
			}
			if repository.saveCalls != 0 {
				t.Fatal("invalid input must not be saved")
			}
		})
	}
}

func TestServiceRepositoryErrors(t *testing.T) {
	failure := errors.New("storage unavailable")
	operations := []struct {
		name    string
		mutates bool
		run     func(*Service) error
	}{
		{name: "add student", mutates: true, run: func(s *Service) error { return s.AddStudent("Alex", "G1") }},
		{name: "add criterion", mutates: true, run: func(s *Service) error { return s.AddCriterion("Style", 10, 50) }},
		{name: "list students", run: func(s *Service) error { _, err := s.ListStudents(); return err }},
	}
	for _, operation := range operations {
		t.Run(operation.name, func(t *testing.T) {
			repository := testRepository{loadErr: failure}
			service := NewService(&repository)
			if err := operation.run(service); !errors.Is(err, failure) {
				t.Fatalf("expected load error, got %v", err)
			}
			if repository.saveCalls != 0 {
				t.Fatal("must not save after load failure")
			}
			if operation.mutates {
				repository.loadErr = nil
				repository.saveErr = failure
				if err := operation.run(service); !errors.Is(err, failure) {
					t.Fatalf("expected save error, got %v", err)
				}
				if repository.saveCalls != 1 {
					t.Fatal("expected save attempt")
				}
			}
		})
	}
}
