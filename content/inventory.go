// Package content provides typed physical inventories and explicit Markdown
// document views. Discovering a file never transforms its original payload.
package content

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type Kind string

const (
	Directory   Kind = "dir"
	File        Kind = "file"
	HeadingKind Kind = "heading"
)

// Node describes one physical object. Path is relative to the inventory root;
// the root itself is ".". Bundle identifies the enclosing inferred skill root.
type Node struct {
	Path     string
	Kind     Kind
	Source   string
	Mode     os.FileMode
	Role     string
	Bundle   string
	Children []string
}

type Inventory struct {
	Root  string
	Nodes map[string]*Node
}

// Discover imports all regular files and directories, rejecting symlinks and
// special files. SKILL.md identifies a skill boundary without parsing or changing
// its bytes; other files remain opaque regardless of extension.
func Discover(root string) (*Inventory, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	info, err := os.Lstat(root)
	if err != nil {
		return nil, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("content root is a symlink: %s", root)
	}
	inventory := &Inventory{Root: root, Nodes: make(map[string]*Node)}
	err = filepath.WalkDir(root, func(file string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if !info.IsDir() && !info.Mode().IsRegular() {
			return fmt.Errorf("content contains symlink or non-regular file: %s", file)
		}
		relative, err := filepath.Rel(root, file)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		kind := File
		if info.IsDir() {
			kind = Directory
		}
		inventory.Nodes[relative] = &Node{Path: relative, Kind: kind, Source: file, Mode: info.Mode().Perm()}
		if relative != "." {
			parent := inventory.Nodes[path.Dir(relative)]
			parent.Children = append(parent.Children, relative)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	for _, node := range inventory.Nodes {
		if node.Kind == Directory {
			if skill := inventory.Nodes[path.Join(node.Path, "SKILL.md")]; skill != nil && skill.Kind == File {
				node.Role = "skill"
			}
		}
	}
	for _, node := range inventory.Nodes {
		for current := node.Path; ; current = path.Dir(current) {
			if ancestor := inventory.Nodes[current]; ancestor != nil && ancestor.Role == "skill" {
				node.Bundle = current
				break
			}
			if current == "." {
				break
			}
		}
	}
	return inventory, nil
}

// Paths returns a deterministic physical inventory, including the root.
func (i *Inventory) Paths() []string {
	paths := make([]string, 0, len(i.Nodes))
	for p := range i.Nodes {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	return paths
}

// Document opens an optional structural view. It does not turn headings into
// physical children or change the copied file/bundle view.
func (i *Inventory) Document(relative string) (*Document, error) {
	node, found := i.Nodes[relative]
	if !found || node.Kind != File {
		return nil, fmt.Errorf("document view requires a file: %q", relative)
	}
	if !strings.EqualFold(filepath.Ext(relative), ".md") && !strings.EqualFold(filepath.Ext(node.Source), ".md") {
		return nil, fmt.Errorf("document view requires a Markdown file: %q", relative)
	}
	bytes, err := os.ReadFile(node.Source)
	if err != nil {
		return nil, err
	}
	return ParseDocument(bytes)
}

// ValidatePath rejects ambiguous or escaping slash paths. Root selection is
// allowed only when allowRoot is true; destination files cannot be the root.
func ValidatePath(value string, allowRoot bool) error {
	if value == "." && allowRoot {
		return nil
	}
	if value == "" || value == "." || value == ".." || filepath.IsAbs(value) || strings.ContainsRune(value, 0) || strings.HasPrefix(value, "/") || strings.Contains(value, "\\") || strings.HasPrefix(value, "../") || path.Clean(value) != value || strings.Contains(value, "/../") {
		return fmt.Errorf("%q must be a safe relative slash path", value)
	}
	return nil
}
