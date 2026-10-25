package di

import (
	"bytes"
	"strings"
	"testing"
	"vyshkat/internal/fleet"
)

type rejectingCenter struct {
	calls int
}

func (service *rejectingCenter) Inspect(fleet.Transport) (bool, string) {
	service.calls++
	return false, "тестовый отказ"
}

func TestContainerInjectsServiceIntoConsole(t *testing.T) {
	service := &rejectingCenter{}
	var output bytes.Buffer
	container, err := NewWithServiceCenter(strings.NewReader("1\n3\nВелосипед\nB1\n6\nда\n3\n0\n"), &output, service)
	if err != nil {
		t.Fatal(err)
	}
	if err = container.App.Run(); err != nil {
		t.Fatal(err)
	}
	if service.calls != 1 {
		t.Fatalf("число вызовов техосмотра: %d, ожидалось 1", service.calls)
	}
	if !strings.Contains(output.String(), "Всего на балансе: 0;") {
		t.Fatal("отклонённый транспорт попал на баланс")
	}
	if !strings.Contains(output.String(), "тестовый отказ") {
		t.Fatal("приложение не вывело причину отказа")
	}
	if _, err = NewWithServiceCenter(strings.NewReader(""), &output, nil); err == nil {
		t.Fatal("nil сервис принят")
	}
}

func TestDefaultContainer(t *testing.T) {
	container, err := New(strings.NewReader("0\n"), &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if err = container.App.Run(); err != nil {
		t.Fatal(err)
	}
}
