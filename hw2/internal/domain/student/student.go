package student

import (
	"errors"
	"strings"
)

type StudentID string

type Student struct {
	id    StudentID
	name  string
	group string
}

func NewStudent(id StudentID, name string, group string) (Student, error) {

	if strings.TrimSpace(string(id)) == "" {
		return Student{}, errors.New("id must not be empty")
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return Student{}, errors.New("name must not be empty")
	}

	group = strings.TrimSpace(group)
	if group == "" {
		return Student{}, errors.New("group must not be empty")
	}

	return Student{
		id:    id,
		name:  name,
		group: group,
	}, nil
}

func (s Student) ID() StudentID { return s.id }

func (s Student) Name() string { return s.name }

func (s Student) Group() string { return s.group }
