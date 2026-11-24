package console

import (
	"bytes"
	"errors"
	"pos/internal/domain/grading"
	"strings"
	"testing"
)

func TestConsoleAddCriterion(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"success", nil, "Критерий добавлен."},
		{"service error", errors.New("invalid weight"), "Не удалось добавить критерий: invalid weight"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &testService{err: tt.err}
			var output bytes.Buffer
			input := strings.NewReader("4\n Homework \nabc\n20\n1.5\n30\n3\n0\n")
			if err := NewConsole(service, input, &output).Run(); err != nil {
				t.Fatal(err)
			}
			if service.criterionCalls != 1 || service.criterionName != "Homework" || service.maxScore != 20 || service.weight != 30 {
				t.Fatalf("unexpected criterion call: %#v", service)
			}
			if service.listCriteriaCalls != 1 || !strings.Contains(output.String(), tt.want) || strings.Count(output.String(), "Введите целое число.") != 2 {
				t.Fatalf("unexpected output: %s", output.String())
			}
		})
	}
}

func TestConsoleListCriteria(t *testing.T) {
	rubric := grading.NewRubric()
	if err := rubric.AddCriterion("c1", "Homework", 20, 30); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name     string
		criteria []grading.Criterion
		err      error
		want     string
	}{
		{"populated", rubric.Criteria(), nil, "Homework | максимум: 20 | вес: 30"},
		{"empty", nil, nil, "Критериев пока нет."},
		{"error", nil, errors.New("load failed"), "Не удалось получить список критериев: load failed"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &testService{criteria: tt.criteria, err: tt.err}
			var output bytes.Buffer
			if err := NewConsole(service, strings.NewReader("3\n0\n"), &output).Run(); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(output.String(), tt.want) {
				t.Fatalf("missing %q: %s", tt.want, output.String())
			}
		})
	}
}

func TestConsoleAddCriterionInterruptedInput(t *testing.T) {
	inputs := []string{"4\n", "4\nHomework\n", "4\nHomework\n20\n"}
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			service := &testService{}
			var output bytes.Buffer
			if err := NewConsole(service, strings.NewReader(input), &output).Run(); err != nil {
				t.Fatal(err)
			}
			if service.criterionCalls != 0 {
				t.Fatal("criterion added with incomplete input")
			}
		})
	}
}
