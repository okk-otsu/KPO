package jsonrepo

import (
	"fmt"
	"pos/internal/application"
	"pos/internal/domain/grading"
	"pos/internal/domain/student"
)

func fromDTO(data stateDTO) (application.State, error) {
	state := application.State{
		Students: make([]student.Student, 0, len(data.Students)),
		Rubric:   grading.NewRubric(),
		Grades:   make([]grading.Grade, 0, len(data.Grades)),
	}

	studentIDs := make(map[student.StudentID]struct{}, len(data.Students))
	for _, item := range data.Students {
		learner, err := student.NewStudent(student.StudentID(item.ID), item.Name, item.Group)
		if err != nil {
			return application.State{}, fmt.Errorf("student %q: %w", item.ID, err)
		}

		if _, exists := studentIDs[learner.ID()]; exists {
			return application.State{}, fmt.Errorf("duplicate student id: %s", learner.ID())
		}

		studentIDs[learner.ID()] = struct{}{}
		state.Students = append(state.Students, learner)
	}

	for _, item := range data.Criteria {
		if err := state.Rubric.AddCriterion(grading.CriterionID(item.ID), item.Name, item.MaxScore, item.Weight); err != nil {
			return application.State{}, fmt.Errorf("criterion %q: %w", item.ID, err)
		}
	}

	type gradeKey struct {
		studentID   student.StudentID
		criterionID grading.CriterionID
	}

	gradeKeys := make(map[gradeKey]struct{}, len(data.Grades))
	for _, item := range data.Grades {

		studentID := student.StudentID(item.StudentID)
		if _, exists := studentIDs[studentID]; !exists {
			return application.State{}, fmt.Errorf("student not found: %s", studentID)
		}

		grade, err := grading.NewGrade(studentID, state.Rubric, grading.CriterionID(item.CriterionID), item.Score)
		if err != nil {
			return application.State{}, fmt.Errorf("grade for student %q: %w", studentID, err)
		}

		key := gradeKey{studentID: studentID, criterionID: grade.CriterionID()}
		if _, exists := gradeKeys[key]; exists {
			return application.State{}, fmt.Errorf("duplicate grade for student %s and criterion %s", studentID, grade.CriterionID())
		}

		gradeKeys[key] = struct{}{}
		state.Grades = append(state.Grades, grade)
	}

	return state, nil
}

func toDTO(state application.State) stateDTO {
	criteria := state.Rubric.Criteria()
	data := stateDTO{
		Students: make([]studentDTO, 0, len(state.Students)),
		Criteria: make([]criterionDTO, 0, len(criteria)),
		Grades:   make([]gradeDTO, 0, len(state.Grades)),
	}

	for _, student := range state.Students {
		data.Students = append(data.Students, studentDTO{
			ID:    string(student.ID()),
			Name:  student.Name(),
			Group: student.Group(),
		})
	}

	for _, criterion := range criteria {
		data.Criteria = append(data.Criteria, criterionDTO{
			ID:       string(criterion.ID()),
			Name:     criterion.Name(),
			MaxScore: criterion.MaxScore(),
			Weight:   criterion.Weight(),
		})
	}

	for _, grade := range state.Grades {
		data.Grades = append(data.Grades, gradeDTO{
			StudentID:   string(grade.StudentID()),
			CriterionID: string(grade.CriterionID()),
			Score:       grade.Score(),
		})
	}
	return data
}
