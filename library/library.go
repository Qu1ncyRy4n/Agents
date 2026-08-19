// Package library indexes Markdown heading trees from local source directories.
package library

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

// Node is a selectable source directory, Markdown heading, or combined
// directory/heading identity. Only heading kinds own Markdown body content.
type Node struct {
	Kind     NodeKind
	Path     string
	Heading  string
	Body     string
	File     string
	Line     int
	Metadata Metadata
	Children []*Node
}

// NodeKind distinguishes selectable organizational directories from Markdown
// headings. Both have stable source paths, but only headings own content.
type NodeKind string

const (
	NodeDirectory        NodeKind = "directory"
	NodeHeading          NodeKind = "heading"
	NodeDirectoryHeading NodeKind = "directory-heading"
)

// Metadata is tool-only source metadata. It is useful for browsing and
// filtering source modules, but is never rendered into AGENTS.md.
type Metadata struct {
	Tags          []string `yaml:"tags,omitempty"`
	TLDR          string   `yaml:"tldr,omitempty"`
	Priority      *float64 `yaml:"priority,omitempty"`
	Scope         string   `yaml:"scope,omitempty"`
	Requires      []string `yaml:"requires,omitempty"`
	ConflictsWith []string `yaml:"conflicts_with,omitempty"`
}

// Index permits exact source-path lookup within one library.
type Index struct {
	Roots  []*Node
	ByPath map[string]*Node
}

// Load recursively reads Markdown files. A source must be a directory, and
// duplicate heading paths are errors rather than arbitrary source precedence.
func Load(root string) (*Index, error) {
	info, err := os.Lstat(root)
	if err != nil {
		return nil, fmt.Errorf("read source %q: %w", root, err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return nil, symlinkError(root, root)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("source %q is not a directory", root)
	}
	var files []string
	err = filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return symlinkError(root, path)
		}
		if entry.IsDir() {
			return nil
		}
		if strings.EqualFold(filepath.Ext(path), ".md") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk source %q: %w", root, err)
	}
	sort.Strings(files)
	if len(files) == 0 {
		return nil, fmt.Errorf("source %q contains no Markdown files", root)
	}
	index := &Index{ByPath: make(map[string]*Node)}
	directories := make(map[string]bool)
	for _, path := range files {
		prefix, err := directoryPrefix(root, path)
		if err != nil {
			return nil, err
		}
		parts := strings.Split(prefix, "/")
		for end := 1; prefix != "" && end <= len(parts); end++ {
			directories[strings.Join(parts[:end], "/")] = true
		}
		roots, err := parseFile(root, path)
		if err != nil {
			return nil, err
		}
		for _, node := range roots {
			if err := index.addHeadings(node); err != nil {
				return nil, err
			}
		}
	}
	index.buildTree(directories)
	return index, nil
}

func symlinkError(root, path string) error {
	target, err := os.Readlink(path)
	if err != nil {
		return fmt.Errorf("source %q contains unsupported symlink %q: inspect or replace the link with ordinary source content: %w", root, path, err)
	}
	display := path
	if relative, err := filepath.Rel(root, path); err == nil && relative != "." {
		display = relative
	}
	if path == root {
		return fmt.Errorf("source %q is an unsupported symlink to %q; declare the target directory directly instead", root, target)
	}
	return fmt.Errorf("source %q contains unsupported symlink %q -> %q; replace it with ordinary source content or declare the target directory as a separate source", root, display, target)
}

func (i *Index) addHeadings(node *Node) error {
	if _, found := i.ByPath[node.Path]; found {
		return fmt.Errorf("duplicate heading path %q", node.Path)
	}
	node.Kind = NodeHeading
	i.ByPath[node.Path] = node
	for _, child := range node.Children {
		if err := i.addHeadings(child); err != nil {
			return err
		}
	}
	return nil
}

