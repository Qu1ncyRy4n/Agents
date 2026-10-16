package assembly

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/content"
	"github.com/Qu1ncyRy4n/Agents/syntax"
)

// LoadLibrary parses and validates one mogent.lib library declaration.
func LoadLibrary(root string) (*Library, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	contents, err := os.ReadFile(filepath.Join(root, LibraryFile))
	if os.IsNotExist(err) {
		return loadRawLibrary(root)
	}
	if err != nil {
		return nil, err
	}
	f, err := syntax.Parse(filepath.Join(root, LibraryFile), contents)
	if err != nil {
		return nil, err
	}
	if f.Kind != syntax.Library {
		return nil, fmt.Errorf("%s must begin with a library header", LibraryFile)
	}
	name, found, err := fieldString(f.Entries[0], "name")
	if err != nil || !found || strings.TrimSpace(name) == "" {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("library header requires a non-empty name")
	}
	l := &Library{Path: filepath.Join(root, LibraryFile), Root: root, Name: name, ByPath: make(map[string]*Node), Presets: make(map[string]Preset), Recipes: make(map[string]Preset), Groups: make(map[string]Group)}
	for _, entry := range f.Entries[1:] {
		switch {
		case isNodeEntry(entry):
			node, err := loadNode(l, nil, entry, "")
			if err != nil {
				return nil, err
			}
			l.Nodes = append(l.Nodes, node)
		case entry.Name == "preset":
			preset, err := loadPreset(l, entry)
			if err != nil {
				return nil, err
			}
			if _, duplicate := l.Presets[preset.Name]; duplicate || l.ByPath[preset.Name] != nil {
				return nil, fmt.Errorf("%s: duplicate public library name %q", l.Path, preset.Name)
			}
			l.Presets[preset.Name] = preset
		case entry.Name == "recipe":
			recipe, err := loadPreset(l, entry)
			if err != nil {
				return nil, err
			}
			for _, member := range recipe.Entries {
				label, _ := member.Label.ResolveString(nil)
				if !strings.HasPrefix(label, "self:") {
					return nil, fmt.Errorf("%s: recipe %q may include only self: references", l.Path, recipe.Name)
				}
			}
			if _, duplicate := l.Recipes[recipe.Name]; duplicate || l.ByPath[recipe.Name] != nil || l.Presets[recipe.Name].Name != "" {
				return nil, fmt.Errorf("%s: duplicate public library name %q", l.Path, recipe.Name)
			}
			l.Recipes[recipe.Name] = recipe
		case entry.Name == "group":
			group, err := loadGroup(entry)
			if err != nil {
				return nil, err
			}
			if _, exists := l.Groups[group.Name]; exists {
				return nil, fmt.Errorf("%s: duplicate group %q", l.Path, group.Name)
			}
			l.Groups[group.Name] = group
		default:
			return nil, fmt.Errorf("%s: unsupported library declaration %q", l.Path, entry.Name)
		}
	}
	if err := validateRelationshipsCatalog(l); err != nil {
		return nil, err
	}
	return l, nil
}

func loadRawLibrary(root string) (*Library, error) {
	l := &Library{Path: filepath.Join(root, LibraryFile), Root: root, Name: filepath.Base(root), ByPath: make(map[string]*Node), Presets: make(map[string]Preset), Recipes: make(map[string]Preset), Groups: make(map[string]Group)}
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == ".git" || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		inventory, err := content.Discover(filepath.Join(root, entry.Name()))
		if err != nil {
			return nil, err
		}
		node := &Node{Name: entry.Name(), Path: entry.Name(), Type: Directory, Source: entry.Name(), Library: l, Physical: inventory}
		l.ByPath[node.Path] = node
		l.Nodes = append(l.Nodes, node)
		for _, relative := range inventory.Paths() {
			if relative == "." {
				continue
			}
			parts := strings.Split(relative, "/")
			parent := node
			prefix := entry.Name()
			for index, part := range parts {
				prefix += "/" + part
				if existing := l.ByPath[prefix]; existing != nil {
					parent = existing
					continue
				}
				kind := inventory.Nodes[strings.Join(parts[:index+1], "/")].Kind
				childType := File
				if kind == content.Directory {
					childType = Directory
				}
				child := &Node{Name: part, Path: prefix, Type: childType, Parent: parent, Library: l, Physical: inventory}
				parent.Children = append(parent.Children, child)
				l.ByPath[prefix] = child
				parent = child
			}
		}
	}
	return l, nil
}

