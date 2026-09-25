// Package render resolves manifests against libraries and renders AGENTS.md.
package render

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/Qu1ncyRy4n/Agents/library"
	"github.com/Qu1ncyRy4n/Agents/manifest"
	"github.com/Qu1ncyRy4n/Agents/renderfs"
	"github.com/Qu1ncyRy4n/Agents/sourcecache"
	"github.com/Qu1ncyRy4n/Agents/sourcepath"
)

// Result is the fully validated build result. Warnings are non-fatal facts that
// the future confirmation UI must present before it writes an output file.
type Result struct {
	Content  string
	Warnings []string
}

// Options controls Markdown-only rendering behavior. Raw directory outputs are
// copied by workspace unchanged.
type Options struct{ PreserveHTMLComments bool }

// Build loads all named local sources and renders the document tree.
// Intent: refuse partial or ambiguous documents before a caller can replace an
// existing AGENTS.md. Source: DI-vukam, DI-sufok.
func Build(value *manifest.Manifest, manifestPath string) (*Result, error) {
	return BuildWithOptions(value, manifestPath, Options{})
}

// BuildWithOptions renders a legacy document outline with explicit options.
func BuildWithOptions(value *manifest.Manifest, manifestPath string, options Options) (*Result, error) {
	sources, err := LoadDocumentSources(value, manifestPath)
	if err != nil {
		return nil, err
	}
	result := &Result{}
	var output strings.Builder
	for _, entry := range value.Doc {
		if err := renderEntry(&output, entry, 1, value.Vars, sources, result, options); err != nil {
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
	directoryOnly := value.RawDirectorySourceAliases()
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
		path, err := ResolveSourcePath(value, manifestPath, alias)
		if err != nil {
			return nil, fmt.Errorf("source %q: %w", alias, err)
		}
		index, _, err := library.LoadWithSidecar(filepath.Clean(path))
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
	return manifest.EntrySourceAliases(entries)
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
	if !filepath.IsAbs(path) {
		path = filepath.Join(filepath.Dir(manifestPath), path)
	}
	return sourcepath.Resolve(filepath.Clean(path), source.Subdir)
}

func isURL(value string) bool {
	return strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://")
}

func renderEntry(output *strings.Builder, entry manifest.Entry, level int, vars map[string]any, sources map[string]*library.Index, result *Result, options Options) error {
	if level > 6 {
		return fmt.Errorf("manifest nesting under %q exceeds Markdown heading level 6", entry.Heading)
	}
	output.WriteString(strings.Repeat("#", level))
	output.WriteByte(' ')
	output.WriteString(entry.Heading)
	output.WriteString("\n\n")
	if len(entry.Children) > 0 {
		for _, child := range entry.Children {
			if err := renderEntry(output, child, level+1, vars, sources, result, options); err != nil {
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
		if err := renderNode(output, node, level, alias, vars, excluded, options); err != nil {
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

func renderNode(output *strings.Builder, node *library.Node, level int, alias string, vars map[string]any, excluded map[string]bool, options Options) error {
	if excluded[alias+":"+node.Path] {
		return nil
	}
	body, err := executeTemplate(node.Path, node.Body, vars)
	if err != nil {
		return err
	}
	if strings.TrimSpace(body) != "" {
		if !options.PreserveHTMLComments {
			body = stripHTMLComments(body)
		}
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
		if err := renderNode(output, child, level+1, alias, vars, excluded, options); err != nil {
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

// RenderOutput renders one canonical Markdown output. Outputs without selectors
// retain the legacy authored document outline.
func RenderOutput(value *manifest.Manifest, manifestPath string, output manifest.Output) (*Result, error) {
	return RenderOutputWithOptions(value, manifestPath, output, Options{})
}

// RenderOutputWithOptions renders one canonical Markdown output with explicit options.
func RenderOutputWithOptions(value *manifest.Manifest, manifestPath string, output manifest.Output, options Options) (*Result, error) {
	if len(output.Include) == 0 {
		return BuildWithOptions(value, manifestPath, options)
	}
	sources, err := loadSelectorSources(value, manifestPath, output)
	if err != nil {
		return nil, err
	}
	selected, err := selectorNodes(output.Include, sources)
	if err != nil {
		return nil, err
	}
	excluded, err := selectorPaths(output.Exclude, sources)
	if err != nil {
		return nil, err
	}
	var text strings.Builder
	for _, item := range selected {
		if excluded[item.alias+":"+item.node.Path] {
			continue
		}
		text.WriteString("# " + item.node.Heading + "\n\n")
		if err := renderNode(&text, item.node, 1, item.alias, value.Vars, excluded, options); err != nil {
			return nil, err
		}
	}
	content := strings.TrimSpace(text.String())
	if content == "" {
		return nil, fmt.Errorf("output %q selects no renderable content", output.Path)
	}
	return &Result{Content: content + "\n"}, nil
}

// stripHTMLComments removes standard comments outside fenced code blocks while
// preserving source line structure well enough for ordinary Markdown spacing.
func stripHTMLComments(content string) string {
	var output strings.Builder
	inComment := false
	var fence library.Fence
	inFence := false
	for _, line := range strings.SplitAfter(content, "\n") {
		bare := strings.TrimSuffix(line, "\n")
		if inFence {
			output.WriteString(line)
			if fence.Closes(bare) {
				inFence = false
			}
			continue
		}
		if opened, ok := library.OpenFence(bare); ok {
			fence, inFence = opened, true
			output.WriteString(line)
			continue
		}
		remaining := bare
		for {
			if inComment {
				_, after, found := strings.Cut(remaining, "-->")
				if !found {
					remaining = ""
					break
				}
				inComment = false
				remaining = after
				continue
			}
			before, after, found := strings.Cut(remaining, "<!--")
			if !found {
				output.WriteString(remaining)
				break
			}
			output.WriteString(before)
			inComment = true
			remaining = after
		}
		if strings.HasSuffix(line, "\n") {
			output.WriteByte('\n')
		}
	}
	return output.String()
}

type selectorNode struct {
	alias string
	node  *library.Node
}

func loadSelectorSources(value *manifest.Manifest, manifestPath string, output manifest.Output) (map[string]*library.Index, error) {
	used := map[string]bool{}
	for _, selector := range append(append([]manifest.Selector(nil), output.Include...), output.Exclude...) {
		if selector.All != "" {
			used[selector.All] = true
		}
		if selector.Source != "" {
			alias, _, _ := manifest.SplitReference(selector.Source)
			used[alias] = true
		}
		if len(selector.Tags) > 0 {
			for alias := range value.Sources {
				used[alias] = true
			}
		}
	}
	aliases := make([]string, 0, len(used))
	for alias := range used {
		aliases = append(aliases, alias)
	}
	return loadSourceAliases(value, manifestPath, aliases)
}

func selectorNodes(selectors []manifest.Selector, sources map[string]*library.Index) ([]selectorNode, error) {
	var result []selectorNode
	for _, selector := range selectors {
		if selector.All != "" {
			index, found := sources[selector.All]
			if !found {
				return nil, fmt.Errorf("selector uses undeclared source %q", selector.All)
			}
			for _, node := range index.Roots {
				result = appendSelectorNode(result, selectorNode{selector.All, node})
			}
			continue
		}
		if selector.Source != "" {
			alias, path, _ := manifest.SplitReference(selector.Source)
			node, err := lookup(sources, alias, path)
			if err != nil {
				return nil, err
			}
			result = appendSelectorNode(result, selectorNode{alias, node})
			continue
		}
		aliases := make([]string, 0, len(sources))
		for alias := range sources {
			aliases = append(aliases, alias)
		}
		sort.Strings(aliases)
		for _, alias := range aliases {
			index := sources[alias]
			paths := make([]string, 0, len(index.ByPath))
			for path := range index.ByPath {
				paths = append(paths, path)
			}
			sort.Strings(paths)
			for _, path := range paths {
				node := index.ByPath[path]
				if tagsAny(node.Metadata.Tags, selector.Tags) {
					result = appendSelectorNode(result, selectorNode{alias, node})
				}
			}
		}
	}
	return result, nil
}

func appendSelectorNode(nodes []selectorNode, candidate selectorNode) []selectorNode {
	for _, existing := range nodes {
		if existing.alias == candidate.alias && (existing.node.Path == candidate.node.Path || strings.HasPrefix(candidate.node.Path, existing.node.Path+"/")) {
			return nodes
		}
	}
	kept := nodes[:0]
	for _, existing := range nodes {
		if existing.alias == candidate.alias && strings.HasPrefix(existing.node.Path, candidate.node.Path+"/") {
			continue
		}
		kept = append(kept, existing)
	}
	return append(kept, candidate)
}

func selectorPaths(selectors []manifest.Selector, sources map[string]*library.Index) (map[string]bool, error) {
	selected, err := selectorNodes(selectors, sources)
	if err != nil {
		return nil, err
	}
	paths := map[string]bool{}
	for _, item := range selected {
		markSelectorSubtree(paths, item.alias, item.node)
	}
	return paths, nil
}

func markSelectorSubtree(paths map[string]bool, alias string, node *library.Node) {
	paths[alias+":"+node.Path] = true
	for _, child := range node.Children {
		markSelectorSubtree(paths, alias, child)
	}
}
func tagsAny(have, wanted []string) bool {
	for _, tag := range wanted {
		for _, candidate := range have {
			if candidate == tag {
				return true
			}
		}
	}
	return false
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
