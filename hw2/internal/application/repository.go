package application

import (
	"pos/internal/domain/grading"
	"pos/internal/domain/student"
)

type State struct {
	Students []student.Student
	Rubric   grading.Rubric
	Grades   []grading.Grade
}

type Repository interface {
	Load() (State, error)
	Save(state State) error
}
