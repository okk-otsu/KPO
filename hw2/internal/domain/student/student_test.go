package student

import "testing"

func TestNewStudent(t *testing.T) {
	got, err := NewStudent("student-1", " \tStudent A\n", "\t G1 \n")
	if err != nil {
		t.Fatal(err)
	}
	if got.ID() != "student-1" {
		t.Errorf("ID() = %q, want %q", got.ID(), "student-1")
	}
	if got.Name() != "Student A" {
		t.Errorf("Name() = %q, want %q", got.Name(), "Student A")
	}
	if got.Group() != "G1" {
		t.Errorf("Group() = %q, want %q", got.Group(), "G1")
	}
}

func TestNewStudentRejectsInvalidInput(t *testing.T) {
	cases := []struct {
		name  string
		id    StudentID
		label string
		group string
	}{
		{name: "empty ID", id: "", label: "Student A", group: "G1"},
		{name: "blank ID", id: " \t\n", label: "Student A", group: "G1"},
		{name: "empty name", id: "student-1", label: "", group: "G1"},
		{name: "blank name", id: "student-1", label: " \t\n", group: "G1"},
		{name: "empty group", id: "student-1", label: "Student A", group: ""},
		{name: "blank group", id: "student-1", label: "Student A", group: " \t\n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewStudent(tc.id, tc.label, tc.group)
			if err == nil {
				t.Fatal("expected validation error")
			}
			zero := Student{}
			if got != zero {
				t.Errorf("returned student on error: %#v", got)
			}
		})
	}
}
