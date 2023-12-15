package pkg

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

const (
	writePermission = 0o600
)

func FromYaml[T any](path string) (*T, error) {
	raw, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var t T
	if err = yaml.Unmarshal(raw, &t); err != nil {
		return nil, fmt.Errorf("failed to unmarshal file: %w", err)
	}

	return &t, nil
}

func ToYaml(path string, content any) error {
	raw, err := yaml.Marshal(content)
	if err != nil {
		return fmt.Errorf("failed to marshal content: %w", err)
	}

	if err = os.WriteFile(filepath.Clean(path), raw, writePermission); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
