// Package render resolves manifests against libraries and renders AGENTS.md.
package render

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/Qu1ncyRy4n/Agents/internal/library"
	"github.com/Qu1ncyRy4n/Agents/internal/manifest"
	"github.com/Qu1ncyRy4n/Agents/internal/sourcecache"
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
	sources, err := LoadSources(value, manifestPath)
	if err != nil {
		return nil, err
	}
	result := &Result{}
	vars := toolVars(value.Vars, filepath.Dir(manifestPath))
	var output strings.Builder
	for _, entry := range value.Doc {
		if err := renderEntry(&output, entry, 1, vars, sources, result); err != nil {
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

func toolVars(user map[string]any, repository string) map[string]any {
	vars := make(map[string]any, len(user)+2)
	for name, value := range user {
		vars[name] = value
	}
	if _, exists := vars["repo_name"]; !exists {
		vars["repo_name"] = filepath.Base(repository)
	}
	if _, exists := vars["repo_url"]; !exists {
		vars["repo_url"] = ""
		command := exec.Command("git", "-C", repository, "config", "--get", "remote.origin.url")
		if output, err := command.Output(); err == nil {
			vars["repo_url"] = strings.TrimSpace(string(output))
		}
	}
	return vars
}

// LoadSources resolves the manifest's named local libraries using the same
// rules as Build. Interactive callers use it to inspect source provenance
// before writing anything.
func LoadSources(value *manifest.Manifest, manifestPath string) (map[string]*library.Index, error) {
	aliases := make([]string, 0, len(value.Sources))
	for alias := range value.Sources {
		aliases = append(aliases, alias)
	}
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
	directory := filepath.Dir(path)
	temporary, err := os.CreateTemp(directory, ".mogent-*")
	if err != nil {
		return fmt.Errorf("create temporary output: %w", err)
	}
	temporaryName := temporary.Name()
	if _, err := temporary.WriteString(content); err != nil {
		return cleanupTemporary(fmt.Errorf("write temporary output: %w", err), temporary, temporaryName)
	}
	if err := temporary.Chmod(0o644); err != nil {
		return cleanupTemporary(fmt.Errorf("set output permissions: %w", err), temporary, temporaryName)
	}
	if err := temporary.Close(); err != nil {
		if removeErr := os.Remove(temporaryName); removeErr != nil {
			return errors.Join(fmt.Errorf("close temporary output: %w", err), fmt.Errorf("remove temporary output: %w", removeErr))
		}
		return fmt.Errorf("close temporary output: %w", err)
	}
	if err := os.Rename(temporaryName, path); err != nil {
		if removeErr := os.Remove(temporaryName); removeErr != nil {
			return errors.Join(fmt.Errorf("replace output: %w", err), fmt.Errorf("remove temporary output: %w", removeErr))
		}
		return fmt.Errorf("replace output: %w", err)
	}
	return nil
}

func cleanupTemporary(buildErr error, temporary *os.File, path string) error {
	closeErr := temporary.Close()
	removeErr := os.Remove(path)
	if closeErr != nil || removeErr != nil {
		return errors.Join(buildErr, closeErr, removeErr)
	}
	return buildErr
}
