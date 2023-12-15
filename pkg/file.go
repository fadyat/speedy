package pkg

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	emptyFn = func() {}
)

func CopyFile(from, to string) error {
	f, err := os.ReadFile(filepath.Clean(from))
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if err = os.WriteFile(to, f, writePermission); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func CreateTemporaryFile(from, temporary string) (cleanup func(), err error) {
	f, err := os.Create(filepath.Clean(temporary))
	if err != nil {
		return emptyFn, fmt.Errorf("failed to create temporary file: %w", err)
	}

	if err = CopyFile(from, f.Name()); err != nil {
		return emptyFn, fmt.Errorf("failed to copy file: %w", err)
	}

	return func() { _ = os.Remove(f.Name()) }, nil
}