func (i *Index) buildTree(directories map[string]bool) {
	for path := range directories {
		if existing, found := i.ByPath[path]; found {
			if existing.Kind == NodeHeading {
				existing.Kind = NodeDirectoryHeading
			}
		} else {
			i.ByPath[path] = &Node{
				Kind:    NodeDirectory,
				Path:    path,
				Heading: lastPathPart(path),
			}
		}
	}
	for _, node := range i.ByPath {
		node.Children = nil
	}
	paths := make([]string, 0, len(i.ByPath))
	for path := range i.ByPath {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		node := i.ByPath[path]
		parentPath, found := strings.CutSuffix(path, "/"+lastPathPart(path))
		if found && parentPath != "" {
			i.ByPath[parentPath].Children = append(i.ByPath[parentPath].Children, node)
			continue
		}
		i.Roots = append(i.Roots, node)
	}
}

func lastPathPart(path string) string {
	if index := strings.LastIndexByte(path, '/'); index >= 0 {
		return path[index+1:]
	}
	return path
}

type parsedNode struct {
	node  *Node
	level int
}

func parseFile(root string, path string) ([]*Node, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open Markdown %q: %w", path, err)
	}
	metadata := Metadata{}
	inFrontmatter := false
	frontmatterClosed := false
	var frontmatter bytes.Buffer
	var roots []*Node
	var stack []*parsedNode
	var current *Node
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if lineNumber == 1 && strings.TrimSpace(line) == "---" {
			inFrontmatter = true
			continue
		}
		if inFrontmatter {
			if strings.TrimSpace(line) == "---" {
				frontmatterClosed = true
				inFrontmatter = false
				parsed, err := parseFrontmatter(frontmatter.Bytes(), path)
				if err != nil {
					return nil, err
				}
				metadata = parsed
				continue
			}
			frontmatter.WriteString(line)
			frontmatter.WriteByte('\n')
			continue
		}
		level, heading, identifier, ok := headingLine(line)
		if !ok {
			if current != nil {
				current.Body += line + "\n"
			}
			continue
		}
		for len(stack) > 0 && stack[len(stack)-1].level >= level {
			stack = stack[:len(stack)-1]
		}
		parentPath := ""
		if len(stack) > 0 {
			parent := stack[len(stack)-1].node
			parentPath = parent.Path
		}
		slug := Slug(heading)
		if identifier != "" {
			slug = identifier
		}
		if slug == "" {
			return nil, fmt.Errorf("Markdown %q has an empty heading slug", path)
		}
		node := &Node{Kind: NodeHeading, Heading: heading, Path: slug, File: path, Line: lineNumber, Metadata: metadata}
		if parentPath != "" {
			node.Path = parentPath + "/" + slug
			stack[len(stack)-1].node.Children = append(stack[len(stack)-1].node.Children, node)
		} else {
			roots = append(roots, node)
		}
		stack = append(stack, &parsedNode{node: node, level: level})
		current = node
	}
	if err := scanner.Err(); err != nil {
		if closeErr := file.Close(); closeErr != nil {
			return nil, errors.Join(fmt.Errorf("read Markdown %q: %w", path, err), fmt.Errorf("close Markdown %q: %w", path, closeErr))
		}
		return nil, fmt.Errorf("read Markdown %q: %w", path, err)
	}
	if inFrontmatter || (lineNumber > 0 && !frontmatterClosed && frontmatter.Len() > 0) {
		if closeErr := file.Close(); closeErr != nil {
			return nil, errors.Join(fmt.Errorf("Markdown %q has unclosed YAML frontmatter", path), fmt.Errorf("close Markdown %q: %w", path, closeErr))
		}
		return nil, fmt.Errorf("Markdown %q has unclosed YAML frontmatter", path)
	}
	if err := file.Close(); err != nil {
		return nil, fmt.Errorf("close Markdown %q: %w", path, err)
	}
	if len(roots) == 0 {
		return nil, fmt.Errorf("Markdown %q contains no headings", path)
	}
	prefix, err := directoryPrefix(root, path)
	if err != nil {
		return nil, err
	}
	for _, root := range roots {
		if err := applyIDs(root, prefix); err != nil {
			return nil, fmt.Errorf("Markdown %q: %w", path, err)
		}
	}
	return roots, nil
}

