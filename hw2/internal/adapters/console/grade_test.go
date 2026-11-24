package console

import (
	"bytes"
	"path/filepath"
	"pos/internal/application"
	"pos/internal/infrastructure/jsonrepo"
	"strings"
	"testing"
)

func TestConsoleGradeWorkflow(t *testing.T) {
	repository := jsonrepo.NewRepository(filepath.Join(t.TempDir(), "state.json"))
	service := application.NewService(repository)
	if err := service.AddStudent("Student A", "group-1"); err != nil {
		t.Fatal(err)
	}
	if err := service.AddCriterion("Homework", 20, 30); err != nil {
		t.Fatal(err)
	}
	state, err := repository.Load()
	if err != nil {
		t.Fatal(err)
	}
	studentID := string(state.Students[0].ID())
	criterionName := " homework "
	input := "5\n" + studentID + "\n" + criterionName + "\ninvalid\n15\n" +
		"5\n" + studentID + "\n" + criterionName + "\n18\n" +
		"6\n" + studentID + "\n" +
		"5\n" + studentID + "\n" + criterionName + "\n21\n" +
		"7\n" + studentID + "\n" + criterionName + "\n" +
		"6\n" + studentID + "\n0\n"
	var output bytes.Buffer
	if err := NewConsole(service, strings.NewReader(input), &output).Run(); err != nil {
		t.Fatal(err)
	}
	want := []string{"Введите целое число.", "Оценка сохранена.", "Критерий: Homework | балл: 18", "Не удалось выставить оценку:", "Оценка отменена.", "Оценок пока нет."}
	for _, text := range want {
		if !strings.Contains(output.String(), text) {
			t.Fatalf("missing %q: %s", text, output.String())
		}
	}
	state, err = repository.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(state.Grades) != 0 {
		t.Fatal("removed grade remains in storage")
	}
}

func TestConsoleGradeErrorsAndInterruptedInput(t *testing.T) {
	inputs := []string{
		"5\n", "5\ns1\n", "5\ns1\nc1\n", "6\n", "7\n", "7\ns1\n",
	}
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			repository := jsonrepo.NewRepository(filepath.Join(t.TempDir(), "state.json"))
			service := application.NewService(repository)
			var output bytes.Buffer
			if err := NewConsole(service, strings.NewReader(input), &output).Run(); err != nil {
				t.Fatal(err)
			}
			state, err := repository.Load()
			if err != nil || len(state.Grades) != 0 {
				t.Fatalf("unexpected state or error: %#v, %v", state, err)
			}
		})
	}
	repository := jsonrepo.NewRepository(filepath.Join(t.TempDir(), "state.json"))
	service := application.NewService(repository)
	var output bytes.Buffer
	input := strings.NewReader("5\nmissing\nc1\n1\n6\nmissing\n7\nmissing\nc1\n0\n")
	if err := NewConsole(service, input, &output).Run(); err != nil {
		t.Fatal(err)
	}
	want := []string{"Не удалось выставить оценку:", "Не удалось получить оценки:", "Не удалось отменить оценку:"}
	for _, text := range want {
		if !strings.Contains(output.String(), text) {
			t.Fatalf("missing %q: %s", text, output.String())
		}
	}
}
