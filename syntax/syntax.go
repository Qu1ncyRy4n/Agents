// Package syntax parses Mogent's small library, plan, and build language.
// It preserves ordered declarations, comments, and source ranges without
// accessing files, resolving libraries, fetching sources, or installing outputs.
package syntax

import (
	"fmt"
	"strings"
)

// Position uses a zero-based byte offset and one-based line and rune column.
type Position struct {
	Offset int
	Line   int
	Column int
}

// Range is a half-open source range, including its file name.
type Range struct {
	Filename string
	Start    Position
	End      Position
}

// Diagnostic is a located syntax or string-input error.
type Diagnostic struct {
	Range   Range
	Message string
}

func (d *Diagnostic) Error() string {
	return fmt.Sprintf("%s:%d:%d: %s", d.Range.Filename, d.Range.Start.Line, d.Range.Start.Column, d.Message)
}

type FileKind string

const (
	Library FileKind = "library"
	Plan    FileKind = "plan"
	Build   FileKind = "build"
)

// File contains syntax only. Semantic selection and source validation belong
// to the catalog/resolver, not to this parser.
type File struct {
	Kind     FileKind
	Entries  []*Entry
	Comments []Comment
	original string
}

// Original returns exactly the parsed source, including whitespace and comments.
// It is not a serialization of subsequent changes to Entries.
func (f *File) Original() string { return f.original }

type Comment struct {
	Text  string
	Range Range
}

type EntryKind string

const (
	Block   EntryKind = "block"
	Field   EntryKind = "field"
	Binding EntryKind = "binding"
)

// Entry retains declaration order; option-field order does not imply execution.
// Label and Selection are optional. HasBody distinguishes a bare include from
// an explicitly empty options block.
type Entry struct {
	Kind      EntryKind
	Name      string
	Label     *Value
	Selection *Value
	Value     *Value
	Body      []*Entry
	HasBody   bool
	Range     Range
}

type ValueKind string

const (
	String  ValueKind = "string"
	Boolean ValueKind = "boolean"
	Integer ValueKind = "integer"
	List    ValueKind = "list"
)

// StringPart is either literal Text or a named Reference. Inserted values are
// never reinterpreted as another interpolation expression.
type StringPart struct {
	Text      string
	Reference string
}

type Value struct {
	Kind    ValueKind
	Parts   []StringPart
	Bool    bool
	Int     int64
	Items   []Value
	Literal bool
	Range   Range
}

// ResolveString substitutes only explicit string bindings. Literal triple-quoted
// strings contain no references and require no special evaluation context.
func (v Value) ResolveString(bindings map[string]string) (string, error) {
	if v.Kind != String {
		return "", &Diagnostic{Range: v.Range, Message: "expected a string"}
	}
	var out strings.Builder
	for _, part := range v.Parts {
		if part.Reference == "" {
			out.WriteString(part.Text)
			continue
		}
		value, found := bindings[part.Reference]
		if !found {
			return "", &Diagnostic{Range: v.Range, Message: fmt.Sprintf("undefined binding %q", part.Reference)}
		}
		out.WriteString(value)
	}
	return out.String(), nil
}

// field returns a field without introducing an unordered representation of Body.
func (e *Entry) field(name string) *Entry {
	for _, child := range e.Body {
		if child != nil && child.Kind == Field && child.Name == name {
			return child
		}
	}
	return nil
}
