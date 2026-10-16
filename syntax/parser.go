package syntax

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Parse parses a complete format-1 library, plan, or build. It validates syntax,
// header identity, declaration shape, and closed option names. It does not check
// filesystem existence, library relationships, unused names, or source integrity.
func Parse(filename string, source []byte) (*File, error) {
	if !utf8.Valid(source) {
		return nil, &Diagnostic{Range: Range{Filename: filename, Start: Position{Line: 1, Column: 1}}, Message: "source must be UTF-8"}
	}
	l := newLexer(filename, string(source))
	p := &parser{lexer: l}
	if err := p.advance(); err != nil {
		return nil, err
	}
	entries, err := p.entries(false, 0)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, l.fail(0, 0, "require a library, plan, or build header")
	}
	f := &File{Kind: FileKind(entries[0].Name), Entries: entries, Comments: l.comments, original: string(source)}
	if err := validate(f); err != nil {
		return nil, err
	}
	return f, nil
}

type parser struct {
	lexer *lexer
	next  token
}

func (p *parser) advance() error {
	t, err := p.lexer.next()
	if err == nil {
		p.next = t
	}
	return err
}

func (p *parser) error(message string) error {
	return p.lexer.fail(p.next.start, p.next.end, message)
}

func (p *parser) entries(nested bool, depth int) ([]*Entry, error) {
	if depth > 128 {
		return nil, p.error("nesting exceeds 128 blocks")
	}
	var entries []*Entry
	for {
		for p.next.kind == tokenNewline || p.next.kind == ';' {
			if err := p.advance(); err != nil {
				return nil, err
			}
		}
		if p.next.kind == tokenEOF {
			if nested {
				return nil, p.error("expected closing brace")
			}
			return entries, nil
		}
		if p.next.kind == '}' {
			if !nested {
				return nil, p.error("unexpected closing brace")
			}
			return entries, nil
		}
		entry, err := p.entry(depth)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
		switch p.next.kind {
		case tokenNewline, ';', tokenEOF, '}':
		default:
			return nil, p.error("expected newline or ';' between statements; commas only separate array items")
		}
	}
}

func (p *parser) entry(depth int) (*Entry, error) {
	if p.next.kind != tokenIdentifier {
		return nil, p.error("expected a declaration or field name")
	}
	start, name := p.next.start, p.next.text
	if err := p.advance(); err != nil {
		return nil, err
	}
	e := &Entry{Name: name, Kind: Block}
	if name == "let" {
		if p.next.kind != tokenIdentifier {
			return nil, p.error("let requires a binding name")
		}
		e.Kind, e.Name = Binding, p.next.text
		if err := p.advance(); err != nil {
			return nil, err
		}
		if p.next.kind != '=' {
			return nil, p.error("binding requires '='")
		}
	}
	if p.next.kind == '=' {
		if e.Kind != Binding {
			e.Kind = Field
		}
		if err := p.advance(); err != nil {
			return nil, err
		}
		value, err := p.value()
		if err != nil {
			return nil, err
		}
		e.Value = &value
		e.Range = p.lexer.span(start, value.Range.End.Offset)
		return e, nil
	}
	end := start + len(name)
	if p.next.kind == tokenString {
		e.Label = p.next.value
		end = p.next.end
		if err := p.advance(); err != nil {
			return nil, err
		}
	}
	if p.next.kind == '[' {
		value, err := p.value()
		if err != nil {
			return nil, err
		}
		e.Selection = &value
		end = value.Range.End.Offset
	}
	if p.next.kind == '{' {
		e.HasBody = true
		if err := p.advance(); err != nil {
			return nil, err
		}
		body, err := p.entries(true, depth+1)
		if err != nil {
			return nil, err
		}
		e.Body = body
		end = p.next.end
		if err := p.advance(); err != nil {
			return nil, err
		}
	}
	e.Range = p.lexer.span(start, end)
	return e, nil
}

func (p *parser) value() (Value, error) {
	t := p.next
	if t.kind == tokenString || t.kind == tokenInteger {
		if err := p.advance(); err != nil {
			return Value{}, err
		}
		return *t.value, nil
	}
	if t.kind == tokenIdentifier && (t.text == "true" || t.text == "false") {
		if err := p.advance(); err != nil {
			return Value{}, err
		}
		return Value{Kind: Boolean, Bool: t.text == "true", Range: p.lexer.span(t.start, t.end)}, nil
	}
	if t.kind != '[' {
		return Value{}, p.error("expected quoted string, boolean, integer, or string array; expressions and object assignments are not supported")
	}
	if err := p.advance(); err != nil {
		return Value{}, err
	}
	v := Value{Kind: List}
	for {
		for p.next.kind == tokenNewline {
			if err := p.advance(); err != nil {
				return Value{}, err
			}
		}
		if p.next.kind == ']' {
			v.Range = p.lexer.span(t.start, p.next.end)
			return v, p.advance()
		}
		if p.next.kind != tokenString {
			return Value{}, p.error("array items must be quoted strings")
		}
		v.Items = append(v.Items, *p.next.value)
		if err := p.advance(); err != nil {
			return Value{}, err
		}
		for p.next.kind == tokenNewline {
			if err := p.advance(); err != nil {
				return Value{}, err
			}
		}
		if p.next.kind == ']' {
			v.Range = p.lexer.span(t.start, p.next.end)
			return v, p.advance()
		}
		if p.next.kind != ',' {
			return Value{}, p.error("expected ',' between array items")
		}
		if err := p.advance(); err != nil {
			return Value{}, err
		}
	}
}

func located(e *Entry, message string) error {
	return &Diagnostic{Range: e.Range, Message: message}
}

func staticLabel(e *Entry) (string, error) {
	if e.Label == nil {
		return "", located(e, fmt.Sprintf("%s requires a quoted name", e.Name))
	}
	label, err := e.Label.ResolveString(nil)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(label) == "" {
		return "", located(e, "name cannot be empty")
	}
	return label, nil
}
