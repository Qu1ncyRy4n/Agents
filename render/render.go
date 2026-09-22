// Package render resolves manifests against libraries and renders AGENTS.md.
package render

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/Qu1ncyRy4n/Agents/library"
	"github.com/Qu1ncyRy4n/Agents/manifest"
	"github.com/Qu1ncyRy4n/Agents/renderfs"
	"github.com/Qu1ncyRy4n/Agents/sourcecache"
)

// Result is the fully validated build result. Warnings are non-fatal facts that
// the future confirmation UI must present before it writes an output file.
type Result struct {
	Content  string
	Warnings []string
}

// Build loads all named local sources and renders the document tree.
// Intent: refuse partial or ambiguous documents before a caller can replace an
// existing AGENTS.md. Source: DI-vukam, DI-sufok.
func Build(value *manifest.Manifest, manifestPath string) (*Result, error) {
	sources, err := LoadDocumentSources(value, manifestPath)
	if err != nil {
		return nil, err
	}
	result := &Result{}
	var output strings.Builder
	for _, entry := range value.Doc {
		if err := renderEntry(&output, entry, 1, value.Vars, sources, result); err != nil {
			return nil, err
		}
	}
	trimmed := strings.TrimSpace(output.String())
	if trimmed == "" {
		return nil, fmt.Errorf("rendered output is empty")
	}
	result.Content = trimmed + "\n"
	return result, nil
}

// LoadSources resolves the manifest's named local libraries using the same
// rules as Build. Interactive callers use it to inspect source provenance
// before writing anything.
func LoadSources(value *manifest.Manifest, manifestPath string) (map[string]*library.Index, error) {
	aliases := make([]string, 0, len(value.Sources))
	for alias := range value.Sources {
		aliases = append(aliases, alias)
	}
	return loadSourceAliases(value, manifestPath, aliases)
}

// LoadDocumentSources indexes only aliases that contribute rendered Markdown.
func LoadDocumentSources(value *manifest.Manifest, manifestPath string) (map[string]*library.Index, error) {
	used := documentSourceAliases(value.Doc)
	aliases := make([]string, 0, len(used))
	for alias := range used {
		aliases = append(aliases, alias)
	}
	return loadSourceAliases(value, manifestPath, aliases)
}

// LoadWorkspaceSources indexes ordinary libraries for interactive browsing while
// skipping aliases used only by raw directory outputs. Those outputs may contain
// non-Mogent Markdown such as Agent Skills frontmatter.
func LoadWorkspaceSources(value *manifest.Manifest, manifestPath string) (map[string]*library.Index, error) {
	usedByDocument := documentSourceAliases(value.Doc)
	directoryOnly := make(map[string]bool)
	for _, output := range value.Outputs {
		if !output.Directory() {
			continue
		}
		alias, _, err := manifest.SplitReference(output.From)
		if err == nil && !usedByDocument[alias] {
			directoryOnly[alias] = true
		}
	}
	aliases := make([]string, 0, len(value.Sources))
	for alias := range value.Sources {
		if !directoryOnly[alias] {
			aliases = append(aliases, alias)
		}
	}
	return loadSourceAliases(value, manifestPath, aliases)
}

func loadSourceAliases(value *manifest.Manifest, manifestPath string, aliases []string) (map[string]*library.Index, error) {
	sort.Strings(aliases)
	indexes := make(map[string]*library.Index, len(aliases))
	for _, alias := range aliases {
		source := value.Sources[alias]
		path := source.Location
		if isURL(path) {
			resolved, err := sourcecache.Resolve(manifestPath, alias, path, source.Subdir)
			if err != nil {
				return nil, err
			}
			index, err := library.Load(resolved)
			if err != nil {
				return nil, fmt.Errorf("source %q: %w", alias, err)
			}
			indexes[alias] = index
			continue
		}
		if source.Subdir != "" {
			return nil, fmt.Errorf("source %q: subdir is currently supported only for URL locations", alias)
		}
		path, err := expandHome(path)
		if err != nil {
			return nil, fmt.Errorf("source %q: %w", alias, err)
		}
		if !filepath.IsAbs(path) {
			path = filepath.Join(filepath.Dir(manifestPath), path)
		}
		index, err := library.Load(filepath.Clean(path))
		if err != nil {
			return nil, fmt.Errorf("source %q: %w", alias, err)
		}
		indexes[alias] = index
	}
	return indexes, nil
}

