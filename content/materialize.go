package content

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/Qu1ncyRy4n/Agents/renderfs"
)

// Materialize writes captured payloads into an existing empty staging directory.
// Installing/rolling back the staging tree and tracking state belong to the
// workspace transaction. No source bytes are reread here.
func (m *Manifest) Materialize(directory string) error {
	info, err := os.Lstat(directory)
	if err != nil {
		return err
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("materialization requires a non-symlink staging directory")
	}
	entries, err := os.ReadDir(directory)
	if err != nil {
		return err
	}
	if len(entries) != 0 {
		return fmt.Errorf("materialization requires an empty staging directory")
	}
	dirs := make([]string, 0, len(m.Directories))
	for relative := range m.Directories {
		if err := ValidatePath(relative, true); err != nil {
			return err
		}
		dirs = append(dirs, relative)
	}
	files := make([]string, 0, len(m.Files))
	for relative := range m.Files {
		if err := ValidatePath(relative, false); err != nil {
			return err
		}
		files = append(files, relative)
	}
	sort.Strings(dirs)
	sort.Strings(files)
	for _, relative := range dirs {
		if err := os.MkdirAll(filepath.Join(directory, filepath.FromSlash(relative)), 0755); err != nil {
			return err
		}
	}
	for _, relative := range files {
		payload := m.Files[relative]
		if err := renderfs.WriteAtomicallyMode(filepath.Join(directory, filepath.FromSlash(relative)), payload.Bytes, payload.Mode); err != nil {
			return err
		}
	}
	// Apply restrictive directory modes only after creating every descendant.
	for index := len(dirs) - 1; index >= 0; index-- {
		relative := dirs[index]
		if err := os.Chmod(filepath.Join(directory, filepath.FromSlash(relative)), m.Directories[relative]); err != nil {
			return err
		}
	}
	return nil
}
