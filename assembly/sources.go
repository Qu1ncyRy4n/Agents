package assembly

import (
	"fmt"
	"os"

	"github.com/Qu1ncyRy4n/Agents/renderfs"
	"github.com/Qu1ncyRy4n/Agents/syntax"
)

// SourceDeclarations reads project source locations without opening libraries or
// fetching Git repositories. Bindings are evaluated in declaration order.
func SourceDeclarations(file *syntax.File) ([]Source, error) {
	if file.Kind != syntax.Plan && file.Kind != syntax.Build {
		return nil, fmt.Errorf("sources require a plan or build file")
	}
	bindings := make(map[string]string)
	seen := make(map[string]bool)
	var sources []Source
	for _, entry := range file.Entries[1:] {
		if entry.Kind == syntax.Binding {
			value, err := entry.Value.ResolveString(bindings)
			if err != nil {
				return nil, err
			}
			bindings[entry.Name] = value
		}
		if entry.Kind != syntax.Block || entry.Name != "source" {
			continue
		}
		alias, err := entry.Label.ResolveString(bindings)
		if err != nil {
			return nil, err
		}
		if seen[alias] {
			return nil, fmt.Errorf("duplicate source alias %q", alias)
		}
		seen[alias] = true
		source := Source{Alias: alias}
		for name, target := range map[string]*string{"local": &source.Local, "git": &source.Git, "ref": &source.Ref, "commit": &source.Commit} {
			value, _, err := stringField(entry, name, bindings)
			if err != nil {
				return nil, err
			}
			*target = value
		}
		if (source.Local == "") == (source.Git == "") {
			return nil, fmt.Errorf("source %q requires exactly one of local or git", alias)
		}
		sources = append(sources, source)
	}
	return sources, nil
}

// AddSource preserves the original bytes and appends one validated declaration.
// Relative locations are interpreted relative to the entry file, like init.
// Generated builds must be changed through their owning plan instead.
func AddSource(path, alias, location string, dryRun bool) ([]byte, error) {
	file, err := LoadFile(path)
	if err != nil {
		return nil, err
	}
	if file.Kind == syntax.Build {
		origin, _, err := stringField(file.Entries[0], "origin", nil)
		if err != nil {
			return nil, err
		}
		if origin != "authored" {
			return nil, fmt.Errorf("cannot edit generated %s; add the source to its plan", path)
		}
	}
	declarations, err := SourceDeclarations(file)
	if err != nil {
		return nil, err
	}
	for _, source := range declarations {
		if source.Alias == alias {
			return nil, fmt.Errorf("source alias %q already exists", alias)
		}
	}
	entry, err := sourceEntry(path, alias, location)
	if err != nil {
		return nil, err
	}
	fragment, err := syntax.Format(&syntax.File{Kind: file.Kind, Entries: []*syntax.Entry{file.Entries[0], entry}}, syntax.FormatOptions{})
	if err != nil {
		return nil, err
	}
	fragmentFile, err := syntax.Parse("source fragment", fragment)
	if err != nil {
		return nil, err
	}
	start := fragmentFile.Entries[1].Range.Start.Offset
	contents := append([]byte(file.Original()+"\n"), fragment[start:]...)
	if _, err := syntax.Parse(path, contents); err != nil {
		return nil, err
	}
	if !dryRun {
		if info, err := os.Lstat(path); err != nil {
			return nil, err
		} else if !info.Mode().IsRegular() {
			return nil, fmt.Errorf("source file must be a regular file")
		}
		if err := renderfs.WriteAtomically(path, contents); err != nil {
			return nil, err
		}
	}
	return contents, nil
}
