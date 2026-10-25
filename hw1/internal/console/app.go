package console

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"vyshkat/internal/fleet"
)

type App struct {
	fleet  *fleet.Fleet
	input  *bufio.Scanner
	output io.Writer
	echo   bool
}

func New(park *fleet.Fleet, input io.Reader, output io.Writer) *App {
	app := &App{
		fleet:  park,
		input:  bufio.NewScanner(input),
		output: output,
	}
	return app
}

func (a *App) read(prompt string) (string, error) {
	fmt.Fprint(a.output, prompt)
	if !a.input.Scan() {
		if err := a.input.Err(); err != nil {
			return "", err
		}
		return "", io.EOF
	}
	value := strings.TrimSpace(a.input.Text())
	if a.echo {
		fmt.Fprintln(a.output, value)
	}
	return value, nil
}

func (a *App) Run() error {
	fmt.Fprintln(a.output, "ВышКат — учёт парка. Энергия указывается в кВт·ч/сутки.")
	for {
		fmt.Fprintln(a.output, "\n1 — добавить транспорт; 2 — добавить вещь; 3 — отчёт; 4 — для новичков; 5 — инвентарь; 0 — выход")
		command, err := a.read("> ")
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		switch command {
		case "0":
			fmt.Fprintln(a.output, "До свидания!")
			return nil
		case "1", "2":
			if command == "1" {
				err = a.addTransport()
			} else {
				err = a.addThing()
			}
			if errors.Is(err, io.EOF) {
				return nil
			}
			if err != nil {
				if a.input.Err() != nil {
					return err
				}
				fmt.Fprintln(a.output, "Ошибка:", err)
			} else {
				fmt.Fprintln(a.output, "Добавлено на баланс.")
			}
		case "3":
			a.report()
		case "4":
			a.beginners()
		case "5":
			a.inventory()
		default:
			fmt.Fprintln(a.output, "Неизвестная команда. Введите число из меню.")
		}
	}
}

func (a *App) integer(prompt string) (int, error) {
	text, err := a.read(prompt)
	if err != nil {
		return 0, err
	}
	value, err := strconv.Atoi(text)
	if err != nil {
		return 0, fmt.Errorf("ожидалось целое число")
	}
	return value, nil
}

func (a *App) energy() (float64, error) {
	text, err := a.read("Суточная энергия (кВт·ч): ")
	if err != nil {
		return 0, err
	}
	text = strings.ReplaceAll(text, ",", ".")
	value, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0, fmt.Errorf("ожидалось число, например 2.5")
	}
	return value, nil
}

func (a *App) readIdentity() (string, string, error) {
	name, err := a.read("Название: ")
	if err != nil {
		return "", "", err
	}
	number, err := a.read("Инвентарный номер: ")
	if err != nil {
		return "", "", err
	}
	if name == "" || number == "" {
		return "", "", fmt.Errorf("название и инвентарный номер обязательны")
	}
	return name, number, nil
}

func (a *App) addTransport() error {
	kind, err := a.read("Транспорт: 1 — самокат, 2 — электровелосипед, 3 — велосипед: ")
	if err != nil {
		return err
	}
	if kind != "1" && kind != "2" && kind != "3" {
		return fmt.Errorf("неизвестный вид объекта")
	}
	name, number, err := a.readIdentity()
	if err != nil {
		return err
	}
	ease, err := a.integer("Простота для новичка (1–10): ")
	if err != nil {
		return err
	}
	if ease < 1 || ease > 10 {
		return fmt.Errorf("простота должна быть от 1 до 10")
	}
	state, err := a.read("Техника исправна? (да/нет): ")
	if err != nil {
		return err
	}
	state = strings.ToLower(state)
	if state != "да" && state != "нет" {
		return fmt.Errorf("введите да или нет")
	}
	serviceable := state == "да"
	var transport fleet.Transport
	switch kind {
	case "1":
		energy, readErr := a.energy()
		if readErr != nil {
			return readErr
		}
		transport, err = fleet.NewScooter(name, number, ease, energy, serviceable)
	case "2":
		energy, readErr := a.energy()
		if readErr != nil {
			return readErr
		}
		transport, err = fleet.NewElectricBike(name, number, ease, energy, serviceable)
	case "3":
		transport, err = fleet.NewBicycle(name, number, ease, serviceable)
	}
	if err != nil {
		return err
	}
	return a.fleet.Add(transport)
}

func (a *App) addThing() error {
	kind, err := a.read("Вещь: 1 — шлем, 2 — зарядная станция: ")
	if err != nil {
		return err
	}
	if kind != "1" && kind != "2" {
		return fmt.Errorf("неизвестный вид объекта")
	}
	name, number, err := a.readIdentity()
	if err != nil {
		return err
	}
	var object fleet.InventoryItem
	switch kind {
	case "1":
		object, err = fleet.NewHelmet(name, number)
	case "2":
		energy, readErr := a.energy()
		if readErr != nil {
			return readErr
		}
		object, err = fleet.NewChargingStation(name, number, energy)
	}
	if err != nil {
		return err
	}
	return a.fleet.Add(object)
}

func (a *App) report() {
	report := a.fleet.Report()
	fmt.Fprintf(a.output, "Всего на балансе: %d; транспорт: %d; вещи: %d\n", report.Total, report.Transports, report.Things)
	fmt.Fprintf(a.output, "Энергия транспорта: %.2f кВт·ч/сутки\nЭнергия всего баланса: %.2f кВт·ч/сутки\n", report.TransportEnergy, report.TotalEnergy)
}

func (a *App) beginners() {
	items := a.fleet.Beginners()
	fmt.Fprintln(a.output, "Транспорт для новичков (простота >= 6):")
	if len(items) == 0 {
		fmt.Fprintln(a.output, "Пока нет.")
	}
	for _, object := range items {
		fmt.Fprintf(a.output, "%s — %s (простота %d)\n", object.Number(), object.Name(), object.BeginnerEase())
	}
}

func (a *App) inventory() {
	items := a.fleet.Inventory()
	fmt.Fprintln(a.output, "Инвентарь:")
	if len(items) == 0 {
		fmt.Fprintln(a.output, "Пока пусто.")
	}
	for _, object := range items {
		fmt.Fprintf(a.output, "%s — %s\n", object.Number(), object.Name())
	}
}
