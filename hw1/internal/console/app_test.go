package console

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
	"vyshkat/internal/fleet"
)

func newApp(t *testing.T, input io.Reader) (*App, *fleet.Fleet, *bytes.Buffer) {
	t.Helper()
	park, err := fleet.New(fleet.TechnicalServiceCenter{})
	if err != nil {
		t.Fatal(err)
	}
	output := &bytes.Buffer{}
	return New(park, input, output), park, output
}

func TestDemo(t *testing.T) {
	app, park, output := newApp(t, strings.NewReader(""))
	if err := app.Demo(); err != nil {
		t.Fatal(err)
	}
	expectedReport := fleet.Report{
		Total:           5,
		Transports:      3,
		Things:          2,
		TransportEnergy: 6.5,
		TotalEnergy:     8,
	}
	if park.Report() != expectedReport {
		t.Fatalf("%+v", park.Report())
	}
	expectedTexts := []string{
		"Пока нет.",
		"Пока пусто.",
		"техосмотр не пройден",
		"инвентарный номер уже занят",
		"простота должна быть",
		"конечным неотрицательным",
		"Неизвестная команда",
		"неизвестный вид",
		"S-001 — Городской самокат",
		"B-001 — Обычный велосипед",
		"8.00",
		"До свидания!",
	}

	for _, text := range expectedTexts {
		if !strings.Contains(output.String(), text) {
			t.Errorf("нет %q", text)
		}
	}
}

func TestInvalidInputAndRecovery(t *testing.T) {
	invalidInputs := []string{
		"1\n0\n",
		"2\n3\n",
		"2\n1\n\nH1\n",
		"2\n1\nШлем\n\n",
		"1\n1\nСамокат\nS1\nabc\n",
		"1\n1\nСамокат\nS1\n0\n",
		"1\n1\nСамокат\nS1\n6\nможет быть\n",
		"1\n1\nСамокат\nS1\n6\nда\nabc\n",
		"1\n1\nСамокат\nS1\n6\nда\nNaN\n",
		"2\n2\nСтанция\nC1\nInf\n",
		"1\n2\nЭлектро\nE1\n6\nда\n-3\n",
		"2\n2\nСтанция\nC1\nabc\n",
	}

	for _, input := range invalidInputs {
		t.Run(input, func(t *testing.T) {
			app, park, output := newApp(t, strings.NewReader(input+"2\n1\nШлем\nOK\n0\n"))
			if err := app.Run(); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(output.String(), "Ошибка:") {
				t.Fatalf("нет сообщения об ошибке: %s", output.String())
			}
			if park.Report().Total != 1 {
				t.Fatalf("после ошибки должен добавиться только корректный объект: %s", output.String())
			}
		})
	}
}

func TestEOFDuringEveryInputStep(t *testing.T) {
	scripts := []string{
		"1\n1\nСамокат\nS1\n6\nда\n2.5\n",
		"2\n2\nСтанция\nC1\n2\n",
	}

	for _, script := range scripts {
		lines := strings.Split(script, "\n")
		for lineCount := 0; lineCount < len(lines)-1; lineCount++ {
			input := strings.Join(lines[:lineCount], "\n")
			if lineCount > 0 {
				input += "\n"
			}
			app, park, _ := newApp(t, strings.NewReader(input))
			if err := app.Run(); err != nil {
				t.Fatal(err)
			}
			if park.Report().Total != 0 {
				t.Fatal("частичный ввод изменил парк")
			}
		}
	}
}

type brokenReader struct{}

func (brokenReader) Read([]byte) (int, error) {
	return 0, errors.New("ошибка чтения")
}

func TestReadFailure(t *testing.T) {
	prefixes := []string{
		"",
		"1\n",
		"1\n1\n",
	}

	for _, prefix := range prefixes {
		app, _, _ := newApp(t, io.MultiReader(strings.NewReader(prefix), brokenReader{}))
		if err := app.Run(); err == nil {
			t.Fatal("потеряна ошибка чтения")
		}
	}
}

func TestDecimalCommaAndCase(t *testing.T) {
	app, park, _ := newApp(t, strings.NewReader("1\n1\nСамокат\nS1\n6\nДА\n2,5\n0\n"))
	if err := app.Run(); err != nil {
		t.Fatal(err)
	}
	if park.Report().TotalEnergy != 2.5 {
		t.Fatal("не разобрана запятая")
	}
}