func loadGroup(entry *syntax.Entry) (Group, error) {
	if entry.Label == nil {
		return Group{}, fmt.Errorf("group requires a name")
	}
	name, err := entry.Label.ResolveString(nil)
	if err != nil {
		return Group{}, err
	}
	mode, found, err := fieldString(entry, "mode")
	if err != nil || !found || mode != "at-most-one" {
		return Group{}, fmt.Errorf("group %q requires mode = %q", name, "at-most-one")
	}
	return Group{Name: name, Mode: mode, Range: entry.Range}, nil
}

func isNodeEntry(e *syntax.Entry) bool {
	return e.Name == "node" || e.Name == "dir" || e.Name == "file" || e.Name == "md-file" || e.Name == "md-section"
}

func loadNode(l *Library, parent *Node, entry *syntax.Entry, prefix string) (*Node, error) {
	name, err := entry.Label.ResolveString(nil)
	if err != nil {
		return nil, err
	}
	typeName := entry.Name
	if entry.Name == "node" {
		typeName, _, err = fieldString(entry, "type")
		if err != nil {
			return nil, err
		}
		if typeName == "" && parent != nil && parent.Type == Directory {
			// Child annotations in a physical root address discovered physical
			// paths; inventory resolution determines file versus directory later.
			typeName = string(File)
		}
		if typeName == "" {
			return nil, fmt.Errorf("%s: node %q requires type", l.Path, name)
		}
	}
	typeValue := NodeType(typeName)
	if typeValue != Directory && typeValue != File && typeValue != Markdown && typeValue != Section {
		return nil, fmt.Errorf("%s: node %q has unknown type %q", l.Path, name, typeName)
	}
	nodePath := name
	if prefix != "" {
		nodePath = prefix + "/" + name
	}
	if _, duplicate := l.ByPath[nodePath]; duplicate {
		return nil, fmt.Errorf("%s: duplicate node path %q", l.Path, nodePath)
	}
	n := &Node{Name: name, Path: nodePath, Type: typeValue, Parent: parent, Library: l, Range: entry.Range}
	if n.Source, _, err = fieldString(entry, "source"); err != nil {
		return nil, err
	}
	if n.Title, _, err = fieldString(entry, "title"); err != nil {
		return nil, err
	}
	if n.Role, _, err = fieldString(entry, "role"); err != nil {
		return nil, err
	}
	if n.Content, _, err = fieldString(entry, "content"); err != nil {
		return nil, err
	}
	if n.Curate, _, err = fieldCurate(entry); err != nil {
		return nil, err
	}
	if n.Into, _, err = fieldString(entry, "into"); err != nil {
		return nil, err
	}
	if n.Tags, _, err = fieldStrings(entry, "tags"); err != nil {
		return nil, err
	}
	if parent != nil {
		n.Tags = appendUnique(parent.Tags, n.Tags)
	}
	if n.Variables, _, err = fieldStrings(entry, "vars"); err != nil {
		return nil, err
	}
	if n.Heading, _, err = fieldStrings(entry, "heading"); err != nil {
		return nil, err
	}
	if n.Requires, _, err = fieldStrings(entry, "requires"); err != nil {
		return nil, err
	}
	if n.ConflictsWith, _, err = fieldStrings(entry, "conflicts_with"); err != nil {
		return nil, err
	}
	if n.ExclusiveGroup, _, err = fieldString(entry, "exclusive_group"); err != nil {
		return nil, err
	}
	if value, found := fieldBool(entry, "default"); found {
		n.Default = &value
	}
	n.ReasonRequired, _ = fieldBool(entry, "reason_required")
	if err := validateNode(l, n); err != nil {
		return nil, err
	}
	l.ByPath[nodePath] = n
	for _, child := range entry.Body {
		if !isNodeEntry(child) {
			continue
		}
		loaded, err := loadNode(l, n, child, nodePath)
		if err != nil {
			return nil, err
		}
		n.Children = append(n.Children, loaded)
	}
	return n, nil
}

func appendUnique(values ...[]string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, group := range values {
		for _, value := range group {
			if !seen[value] {
				result, seen[value] = append(result, value), true
			}
		}
	}
	return result
}

