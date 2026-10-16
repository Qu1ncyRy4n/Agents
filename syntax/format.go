package syntax

import (
	"strconv"
	"strings"
)

// FormatOptions affect spelling, never selection or execution semantics.
type FormatOptions struct {
	ExpandAliases bool
}

// Format writes deterministic, multiline syntax, preserving all parsed comments.
// Comments inside values may move to their containing statement's end. Use
// Original for exact no-op source preservation; this is not an authored editor.
func Format(file *File, options FormatOptions) ([]byte, error) {
	if err := validate(file); err != nil {
		return nil, err
	}
	p := &formatter{comments: file.Comments, options: options}
	for _, e := range file.Entries {
		p.entry(e, 0)
	}
	p.commentsBefore(int(^uint(0)>>1), 0)
	return []byte(p.out.String()), nil
}

type formatter struct {
	out      strings.Builder
	comments []Comment
	comment  int
	options  FormatOptions
}

func (p *formatter) commentsBefore(offset, depth int) {
	for p.comment < len(p.comments) && p.comments[p.comment].Range.Start.Offset < offset {
		p.out.WriteString(strings.Repeat("  ", depth))
		p.out.WriteString(p.comments[p.comment].Text)
		p.out.WriteByte('\n')
		p.comment++
	}
}

func (p *formatter) entry(e *Entry, depth int) {
	p.commentsBefore(e.Range.Start.Offset, depth)
	indent := strings.Repeat("  ", depth)
	p.out.WriteString(indent)
	if e.Kind == Binding {
		p.out.WriteString("let ")
	}
	name := e.Name
	if e.Kind == Block && p.options.ExpandAliases {
		if name == "replace" {
			name = "include"
		} else if isNode(name) && name != "node" {
			name = "node"
		}
	}
	p.out.WriteString(name)
	if e.Kind == Binding || e.Kind == Field {
		p.out.WriteString(" = ")
		p.value(*e.Value)
		p.out.WriteByte('\n')
		p.commentsBefore(e.Range.End.Offset, depth)
		return
	}
	if e.Label != nil {
		p.out.WriteByte(' ')
		p.value(*e.Label)
	}
	if e.Selection != nil {
		p.out.WriteByte(' ')
		p.value(*e.Selection)
	}
	if !e.HasBody && name == e.Name {
		p.out.WriteByte('\n')
		return
	}
	p.out.WriteString(" {\n")
	if p.options.ExpandAliases && e.Name != name {
		if e.Name == "replace" && e.field("replace") == nil {
			p.out.WriteString(indent + "  replace = true\n")
		} else if isNode(e.Name) && e.field("type") == nil {
			p.out.WriteString(indent + "  type = " + strconv.Quote(e.Name) + "\n")
		}
	}
	for _, child := range e.Body {
		p.entry(child, depth+1)
	}
	p.commentsBefore(e.Range.End.Offset, depth+1)
	p.out.WriteString(indent + "}\n")
}

func (p *formatter) value(value Value) {
	switch value.Kind {
	case String:
		p.out.WriteByte('"')
		for _, part := range value.Parts {
			if part.Reference != "" {
				p.out.WriteString("${" + part.Reference + "}")
				continue
			}
			for i := 0; i < len(part.Text); i++ {
				switch part.Text[i] {
				case '\\':
					p.out.WriteString(`\\`)
				case '"':
					p.out.WriteString(`\"`)
				case '\n':
					p.out.WriteString(`\n`)
				case '\r':
					p.out.WriteString(`\r`)
				case '\t':
					p.out.WriteString(`\t`)
				case '$':
					if i+1 < len(part.Text) && part.Text[i+1] == '{' {
						p.out.WriteByte('\\')
					}
					p.out.WriteByte('$')
				default:
					p.out.WriteByte(part.Text[i])
				}
			}
		}
		p.out.WriteByte('"')
	case Boolean:
		p.out.WriteString(strconv.FormatBool(value.Bool))
	case Integer:
		p.out.WriteString(strconv.FormatInt(value.Int, 10))
	case List:
		p.out.WriteByte('[')
		for i, item := range value.Items {
			if i > 0 {
				p.out.WriteString(", ")
			}
			p.value(item)
		}
		p.out.WriteByte(']')
	}
}
