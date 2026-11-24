package main

import (
	"fmt"
	"os"
	"path/filepath"
	"pos/internal/adapters/console"
	"pos/internal/application"
	"pos/internal/domain/grading"
	"pos/internal/infrastructure/jsonrepo"
	"runtime"
)

func main() {
	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok || !filepath.IsAbs(sourceFile) {
		fmt.Fprintln(os.Stderr, "Не удалось определить корень проекта: соберите программу без -trimpath")
		os.Exit(1)
	}
	projectRoot := filepath.Join(filepath.Dir(sourceFile), "..", "..")
	repository := jsonrepo.NewRepository(filepath.Join(projectRoot, "data.json"))
	service := application.NewServiceWithScoreStrategies(repository, grading.FinalScoreStrategy{}, grading.CurrentScoreStrategy{})
	app := console.NewConsole(service, os.Stdin, os.Stdout)
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка приложения: %v\n", err)
		os.Exit(1)
	}
}
