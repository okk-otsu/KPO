package console

import (
	"fmt"
	"pos/internal/domain/student"
)

func (c *Console) calculateFinalScore() error {
	id, err := c.readLine("ID студента (из списка студентов): ")
	if err != nil {
		return err
	}
	score, err := c.service.CalculateFinalScore(student.StudentID(id))
	if err != nil {
		fmt.Fprintf(c.output, "Не удалось рассчитать итог: %v\n", err)
		return nil
	}
	fmt.Fprintf(c.output, "Итоговый балл: %d из 10.\n", score)
	return nil
}

func (c *Console) calculateCurrentScore() error {
	id, err := c.readLine("ID студента (из списка студентов): ")
	if err != nil {
		return err
	}
	score, err := c.service.CalculateCurrentScore(student.StudentID(id))
	if err != nil {
		fmt.Fprintf(c.output, "Не удалось рассчитать текущую успеваемость: %v\n", err)
		return nil
	}
	fmt.Fprintf(c.output, "Текущая успеваемость: %d из 10.\n", score)
	return nil
}
