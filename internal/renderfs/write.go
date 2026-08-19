// Package renderfs provides small atomic file operations without importing the renderer.
package renderfs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func WriteAtomically(path string, content []byte) error {
	directory := filepath.Dir(path)
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	temporary, err := os.CreateTemp(directory, ".mogent-write-")
	if err != nil {
		return fmt.Errorf("create temporary output: %w", err)
	}
	name := temporary.Name()
	open := true
	cleanup := func(original error) error {
		var closeErr error
		if open {
			closeErr = temporary.Close()
			open = false
		}
		removeErr := os.Remove(name)
		if errors.Is(removeErr, os.ErrNotExist) {
			removeErr = nil
		}
		return errors.Join(original, closeErr, removeErr)
	}
	if _, err := temporary.Write(content); err != nil {
		return cleanup(fmt.Errorf("write temporary output: %w", err))
	}
	if err := temporary.Chmod(0o644); err != nil {
		return cleanup(fmt.Errorf("set output permissions: %w", err))
	}
	if err := temporary.Close(); err != nil {
		open = false
		return cleanup(fmt.Errorf("close temporary output: %w", err))
	}
	open = false
	if err := os.Rename(name, path); err != nil {
		return cleanup(fmt.Errorf("replace output: %w", err))
	}
	return nil
}
