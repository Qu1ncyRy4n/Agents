package content

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// Heading describes a parsed Markdown subtree. Byte spans address the original
// file, while Titles are display text, not globally unique IDs.
type Heading struct {
	Kind      Kind
	Title     string
	Level     int
	Start     int
	End       int
	HeaderEnd int
	Children  []*Heading
}

type Document struct {
	Source    []byte
	Headings  []*Heading
	bodyStart int
	ordered   []*Heading
}

type Rendered struct {
	Bytes      []byte
	References []string
}

// ParseDocument builds a CommonMark heading view, ignoring fenced code and
// keeping leading YAML frontmatter outside the document body. Copy callers
// retain Source unchanged. Unterminated frontmatter is an explicit error.
func ParseDocument(source []byte) (*Document, error) {
	if !utf8.Valid(source) {
		return nil, fmt.Errorf("Markdown document is not UTF-8")
	}
	doc := &Document{Source: append([]byte(nil), source...)}
	if first := line(source, 0); strings.TrimSuffix(first, "\r") == "---" {
		position := lineEnd(source, 0)
		closed := false
		for position < len(source) {
			value := strings.TrimSuffix(line(source, position), "\r")
			end := lineEnd(source, position)
			if value == "---" || value == "..." {
				doc.bodyStart = end
				closed = true
				break
			}
			position = end
		}
		if !closed {
			return nil, fmt.Errorf("unterminated Markdown frontmatter")
		}
	}
	body := doc.Source[doc.bodyStart:]
	root := goldmark.DefaultParser().Parse(text.NewReader(body))
	var stack []*Heading
	err := ast.Walk(root, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		h, ok := node.(*ast.Heading)
		if !entering || !ok {
			return ast.WalkContinue, nil
		}
		start := h.Lines().At(0).Start
		for start > 0 && body[start-1] != '\n' {
			start--
		}
		last := h.Lines().At(h.Lines().Len() - 1)
		end := lineEnd(body, last.Stop)
		// Setext heading segments contain the text but not the underline.
		if !isATX(line(body, start)) && end < len(body) {
			end = lineEnd(body, end)
		}
		heading := &Heading{Kind: HeadingKind, Title: strings.TrimSpace(string(h.Text(body))), Level: h.Level, Start: doc.bodyStart + start, End: len(doc.Source), HeaderEnd: doc.bodyStart + end}
		for len(stack) > 0 && stack[len(stack)-1].Level >= heading.Level {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			doc.Headings = append(doc.Headings, heading)
		} else {
			parent := stack[len(stack)-1]
			parent.Children = append(parent.Children, heading)
		}
		stack = append(stack, heading)
		doc.ordered = append(doc.ordered, heading)
		return ast.WalkContinue, nil
	})
	if err != nil {
		return nil, err
	}
	for index, heading := range doc.ordered {
		for _, next := range doc.ordered[index+1:] {
			if next.Level <= heading.Level {
				heading.End = next.Start
				break
			}
		}
	}
	return doc, nil
}

// Resolve matches exact title components within this file's document view.
// Duplicate sibling titles are diagnosed instead of choosing an occurrence.
func (d *Document) Resolve(titles []string) (*Heading, error) {
	if len(titles) == 0 {
		return nil, fmt.Errorf("heading address must have at least one title")
	}
	candidates := d.Headings
	var selected *Heading
	for _, title := range titles {
		selected = nil
		for _, candidate := range candidates {
			if candidate.Title == title {
				if selected != nil {
					return nil, fmt.Errorf("ambiguous heading %q in address %q", title, strings.Join(titles, " / "))
				}
				selected = candidate
			}
		}
		if selected == nil {
			return nil, fmt.Errorf("heading %q not found in address %q", title, strings.Join(titles, " / "))
		}
		candidates = selected.Children
	}
	return selected, nil
}

// Render explicitly renders one subtree at level, or the complete document
// body when titles is empty. It removes frontmatter from the derivative only,
// rebases real headings, and reports link/image references for placement checks.
func (d *Document) Render(titles []string, level int) (Rendered, error) {
	if level < 1 || level > 6 {
		return Rendered{}, fmt.Errorf("render heading level must be between 1 and 6")
	}
	start, end, base := d.bodyStart, len(d.Source), 1
	if len(titles) > 0 {
		h, err := d.Resolve(titles)
		if err != nil {
			return Rendered{}, err
		}
		start, end, base = h.Start, h.End, h.Level
	} else if len(d.ordered) > 0 {
		base = 6
		for _, h := range d.ordered {
			if h.Level < base {
				base = h.Level
			}
		}
	}
	var out bytes.Buffer
	position := start
	for _, h := range d.ordered {
		if h.Start < start || h.Start >= end {
			continue
		}
		depth := h.Level - base + level
		if depth > 6 {
			return Rendered{}, fmt.Errorf("heading %q exceeds level 6 after rendering", h.Title)
		}
		out.Write(d.Source[position:h.Start])
		out.WriteString(strings.Repeat("#", depth) + " " + h.Title + "\n")
		position = h.HeaderEnd
	}
	out.Write(d.Source[position:end])
	rendered := Rendered{Bytes: []byte(strings.TrimSpace(out.String()) + "\n")}
	root := goldmark.DefaultParser().Parse(text.NewReader(rendered.Bytes))
	err := ast.Walk(root, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			switch n := node.(type) {
			case *ast.Link:
				rendered.References = append(rendered.References, string(n.Destination))
			case *ast.Image:
				rendered.References = append(rendered.References, string(n.Destination))
			}
		}
		return ast.WalkContinue, nil
	})
	return rendered, err
}

func line(source []byte, start int) string {
	end := lineEnd(source, start)
	return strings.TrimSuffix(string(source[start:end]), "\n")
}

func isATX(value string) bool {
	value = strings.TrimLeft(value, " ")
	count := 0
	for count < len(value) && value[count] == '#' {
		count++
	}
	return count >= 1 && count <= 6 && (count == len(value) || value[count] == ' ' || value[count] == '\t' || value[count] == '\r')
}

func lineEnd(source []byte, start int) int {
	if start >= len(source) {
		return len(source)
	}
	if offset := bytes.IndexByte(source[start:], '\n'); offset >= 0 {
		return start + offset + 1
	}
	return len(source)
}
