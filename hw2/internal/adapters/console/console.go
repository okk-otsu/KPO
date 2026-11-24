package console

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"pos/internal/application"
	"pos/internal/domain/grading"
	"pos/internal/domain/student"
	"strings"
)

type Service interface {
	AddStudent(name, group string) error
	ListStudents() ([]student.Student, error)
	AddCriterion(name string, maxScore, weight int) error
	ListCriteria() ([]grading.Criterion, error)
	SetGrade(studentID student.StudentID, criterionID grading.CriterionID, score int) error
	ListGrades(studentID student.StudentID) ([]grading.Grade, error)
	RemoveGrade(studentID student.StudentID, criterionID grading.CriterionID) error
	CalculateFinalScore(studentID student.StudentID) (int, error)
	CalculateCurrentScore(studentID student.StudentID) (int, error)
	ListResults() ([]application.StudentResult, error)
}

type Console struct {
	service Service
	input   *bufio.Scanner
	output  io.Writer
}

func NewConsole(service Service, input io.Reader, output io.Writer) *Console {
	return &Console{service: service, input: bufio.NewScanner(input), output: output}
}

func (c *Console) Run() error {

	for {
		fmt.Fprintln(c.output, "")

		fmt.Fprintln(c.output, "1. Показать студентов")
		fmt.Fprintln(c.output, "2. Добавить студента")
		fmt.Fprintln(c.output, "3. Показать критерии")
		fmt.Fprintln(c.output, "4. Добавить критерий")
		fmt.Fprintln(c.output, "5. Выставить или исправить оценку")
		fmt.Fprintln(c.output, "6. Показать оценки студента")
		fmt.Fprintln(c.output, "7. Отменить оценку")
		fmt.Fprintln(c.output, "8. Рассчитать итог")
		fmt.Fprintln(c.output, "9. Экспортировать итоги в CSV (новый файл)")
		fmt.Fprintln(c.output, "10. Показать текущую успеваемость")
		fmt.Fprintln(c.output, "0. Выход")
		fmt.Fprint(c.output, "Выберите команду: ")
		if !c.input.Scan() {
			if err := c.input.Err(); err != nil {
				return fmt.Errorf("read command: %w", err)
			}
			return nil
		}
		var err error
		switch strings.TrimSpace(c.input.Text()) {
		case "0":
			return nil
		case "1":
			c.listStudents()
		case "2":
			err = c.addStudent()
		case "3":
			c.listCriteria()
		case "4":
			err = c.addCriterion()
		case "5":
			err = c.setGrade()
		case "6":
			err = c.listGrades()
		case "7":
			err = c.removeGrade()
		case "8":
			err = c.calculateFinalScore()
		case "10":
			err = c.calculateCurrentScore()
		case "9":
			err = c.exportResults()
		default:
			fmt.Fprintln(c.output, "Неизвестная команда. Выберите пункт меню.")
		}
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

func (c *Console) readLine(prompt string) (string, error) {
	fmt.Fprint(c.output, prompt)
	if !c.input.Scan() {
		if err := c.input.Err(); err != nil {
			return "", fmt.Errorf("read input: %w", err)
		}
		return "", io.EOF
	}
	return strings.TrimSpace(c.input.Text()), nil
}

func (c *Console) addStudent() error {
	name, err := c.readLine("Имя студента: ")
	if err != nil {
		return err
	}
	group, err := c.readLine("Группа: ")
	if err != nil {
		return err
	}
	if err := c.service.AddStudent(name, group); err != nil {
		fmt.Fprintf(c.output, "Не удалось добавить студента: %v\n", err)
		return nil
	}
	fmt.Fprintln(c.output, "Студент добавлен.")
	return nil
}

func (c *Console) listStudents() {
	students, err := c.service.ListStudents()
	if err != nil {
		fmt.Fprintf(c.output, "Не удалось получить список студентов: %v\n", err)
		return
	}
	if len(students) == 0 {
		fmt.Fprintln(c.output, "Студентов пока нет.")
		return
	}
	for _, learner := range students {
		fmt.Fprintf(c.output, "%s | %s | %s\n", learner.ID(), learner.Name(), learner.Group())
	}
}
