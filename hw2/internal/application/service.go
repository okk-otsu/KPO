package application

import (
	"fmt"
	"pos/internal/domain/grading"
	"pos/internal/domain/student"
	"slices"

	"github.com/google/uuid"
)

type Service struct {
	repository   Repository
	finalScore   grading.ScoreStrategy
	currentScore grading.ScoreStrategy
}

func (s *Service) CalculateFinalScore(studentID student.StudentID) (int, error) {
	return s.calculateScore(studentID, s.finalScore)
}

func (s *Service) CalculateCurrentScore(studentID student.StudentID) (int, error) {
	return s.calculateScore(studentID, s.currentScore)
}

func (s *Service) calculateScore(studentID student.StudentID, strategy grading.ScoreStrategy) (int, error) {
	state, err := s.repository.Load()
	if err != nil {
		return 0, err
	}

	if !slices.ContainsFunc(state.Students, func(existing student.Student) bool {
		return existing.ID() == studentID
	}) {
		return 0, fmt.Errorf("student not found: %s", studentID)
	}

	grades := make([]grading.Grade, 0)
	for _, grade := range state.Grades {
		if grade.StudentID() == studentID {
			grades = append(grades, grade)
		}
	}

	return strategy.Calculate(state.Rubric, grades)
}

func (s *Service) RemoveGrade(studentID student.StudentID, criterionID grading.CriterionID) error {
	state, err := s.repository.Load()
	if err != nil {
		return err
	}

	if !slices.ContainsFunc(state.Students, func(existing student.Student) bool {
		return existing.ID() == studentID
	}) {
		return fmt.Errorf("student not found: %s", studentID)
	}

	if _, err := state.Rubric.FindCriterion(criterionID); err != nil {
		return err
	}

	for i, grade := range state.Grades {
		if grade.StudentID() == studentID && grade.CriterionID() == criterionID {
			state.Grades = slices.Delete(slices.Clone(state.Grades), i, i+1)
			return s.repository.Save(state)
		}
	}
	return fmt.Errorf("grade not found for student %s and criterion %s", studentID, criterionID)
}

func (s *Service) ListGrades(studentID student.StudentID) ([]grading.Grade, error) {
	state, err := s.repository.Load()
	if err != nil {
		return nil, err
	}

	if !slices.ContainsFunc(state.Students, func(existing student.Student) bool {
		return existing.ID() == studentID
	}) {
		return nil, fmt.Errorf("student not found: %s", studentID)
	}

	grades := make([]grading.Grade, 0)
	for _, grade := range state.Grades {
		if grade.StudentID() == studentID {
			grades = append(grades, grade)
		}
	}
	return grades, nil
}

func (s *Service) SetGrade(studentID student.StudentID, criterionID grading.CriterionID, score int) error {
	state, err := s.repository.Load()
	if err != nil {
		return err
	}

	if !slices.ContainsFunc(state.Students, func(existing student.Student) bool {
		return existing.ID() == studentID
	}) {
		return fmt.Errorf("student not found: %s", studentID)
	}

	grade, err := grading.NewGrade(studentID, state.Rubric, criterionID, score)
	if err != nil {
		return err
	}

	state.Grades = slices.Clone(state.Grades)
	for i, existing := range state.Grades {
		if existing.StudentID() == studentID && existing.CriterionID() == criterionID {
			state.Grades[i] = grade
			return s.repository.Save(state)
		}
	}

	state.Grades = append(state.Grades, grade)
	return s.repository.Save(state)
}

func NewService(repository Repository) *Service {
	return NewServiceWithScoreStrategies(repository, grading.FinalScoreStrategy{}, grading.CurrentScoreStrategy{})
}

func NewServiceWithScoreStrategies(repository Repository, finalScore, currentScore grading.ScoreStrategy) *Service {
	return &Service{repository: repository, finalScore: finalScore, currentScore: currentScore}
}

func (s *Service) AddStudent(name, group string) error {

	state, err := s.repository.Load()
	if err != nil {
		return err
	}

	id := student.StudentID(uuid.New().String())
	newStudent, err := student.NewStudent(id, name, group)
	if err != nil {
		return err
	}

	state.Students = append(state.Students, newStudent)
	if err := s.repository.Save(state); err != nil {
		return err
	}

	return nil
}

func (s *Service) ListStudents() ([]student.Student, error) {
	state, err := s.repository.Load()
	if err != nil {
		return []student.Student{}, err
	}

	return state.Students, nil
}

func (s *Service) AddCriterion(name string, maxScore, weight int) error {
	state, err := s.repository.Load()
	if err != nil {
		return err
	}

	id := grading.CriterionID(uuid.New().String())
	if err := state.Rubric.AddCriterion(id, name, maxScore, weight); err != nil {
		return err
	}

	return s.repository.Save(state)
}

func (s *Service) ListCriteria() ([]grading.Criterion, error) {
	state, err := s.repository.Load()

	if err != nil {
		return []grading.Criterion{}, err
	}

	return state.Rubric.Criteria(), nil
}
