package workspace

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/library"
	"github.com/Qu1ncyRy4n/Agents/manifest"
)

type SourceNode struct {
	Kind      library.NodeKind
	Alias     string
	Reference string
	Path      string
	File      string
	Line      int
	Heading   string
	Metadata  library.Metadata
	Content   string
}

type SourceListOptions struct {
	SourceAlias string
	Tag         string
	TagSearch   string
	Search      string
	Sort        string
}

// SourceNode returns one source node by explicit alias:path reference.
func (s *Session) SourceNode(reference string) (*SourceNode, error) {
	alias, path, err := manifest.SplitReference(reference)
	if err != nil {
		return nil, err
	}
	index, found := s.Sources[alias]
	if !found {
		return nil, s.undeclaredSourceError(alias)
	}
	node, found := index.ByPath[path]
	if !found {
		return nil, missingSourcePathError(alias, path, index)
	}
	return &SourceNode{
		Kind:      node.Kind,
		Alias:     alias,
		Reference: alias + ":" + path,
		Path:      path,
		File:      node.File,
		Line:      node.Line,
		Heading:   node.Heading,
		Metadata:  node.Metadata,
		Content:   sourceContent(node),
	}, nil
}

// SourceNodes lists source nodes across the loaded manifest's declared sources.
func (s *Session) SourceNodes(options SourceListOptions) ([]SourceNode, error) {
	if options.SourceAlias != "" {
		if _, found := s.Sources[options.SourceAlias]; !found {
			return nil, s.undeclaredSourceError(options.SourceAlias)
		}
	}
	var nodes []SourceNode
	aliases := make([]string, 0, len(s.Sources))
	for alias := range s.Sources {
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)
	for _, alias := range aliases {
		if options.SourceAlias != "" && alias != options.SourceAlias {
			continue
		}
		index := s.Sources[alias]
		paths := make([]string, 0, len(index.ByPath))
		for path := range index.ByPath {
			paths = append(paths, path)
		}
		sort.Strings(paths)
		for _, path := range paths {
			node := index.ByPath[path]
			if options.Tag != "" && !hasExactTag(node.Metadata.Tags, options.Tag) {
				continue
			}
			if options.TagSearch != "" && !hasMatchingTag(node.Metadata.Tags, options.TagSearch) {
				continue
			}
			if options.Search != "" && !matchesSourceSearch(alias, path, node, options.Search) {
				continue
			}
			nodes = append(nodes, SourceNode{
				Kind:      node.Kind,
				Alias:     alias,
				Reference: alias + ":" + path,
				Path:      path,
				File:      node.File,
				Line:      node.Line,
				Heading:   node.Heading,
				Metadata:  node.Metadata,
				Content:   sourceContent(node),
			})
		}
	}
	switch options.Sort {
	case "", "path":
		sort.SliceStable(nodes, func(first, second int) bool {
			return nodes[first].Reference < nodes[second].Reference
		})
	case "priority":
		sort.SliceStable(nodes, func(first, second int) bool {
			firstPriority := priorityValue(nodes[first].Metadata.Priority)
			secondPriority := priorityValue(nodes[second].Metadata.Priority)
			if firstPriority == secondPriority {
				return nodes[first].Reference < nodes[second].Reference
			}
			return firstPriority > secondPriority
		})
	default:
		return nil, fmt.Errorf("unknown source list sort %q; use path or priority", options.Sort)
	}
	return nodes, nil
}

func (s *Session) undeclaredSourceError(alias string) error {
	aliases := s.SavedSourceAliases()
	if suggestion := closestAlias(alias, aliases); suggestion != "" {
		return fmt.Errorf("source %q is not declared in manifest; did you mean %q?", alias, suggestion)
	}
	return fmt.Errorf("source %q is not declared in manifest", alias)
}

func closestAlias(alias string, aliases []string) string {
	best := ""
	bestDistance := 4
	for _, candidate := range aliases {
		distance := levenshtein(strings.ToLower(alias), strings.ToLower(candidate))
		if distance < bestDistance {
			best = candidate
			bestDistance = distance
		}
	}
	return best
}

func (s *Session) SourceReferences(sourceAlias string) ([]string, error) {
	nodes, err := s.SourceNodes(SourceListOptions{SourceAlias: sourceAlias, Sort: "path"})
	if err != nil {
		return nil, err
	}
	references := make([]string, 0, len(nodes))
	for _, node := range nodes {
		references = append(references, node.Reference)
	}
	return references, nil
}

func (s *Session) ManifestHeadingPaths() []string {
	return manifestHeadingPaths(s.Draft.Doc, nil)
}

func manifestHeadingPaths(entries []manifest.Entry, parents []string) []string {
	var paths []string
	for _, entry := range entries {
		current := append(append([]string(nil), parents...), entry.Heading)
		paths = append(paths, strings.Join(current, "/"))
		paths = append(paths, manifestHeadingPaths(entry.Children, current)...)
	}
	return paths
}

