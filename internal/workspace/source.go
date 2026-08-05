package workspace

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/internal/library"
	"github.com/Qu1ncyRy4n/Agents/internal/manifest"
)

type SourceNode struct {
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
		return nil, fmt.Errorf("reference uses undeclared source %q", alias)
	}
	node, found := index.ByPath[path]
	if !found {
		return nil, fmt.Errorf("source %q has no heading path %q", alias, path)
	}
	return &SourceNode{
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
			return nil, fmt.Errorf("source %q is not declared in manifest", options.SourceAlias)
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
			nodes = append(nodes, SourceNode{
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
