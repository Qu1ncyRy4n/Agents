// Package sourcepath validates and resolves source roots below a checkout.
package sourcepath

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// ValidateSubdir accepts an empty root or a portable relative slash path.
func ValidateSubdir(value string) error {
	if value == "" {
		return nil
	}
	if strings.Contains(value, "\\") || path.IsAbs(value) || path.Clean(value) != value {
		return fmt.Errorf("source subdir %q must be a clean relative slash path", value)
	}
	for _, part := range strings.Split(value, "/") {
		if part == "" || part == "." || part == ".." {
			return fmt.Errorf("source subdir %q contains an unsafe path component", value)
		}
	}
	return nil
}

// Resolve returns a directory below root and rejects symlinked components.
func Resolve(root, subdir string) (string, error) {
	if err := ValidateSubdir(subdir); err != nil {
		return "", err
	}
	current := root
	rootInfo, err := os.Lstat(current)
	if err != nil {
		return "", err
	}
	if rootInfo.Mode()&os.ModeSymlink != 0 {
		return "", fmt.Errorf("source root %q is a symlink", root)
	}
	if !rootInfo.IsDir() {
		return "", fmt.Errorf("source root %q is not a directory", root)
	}
	if subdir == "" {
		return current, nil
	}
	for _, part := range strings.Split(subdir, "/") {
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if err != nil {
			return "", fmt.Errorf("resolve source subdir %q: %w", subdir, err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("source subdir %q contains symlinked component %q", subdir, part)
		}
		if !info.IsDir() {
			return "", fmt.Errorf("source subdir %q component %q is not a directory", subdir, part)
		}
	}
	return current, nil
}
