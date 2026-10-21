package assembly

import (
	"fmt"

	"github.com/Qu1ncyRy4n/Agents/renderfs"
	"github.com/Qu1ncyRy4n/Agents/syntax"
)

// Reorder groups bindings and source declarations before ordered includes. It
// preserves order within each group, so it cannot change assembly order. It is
// an explicit formatting rewrite, not overlap optimization or selection repair.
func Reorder(path string, dryRun bool) ([]byte, bool, error) {
	file, err := LoadFile(path)
	if err != nil {
		return nil, false, err
	}
	if file.Kind == syntax.Build {
		origin, _, err := stringField(file.Entries[0], "origin", nil)
		if err != nil {
			return nil, false, err
		}
		if origin != "authored" {
			return nil, false, fmt.Errorf("cannot reorder generated %s; reorder its plan", path)
		}
	}
	ordered := make([]*syntax.Entry, 0, len(file.Entries))
	ordered = append(ordered, file.Entries[0])
	for _, category := range []int{0, 1, 2} {
		for _, entry := range file.Entries[1:] {
			if reorderCategory(entry) == category {
				ordered = append(ordered, entry)
			}
		}
	}
	copy := *file
	copy.Entries = ordered
	contents, err := syntax.Format(&copy, syntax.FormatOptions{})
	if err != nil {
		return nil, false, err
	}
	changed := string(contents) != file.Original()
	if changed && !dryRun {
		if err := renderfs.WriteAtomically(path, contents); err != nil {
			return nil, false, err
		}
	}
	return contents, changed, nil
}

func reorderCategory(entry *syntax.Entry) int {
	if entry.Kind == syntax.Binding {
		return 0
	}
	if entry.Kind == syntax.Block && entry.Name == "source" {
		return 1
	}
	return 2
}
