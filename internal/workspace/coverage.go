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
	Tag         string
	ContentOnly bool
	LeavesOnly  bool
	LimitDepth  bool
	MaxDepth    int
}

type SourceCoverage struct {
	Alias    string
	Path     string
	Total    int
	Included int
	Nodes    []CoverageNode
	Unused   []CoverageNode
}

type CoverageNode struct {
	Path    string
	Heading string
	File    string
	TLDR    string
	Depth   int
	State   CoverageState
}

type CoverageState string

const (
	CoverageUnused    CoverageState = "unused"
	CoverageIncluded  CoverageState = "included"
	CoverageInherited CoverageState = "inherited"
	CoverageExcluded  CoverageState = "excluded"
	CoveragePartial   CoverageState = "partial"
)

// Coverage reports which source nodes are included by the current draft and
// which source nodes are available but unused.
func (s *Session) Coverage() Coverage {
	return s.CoverageWithOptions(CoverageOptions{})
}

func (s *Session) CoverageWithOptions(options CoverageOptions) Coverage {
	states := make(map[string]CoverageState)
	markEntries(s.Draft.Doc, s.Sources, states)
	markPartialAncestors(states)

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
			Path:  s.Saved.Sources[alias].Display(),
		}
		for _, path := range paths {
			node := index.ByPath[path]
			depth := pathDepth(path)
			if options.ContentOnly && strings.TrimSpace(node.Body) == "" {
				continue
			}
			if options.Tag != "" && !hasTag(node.Metadata.Tags, options.Tag) {
				continue
			}
			if options.LeavesOnly && len(node.Children) > 0 {
				continue
			}
			if options.LimitDepth && depth > options.MaxDepth {
				continue
			}
			source.Total++
			state := states[alias+":"+path]
			if state == "" {
				state = CoverageUnused
			}
			coverageNode := CoverageNode{
				Path:    path,
				Heading: node.Heading,
				File:    node.File,
				TLDR:    node.Metadata.TLDR,
				Depth:   depth,
				State:   state,
			}
			source.Nodes = append(source.Nodes, coverageNode)
			if state == CoverageIncluded || state == CoverageInherited {
				source.Included++
				continue
			}
			if state == CoverageUnused {
				source.Unused = append(source.Unused, coverageNode)
			}
		}
		result.Sources = append(result.Sources, source)
	}
	return result
}

func markPartialAncestors(states map[string]CoverageState) {
	for reference, state := range cloneCoverageStates(states) {
		if state != CoverageIncluded && state != CoverageInherited && state != CoverageExcluded {
			continue
		}
		alias, path, found := strings.Cut(reference, ":")
		if !found {
			continue
		}
		parts := strings.Split(path, "/")
		for index := 1; index < len(parts); index++ {
			parentReference := alias + ":" + strings.Join(parts[:index], "/")
			if states[parentReference] == "" {
				states[parentReference] = CoveragePartial
			}
		}
	}
}

func cloneCoverageStates(states map[string]CoverageState) map[string]CoverageState {
	clone := make(map[string]CoverageState, len(states))
	for key, value := range states {
		clone[key] = value
	}
	return clone
}

func hasTag(tags []string, target string) bool {
	for _, tag := range tags {
		if tag == target {
			return true
		}
	}
	return false
}

func markEntries(entries []manifest.Entry, sources map[string]*library.Index, states map[string]CoverageState) {
	for _, entry := range entries {
		if len(entry.Children) > 0 {
			markEntries(entry.Children, sources, states)
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
			markNode(alias, node, excluded, states, true)
		}
	}
}

func markNode(alias string, node *library.Node, excluded map[string]bool, states map[string]CoverageState, direct bool) {
	reference := alias + ":" + node.Path
	if excluded[reference] {
		markExcludedNode(alias, node, states)
		return
	}
	if direct {
		states[reference] = CoverageIncluded
	} else {
		states[reference] = CoverageInherited
	}
	for _, child := range node.Children {
		markNode(alias, child, excluded, states, false)
	}
}

func markExcludedNode(alias string, node *library.Node, states map[string]CoverageState) {
	states[alias+":"+node.Path] = CoverageExcluded
	for _, child := range node.Children {
		markExcludedNode(alias, child, states)
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
