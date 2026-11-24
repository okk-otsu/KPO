package console

import (
	"fmt"
	"pos/internal/domain/grading"
	"pos/internal/domain/student"
	"strings"
)

func (c *Console) readGradeSelection() (student.StudentID, string, error) {
	studentID, err := c.readLine("ID студента (из списка студентов): ")
	if err != nil {
		return "", "", err
	}
	name, err := c.readLine("Название критерия: ")
	return student.StudentID(studentID), name, err
}

func (c *Console) setGrade() error {
	studentID, name, err := c.readGradeSelection()
	if err != nil {
		return err
	}
	score, err := c.readInt("Балл: ")
	if err != nil {
		return err
	}
	criterionID, err := c.findCriterionID(name)
	if err == nil {
		err = c.service.SetGrade(studentID, criterionID, score)
	}
	if err != nil {
		fmt.Fprintf(c.output, "Не удалось выставить оценку: %v\n", err)
		return nil
	}
	fmt.Fprintln(c.output, "Оценка сохранена.")
	return nil
}

func (c *Console) listGrades() error {
	id, err := c.readLine("ID студента (из списка студентов): ")
	if err != nil {
		return err
	}
	grades, err := c.service.ListGrades(student.StudentID(id))
	if err != nil {
		fmt.Fprintf(c.output, "Не удалось получить оценки: %v\n", err)
		return nil
	}
	if len(grades) == 0 {
		fmt.Fprintln(c.output, "Оценок пока нет.")
		return nil
	}
	criteria, err := c.service.ListCriteria()
	if err != nil {
		fmt.Fprintf(c.output, "Не удалось получить критерии: %v\n", err)
		return nil
	}
	names := make(map[grading.CriterionID]string, len(criteria))
	for _, criterion := range criteria {
		names[criterion.ID()] = criterion.Name()
	}
	for _, grade := range grades {
		fmt.Fprintf(c.output, "Критерий: %s | балл: %d\n", names[grade.CriterionID()], grade.Score())
	}
	return nil
}

func (c *Console) removeGrade() error {
	studentID, name, err := c.readGradeSelection()
	if err != nil {
		return err
	}
	criterionID, err := c.findCriterionID(name)
	if err == nil {
		err = c.service.RemoveGrade(studentID, criterionID)
	}
	if err != nil {
		fmt.Fprintf(c.output, "Не удалось отменить оценку: %v\n", err)
		return nil
	}
	fmt.Fprintln(c.output, "Оценка отменена.")
	return nil
}

func (c *Console) findCriterionID(name string) (grading.CriterionID, error) {
	criteria, err := c.service.ListCriteria()
	if err != nil {
		return "", err
	}
	for _, criterion := range criteria {
		if strings.EqualFold(criterion.Name(), name) {
			return criterion.ID(), nil
		}
	}
	return "", fmt.Errorf("критерий %q не найден", name)
}
