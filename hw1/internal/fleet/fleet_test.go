package fleet

import (
	"errors"
	"math"
	"reflect"
	"testing"
)

type inspectionStub struct {
	accepted bool
	calls    int
}

func (service *inspectionStub) Inspect(Transport) (bool, string) {
	service.calls++
	return service.accepted, "решение тестового сервиса"
}

func testFleet(t *testing.T, service ServiceCenter) *Fleet {
	t.Helper()
	park, err := New(service)
	if err != nil {
		t.Fatal(err)
	}
	return park
}

func add(t *testing.T, park *Fleet, object InventoryItem) {
	t.Helper()
	if err := park.Add(object); err != nil {
		t.Fatal(err)
	}
}

func TestInspectionDecisionThroughDI(t *testing.T) {
	inspectionDecisions := []bool{
		false,
		true,
	}

	for _, accepted := range inspectionDecisions {
		service := &inspectionStub{accepted: accepted}
		park := testFleet(t, service)
		transport, _ := NewScooter("Самокат", "S1", 6, 2, false)
		err := park.Add(transport)
		if accepted && err != nil {
			t.Fatal(err)
		}
		if !accepted && !errors.Is(err, ErrRejected) {
			t.Fatal("ожидался отказ", err)
		}
		expected := 0
		if accepted {
			expected = 1
		}
		if service.calls != 1 {
			t.Fatalf("число осмотров: %d, ожидалось 1", service.calls)
		}
		if park.Report().Total != expected {
			t.Fatalf("на балансе %d объектов, ожидалось %d", park.Report().Total, expected)
		}
		if len(park.Beginners()) != expected {
			t.Fatalf("для новичков доступно %d объектов, ожидалось %d", len(park.Beginners()), expected)
		}
	}
}

func TestRejectedNumberCanBeRetried(t *testing.T) {
	service := &inspectionStub{}
	park := testFleet(t, service)
	transport, _ := NewScooter("Самокат", "S1", 6, 2, true)
	if !errors.Is(park.Add(transport), ErrRejected) {
		t.Fatal("ожидался отказ")
	}
	if park.Report().TotalEnergy != 0 {
		t.Fatal("отказ изменил энергию")
	}
	service.accepted = true
	add(t, park, transport)
}

func TestTechnicalServiceCenter(t *testing.T) {
	technicalStates := []bool{
		false,
		true,
	}

	for _, healthy := range technicalStates {
		transport, _ := NewBicycle("Велосипед", "B1", 5, healthy)
		accepted, reason := (TechnicalServiceCenter{}).Inspect(transport)
		if accepted != healthy {
			t.Fatalf("решение техосмотра: %v, исправность: %v", accepted, healthy)
		}
		if reason == "" {
			t.Fatal("сервис не указал причину решения")
		}
	}
}

func TestBeginnersBoundary(t *testing.T) {
	park := testFleet(t, TechnicalServiceCenter{})
	cases := []struct {
		number string
		ease   int
	}{
		{number: "B1", ease: 1},
		{number: "B5", ease: 5},
		{number: "B6", ease: 6},
		{number: "B9", ease: 10},
	}
	for _, testCase := range cases {
		transport, _ := NewBicycle("Велосипед", testCase.number, testCase.ease, true)
		add(t, park, transport)
	}
	helmet, _ := NewHelmet("Шлем", "H1")
	add(t, park, helmet)
	var numbers []string
	for _, transport := range park.Beginners() {
		numbers = append(numbers, transport.Number())
	}
	if !reflect.DeepEqual(numbers, []string{"B6", "B9"}) {
		t.Fatal(numbers)
	}
}

func TestReportAndInventory(t *testing.T) {
	park := testFleet(t, TechnicalServiceCenter{})
	emptyReport := Report{}
	if park.Report() != emptyReport {
		t.Fatal("отчёт нового парка не пуст")
	}
	if len(park.Inventory()) != 0 {
		t.Fatal("инвентарь нового парка не пуст")
	}
	if len(park.Beginners()) != 0 {
		t.Fatal("список для новичков нового парка не пуст")
	}
	scooter, _ := NewScooter("Самокат", "S1", 6, 2.5, true)
	electricBike, _ := NewElectricBike("Электро", "E1", 5, 4, true)
	bicycle, _ := NewBicycle("Велосипед", "B1", 10, true)
	helmet, _ := NewHelmet("Шлем", "H1")
	station, _ := NewChargingStation("Станция", "C1", 1.5)
	objects := []InventoryItem{
		scooter,
		electricBike,
		bicycle,
		helmet,
		station,
	}

	for _, object := range objects {
		add(t, park, object)
	}
	expected := Report{
		Total:           5,
		Transports:      3,
		Things:          2,
		TransportEnergy: 6.5,
		TotalEnergy:     8,
	}
	if park.Report() != expected {
		t.Fatalf("%+v", park.Report())
	}
	inventory := park.Inventory()
	var numbers []string
	for _, object := range inventory {
		numbers = append(numbers, object.Number())
	}
	if !reflect.DeepEqual(numbers, []string{"B1", "C1", "E1", "H1", "S1"}) {
		t.Fatal(numbers)
	}
	inventory[0] = nil
	if park.Inventory()[0] == nil {
		t.Fatal("выдан внутренний срез")
	}
	if _, consumesEnergy := any(bicycle).(EnergyConsumer); consumesEnergy {
		t.Fatal("велосипед потребляет электричество")
	}
	if _, consumesEnergy := any(helmet).(EnergyConsumer); consumesEnergy {
		t.Fatal("шлем потребляет электричество")
	}
}

