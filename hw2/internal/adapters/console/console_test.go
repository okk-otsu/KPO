package console

import (
	"bytes"
	"errors"
	"io"
	"pos/internal/application"
	"pos/internal/domain/grading"
	"pos/internal/domain/student"
	"strings"
	"testing"
)

type testService struct {
	students          []student.Student
	err               error
	calls             int
	addCalls          int
	name              string
	group             string
	criteria          []grading.Criterion
	criterionCalls    int
	listCriteriaCalls int
	criterionName     string
	maxScore          int
	weight            int
}

func (s *testService) ListResults() ([]application.StudentResult, error) {
	return nil, s.err
}

func (s *testService) CalculateFinalScore(student.StudentID) (int, error) {
	return 0, s.err
}

func (s *testService) SetGrade(student.StudentID, grading.CriterionID, int) error {
	return s.err
}

func (s *testService) ListGrades(student.StudentID) ([]grading.Grade, error) {
	return nil, s.err
}

func (s *testService) RemoveGrade(student.StudentID, grading.CriterionID) error {
	return s.err
}

func (s *testService) AddCriterion(name string, maxScore, weight int) error {
	s.criterionCalls++
	s.criterionName, s.maxScore, s.weight = name, maxScore, weight
	return s.err
}

func (s *testService) ListCriteria() ([]grading.Criterion, error) {
	s.listCriteriaCalls++
	return s.criteria, s.err
}

func (s *testService) AddStudent(name, group string) error {
	s.addCalls++
	s.name, s.group = name, group
	return s.err
}

func (s *testService) ListStudents() ([]student.Student, error) {
	s.calls++
	return s.students, s.err
}

func TestConsoleListsStudentsAndContinues(t *testing.T) {
	learner, err := student.NewStudent("s1", "Student A", "group-1")
	if err != nil {
		t.Fatal(err)
	}
	students := []student.Student{learner}
	service := &testService{students: students}
	var output bytes.Buffer
	console := NewConsole(service, strings.NewReader("invalid\n 1 \n0\n1\n"), &output)
	if err := console.Run(); err != nil {
		t.Fatal(err)
	}
	if service.calls != 1 {
		t.Fatalf("service calls = %d, want 1", service.calls)
	}
	want := []string{"Неизвестная команда", "s1 | Student A | group-1"}
	for _, text := range want {
		if !strings.Contains(output.String(), text) {
			t.Fatalf("output does not contain %q: %s", text, output.String())
		}
	}
}

func TestConsoleEmptyListAndServiceError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"empty list", nil, "Студентов пока нет."},
		{"service error", errors.New("storage unavailable"), "Не удалось получить список студентов: storage unavailable"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &testService{err: tt.err}
			var output bytes.Buffer
			console := NewConsole(service, strings.NewReader("1\n1\n0\n"), &output)
			if err := console.Run(); err != nil {
				t.Fatal(err)
			}
			if service.calls != 2 || !strings.Contains(output.String(), tt.want) {
				t.Fatalf("calls = %d, output = %s", service.calls, output.String())
			}
		})
	}
}

type failingReader struct{ err error }

func (r failingReader) Read([]byte) (int, error) { return 0, r.err }

func TestConsoleInputEndAndError(t *testing.T) {
	service := &testService{}
	if err := NewConsole(service, strings.NewReader(""), io.Discard).Run(); err != nil {
		t.Fatal(err)
	}
	want := errors.New("input unavailable")
	input := failingReader{err: want}
	if err := NewConsole(service, input, io.Discard).Run(); !errors.Is(err, want) {
		t.Fatalf("error = %v, want %v", err, want)
	}
}

func (s *testService) CalculateCurrentScore(student.StudentID) (int, error) { return 0, s.err }