// documentSourceAliases returns only sources whose Markdown headings participate
// in the rendered document. Directory outputs copy raw file trees and must not
// require unrelated Markdown elsewhere in their source roots to parse as a
// Mogent library.
func documentSourceAliases(entries []manifest.Entry) map[string]bool {
	aliases := make(map[string]bool)
	var visit func([]manifest.Entry)
	visit = func(entries []manifest.Entry) {
		for _, entry := range entries {
			for _, reference := range append(append([]string(nil), entry.From...), entry.Exclude...) {
				alias, _, err := manifest.SplitReference(reference)
				if err == nil {
					aliases[alias] = true
				}
			}
			visit(entry.Children)
		}
	}
	visit(entries)
	return aliases
}

// ResolveSourcePath returns the physical root used for a source. Directory
// outputs use this rather than the Markdown index so assets are copied intact.
func ResolveSourcePath(value *manifest.Manifest, manifestPath, alias string) (string, error) {
	source, found := value.Sources[alias]
	if !found {
		return "", fmt.Errorf("reference uses undeclared source %q", alias)
	}
	path := source.Location
	if isURL(path) {
		resolved, err := sourcecache.Resolve(manifestPath, alias, path, source.Subdir)
		if err != nil {
			return "", err
		}
		return resolved, nil
	}
	if source.Subdir != "" {
		return "", fmt.Errorf("source %q: subdir is currently supported only for URL locations", alias)
	}
	path, err := expandHome(path)
	if err != nil {
		return "", err
	}
	if !filepath.IsAbs(path) {
		path = filepath.Join(filepath.Dir(manifestPath), path)
	}
	return filepath.Clean(path), nil
}

func isURL(value string) bool {
	return strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://")
}

func expandHome(path string) (string, error) {
	if path != "~" && !strings.HasPrefix(path, "~/") {
		return path, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home directory: %w", err)
	}
	if path == "~" {
		return home, nil
	}
	return filepath.Join(home, strings.TrimPrefix(path, "~/")), nil
}

func renderEntry(output *strings.Builder, entry manifest.Entry, level int, vars map[string]any, sources map[string]*library.Index, result *Result) error {
	if level > 6 {
		return fmt.Errorf("manifest nesting under %q exceeds Markdown heading level 6", entry.Heading)
	}
	output.WriteString(strings.Repeat("#", level))
	output.WriteByte(' ')
	output.WriteString(entry.Heading)
	output.WriteString("\n\n")
	if len(entry.Children) > 0 {
		for _, child := range entry.Children {
			if err := renderEntry(output, child, level+1, vars, sources, result); err != nil {
				return err
			}
		}
		return nil
	}

	excluded, err := excludedPaths(entry, sources)
	if err != nil {
		return err
	}
	seenRelative := make(map[string]string)
	selectedRoots := make(map[string]string)
	for _, reference := range entry.From {
		alias, path, _ := manifest.SplitReference(reference)
		for selectedReference, selectedPath := range selectedRoots {
			selectedAlias, _, _ := manifest.SplitReference(selectedReference)
			if alias == selectedAlias && (path == selectedPath || strings.HasPrefix(path, selectedPath+"/") || strings.HasPrefix(selectedPath, path+"/")) {
				return fmt.Errorf("entry %q composes overlapping source roots %s and %s", entry.Heading, selectedReference, reference)
			}
		}
		selectedRoots[reference] = path
		node, err := lookup(sources, alias, path)
		if err != nil {
			return fmt.Errorf("entry %q: %w", entry.Heading, err)
		}
		if !hasContent(node, alias, excluded) {
			return fmt.Errorf("entry %q selects empty source node %q", entry.Heading, reference)
		}
		for _, relative := range descendantPaths(node, node.Path) {
			if existing, found := seenRelative[relative]; found {
				return fmt.Errorf("entry %q composes overlapping descendant path %q from %s and %s", entry.Heading, relative, existing, reference)
			} else {
				seenRelative[relative] = reference
			}
		}
		if err := renderNode(output, node, level, alias, vars, excluded); err != nil {
			return fmt.Errorf("entry %q: %w", entry.Heading, err)
		}
	}
	return nil
}

