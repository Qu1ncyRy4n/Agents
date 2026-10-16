// Package assembly loads typed Mogent libraries and resolves plan/build files
// into an inspectable ordered result. It intentionally depends on syntax and
// content rather than legacy consumer configuration formats.
package assembly

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/content"
	"github.com/Qu1ncyRy4n/Agents/syntax"
)

const LibraryFile = "mogent.lib"

type NodeType string

const (
	Directory NodeType = "dir"
	File      NodeType = "file"
	Markdown  NodeType = "md-file"
	Section   NodeType = "md-section"
)

type Curate string

const (
	Foundation Curate = "foundation"
	Choose     Curate = "choose"
	Optional   Curate = "optional"
	OptIn      Curate = "opt_in"
)

// Node represents either a logical document section or a physical payload.
// Path is library-relative for physical nodes and logical for sections.
type Node struct {
	Name           string
	Path           string
	Type           NodeType
	Source         string
	Title          string
	Tags           []string
	Role           string
	Content        string
	Curate         Curate
	Default        *bool
	ReasonRequired bool
	Requires       []string
	ConflictsWith  []string
	ExclusiveGroup string
	Into           string
	Variables      []string
	Heading        []string
	Children       []*Node
	Parent         *Node
	Physical       *content.Inventory
	PhysicalNode   string
	Library        *Library
	Range          syntax.Range
}

type Preset struct {
	Name    string
	Entries []*syntax.Entry
	Range   syntax.Range
}

type Library struct {
	Path    string
	Root    string
	Name    string
	Nodes   []*Node
	ByPath  map[string]*Node
	Presets map[string]Preset
}

type Source struct {
	Alias string
	Root  string
	Lib   *Library
}

type Severity string

const (
	Error   Severity = "error"
	Warning Severity = "warning"
	Info    Severity = "info"
)

type Diagnostic struct {
	Severity Severity
	Range    syntax.Range
	Path     string
	Message  string
}

func (d Diagnostic) String() string {
	location := d.Range.Filename
	if location != "" {
		location += fmt.Sprintf(":%d:%d: ", d.Range.Start.Line, d.Range.Start.Column)
	}
	if d.Path != "" {
		location += d.Path + ": "
	}
	return fmt.Sprintf("%s%s: %s", location, d.Severity, d.Message)
}

type Decision struct {
	Path   string
	Action string
	Reason string
	Range  syntax.Range
}

// Contribution is one resolved ordered source or inline document contribution.
type Contribution struct {
	Source   string
	Node     *Node
	Into     string
	Replace  bool
	Inline   string
	Values   map[string]string
	Decision []Decision
	Range    syntax.Range
}

type Result struct {
	Contributions []Contribution
	Diagnostics   []Diagnostic
}

func (r *Result) HasErrors() bool {
	for _, d := range r.Diagnostics {
		if d.Severity == Error {
			return true
		}
	}
	return false
}

func field(e *syntax.Entry, name string) *syntax.Value {
	for _, child := range e.Body {
		if child.Kind == syntax.Field && child.Name == name {
			return child.Value
		}
	}
	return nil
}

func fieldString(e *syntax.Entry, name string) (string, bool, error) {
	v := field(e, name)
	if v == nil {
		return "", false, nil
	}
	value, err := v.ResolveString(nil)
	return value, true, err
}

func fieldBool(e *syntax.Entry, name string) (bool, bool) {
	v := field(e, name)
	if v == nil || v.Kind != syntax.Boolean {
		return false, false
	}
	return v.Bool, true
}

func fieldStrings(e *syntax.Entry, name string) ([]string, bool, error) {
	v := field(e, name)
	if v == nil {
		return nil, false, nil
	}
	if v.Kind != syntax.List {
		return nil, true, fmt.Errorf("%s must be a string list", name)
	}
	result := make([]string, 0, len(v.Items))
	for _, item := range v.Items {
		value, err := item.ResolveString(nil)
		if err != nil {
			return nil, true, err
		}
		result = append(result, value)
	}
	return result, true, nil
}

func safePath(root, value string) (string, error) {
	if value == "" || filepath.IsAbs(value) || strings.Contains(value, "\\") || value == ".." || strings.HasPrefix(value, "../") || strings.Contains(value, "/../") {
		return "", fmt.Errorf("%q must be a safe relative slash path", value)
	}
	path := filepath.Join(root, filepath.FromSlash(value))
	return path, nil
}

func orderedChildren(nodes []*Node) []*Node {
	result := append([]*Node(nil), nodes...)
	sort.SliceStable(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result
}
