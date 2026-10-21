package assembly

import (
	"fmt"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/renderfs"
	"github.com/Qu1ncyRy4n/Agents/syntax"
)

type SelectOptions struct {
	All    bool
	DryRun bool
}

// Select appends one validated include to an authored build or plan. It never
// assembles outputs; callers use show/build afterward. Existing authored bytes
// and comments remain untouched ahead of the new declaration.
func Select(path, reference string, dryRun bool) ([]byte, error) {
	return SelectWithOptions(path, reference, SelectOptions{DryRun: dryRun})
}

// SelectWithOptions additionally supports explicit all-child selection for
// opt-in physical trees such as skill libraries.
func SelectWithOptions(path, reference string, options SelectOptions) ([]byte, error) {
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
			return nil, fmt.Errorf("cannot edit generated %s; select the module in its plan", path)
		}
	}
	alias, nodePath, found := strings.Cut(reference, ":")
	if !found || alias == "" || nodePath == "" || strings.Contains(nodePath, ":") {
		return nil, fmt.Errorf("module reference must use source:path, for example qmr:intro/author-ownership")
	}
	_, sources, err := ResolvePlan(file, path)
	if err != nil {
		return nil, err
	}
	source, found := sources[alias]
	if !found {
		return nil, fmt.Errorf("unknown source %q; list sources with mogent source list", alias)
	}
	node := source.Lib.ByPath[nodePath]
	if node == nil && source.Lib.Presets[nodePath].Name == "" && source.Lib.Recipes[nodePath].Name == "" {
		return nil, fmt.Errorf("unknown module %q; browse available modules with mogent source show %s", reference, alias)
	}
	if options.All && node == nil {
		return nil, fmt.Errorf("--all requires a content node, not preset or recipe %q", reference)
	}
	for _, entry := range file.Entries[1:] {
		if entry.Name != "include" || entry.Label == nil {
			continue
		}
		value, err := entry.Label.ResolveString(nil)
		if err == nil && value == reference {
			if !options.All {
				return nil, fmt.Errorf("module %q is already selected", reference)
			}
			return selectAllExisting(file, path, entry, node, reference, options.DryRun)
		}
	}
	if options.All && len(node.Children) == 0 {
		return nil, fmt.Errorf("module %q has no children to select with --all", reference)
	}
	declaration := "\ninclude \"" + reference + "\"\n"
	if options.All {
		declaration = "\ninclude \"" + reference + "\" [" + selectionNames(node) + "]\n"
	}
	contents := append([]byte(file.Original()), []byte(declaration)...)
	if _, err := syntax.Parse(path, contents); err != nil {
		return nil, err
	}
	if !options.DryRun {
		if err := renderfs.WriteAtomically(path, contents); err != nil {
			return nil, err
		}
	}
	return contents, nil
}

func selectAllExisting(file *syntax.File, path string, entry *syntax.Entry, node *Node, reference string, dryRun bool) ([]byte, error) {
	if node == nil || len(node.Children) == 0 {
		return nil, fmt.Errorf("module %q has no children to select with --all", reference)
	}
	if entry.Selection != nil {
		return nil, fmt.Errorf("module %q already has an explicit child selection", reference)
	}
	entry.Selection = &syntax.Value{Kind: syntax.List, Items: selectionValues(node)}
	contents, err := syntax.Format(file, syntax.FormatOptions{})
	if err != nil {
		return nil, err
	}
	if !dryRun {
		if err := renderfs.WriteAtomically(path, contents); err != nil {
			return nil, err
		}
	}
	return contents, nil
}

func selectionValues(node *Node) []syntax.Value {
	values := make([]syntax.Value, 0, len(node.Children))
	for _, child := range node.Children {
		values = append(values, syntax.Value{Kind: syntax.String, Parts: []syntax.StringPart{{Text: child.Name}}})
	}
	return values
}

func selectionNames(node *Node) string {
	values := selectionValues(node)
	parts := make([]string, 0, len(values))
	for _, value := range values {
		name, _ := value.ResolveString(nil)
		parts = append(parts, fmt.Sprintf("%q", name))
	}
	return strings.Join(parts, ", ")
}
