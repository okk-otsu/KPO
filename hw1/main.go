package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"vyshkat/internal/di"
)

func run(args []string, input io.Reader, output io.Writer) error {
	flags := flag.NewFlagSet("vyshkat", flag.ContinueOnError)
	flags.SetOutput(output)
	demo := flags.Bool("demo", false, "показать все основные сценарии")
	err := flags.Parse(args)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("неизвестные аргументы: %v", flags.Args())
	}
	container, err := di.New(input, output)
	if err != nil {
		return err
	}
	if *demo {
		return container.App.Demo()
	}
	return container.App.Run()
}

func main() {
	args := os.Args[1:]
	err := run(args, os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка:", err)
		os.Exit(1)
	}
}