func excludedPaths(entry manifest.Entry, sources map[string]*library.Index) (map[string]bool, error) {
	excluded := make(map[string]bool, len(entry.Exclude))
	for _, reference := range entry.Exclude {
		alias, path, _ := manifest.SplitReference(reference)
		if err := excludePath(excluded, sources, entry.From, alias, path, reference); err != nil {
			return nil, err
		}
	}
	return excluded, nil
}

func excludePath(excluded map[string]bool, sources map[string]*library.Index, from []string, alias, path, reference string) error {
	if _, err := lookup(sources, alias, path); err != nil {
		return fmt.Errorf("exclude %q: %w", reference, err)
	}
	for _, sourceReference := range from {
		sourceAlias, sourcePath, _ := manifest.SplitReference(sourceReference)
		if alias == sourceAlias && path == sourcePath {
			return fmt.Errorf("exclude %q cannot exclude an entry source root", reference)
		}
		if alias == sourceAlias && strings.HasPrefix(path, sourcePath+"/") {
			excluded[alias+":"+path] = true
			return nil
		}
	}
	return fmt.Errorf("exclude %q is not inside an entry source pull", reference)
}

func lookup(sources map[string]*library.Index, alias, path string) (*library.Node, error) {
	index, found := sources[alias]
	if !found {
		return nil, fmt.Errorf("reference uses undeclared source %q", alias)
	}
	node, found := index.ByPath[path]
	if !found {
		return nil, fmt.Errorf("source %q has no heading path %q", alias, path)
	}
	return node, nil
}

func renderNode(output *strings.Builder, node *library.Node, level int, alias string, vars map[string]any, excluded map[string]bool) error {
	if excluded[alias+":"+node.Path] {
		return nil
	}
	body, err := executeTemplate(node.Path, node.Body, vars)
	if err != nil {
		return err
	}
	if strings.TrimSpace(body) != "" {
		output.WriteString(strings.TrimSpace(body))
		output.WriteString("\n\n")
	}
	for _, child := range node.Children {
		if excluded[alias+":"+child.Path] {
			continue
		}
		if level+1 > 6 {
			return fmt.Errorf("source path %q exceeds Markdown heading level 6", child.Path)
		}
		output.WriteString(strings.Repeat("#", level+1))
		output.WriteByte(' ')
		output.WriteString(child.Heading)
		output.WriteString("\n\n")
		if err := renderNode(output, child, level+1, alias, vars, excluded); err != nil {
			return err
		}
	}
	return nil
}

func executeTemplate(name, content string, vars map[string]any) (string, error) {
	tmpl, err := template.New(name).Option("missingkey=error").Parse(content)
	if err != nil {
		return "", fmt.Errorf("parse template %q: %w", name, err)
	}
	var output bytes.Buffer
	if err := tmpl.Execute(&output, vars); err != nil {
		if strings.Contains(err.Error(), "map has no entry for key") {
			return "", fmt.Errorf("render template %q: %w; declare the missing value under manifest vars (init materializes repo_name and repo_url)", name, err)
		}
		return "", fmt.Errorf("render template %q: %w", name, err)
	}
	return output.String(), nil
}

func hasContent(node *library.Node, alias string, excluded map[string]bool) bool {
	if excluded[alias+":"+node.Path] {
		return false
	}
	if strings.TrimSpace(node.Body) != "" {
		return true
	}
	for _, child := range node.Children {
		if hasContent(child, alias, excluded) {
			return true
		}
	}
	return false
}

func descendantPaths(node *library.Node, rootPath string) []string {
	var paths []string
	for _, child := range node.Children {
		relative := strings.TrimPrefix(child.Path, rootPath+"/")
		paths = append(paths, relative)
		paths = append(paths, descendantPaths(child, rootPath)...)
	}
	return paths
}

// WriteAtomically replaces output only after a successful full build.
func WriteAtomically(path, content string) error {
	return renderfs.WriteAtomically(path, []byte(content))
}
