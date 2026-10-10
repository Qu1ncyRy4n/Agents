package state

import (
	"fmt"
	"os"
	"path/filepath"
)

// Layout records POSIX permission bits and every directory, including ".".
// It complements content hashes, so empty directories and chmod edits are not
// invisible to typed artifact ownership.
type Layout struct {
	Files       map[string]uint32 `json:"files"`
	Directories map[string]uint32 `json:"directories"`
}

func InspectLayout(root string) (Layout, error) {
	layout := Layout{Files: make(map[string]uint32), Directories: make(map[string]uint32)}
	err := filepath.WalkDir(root, func(file string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if !info.IsDir() && !info.Mode().IsRegular() {
			return fmt.Errorf("symlink or non-regular file in directory layout %q", file)
		}
		relative, err := filepath.Rel(root, file)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		if info.IsDir() {
			layout.Directories[relative] = uint32(info.Mode().Perm())
		} else {
			layout.Files[relative] = uint32(info.Mode().Perm())
		}
		return nil
	})
	return layout, err
}

func cloneLayout(layout *Layout) *Layout {
	if layout == nil {
		return nil
	}
	copy := &Layout{Files: make(map[string]uint32, len(layout.Files)), Directories: make(map[string]uint32, len(layout.Directories))}
	for name, mode := range layout.Files {
		copy.Files[name] = mode
	}
	for name, mode := range layout.Directories {
		copy.Directories[name] = mode
	}
	return copy
}

func equalModes(left, right map[string]uint32) bool {
	if len(left) != len(right) {
		return false
	}
	for name, mode := range left {
		if other, found := right[name]; !found || other != mode {
			return false
		}
	}
	return true
}
