package grading

import (
	"fmt"
	"pos/internal/domain/student"
	"strings"
)

type Grade struct {
	studentID   student.StudentID
	criterionID CriterionID
	score       int
}

func NewGrade(
	studentID student.StudentID,
	rubric Rubric,
	criterionID CriterionID,
	score int,
) (Grade, error) {
	if strings.TrimSpace(string(studentID)) == "" {
		return Grade{}, fmt.Errorf("studentID is empty")
	}
	criterion, err := rubric.FindCriterion(criterionID)
	if err != nil {
		return Grade{}, err
	}

	if score < 0 || criterion.MaxScore() < score {
		return Grade{}, fmt.Errorf("score is not valid: %v", score)
	}

	return Grade{
		studentID:   studentID,
		criterionID: criterion.ID(),
		score:       score,
	}, nil
}

func (g Grade) StudentID() student.StudentID { return g.studentID }

func (g Grade) CriterionID() CriterionID { return g.criterionID }

func (g Grade) Score() int { return g.score }
