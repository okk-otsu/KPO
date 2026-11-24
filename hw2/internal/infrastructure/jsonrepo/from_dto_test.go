package jsonrepo

import (
	"encoding/json"
	"pos/internal/application"
	"reflect"
	"strings"
	"testing"
)

func TestFromDTORoundTrip(t *testing.T) {
	input := `{"students":[{"id":"s1","name":"Student A","group":"group-1"}],"criteria":[{"id":"c1","name":"Homework","max_score":20,"weight":30}],"grades":[{"student_id":"s1","criterion_id":"c1","score":15}]}`
	var data stateDTO
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		t.Fatal(err)
	}
	state, err := fromDTO(data)
	if err != nil {
		t.Fatal(err)
	}
	if got := toDTO(state); !reflect.DeepEqual(got, data) {
		t.Fatalf("restored data = %#v, want %#v", got, data)
	}
}

func TestFromDTOEmpty(t *testing.T) {
	data := stateDTO{}
	state, err := fromDTO(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(state.Students) != 0 || len(state.Rubric.Criteria()) != 0 || len(state.Grades) != 0 {
		t.Fatalf("expected empty state, got %#v", state)
	}
}

func TestFromDTORejectsInvalidData(t *testing.T) {
	tests := []struct {
		name   string
		change func(*stateDTO)
		want   string
	}{
		{"invalid student", func(d *stateDTO) { d.Students[0].Name = " " }, "name must not be empty"},
		{"duplicate student", func(d *stateDTO) { d.Students = append(d.Students, d.Students[0]) }, "duplicate student id"},
		{"invalid maximum", func(d *stateDTO) { d.Criteria[0].MaxScore = 0 }, "maxScore must be positive"},
		{"invalid weight", func(d *stateDTO) { d.Criteria[0].Weight = 0 }, "weight must be"},
		{"duplicate criterion", func(d *stateDTO) { d.Criteria = append(d.Criteria, d.Criteria[0]) }, "duplicate criterion id"},
		{"duplicate criterion name", func(d *stateDTO) {
			other := d.Criteria[0]
			other.ID = "c2"
			other.Name = "HOMEWORK"
			d.Criteria = append(d.Criteria, other)
		}, "duplicate criterion name"},
		{"unknown student", func(d *stateDTO) { d.Grades[0].StudentID = "missing" }, "student not found"},
		{"unknown criterion", func(d *stateDTO) { d.Grades[0].CriterionID = "missing" }, "criterion not found"},
		{"negative score", func(d *stateDTO) { d.Grades[0].Score = -1 }, "score is not valid"},
		{"excessive score", func(d *stateDTO) { d.Grades[0].Score = 21 }, "score is not valid"},
		{"duplicate grade", func(d *stateDTO) { d.Grades = append(d.Grades, d.Grades[0]) }, "duplicate grade"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			learner := studentDTO{ID: "s1", Name: "Student A", Group: "group-1"}
			criterion := criterionDTO{ID: "c1", Name: "Homework", MaxScore: 20, Weight: 30}
			grade := gradeDTO{StudentID: "s1", CriterionID: "c1", Score: 15}
			students := []studentDTO{learner}
			criteria := []criterionDTO{criterion}
			grades := []gradeDTO{grade}
			data := stateDTO{Students: students, Criteria: criteria, Grades: grades}
			tt.change(&data)
			state, err := fromDTO(data)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %v, want containing %q", err, tt.want)
			}
			zero := application.State{}
			if !reflect.DeepEqual(state, zero) {
				t.Fatalf("returned partial state on error: %#v", state)
			}
		})
	}
}
