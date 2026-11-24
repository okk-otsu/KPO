package console

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestConsoleAddStudent(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"success", nil, "Студент добавлен."},
		{"error", errors.New("save failed"), "Не удалось добавить студента: save failed"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &testService{err: tt.err}
			var output bytes.Buffer
			input := strings.NewReader("2\n Student A \n group-1 \n1\n0\n")
			if err := NewConsole(service, input, &output).Run(); err != nil {
				t.Fatal(err)
			}
			if service.addCalls != 1 || service.name != "Student A" || service.group != "group-1" {
				t.Fatalf("unexpected add call: %#v", service)
			}
			if service.calls != 1 || !strings.Contains(output.String(), tt.want) {
				t.Fatalf("menu did not continue or missing result: %s", output.String())
			}
		})
	}
}

func TestConsoleAddStudentInterruptedInput(t *testing.T) {
	inputs := []string{"2\n", "2\nStudent A\n"}
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			service := &testService{}
			var output bytes.Buffer
			if err := NewConsole(service, strings.NewReader(input), &output).Run(); err != nil {
				t.Fatal(err)
			}
			if service.addCalls != 0 {
				t.Fatal("student added with incomplete input")
			}
		})
	}
}
