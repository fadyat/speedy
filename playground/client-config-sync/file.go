package main

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

func copyFile(from, to string) error {
	f, err := os.ReadFile(filepath.Clean(from))
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if err = os.WriteFile(to, f, 0o600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func withTemporaryFile(from, to string) (cleanup func()) {
	f, err := os.Create(filepath.Clean(to))
	if err != nil {
		zap.S().Fatalf("failed to create temporary file: %v", err)
	}

	if err = copyFile(from, f.Name()); err != nil {
		zap.S().Fatalf("failed to copy file: %v", err)
	}

	zap.S().Infof("copied file from %s to %s", from, f.Name())
	return func() {
		if err = os.Remove(f.Name()); err != nil {
			zap.S().Fatalf("failed to remove file: %v", err)
		}
	}
}