func validateRelationshipsCatalog(l *Library) error {
	for path, node := range l.ByPath {
		for _, target := range append(append([]string(nil), node.Requires...), node.ConflictsWith...) {
			if l.ByPath[target] == nil {
				return fmt.Errorf("%s: node %q names missing relationship target %q", l.Path, path, target)
			}
		}
		if node.ExclusiveGroup != "" && l.Groups[node.ExclusiveGroup].Name == "" {
			return fmt.Errorf("%s: node %q names unknown exclusive group %q", l.Path, path, node.ExclusiveGroup)
		}
	}
	visiting, visited := make(map[string]bool), make(map[string]bool)
	var visit func(string) error
	visit = func(path string) error {
		if visiting[path] {
			return fmt.Errorf("%s: hard requirement cycle includes %q", l.Path, path)
		}
		if visited[path] {
			return nil
		}
		visiting[path] = true
		for _, target := range l.ByPath[path].Requires {
			if err := visit(target); err != nil {
				return err
			}
		}
		visiting[path], visited[path] = false, true
		return nil
	}
	for path := range l.ByPath {
		if err := visit(path); err != nil {
			return err
		}
	}
	return nil
}

func fieldCurate(e *syntax.Entry) (Curate, bool, error) {
	value, found, err := fieldString(e, "curate")
	if err != nil || !found {
		return "", found, err
	}
	curate := Curate(value)
	if curate != Foundation && curate != Choose && curate != Optional && curate != OptIn {
		return "", true, fmt.Errorf("%s: unknown curate value %q", e.Name, value)
	}
	return curate, true, nil
}

func validateNode(l *Library, n *Node) error {
	if n.Type == Section {
		if n.Source != "" {
			if _, err := safePath(l.Root, n.Source); err != nil {
				return fmt.Errorf("%s: section %q source: %w", l.Path, n.Path, err)
			}
			if filepath.Ext(n.Source) != ".md" {
				return fmt.Errorf("%s: section %q source must be Markdown", l.Path, n.Path)
			}
		}
		if n.Source == "" && len(n.Heading) > 0 {
			return fmt.Errorf("%s: section %q heading needs a source", l.Path, n.Path)
		}
		if n.Content != "" || n.Role != "" {
			return fmt.Errorf("%s: section %q cannot set content or role", l.Path, n.Path)
		}
		return nil
	}
	if len(n.Variables) > 0 || len(n.Heading) > 0 {
		return fmt.Errorf("%s: physical node %q cannot declare vars or heading", l.Path, n.Path)
	}
	if n.Source == "" && n.Parent == nil {
		return fmt.Errorf("%s: physical root %q requires source", l.Path, n.Path)
	}
	if n.Source != "" {
		if _, err := safePath(l.Root, n.Source); err != nil {
			return fmt.Errorf("%s: node %q source: %w", l.Path, n.Path, err)
		}
	}
	if n.Type == Markdown && n.Source != "" && filepath.Ext(n.Source) != ".md" {
		return fmt.Errorf("%s: md-file %q source must end in .md", l.Path, n.Path)
	}
	return nil
}

func loadPreset(l *Library, entry *syntax.Entry) (Preset, error) {
	name, err := entry.Label.ResolveString(nil)
	if err != nil {
		return Preset{}, err
	}
	for _, child := range entry.Body {
		if child.Name != "include" && child.Name != "replace" {
			return Preset{}, fmt.Errorf("%s: preset %q contains %q, not a selection", l.Path, name, child.Name)
		}
	}
	return Preset{Name: name, Entries: append([]*syntax.Entry(nil), entry.Body...), Range: entry.Range}, nil
}

func (l *Library) ResolvePhysical(node *Node) (*content.Inventory, error) {
	root := node
	for root.Parent != nil && root.Source == "" {
		root = root.Parent
	}
	if root.Type != Directory || root.Source == "" {
		return nil, fmt.Errorf("%s is not under a physical directory root", node.Path)
	}
	if root.Physical != nil {
		return root.Physical, nil
	}
	directory, err := safePath(l.Root, root.Source)
	if err != nil {
		return nil, err
	}
	inventory, err := content.Discover(directory)
	if err != nil {
		return nil, err
	}
	root.Physical = inventory
	return inventory, nil
}

func (l *Library) Node(value string) (*Node, bool) {
	node, found := l.ByPath[path.Clean(value)]
	return node, found
}
