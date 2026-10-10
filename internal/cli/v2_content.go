package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/Qu1ncyRy4n/Agents/content"
	"github.com/Qu1ncyRy4n/Agents/v2"
)

func writeTypedSource(out *printer, library *v2.Library, alias, reference string, headings bool, lines int) (bool, error) {
	root, relative, hasPath := strings.Cut(reference, "/")
	inventory := library.Contents[root]
	if inventory == nil {
		return false, nil
	}
	if library.ByPath[reference] != nil {
		return true, fmt.Errorf("ambiguous authored/content reference %q; use %s:%s/. for the physical root", reference, alias, root)
	}
	if !hasPath {
		relative = "."
	}
	node := inventory.Nodes[relative]
	if node == nil {
		return true, fmt.Errorf("content %q has no node %q", root, relative)
	}
	out.printf("%s:%s\nKind:   %s\nSource: %s\n", alias, reference, node.Kind, node.Source)
	if node.Role != "" {
		out.printf("Role:   %s\n", node.Role)
	}
	if node.Bundle != "" {
		out.printf("Bundle: %s\n", node.Bundle)
	}
	if headings {
		document, err := inventory.Document(relative)
		if err != nil {
			return true, err
		}
		out.line("Heading addresses:")
		writeHeadingAddresses(out, document.Headings, nil)
		return true, out.flush()
	}
	if node.Kind == content.Directory {
		for _, child := range node.Children {
			out.printf("  %s [%s]\n", child, inventory.Nodes[child].Kind)
		}
		return true, out.flush()
	}
	bytes, err := os.ReadFile(node.Source)
	if err != nil {
		return true, err
	}
	if !utf8.Valid(bytes) || strings.ContainsRune(string(bytes), 0) {
		out.printf("Payload: opaque (%d bytes)\n", len(bytes))
		return true, out.flush()
	}
	body := strings.Split(strings.TrimRight(string(bytes), "\n"), "\n")
	if lines > 0 && len(body) > lines {
		body = append(body[:lines], fmt.Sprintf("... (%d more lines)", len(body)-lines))
	}
	out.line("---")
	for _, row := range body {
		out.line(row)
	}
	return true, out.flush()
}

func writeHeadingAddresses(out *printer, headings []*content.Heading, prefix []string) {
	for _, heading := range headings {
		address := append(append([]string(nil), prefix...), heading.Title)
		quoted := make([]string, len(address))
		for index, title := range address {
			quoted[index] = strconv.Quote(title)
		}
		out.printf("  heading = [%s]\n", strings.Join(quoted, ", "))
		writeHeadingAddresses(out, heading.Children, address)
	}
}