func directoryPrefix(root string, path string) (string, error) {
	relative, err := filepath.Rel(root, filepath.Dir(path))
	if err != nil {
		return "", fmt.Errorf("resolve Markdown path %q under source %q: %w", path, root, err)
	}
	if relative == "." {
		return "", nil
	}
	parts := strings.Split(filepath.ToSlash(relative), "/")
	slugs := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == "." || part == "" {
			continue
		}
		slug := Slug(part)
		if slug == "" {
			return "", fmt.Errorf("Markdown %q has an empty directory slug for %q", path, part)
		}
		slugs = append(slugs, slug)
	}
	return strings.Join(slugs, "/"), nil
}

func parseFrontmatter(contents []byte, path string) (Metadata, error) {
	if len(bytes.TrimSpace(contents)) == 0 {
		return Metadata{}, nil
	}
	var metadata Metadata
	decoder := yaml.NewDecoder(bytes.NewReader(contents))
	decoder.KnownFields(true)
	if err := decoder.Decode(&metadata); err != nil {
		return Metadata{}, fmt.Errorf("parse metadata in %q: %w", path, err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != nil && !errors.Is(err, io.EOF) {
		return Metadata{}, fmt.Errorf("parse metadata in %q: %w", path, err)
	} else if err == nil {
		return Metadata{}, fmt.Errorf("parse metadata in %q: frontmatter must contain one YAML document", path)
	}
	if metadata.Priority != nil && (*metadata.Priority < 0 || *metadata.Priority > 1) {
		return Metadata{}, fmt.Errorf("parse metadata in %q: priority must be between 0.0 and 1.0", path)
	}
	if err := validateMetadataStrings(path, metadata); err != nil {
		return Metadata{}, err
	}
	return metadata, nil
}

func validateMetadataStrings(path string, metadata Metadata) error {
	for name, values := range map[string][]string{
		"tags":           metadata.Tags,
		"requires":       metadata.Requires,
		"conflicts_with": metadata.ConflictsWith,
	} {
		for _, value := range values {
			if strings.TrimSpace(value) == "" {
				return fmt.Errorf("parse metadata in %q: %s contains an empty value", path, name)
			}
		}
	}
	return nil
}

// applyIDs propagates already-parsed inline IDs into complete heading paths.
func applyIDs(node *Node, parentPath string) error {
	segment := Slug(node.Heading)
	if suffix := strings.LastIndex(node.Path, "/"); suffix >= 0 {
		segment = node.Path[suffix+1:]
	} else {
		segment = node.Path
	}
	if parentPath == "" {
		node.Path = segment
	} else {
		node.Path = parentPath + "/" + segment
	}
	for _, child := range node.Children {
		if err := applyIDs(child, node.Path); err != nil {
			return err
		}
	}
	return nil
}

func headingLine(line string) (int, string, string, bool) {
	if !strings.HasPrefix(line, "#") {
		return 0, "", "", false
	}
	level := 0
	for level < len(line) && line[level] == '#' {
		level++
	}
	if level == 0 || level > 6 || level == len(line) || line[level] != ' ' {
		return 0, "", "", false
	}
	heading := strings.TrimSpace(line[level:])
	heading = strings.TrimSpace(strings.TrimSuffix(heading, "#"))
	identifier := ""
	if before, comment, found := strings.Cut(heading, "<!-- id:"); found {
		if !strings.HasSuffix(strings.TrimSpace(comment), "-->") {
			return 0, "", "", false
		}
		identifier = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(comment), "-->"))
		if identifier == "" || identifier != Slug(identifier) {
			return 0, "", "", false
		}
		heading = strings.TrimSpace(before)
	}
	return level, heading, identifier, heading != ""
}

// Slug derives the stable, human-readable address documented by mogent.
func Slug(value string) string {
	var out strings.Builder
	dash := false
	for _, r := range strings.ToLower(strings.TrimSpace(value)) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			out.WriteRune(r)
			dash = false
		} else if !dash && out.Len() > 0 {
			out.WriteByte('-')
			dash = true
		}
	}
	return strings.Trim(out.String(), "-")
}
