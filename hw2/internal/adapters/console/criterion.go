package console

import (
	"fmt"
	"strconv"
)

func (c *Console) readInt(prompt string) (int, error) {
	for {
		text, err := c.readLine(prompt)
		if err != nil {
			return 0, err
		}
		value, err := strconv.Atoi(text)
		if err != nil {
			fmt.Fprintln(c.output, "Введите целое число.")
			continue
		}
		return value, nil
	}
}

func (c *Console) addCriterion() error {
	name, err := c.readLine("Название критерия: ")
	if err != nil {
		return err
	}
	maxScore, err := c.readInt("Максимальный балл: ")
	if err != nil {
		return err
	}
	weight, err := c.readInt("Вес (1–100): ")
	if err != nil {
		return err
	}
	if err := c.service.AddCriterion(name, maxScore, weight); err != nil {
		fmt.Fprintf(c.output, "Не удалось добавить критерий: %v\n", err)
		return nil
	}
	fmt.Fprintln(c.output, "Критерий добавлен.")
	return nil
}

func (c *Console) listCriteria() {
	criteria, err := c.service.ListCriteria()
	if err != nil {
		fmt.Fprintf(c.output, "Не удалось получить список критериев: %v\n", err)
		return
	}
	if len(criteria) == 0 {
		fmt.Fprintln(c.output, "Критериев пока нет.")
		return
	}
	for _, criterion := range criteria {
		fmt.Fprintf(c.output, "%s | максимум: %d | вес: %d\n",
			criterion.Name(), criterion.MaxScore(), criterion.Weight())
	}
}
