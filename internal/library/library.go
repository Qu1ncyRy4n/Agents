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

// Node is one Markdown heading and the plain content directly beneath it.
type Node struct {
	Path     string
	Heading  string
	Body     string
	File     string
	Line     int
	Metadata Metadata
	Children []*Node
}

// Metadata is tool-only source metadata. It is useful for browsing and
// filtering source modules, but is never rendered into AGENTS.md.
type Metadata struct {
	Tags          []string `yaml:"tags"`
	TLDR          string   `yaml:"tldr"`
	Priority      *float64 `yaml:"priority"`
	Scope         string   `yaml:"scope"`
	Requires      []string `yaml:"requires"`
	ConflictsWith []string `yaml:"conflicts_with"`
}

// Index permits exact heading-path lookup within one library.
type Index struct {
	Roots  []*Node
	ByPath map[string]*Node
}

// Load recursively reads Markdown files. A source must be a directory, and
// duplicate heading paths are errors rather than arbitrary source precedence.
func Load(root string) (*Index, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("read source %q: %w", root, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("source %q is not a directory", root)
	}
	var files []string
	err = filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
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
	for _, path := range files {
		roots, err := parseFile(path)
		if err != nil {
			return nil, err
		}
		for _, node := range roots {
			if err := index.add(node); err != nil {
				return nil, err
			}
			index.Roots = append(index.Roots, node)
		}
	}
	return index, nil
}

func (i *Index) add(node *Node) error {
	if _, found := i.ByPath[node.Path]; found {
		return fmt.Errorf("duplicate heading path %q", node.Path)
	}
	i.ByPath[node.Path] = node
	for _, child := range node.Children {
		if err := i.add(child); err != nil {
			return err
		}
	}
	return nil
}

type parsedNode struct {
	node  *Node
	level int
}

func parseFile(path string) ([]*Node, error) {
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
		node := &Node{Heading: heading, Path: slug, File: path, Line: lineNumber, Metadata: metadata}
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
	for _, root := range roots {
		if err := applyIDs(root, ""); err != nil {
			return nil, fmt.Errorf("Markdown %q: %w", path, err)
		}
	}
	return roots, nil
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