func hasExactTag(tags []string, target string) bool {
	for _, tag := range tags {
		if tag == target {
			return true
		}
	}
	return false
}

func hasMatchingTag(tags []string, search string) bool {
	search = strings.ToLower(search)
	for _, tag := range tags {
		if strings.Contains(strings.ToLower(tag), search) {
			return true
		}
	}
	return false
}

func matchesSourceSearch(alias string, path string, node *library.Node, search string) bool {
	search = strings.ToLower(strings.TrimSpace(search))
	if search == "" {
		return true
	}
	haystack := strings.ToLower(strings.Join([]string{
		alias + ":" + path,
		path,
		node.Heading,
		node.Metadata.TLDR,
		strings.Join(node.Metadata.Tags, " "),
		strings.TrimSpace(node.Body),
	}, "\n"))
	return strings.Contains(haystack, search)
}

func missingSourcePathError(alias string, path string, index *library.Index) error {
	descendants := descendantSuggestions(path, index, 8)
	if len(descendants) > 0 {
		return fmt.Errorf("source %q has no heading path %q; it looks like an organizational path. Try a descendant such as %s", alias, path, strings.Join(descendants, ", "))
	}
	close := closePathSuggestions(path, index, 5)
	if len(close) > 0 {
		return fmt.Errorf("source %q has no heading path %q; did you mean %s?", alias, path, strings.Join(close, ", "))
	}
	return fmt.Errorf("source %q has no heading path %q", alias, path)
}

func descendantSuggestions(path string, index *library.Index, limit int) []string {
	prefix := strings.TrimSuffix(path, "/") + "/"
	var suggestions []string
	for _, candidate := range sortedPaths(index) {
		if strings.HasPrefix(candidate, prefix) {
			suggestions = append(suggestions, candidate)
			if len(suggestions) >= limit {
				break
			}
		}
	}
	return suggestions
}

func closePathSuggestions(path string, index *library.Index, limit int) []string {
	type scoredPath struct {
		path  string
		score int
	}
	var scored []scoredPath
	for _, candidate := range sortedPaths(index) {
		score := levenshtein(strings.ToLower(path), strings.ToLower(candidate))
		if strings.Contains(candidate, path) || strings.Contains(path, candidate) {
			score--
		}
		scored = append(scored, scoredPath{path: candidate, score: score})
	}
	sort.SliceStable(scored, func(first, second int) bool {
		if scored[first].score == scored[second].score {
			return scored[first].path < scored[second].path
		}
		return scored[first].score < scored[second].score
	})
	var suggestions []string
	for _, candidate := range scored {
		if candidate.score > 8 && !strings.Contains(candidate.path, path) {
			continue
		}
		suggestions = append(suggestions, candidate.path)
		if len(suggestions) >= limit {
			break
		}
	}
	return suggestions
}

func levenshtein(first string, second string) int {
	if first == second {
		return 0
	}
	if first == "" {
		return len([]rune(second))
	}
	if second == "" {
		return len([]rune(first))
	}
	a := []rune(first)
	b := []rune(second)
	previous := make([]int, len(b)+1)
	current := make([]int, len(b)+1)
	for index := range previous {
		previous[index] = index
	}
	for i, ar := range a {
		current[0] = i + 1
		for j, br := range b {
			cost := 1
			if ar == br {
				cost = 0
			}
			current[j+1] = minInt(current[j]+1, previous[j+1]+1, previous[j]+cost)
		}
		previous, current = current, previous
	}
	return previous[len(b)]
}

func minInt(values ...int) int {
	best := values[0]
	for _, value := range values[1:] {
		if value < best {
			best = value
		}
	}
	return best
}

func priorityValue(value *float64) float64 {
	if value == nil {
		return -1
	}
	return *value
}

func sourceContent(node *library.Node) string {
	var output strings.Builder
	writeSourceContent(&output, node, 0)
	return strings.TrimSpace(output.String())
}

func writeSourceContent(output *strings.Builder, node *library.Node, depth int) {
	if depth > 0 {
		output.WriteString(strings.Repeat("#", depth+1))
		output.WriteByte(' ')
		output.WriteString(node.Heading)
		output.WriteString("\n\n")
	}
	body := strings.TrimSpace(node.Body)
	if body != "" {
		output.WriteString(body)
		output.WriteString("\n\n")
	}
	for _, child := range node.Children {
		writeSourceContent(output, child, depth+1)
	}
}

func Snippet(content string, lines int) string {
	if lines <= 0 {
		return ""
	}
	parts := strings.Split(strings.TrimSpace(content), "\n")
	if len(parts) <= lines {
		return strings.Join(parts, "\n")
	}
	return strings.Join(parts[:lines], "\n")
}
