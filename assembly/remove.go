package assembly

import (
	"fmt"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/renderfs"
	"github.com/Qu1ncyRy4n/Agents/syntax"
)

// Remove deletes one exact direct include from an authored build or plan. It
// deliberately does not treat an ancestor include as a match: exclusions and
// direct selection are different authoring operations.
func Remove(path, reference string, dryRun bool) ([]byte, error) {
	return RemoveMany(path, []string{reference}, dryRun)
}

// RemoveMany removes exact direct include declarations as one validated edit.
// It is used by optimize only after the user explicitly requests --apply.
func RemoveMany(path string, references []string, dryRun bool) ([]byte, error) {
	file, err := LoadFile(path)
	if err != nil {
		return nil, err
	}
	if err := editableFile(file, path); err != nil {
		return nil, err
	}
	if len(references) == 0 {
		return nil, fmt.Errorf("no module references to remove")
	}
	wanted := make(map[string]bool, len(references))
	for _, reference := range references {
		alias, node, found := strings.Cut(reference, ":")
		if !found || alias == "" || node == "" || strings.Contains(node, ":") {
			return nil, fmt.Errorf("module reference must use source:path, for example qmr:intro/author-ownership")
		}
		if wanted[reference] {
			return nil, fmt.Errorf("module %q was requested more than once", reference)
		}
		wanted[reference] = true
	}
	matches := make(map[string]int, len(wanted))
	var retained []*syntax.Entry
	for index, entry := range file.Entries {
		if index == 0 || entry.Name != "include" || entry.Label == nil {
			retained = append(retained, entry)
			continue
		}
		label, labelErr := entry.Label.ResolveString(nil)
		if labelErr != nil || !wanted[label] {
			retained = append(retained, entry)
			continue
		}
		matches[label]++
	}
	for _, reference := range references {
		switch matches[reference] {
		case 0:
			return nil, fmt.Errorf("module %q is not directly selected", reference)
		case 1:
		default:
			return nil, fmt.Errorf("module %q is selected more than once; remove one declaration manually to preserve deliberate repetition", reference)
		}
	}
	file.Entries = retained
	contents, err := syntax.Format(file, syntax.FormatOptions{})
	if err != nil {
		return nil, err
	}
	updated, err := syntax.Parse(path, contents)
	if err != nil {
		return nil, err
	}
	result, _, err := ResolvePlan(updated, path)
	if err != nil {
		return nil, err
	}
	if result.HasErrors() {
		return nil, fmt.Errorf("removing modules produces unresolved selection errors; inspect with mogent show --issues")
	}
	if !dryRun {
		if err := renderfs.WriteAtomically(path, contents); err != nil {
			return nil, err
		}
	}
	return contents, nil
}

func editableFile(file *syntax.File, path string) error {
	if file.Kind != syntax.Build {
		return nil
	}
	if len(file.Entries) == 0 {
		return fmt.Errorf("build %s has no header", path)
	}
	origin, _, err := stringField(file.Entries[0], "origin", nil)
	if err != nil {
		return err
	}
	if origin != "authored" {
		return fmt.Errorf("cannot edit generated %s; edit its plan", path)
	}
	return nil
}
