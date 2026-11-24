package csvexport

import (
	"bytes"
	"encoding/csv"
	"errors"
	"pos/internal/application"
	"strings"
	"testing"
)

type failingWriter struct{ err error }

func (w failingWriter) Write([]byte) (int, error) { return 0, w.err }

func TestWriteEmptyAndFailure(t *testing.T) {
	results := []application.StudentResult{}
	var output bytes.Buffer
	if err := Write(&output, results); err != nil {
		t.Fatal(err)
	}
	want := "\uFEFFstudent_id,name,group,final_score\n"
	if output.String() != want {
		t.Fatalf("CSV = %q, want %q", output.String(), want)
	}
	wantErr := errors.New("write failed")
	writer := failingWriter{err: wantErr}
	if err := Write(writer, results); !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
}

func TestWriteCyrillic(t *testing.T) {
	result := application.StudentResult{StudentID: "s1", Name: "Студент А", Group: "Группа 1", Score: 8}
	results := []application.StudentResult{result}
	var output bytes.Buffer
	if err := Write(&output, results); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(output.String(), "\uFEFF") {
		t.Fatal("missing UTF-8 BOM")
	}
	rows, err := csv.NewReader(strings.NewReader(strings.TrimPrefix(output.String(), "\uFEFF"))).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 || rows[1][1] != result.Name || rows[1][2] != result.Group {
		t.Fatalf("Cyrillic text changed: %v", rows)
	}
}
