// Package renderfs provides small atomic file operations without importing the renderer.
package renderfs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func WriteAtomically(path string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	temporary, err := os.CreateTemp(filepath.Dir(path), ".mogent-write-")
	if err != nil {
		return fmt.Errorf("create temporary output: %w", err)
	}
	name := temporary.Name()
	cleanup := func(original error) error { return errors.Join(original, temporary.Close(), os.Remove(name)) }
	if _, err := temporary.Write(content); err != nil {
		return cleanup(fmt.Errorf("write temporary output: %w", err))
	}
	if err := temporary.Chmod(0o644); err != nil {
		return cleanup(fmt.Errorf("set output permissions: %w", err))
	}
	if err := temporary.Close(); err != nil {
		_ = os.Remove(name)
		return fmt.Errorf("close temporary output: %w", err)
	}
	if err := os.Rename(name, path); err != nil {
		_ = os.Remove(name)
		return fmt.Errorf("replace output: %w", err)
	}
	return nil
}
