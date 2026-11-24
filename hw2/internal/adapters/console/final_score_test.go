package console

import (
	"bytes"
	"path/filepath"
	"pos/internal/application"
	"pos/internal/infrastructure/jsonrepo"
	"strings"
	"testing"
)

func TestConsoleCalculateFinalScore(t *testing.T) {
	tests := []struct {
		name           string
		withCriterion  bool
		withGrade      bool
		unknownStudent bool
		want           string
		command        string
	}{
		{"score", true, true, false, "Итоговый балл: 8 из 10.", "8"},
		{"missing grades", true, false, false, "Итоговый балл: 0 из 10.", "8"},
		{"empty rubric", false, false, false, "Не удалось рассчитать итог: rubric must not be empty", "8"},
		{"unknown student", true, false, true, "Не удалось рассчитать итог: student not found", "8"},
		{"current score", true, true, false, "Текущая успеваемость: 8 из 10.", "10"},
		{"current without grades", true, false, false, "нет выставленных оценок", "10"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := jsonrepo.NewRepository(filepath.Join(t.TempDir(), "state.json"))
			service := application.NewService(repository)
			if err := service.AddStudent("Student A", "group-1"); err != nil {
				t.Fatal(err)
			}
			if tt.withCriterion {
				if err := service.AddCriterion("Homework", 20, 30); err != nil {
					t.Fatal(err)
				}
			}
			state, err := repository.Load()
			if err != nil {
				t.Fatal(err)
			}
			id := state.Students[0].ID()
			if tt.withGrade {
				if err := service.SetGrade(id, state.Rubric.Criteria()[0].ID(), 15); err != nil {
					t.Fatal(err)
				}
			}
			if tt.unknownStudent {
				id = "missing"
			}
			var output bytes.Buffer
			input := strings.NewReader(tt.command + "\n " + string(id) + " \n1\n0\n")
			if err := NewConsole(service, input, &output).Run(); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(output.String(), tt.want) || !strings.Contains(output.String(), "Student A") {
				t.Fatalf("unexpected output: %s", output.String())
			}
		})
	}
}

func TestConsoleCalculateFinalScoreInterruptedInput(t *testing.T) {
	service := &testService{}
	var output bytes.Buffer
	if err := NewConsole(service, strings.NewReader("8\n"), &output).Run(); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(output.String(), "Итоговый балл:") {
		t.Fatal("calculated score with incomplete input")
	}
}
