package workspace

import (
	"sort"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/internal/library"
	"github.com/Qu1ncyRy4n/Agents/internal/manifest"
)

type Coverage struct {
	Sources []SourceCoverage
}

type CoverageOptions struct {
	SourceAlias string
	ContentOnly bool
}

type SourceCoverage struct {
	Alias    string
	Path     string
	Total    int
	Included int
	Unused   []CoverageNode
}

type CoverageNode struct {
	Path    string
	Heading string
	Depth   int
}

// Coverage reports which source nodes are included by the current draft and
// which source nodes are available but unused.
func (s *Session) Coverage() Coverage {
	return s.CoverageWithOptions(CoverageOptions{})
}

func (s *Session) CoverageWithOptions(options CoverageOptions) Coverage {
	included := make(map[string]bool)
	markEntries(s.Draft.Doc, s.Sources, included)

	aliases := make([]string, 0, len(s.Sources))
	for alias := range s.Sources {
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)

	result := Coverage{Sources: make([]SourceCoverage, 0, len(aliases))}
	for _, alias := range aliases {
		if options.SourceAlias != "" && alias != options.SourceAlias {
			continue
		}
		index := s.Sources[alias]
		paths := sortedPaths(index)
		source := SourceCoverage{
			Alias: alias,
			Path:  s.Saved.Sources[alias],
			Total: len(paths),
		}
		for _, path := range paths {
			node := index.ByPath[path]
			if options.ContentOnly && strings.TrimSpace(node.Body) == "" {
				continue
			}
			if included[alias+":"+path] {
				source.Included++
				continue
			}
			source.Unused = append(source.Unused, CoverageNode{
				Path:    path,
				Heading: node.Heading,
				Depth:   pathDepth(path),
			})
		}
		result.Sources = append(result.Sources, source)
	}
	return result
}

func markEntries(entries []manifest.Entry, sources map[string]*library.Index, included map[string]bool) {
	for _, entry := range entries {
		if len(entry.Children) > 0 {
			markEntries(entry.Children, sources, included)
			continue
		}
		excluded := excludedReferences(entry.Exclude)
		for _, reference := range entry.From {
			alias, path, err := manifest.SplitReference(reference)
			if err != nil {
				continue
			}
			index, found := sources[alias]
			if !found {
				continue
			}
			node, found := index.ByPath[path]
			if !found {
				continue
			}
			markNode(alias, node, excluded, included)
		}
	}
}

func markNode(alias string, node *library.Node, excluded map[string]bool, included map[string]bool) {
	reference := alias + ":" + node.Path
	if excluded[reference] {
		return
	}
	included[reference] = true
	for _, child := range node.Children {
		markNode(alias, child, excluded, included)
	}
}

func excludedReferences(references []string) map[string]bool {
	excluded := make(map[string]bool, len(references))
	for _, reference := range references {
		excluded[reference] = true
	}
	return excluded
}

func sortedPaths(index *library.Index) []string {
	paths := make([]string, 0, len(index.ByPath))
	for path := range index.ByPath {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths
}

func pathDepth(path string) int {
	depth := 0
	for _, char := range path {
		if char == '/' {
			depth++
		}
	}
	return depth
}
