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
	file, err := LoadFile(path)
	if err != nil {
		return nil, err
	}
	if err := editableFile(file, path); err != nil {
		return nil, err
	}
	alias, node, found := strings.Cut(reference, ":")
	if !found || alias == "" || node == "" || strings.Contains(node, ":") {
		return nil, fmt.Errorf("module reference must use source:path, for example qmr:intro/author-ownership")
	}
	match := -1
	for index, entry := range file.Entries[1:] {
		if entry.Name != "include" || entry.Label == nil {
			continue
		}
		label, labelErr := entry.Label.ResolveString(nil)
		if labelErr != nil || label != reference {
			continue
		}
		if match >= 0 {
			return nil, fmt.Errorf("module %q is selected more than once; remove one declaration manually to preserve deliberate repetition", reference)
		}
		match = index + 1
	}
	if match < 0 {
		return nil, fmt.Errorf("module %q is not directly selected", reference)
	}
	file.Entries = append(file.Entries[:match], file.Entries[match+1:]...)
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
		return nil, fmt.Errorf("removing %q produces unresolved selection errors; inspect with mogent show --issues", reference)
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
