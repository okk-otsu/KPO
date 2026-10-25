package fleet

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strings"
)

var (
	ErrDuplicate = errors.New("инвентарный номер уже занят")
	ErrRejected  = errors.New("техосмотр не пройден")
)

type ServiceCenter interface {
	Inspect(Transport) (accepted bool, reason string)
}

type TechnicalServiceCenter struct{}

func (TechnicalServiceCenter) Inspect(transport Transport) (bool, string) {
	if !transport.IsServiceable() {
		return false, "обнаружена неисправность"
	}
	return true, "техника исправна"
}

type Fleet struct {
	service ServiceCenter
	items   map[string]InventoryItem
}

func New(service ServiceCenter) (*Fleet, error) {
	if isNil(service) {
		return nil, fmt.Errorf("сервис техосмотра не задан")
	}
	park := &Fleet{
		service: service,
		items:   make(map[string]InventoryItem),
	}
	return park, nil
}

func isNil(value any) bool {
	if value == nil {
		return true
	}
	reflectedValue := reflect.ValueOf(value)
	switch reflectedValue.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return reflectedValue.IsNil()
	default:
		return false
	}
}

func (f *Fleet) Add(object InventoryItem) error {
	if isNil(object) {
		return fmt.Errorf("объект не задан")
	}
	name := strings.TrimSpace(object.Name())
	number := strings.TrimSpace(object.Number())
	if name == "" || number == "" {
		return fmt.Errorf("некорректное название или инвентарный номер")
	}
	if object.Number() != number {
		return fmt.Errorf("некорректное название или инвентарный номер")
	}
	_, exists := f.items[object.Number()]
	if exists {
		return fmt.Errorf("%w: %s", ErrDuplicate, object.Number())
	}
	consumer, consumesEnergy := object.(EnergyConsumer)
	if consumesEnergy {
		energy := consumer.DailyEnergy()
		if energy < 0 || math.IsNaN(energy) || math.IsInf(energy, 0) {
			return fmt.Errorf("некорректный расход энергии")
		}
	}
	transport, isTransport := object.(Transport)
	if isTransport {
		ease := transport.BeginnerEase()
		if ease < 1 || ease > 10 {
			return fmt.Errorf("простота должна быть от 1 до 10")
		}
		accepted, reason := f.service.Inspect(transport)
		if !accepted {
			return fmt.Errorf("%w: %s", ErrRejected, reason)
		}
	}
	f.items[object.Number()] = object
	return nil
}

func (f *Fleet) Inventory() []InventoryItem {
	result := make([]InventoryItem, 0, len(f.items))
	for _, object := range f.items {
		result = append(result, object)
	}
	sort.Slice(result, func(left, right int) bool {
		return result[left].Number() < result[right].Number()
	})
	return result
}

func (f *Fleet) Beginners() []Transport {
	result := make([]Transport, 0)
	for _, object := range f.Inventory() {
		transport, isTransport := object.(Transport)
		if !isTransport {
			continue
		}
		if transport.BeginnerEase() >= 6 {
			result = append(result, transport)
		}
	}
	return result
}

type Report struct {
	Total           int
	Transports      int
	Things          int
	TransportEnergy float64
	TotalEnergy     float64
}

func (f *Fleet) Report() Report {
	report := Report{Total: len(f.items)}
	for _, object := range f.Inventory() {
		_, isTransport := object.(Transport)
		if isTransport {
			report.Transports++
		} else {
			report.Things++
		}
		consumer, consumesEnergy := object.(EnergyConsumer)
		if consumesEnergy {
			report.TotalEnergy += consumer.DailyEnergy()
			if isTransport {
				report.TransportEnergy += consumer.DailyEnergy()
			}
		}
	}
	return report
}
