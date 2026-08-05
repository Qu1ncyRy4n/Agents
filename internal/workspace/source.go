package workspace

import (
	"fmt"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/internal/library"
	"github.com/Qu1ncyRy4n/Agents/internal/manifest"
)

type SourceNode struct {
	Alias     string
	Reference string
	File      string
	Line      int
	Heading   string
	Metadata  library.Metadata
	Content   string
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
		File:      node.File,
		Line:      node.Line,
		Heading:   node.Heading,
		Metadata:  node.Metadata,
		Content:   sourceContent(node),
	}, nil
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
