package assembly

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/Qu1ncyRy4n/Agents/content"
	"github.com/Qu1ncyRy4n/Agents/syntax"
)

// Preview is a read-only view of one library node. Content retains rendered
// section semantics and physical payload bytes are never transformed.
type Preview struct {
	Reference string
	Type      NodeType
	Source    string
	Content   []byte
	Opaque    bool
	Bytes     int
	Children  []string
	Headings  [][]string
}

// PreviewSource resolves one declared source reference without changing the
// entry file or workspace. lines limits text payload output; zero means full.
func PreviewSource(file *syntax.File, filename, reference string, lines int, headings bool) (*Preview, error) {
	if lines < 0 {
		return nil, fmt.Errorf("preview lines must be zero or greater")
	}
	alias, nodePath, found := strings.Cut(reference, ":")
	if !found || alias == "" || nodePath == "" {
		return nil, fmt.Errorf("content preview requires source:path, for example qmr:intro/assistant-agent")
	}
	_, sources, err := ResolvePlan(file, filename)
	if err != nil {
		return nil, err
	}
	source, found := sources[alias]
	if !found {
		return nil, fmt.Errorf("unknown source %q", alias)
	}
	node := source.Lib.ByPath[nodePath]
	if node == nil {
		return nil, fmt.Errorf("unknown module %q", reference)
	}
	preview := &Preview{Reference: reference, Type: node.Type}
	if node.Type == Section {
		if node.Source != "" {
			preview.Source = node.Source
		}
		bindings, err := previewBindings(file)
		if err != nil {
			return nil, err
		}
		contents, err := renderSection(node, bindings)
		if err != nil {
			return nil, fmt.Errorf("render %s: %w", reference, err)
		}
		if headings {
			document, err := content.ParseDocument(contents)
			if err != nil {
				return nil, err
			}
			collectHeadings(document.Headings, nil, &preview.Headings)
			return preview, nil
		}
		preview.Content = previewLines(contents, lines)
		return preview, nil
	}
	inventory, err := physicalInventory(node)
	if err != nil {
		return nil, err
	}
	item := inventory.Nodes[physicalPath(node)]
	if item == nil {
		return nil, fmt.Errorf("physical node missing for %q", reference)
	}
	preview.Source = item.Source
	if item.Kind == content.Directory {
		preview.Children = append([]string(nil), item.Children...)
		if headings {
			return nil, fmt.Errorf("--headings requires a physical Markdown file")
		}
		return preview, nil
	}
	if headings {
		document, err := inventory.Document(item.Path)
		if err != nil {
			return nil, err
		}
		collectHeadings(document.Headings, nil, &preview.Headings)
		return preview, nil
	}
	bytes, err := os.ReadFile(item.Source)
	if err != nil {
		return nil, err
	}
	preview.Bytes = len(bytes)
	if !utf8.Valid(bytes) || strings.ContainsRune(string(bytes), 0) {
		preview.Opaque = true
		return preview, nil
	}
	preview.Content = previewLines(bytes, lines)
	return preview, nil
}

func previewBindings(file *syntax.File) (map[string]string, error) {
	bindings := make(map[string]string)
	for _, entry := range file.Entries[1:] {
		if entry.Kind != syntax.Binding {
			continue
		}
		value, err := entry.Value.ResolveString(bindings)
		if err != nil {
			return nil, err
		}
		bindings[entry.Name] = value
	}
	return bindings, nil
}

func previewLines(bytes []byte, lines int) []byte {
	if lines == 0 {
		return bytes
	}
	rows := strings.Split(strings.TrimRight(string(bytes), "\n"), "\n")
	if len(rows) <= lines {
		return bytes
	}
	return []byte(strings.Join(rows[:lines], "\n") + fmt.Sprintf("\n... (%d more lines)\n", len(rows)-lines))
}

func collectHeadings(headings []*content.Heading, prefix []string, output *[][]string) {
	for _, heading := range headings {
		path := append(append([]string(nil), prefix...), heading.Title)
		*output = append(*output, path)
		collectHeadings(heading.Children, path, output)
	}
}
