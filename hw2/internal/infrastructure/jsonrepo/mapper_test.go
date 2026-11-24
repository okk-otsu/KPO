package jsonrepo

import (
	"encoding/json"
	"pos/internal/application"
	"pos/internal/domain/grading"
	"pos/internal/domain/student"
	"testing"
)

func TestToDTOJSON(t *testing.T) {
	learner, err := student.NewStudent("student-1", "Student A", "group-1")
	if err != nil {
		t.Fatal(err)
	}
	rubric := grading.NewRubric()
	if err := rubric.AddCriterion("criterion-1", "Homework", 20, 30); err != nil {
		t.Fatal(err)
	}
	grade, err := grading.NewGrade(learner.ID(), rubric, "criterion-1", 15)
	if err != nil {
		t.Fatal(err)
	}
	students := []student.Student{learner}
	grades := []grading.Grade{grade}
	state := application.State{Students: students, Rubric: rubric, Grades: grades}

	encoded, err := json.Marshal(toDTO(state))
	if err != nil {
		t.Fatal(err)
	}
	want := `{"students":[{"id":"student-1","name":"Student A","group":"group-1"}],"criteria":[{"id":"criterion-1","name":"Homework","max_score":20,"weight":30}],"grades":[{"student_id":"student-1","criterion_id":"criterion-1","score":15}]}`
	if string(encoded) != want {
		t.Fatalf("JSON = %s, want %s", encoded, want)
	}
}

func TestToDTOEmptyStateJSON(t *testing.T) {
	state := application.State{}
	encoded, err := json.Marshal(toDTO(state))
	if err != nil {
		t.Fatal(err)
	}
	want := `{"students":[],"criteria":[],"grades":[]}`
	if string(encoded) != want {
		t.Fatalf("JSON = %s, want %s", encoded, want)
	}
}
