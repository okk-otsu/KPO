package fleet

import (
	"fmt"
	"math"
	"strings"
)

type InventoryItem interface {
	Name() string
	Number() string
}

type EnergyConsumer interface {
	DailyEnergy() float64
}

type Transport interface {
	InventoryItem
	BeginnerEase() int
	IsServiceable() bool
}

type item struct {
	name   string
	number string
}

func (object item) Name() string {
	return object.name
}

func (object item) Number() string {
	return object.number
}

func newItem(name, number string) (item, error) {
	name = strings.TrimSpace(name)
	number = strings.TrimSpace(number)
	if name == "" || number == "" {
		return item{}, fmt.Errorf("название и инвентарный номер не должны быть пустыми")
	}
	return item{name: name, number: number}, nil
}

type vehicle struct {
	item
	ease        int
	serviceable bool
}

func (transport vehicle) BeginnerEase() int {
	return transport.ease
}

func (transport vehicle) IsServiceable() bool {
	return transport.serviceable
}

func newVehicle(name, number string, ease int, serviceable bool) (vehicle, error) {
	object, err := newItem(name, number)
	if err != nil {
		return vehicle{}, err
	}
	if ease < 1 || ease > 10 {
		return vehicle{}, fmt.Errorf("простота должна быть от 1 до 10")
	}
	return vehicle{item: object, ease: ease, serviceable: serviceable}, nil
}

type energyUsage struct {
	energy float64
}

func (power energyUsage) DailyEnergy() float64 {
	return power.energy
}

func newEnergyUsage(energy float64) (energyUsage, error) {
	if math.IsNaN(energy) || math.IsInf(energy, 0) || energy < 0 {
		return energyUsage{}, fmt.Errorf("расход энергии должен быть конечным неотрицательным числом")
	}
	return energyUsage{energy: energy}, nil
}

type Scooter struct {
	vehicle
	energyUsage
}
type ElectricBike struct {
	vehicle
	energyUsage
}
type Bicycle struct {
	vehicle
}
type Helmet struct {
	item
}
type ChargingStation struct {
	item
	energyUsage
}

func NewScooter(name, number string, ease int, energy float64, serviceable bool) (Scooter, error) {
	transport, err := newVehicle(name, number, ease, serviceable)
	if err != nil {
		return Scooter{}, err
	}
	power, err := newEnergyUsage(energy)
	if err != nil {
		return Scooter{}, err
	}
	return Scooter{vehicle: transport, energyUsage: power}, nil
}

func NewElectricBike(name, number string, ease int, energy float64, serviceable bool) (ElectricBike, error) {
	transport, err := newVehicle(name, number, ease, serviceable)
	if err != nil {
		return ElectricBike{}, err
	}
	power, err := newEnergyUsage(energy)
	if err != nil {
		return ElectricBike{}, err
	}
	return ElectricBike{vehicle: transport, energyUsage: power}, nil
}

func NewBicycle(name, number string, ease int, serviceable bool) (Bicycle, error) {
	transport, err := newVehicle(name, number, ease, serviceable)
	if err != nil {
		return Bicycle{}, err
	}
	return Bicycle{vehicle: transport}, nil
}

func NewHelmet(name, number string) (Helmet, error) {
	object, err := newItem(name, number)
	if err != nil {
		return Helmet{}, err
	}
	return Helmet{item: object}, nil
}

func NewChargingStation(name, number string, energy float64) (ChargingStation, error) {
	object, err := newItem(name, number)
	if err != nil {
		return ChargingStation{}, err
	}
	power, err := newEnergyUsage(energy)
	if err != nil {
		return ChargingStation{}, err
	}
	return ChargingStation{item: object, energyUsage: power}, nil
}

var (
	_ Transport      = Scooter{}
	_ Transport      = ElectricBike{}
	_ Transport      = Bicycle{}
	_ InventoryItem  = Helmet{}
	_ InventoryItem  = ChargingStation{}
	_ EnergyConsumer = Scooter{}
	_ EnergyConsumer = ElectricBike{}
	_ EnergyConsumer = ChargingStation{}
)