func TestDuplicateAcrossTypesDoesNotInspectOrReplace(t *testing.T) {
	service := &inspectionStub{accepted: true}
	park := testFleet(t, service)
	helmet, _ := NewHelmet("Шлем", "1")
	transport, _ := NewScooter("Самокат", "1", 6, 2, true)
	add(t, park, helmet)
	if !errors.Is(park.Add(transport), ErrDuplicate) {
		t.Fatal("дубликат принят")
	}
	if service.calls != 0 {
		t.Fatal("для дубликата был вызван техосмотр")
	}
	if park.Inventory()[0].Name() != "Шлем" {
		t.Fatal("дубликат заменил исходный объект")
	}
}

func TestConstructorsIdentity(t *testing.T) {
	constructors := map[string]func(string, string) (InventoryItem, error){
		"scooter": func(name, number string) (InventoryItem, error) {
			return NewScooter(name, number, 6, 1, true)
		},
		"electric bike": func(name, number string) (InventoryItem, error) {
			return NewElectricBike(name, number, 6, 1, true)
		},
		"bicycle": func(name, number string) (InventoryItem, error) {
			return NewBicycle(name, number, 6, true)
		},
		"helmet": func(name, number string) (InventoryItem, error) {
			return NewHelmet(name, number)
		},
		"station": func(name, number string) (InventoryItem, error) {
			return NewChargingStation(name, number, 1)
		},
	}
	for name, create := range constructors {
		t.Run(name, func(t *testing.T) {
			for _, fields := range [][2]string{{"", "1"}, {"Техника", " "}} {
				if _, err := create(fields[0], fields[1]); err == nil {
					t.Fatal("пустое поле принято")
				}
			}
			object, err := create(" Название ", " 42 ")
			if err != nil {
				t.Fatal(err)
			}
			if object.Name() != "Название" || object.Number() != "42" {
				t.Fatalf("не удалены пробелы: %q, %q", object.Name(), object.Number())
			}
		})
	}
}

func TestTransportEaseValidation(t *testing.T) {
	for _, ease := range []int{-1, 0, 11} {
		if _, err := NewScooter("Самокат", "S1", ease, 1, true); err == nil {
			t.Errorf("самокат: принята простота %d", ease)
		}
		if _, err := NewElectricBike("Электровелосипед", "E1", ease, 1, true); err == nil {
			t.Errorf("электровелосипед: принята простота %d", ease)
		}
		if _, err := NewBicycle("Велосипед", "B1", ease, true); err == nil {
			t.Errorf("велосипед: принята простота %d", ease)
		}
	}
}

func TestEnergyValidation(t *testing.T) {
	for _, energy := range []float64{-1, math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := NewScooter("Самокат", "S1", 6, energy, true); err == nil {
			t.Errorf("самокат: принят расход %v", energy)
		}
		if _, err := NewElectricBike("Электровелосипед", "E1", 6, energy, true); err == nil {
			t.Errorf("электровелосипед: принят расход %v", energy)
		}
		if _, err := NewChargingStation("Станция", "C1", energy); err == nil {
			t.Errorf("станция: принят расход %v", energy)
		}
	}
	if _, err := NewScooter("Самокат", "S1", 6, 0, true); err != nil {
		t.Fatal(err)
	}
	if _, err := NewElectricBike("Электровелосипед", "E1", 6, 0, true); err != nil {
		t.Fatal(err)
	}
	if _, err := NewChargingStation("Станция", "C1", 0); err != nil {
		t.Fatal(err)
	}
}

type testTransport struct {
	number string
	ease   int
	energy float64
}

func (transport testTransport) Name() string {
	return "Новый транспорт"
}

func (transport testTransport) Number() string {
	return transport.number
}

func (transport testTransport) BeginnerEase() int {
	return transport.ease
}

func (transport testTransport) DailyEnergy() float64 {
	return transport.energy
}

func (transport testTransport) IsServiceable() bool {
	return true
}

func TestExtensionAndInvalidObjects(t *testing.T) {
	park := testFleet(t, TechnicalServiceCenter{})
	add(t, park, testTransport{number: "NEW", ease: 7, energy: 3})
	if park.Report().TransportEnergy != 3 {
		t.Fatal("энергия нового типа не учтена")
	}
	if len(park.Beginners()) != 1 {
		t.Fatal("новый тип не попал в список для новичков")
	}
	var pointer *Scooter
	invalidObjects := []InventoryItem{
		nil,
		pointer,
		Scooter{},
		Helmet{},
		testTransport{number: " X ", ease: 6},
		testTransport{number: "X", ease: 11},
		testTransport{number: "X", ease: 6, energy: math.NaN()},
	}

	for _, object := range invalidObjects {
		if err := park.Add(object); err == nil {
			t.Fatalf("принят %T", object)
		}
	}
	if park.Report().Total != 1 {
		t.Fatal("ошибки изменили парк")
	}
	var nilService *inspectionStub
	services := []ServiceCenter{
		nil,
		nilService,
	}

	for _, service := range services {
		if _, err := New(service); err == nil {
			t.Fatal("nil сервис принят")
		}
	}
}
