package application

import (
	"fmt"
	"pos/internal/domain/grading"
	"pos/internal/domain/student"
)

type StudentResult struct {
	StudentID student.StudentID
	Name      string
	Group     string
	Score     int
}

func (s *Service) ListResults() ([]StudentResult, error) {
	state, err := s.repository.Load()
	if err != nil {
		return nil, err
	}
	gradesByStudent := make(map[student.StudentID][]grading.Grade)
	for _, grade := range state.Grades {
		gradesByStudent[grade.StudentID()] = append(gradesByStudent[grade.StudentID()], grade)
	}
	results := make([]StudentResult, 0, len(state.Students))
	for _, learner := range state.Students {
		score, err := s.finalScore.Calculate(state.Rubric, gradesByStudent[learner.ID()])
		if err != nil {
			return nil, fmt.Errorf("calculate result for student %s: %w", learner.ID(), err)
		}
		results = append(results, StudentResult{
			StudentID: learner.ID(), Name: learner.Name(), Group: learner.Group(), Score: score,
		})
	}
	return results, nil
}
