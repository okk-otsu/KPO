package console

import (
	"bytes"
	"encoding/csv"
	"errors"
	"os"
	"path/filepath"
	"pos/internal/application"
	"pos/internal/infrastructure/jsonrepo"
	"reflect"
	"strings"
	"testing"
)

func TestConsoleExportResults(t *testing.T) {
	dir := t.TempDir()
	repository := jsonrepo.NewRepository(filepath.Join(dir, "state.json"))
	service := application.NewService(repository)
	if err := service.AddStudent(`Student, "A"`, "group-1"); err != nil {
		t.Fatal(err)
	}
	if err := service.AddStudent("Student B", "group-2"); err != nil {
		t.Fatal(err)
	}
	if err := service.AddCriterion("Homework", 20, 30); err != nil {
		t.Fatal(err)
	}
	state, err := repository.Load()
	if err != nil {
		t.Fatal(err)
	}
	if err := service.SetGrade(state.Students[0].ID(), state.Rubric.Criteria()[0].ID(), 15); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "results.csv")
	var output bytes.Buffer
	if err := NewConsole(service, strings.NewReader("9\n"+path+"\n0\n"), &output).Run(); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	got, err := csv.NewReader(bytes.NewReader(bytes.TrimPrefix(content, []byte("\uFEFF")))).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	want := [][]string{
		{"student_id", "name", "group", "final_score"},
		{string(state.Students[0].ID()), `Student, "A"`, "group-1", "8"},
		{string(state.Students[1].ID()), "Student B", "group-2", "0"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("CSV = %#v, want %#v", got, want)
	}
	if !strings.Contains(output.String(), "Итоги экспортированы:") {
		t.Fatal(output.String())
	}
	// A second export must not overwrite an existing file.
	if err := NewConsole(service, strings.NewReader("9\n"+path+"\n0\n"), &output).Run(); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(after, content) {
		t.Fatalf("existing file changed: %v", err)
	}
}

func TestConsoleExportErrorsAndEOF(t *testing.T) {
	service := &testService{err: errors.New("load failed")}
	path := filepath.Join(t.TempDir(), "results.csv")
	var output bytes.Buffer
	if err := NewConsole(service, strings.NewReader("9\n"+path+"\n1\n0\n"), &output).Run(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "Не удалось подготовить итоги: load failed") || service.calls != 1 {
		t.Fatal(output.String())
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("file created on service failure: %v", err)
	}
	if err := NewConsole(service, strings.NewReader("9\n"), &output).Run(); err != nil {
		t.Fatal(err)
	}
	service.err = nil
	badPath := filepath.Join(t.TempDir(), "missing", "out.csv")
	if err := NewConsole(service, strings.NewReader("9\n"+badPath+"\n0\n"), &output).Run(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "Не удалось создать файл") {
		t.Fatal(output.String())
	}
}
