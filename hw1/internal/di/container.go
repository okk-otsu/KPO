package di

import (
	"io"
	"vyshkat/internal/console"
	"vyshkat/internal/fleet"
)

type Container struct {
	App *console.App
}

func New(input io.Reader, output io.Writer) (*Container, error) {
	return NewWithServiceCenter(input, output, fleet.TechnicalServiceCenter{})
}

func NewWithServiceCenter(input io.Reader, output io.Writer, service fleet.ServiceCenter) (*Container, error) {
	park, err := fleet.New(service)
	if err != nil {
		return nil, err
	}
	app := console.New(park, input, output)
	container := &Container{
		App: app,
	}
	return container, nil
}
