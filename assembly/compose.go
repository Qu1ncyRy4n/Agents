package assembly

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/content"
)

// Compose captures every resolved contribution in one physical manifest. It does
// not write files; callers can inspect or stage the returned manifest first.
func (r *Result) Compose() (*content.Manifest, error) {
	if r.HasErrors() {
		return nil, fmt.Errorf("cannot compose a result with errors")
	}
	var contributions []content.Contribution
	documents := make(map[string]bool)
	for _, contribution := range r.Contributions {
		if contribution.Inline != "" {
			contributions = append(contributions, content.Contribution{
				Rendered:  &content.Rendered{Bytes: normalizeDocument(contribution.Inline)},
				Operation: content.RenderMarkdown,
				Into:      contribution.Into,
				Append:    documents[contribution.Into],
				Replace:   contribution.Replace,
				Origin:    "inline",
			})
			documents[contribution.Into] = true
			continue
		}
		if contribution.Node == nil {
			return nil, fmt.Errorf("resolved contribution has neither inline text nor node")
		}
		switch contribution.Node.Type {
		case Section:
			rendered, err := renderSection(contribution.Node, contribution.Values)
			if err != nil {
				return nil, fmt.Errorf("render %s:%s: %w", contribution.Source, contribution.Node.Path, err)
			}
			contributions = append(contributions, content.Contribution{
				Rendered:  &content.Rendered{Bytes: rendered},
				Operation: content.RenderMarkdown,
				Into:      contribution.Into,
				Append:    documents[contribution.Into],
				Replace:   contribution.Replace,
				Origin:    contribution.Source + ":" + contribution.Node.Path,
			})
			documents[contribution.Into] = true
		case Directory, File, Markdown:
			inventory, err := physicalInventory(contribution.Node)
			if err != nil {
				return nil, err
			}
			node := physicalPath(contribution.Node)
			contributions = append(contributions, content.Contribution{
				Inventory: inventory,
				Node:      node,
				Operation: content.Copy,
				Into:      contribution.Into,
				Replace:   contribution.Replace,
				Origin:    contribution.Source + ":" + contribution.Node.Path,
			})
		default:
			return nil, fmt.Errorf("unsupported node type %q", contribution.Node.Type)
		}
	}
	return content.Compose(contributions)
}

func physicalInventory(node *Node) (*content.Inventory, error) {
	root := node
	for root.Parent != nil && root.Source == "" {
		root = root.Parent
	}
	if root.Physical != nil {
		return root.Physical, nil
	}
	if root.Source == "" {
		return nil, fmt.Errorf("physical node %q has no root inventory", node.Path)
	}
	path, err := safePath(filepath.Dir(root.Range.Filename), root.Source)
	if err != nil {
		return nil, err
	}
	return content.Discover(path)
}

func physicalPath(node *Node) string {
	root := node
	for root.Parent != nil && root.Source == "" {
		root = root.Parent
	}
	if node == root {
		return "."
	}
	return strings.TrimPrefix(node.Path, root.Path+"/")
}

func renderSection(node *Node, values map[string]string) ([]byte, error) {
	var out bytes.Buffer
	if node.Source != "" {
		path, err := safePath(filepath.Dir(node.Range.Filename), node.Source)
		if err != nil {
			return nil, err
		}
		contents, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		rendered, err := substituteGuidance(string(contents), node.Variables, values)
		if err != nil {
			return nil, err
		}
		doc, err := content.ParseDocument([]byte(rendered))
		if err != nil {
			return nil, err
		}
		view, err := doc.Render(node.Heading, 1)
		if err != nil {
			return nil, err
		}
		out.Write(view.Bytes)
	}
	for _, child := range node.Children {
		part, err := renderSection(child, values)
		if err != nil {
			return nil, err
		}
		if len(out.Bytes()) > 0 && len(part) > 0 {
			out.WriteByte('\n')
		}
		out.Write(part)
	}
	return normalizeDocument(out.String()), nil
}

func substituteGuidance(source string, variables []string, values map[string]string) (string, error) {
	allowed := make(map[string]bool, len(variables))
	for _, variable := range variables {
		allowed[variable] = true
	}
	for name := range allowed {
		if _, found := values[name]; !found {
			return "", fmt.Errorf("missing declared guidance value %q", name)
		}
	}
	for name, value := range values {
		if !allowed[name] {
			return "", fmt.Errorf("undeclared guidance value %q", name)
		}
		source = strings.ReplaceAll(source, "${"+name+"}", value)
	}
	return strings.ReplaceAll(source, `\${`, "${"), nil
}

func normalizeDocument(value string) []byte {
	trimmed := strings.TrimSpace(strings.ReplaceAll(value, "\r\n", "\n"))
	if trimmed == "" {
		return nil
	}
	return []byte(trimmed + "\n")
}
