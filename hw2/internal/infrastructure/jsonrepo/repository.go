package jsonrepo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"pos/internal/application"
)

type Repository struct {
	path string
}

var _ application.Repository = (*Repository)(nil)

func NewRepository(path string) *Repository {
	return &Repository{path: path}
}

func (r *Repository) Load() (application.State, error) {
	content, err := os.ReadFile(r.path)
	if errors.Is(err, os.ErrNotExist) {
		return fromDTO(stateDTO{})
	}
	if err != nil {
		return application.State{}, fmt.Errorf("read state file %q: %w", r.path, err)
	}

	var data stateDTO
	if err := json.Unmarshal(content, &data); err != nil {
		return application.State{}, fmt.Errorf("decode state file %q: %w", r.path, err)
	}

	state, err := fromDTO(data)
	if err != nil {
		return application.State{}, fmt.Errorf("restore state from %q: %w", r.path, err)
	}
	return state, nil
}

func (r *Repository) Save(state application.State) error {
	content, err := json.MarshalIndent(toDTO(state), "", "  ")
	if err != nil {
		return fmt.Errorf("encode state: %w", err)
	}

	file, err := os.CreateTemp(filepath.Dir(r.path), ".pos-state-*.json")
	if err != nil {
		return fmt.Errorf("create temporary state file: %w", err)
	}
	defer os.Remove(file.Name())

	if _, err := file.Write(content); err != nil {
		file.Close()
		return fmt.Errorf("write temporary state file: %w", err)
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("close temporary state file: %w", err)
	}

	if err := os.Rename(file.Name(), r.path); err != nil {
		return fmt.Errorf("replace state file %q: %w", r.path, err)
	}

	return nil
}
