package console

import (
	"fmt"
	"os"
	"pos/internal/adapters/csvexport"
)

func (c *Console) exportResults() error {
	path, err := c.readLine("Путь к CSV-файлу: ")
	if err != nil {
		return err
	}
	results, err := c.service.ListResults()
	if err != nil {
		fmt.Fprintf(c.output, "Не удалось подготовить итоги: %v\n", err)
		return nil
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		fmt.Fprintf(c.output, "Не удалось создать файл (укажите новый путь): %v\n", err)
		return nil
	}
	writeErr := csvexport.Write(file, results)
	closeErr := file.Close()
	if writeErr != nil {
		os.Remove(path)
		fmt.Fprintf(c.output, "Не удалось экспортировать итоги: %v\n", writeErr)
		return nil
	}
	if closeErr != nil {
		os.Remove(path)
		fmt.Fprintf(c.output, "Не удалось завершить запись: %v\n", closeErr)
		return nil
	}
	fmt.Fprintf(c.output, "Итоги экспортированы: %s\n", path)
	return nil
}
