// Package manifest loads and validates the YAML document outline.
package manifest

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/renderfs"
	"github.com/Qu1ncyRy4n/Agents/sourcepath"
	"gopkg.in/yaml.v3"
)

// Manifest is the complete, repository-local build specification.
type Manifest struct {
	Roots   map[string]string `yaml:"roots,omitempty"`
	Sources map[string]Source `yaml:"sources"`
	Vars    map[string]any    `yaml:"vars,omitempty"`
	Output  string            `yaml:"output,omitempty"`
	Outputs []Output          `yaml:"outputs,omitempty"`
	Doc     []Entry           `yaml:"doc"`
}

// Output is an additional generated target. Kind is inferred from Path: .md
// renders the document, while paths ending in / or .d/ copy a source tree.
type Output struct {
	Path    string     `yaml:"path"`
	From    string     `yaml:"from,omitempty"` // Legacy directory selection.
	Include []Selector `yaml:"include,omitempty"`
	Exclude []Selector `yaml:"exclude,omitempty"`
}

// Selector selects source content for one canonical output.
type Selector struct {
	All    string
	Source string
	Tags   []string
}

func (s *Selector) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode || len(value.Content) != 2 {
		return fmt.Errorf("selector must contain exactly one of all, source, or tags")
	}
	key, item := value.Content[0].Value, value.Content[1]
	switch key {
	case "all", "source":
		if item.Kind != yaml.ScalarNode || item.Tag != "!!str" || strings.TrimSpace(item.Value) == "" {
			return fmt.Errorf("selector %q must be a non-empty string", key)
		}
		if key == "all" {
			s.All = item.Value
		} else {
			s.Source = item.Value
		}
	case "tags":
		if item.Kind != yaml.MappingNode || len(item.Content) != 2 || item.Content[0].Value != "any" {
			return fmt.Errorf("selector tags must contain only any")
		}
		values, err := stringsOnly(item.Content[1], "selector tags.any")
		if err != nil || len(values) == 0 {
			return fmt.Errorf("selector tags.any must be a non-empty sequence of strings")
		}
		s.Tags = values
	default:
		return fmt.Errorf("unknown selector key %q", key)
	}
	return nil
}

func (s Selector) MarshalYAML() (any, error) {
	if s.All != "" {
		return map[string]string{"all": s.All}, nil
	}
	if s.Source != "" {
		return map[string]string{"source": s.Source}, nil
	}
	return struct {
		Tags struct {
			Any []string `yaml:"any"`
		} `yaml:"tags"`
	}{Tags: struct {
		Any []string `yaml:"any"`
	}{Any: s.Tags}}, nil
}

func (o *Output) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode || len(value.Content)%2 != 0 {
		return fmt.Errorf("output must be a mapping")
	}
	for i := 0; i < len(value.Content); i += 2 {
		key, item := value.Content[i].Value, value.Content[i+1]
		switch key {
		case "path":
			if item.Kind != yaml.ScalarNode || item.Tag != "!!str" {
				return fmt.Errorf("output path must be a string")
			}
			o.Path = item.Value
		case "from":
			if item.Kind != yaml.ScalarNode || item.Tag != "!!str" {
				return fmt.Errorf("output from must be a string")
			}
			o.From = item.Value
		case "include":
			if item.Kind != yaml.SequenceNode {
				return fmt.Errorf("output include must be a sequence")
			}
			if err := item.Decode(&o.Include); err != nil {
				return err
			}
		case "exclude":
			if item.Kind != yaml.SequenceNode {
				return fmt.Errorf("output exclude must be a sequence")
			}
			if err := item.Decode(&o.Exclude); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown output key %q", key)
		}
	}
	return nil
}

func (o Output) Directory() bool {
	return strings.HasSuffix(o.Path, "/") || strings.HasSuffix(o.Path, ".d/")
}

// Source normalizes compact scalar locations and explicit source options.
type Source struct {
	Location string `yaml:"location"`
	Root     string `yaml:"-"`
	Subdir   string `yaml:"subdir,omitempty"`
}

func (s *Source) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		if value.Tag != "!!str" {
			return fmt.Errorf("source location must be a string")
		}
		s.Location = value.Value
		return nil
	case yaml.MappingNode:
		if len(value.Content)%2 != 0 {
			return fmt.Errorf("source options must be a mapping")
		}
		seen := make(map[string]bool, len(value.Content)/2)
		for index := 0; index < len(value.Content); index += 2 {
			key, item := value.Content[index].Value, value.Content[index+1]
			if seen[key] {
				return fmt.Errorf("duplicate source option %q", key)
			}
			seen[key] = true
			switch key {
			case "location":
				if item.Kind != yaml.ScalarNode || item.Tag != "!!str" {
					return fmt.Errorf("source option %q must be a string", key)
				}
				s.Location = item.Value
			case "subdir":
				if item.Kind != yaml.ScalarNode || item.Tag != "!!str" {
					return fmt.Errorf("source option %q must be a string", key)
				}
				s.Subdir = item.Value
			case "path":
				if item.Kind == yaml.ScalarNode && item.Tag == "!!str" {
					s.Location = item.Value
					continue
				}
				if item.Kind != yaml.MappingNode || len(item.Content)%2 != 0 {
					return fmt.Errorf("source path must be a string or root/subdir mapping")
				}
				for j := 0; j < len(item.Content); j += 2 {
					pathKey, pathValue := item.Content[j].Value, item.Content[j+1]
					if pathValue.Kind != yaml.ScalarNode || pathValue.Tag != "!!str" {
						return fmt.Errorf("source path %q must be a string", pathKey)
					}
					switch pathKey {
					case "root":
						s.Root = pathValue.Value
					case "subdir":
						s.Subdir = pathValue.Value
					default:
						return fmt.Errorf("unknown source path option %q", pathKey)
					}
				}
			default:
				return fmt.Errorf("unknown source option %q", key)
			}
		}
		return nil
	default:
		return fmt.Errorf("source must be a location string or options mapping")
	}
}

func (s Source) MarshalYAML() (any, error) {
	if s.Root != "" {
		return struct {
			Path struct {
				Root   string `yaml:"root"`
				Subdir string `yaml:"subdir,omitempty"`
			} `yaml:"path"`
		}{Path: struct {
			Root   string `yaml:"root"`
			Subdir string `yaml:"subdir,omitempty"`
		}{Root: s.Root, Subdir: s.Subdir}}, nil
	}
	if s.Subdir == "" {
		return s.Location, nil
	}
	return struct {
		Location string `yaml:"location"`
		Subdir   string `yaml:"subdir"`
	}{Location: s.Location, Subdir: s.Subdir}, nil
}

// Display returns the declared location with its selected root when present.
func (s Source) Display() string {
	if s.Subdir == "" {
		return s.Location
	}
	base := s.Location
	if s.Root != "" {
		base = "root: " + s.Root
	}
	return base + " (subdir: " + s.Subdir + ")"
}

// Entry owns one rendered heading and either composes source subtrees or owns
// an explicitly nested document outline. Source headings never control output.
type Entry struct {
	Heading  string   `yaml:"heading"`
	From     []string `yaml:"from,omitempty"`
	Exclude  []string `yaml:"exclude,omitempty"`
	Children []Entry  `yaml:"children,omitempty"`
}

// UnmarshalYAML accepts both canonical authoring forms and normalizes them to
// Entry. Compact entries keep the manifest readable; object entries provide
// ordered composition and exclusions without creating another content format.
func (e *Entry) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode || len(value.Content)%2 != 0 {
		return fmt.Errorf("entry must be a mapping")
	}
	keys := make(map[string]*yaml.Node, len(value.Content)/2)
	for index := 0; index < len(value.Content); index += 2 {
		key, node := value.Content[index].Value, value.Content[index+1]
		if _, exists := keys[key]; exists {
			return fmt.Errorf("duplicate entry key %q", key)
		}
		keys[key] = node
	}
	if heading, explicit := keys["heading"]; explicit {
		if err := rejectUnknownEntryKeys(keys); err != nil {
			return err
		}
		if heading.Kind != yaml.ScalarNode || heading.Tag != "!!str" {
			return fmt.Errorf("entry heading must be a string")
		}
		e.Heading = heading.Value
		return e.decodeOptions(keys)
	}
	if len(keys) != 1 {
		return fmt.Errorf("compact entry must contain one heading key")
	}
	for heading, node := range keys {
		e.Heading = heading
		switch node.Kind {
		case yaml.ScalarNode:
			if node.Tag != "!!str" {
				return fmt.Errorf("compact entry %q source must be a string", heading)
			}
			e.From = []string{node.Value}
			return nil
		case yaml.SequenceNode:
			return node.Decode(&e.Children)
		default:
			return fmt.Errorf("compact entry %q must contain a source string or children sequence", heading)
		}
	}
	return nil
}

// MarshalYAML keeps saved manifests close to the authoring model: single-source
// entries use compact syntax, while composed or nested entries use explicit
// fields.
func (e Entry) MarshalYAML() (any, error) {
	if len(e.From) == 1 && len(e.Exclude) == 0 && len(e.Children) == 0 {
		return map[string]string{e.Heading: e.From[0]}, nil
	}
	value := struct {
		Heading  string   `yaml:"heading"`
		From     []string `yaml:"from,omitempty"`
		Exclude  []string `yaml:"exclude,omitempty"`
		Children []Entry  `yaml:"children,omitempty"`
	}{
		Heading:  e.Heading,
		From:     e.From,
		Exclude:  e.Exclude,
		Children: e.Children,
	}
	return value, nil
}

func (e *Entry) decodeOptions(options map[string]*yaml.Node) error {
	if from, found := options["from"]; found {
		values, err := stringsOnly(from, "from")
		if err != nil {
			return err
		}
		e.From = values
	}
	if exclude, found := options["exclude"]; found {
		values, err := stringsOnly(exclude, "exclude")
		if err != nil {
			return err
		}
		e.Exclude = values
	}
	if children, found := options["children"]; found {
		if children.Kind != yaml.SequenceNode {
			return fmt.Errorf("children must be a sequence")
		}
		if err := children.Decode(&e.Children); err != nil {
			return err
		}
	}
	return nil
}

func mappingNodes(value *yaml.Node) (map[string]*yaml.Node, error) {
	if value.Kind != yaml.MappingNode || len(value.Content)%2 != 0 {
		return nil, fmt.Errorf("entry options must be a mapping")
	}
	result := make(map[string]*yaml.Node, len(value.Content)/2)
	for index := 0; index < len(value.Content); index += 2 {
		key, node := value.Content[index].Value, value.Content[index+1]
		if _, exists := result[key]; exists {
			return nil, fmt.Errorf("duplicate entry key %q", key)
		}
		result[key] = node
	}
	return result, nil
}

func rejectUnknownEntryKeys(keys map[string]*yaml.Node) error {
	for key := range keys {
		if key != "heading" && key != "from" && key != "exclude" && key != "children" {
			return fmt.Errorf("unknown entry key %q", key)
		}
	}
	return nil
}

func stringsOnly(value *yaml.Node, name string) ([]string, error) {
	if value.Kind != yaml.SequenceNode {
		return nil, fmt.Errorf("%s must be a sequence of strings", name)
	}
	result := make([]string, 0, len(value.Content))
	for _, item := range value.Content {
		if item.Kind != yaml.ScalarNode || item.Tag != "!!str" {
			return nil, fmt.Errorf("%s must be a sequence of strings", name)
		}
		result = append(result, item.Value)
	}
	return result, nil
}

// Load reads one manifest and returns its absolute path for relative source
// resolution. KnownFields and yaml.v3's duplicate-key checks prevent typo and
// last-wins configuration behavior.
func Load(path string) (*Manifest, string, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, "", fmt.Errorf("read manifest: %w", err)
	}
	var tree yaml.Node
	if err := yaml.Unmarshal(contents, &tree); err != nil {
		return nil, "", fmt.Errorf("parse manifest: %w", err)
	}
	if err := rejectAnchors(&tree); err != nil {
		return nil, "", err
	}
	decoder := yaml.NewDecoder(bytes.NewReader(contents))
	decoder.KnownFields(true)
	var value Manifest
	if err := decoder.Decode(&value); err != nil {
		return nil, "", fmt.Errorf("parse manifest: %w", err)
	}
	if err := ensureEOF(decoder); err != nil {
		return nil, "", err
	}
	if err := value.Validate(); err != nil {
		return nil, "", err
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, "", fmt.Errorf("resolve manifest path: %w", err)
	}
	return &value, abs, nil
}

// Clone returns an independent manifest draft.
func (m *Manifest) Clone() *Manifest {
	clone := &Manifest{
		Roots:   make(map[string]string, len(m.Roots)),
		Sources: make(map[string]Source, len(m.Sources)),
		Vars:    make(map[string]any, len(m.Vars)),
		Output:  m.Output,
		Outputs: cloneOutputs(m.Outputs),
		Doc:     cloneEntries(m.Doc),
	}
	for key, value := range m.Roots {
		clone.Roots[key] = value
	}
	for key, value := range m.Sources {
		clone.Sources[key] = value
	}
	for key, value := range m.Vars {
		clone.Vars[key] = value
	}
	return clone
}

func cloneOutputs(outputs []Output) []Output {
	clone := make([]Output, len(outputs))
	for i, output := range outputs {
		clone[i] = output
		clone[i].Include = append([]Selector(nil), output.Include...)
		clone[i].Exclude = append([]Selector(nil), output.Exclude...)
		for j := range clone[i].Include {
			clone[i].Include[j].Tags = append([]string(nil), output.Include[j].Tags...)
		}
		for j := range clone[i].Exclude {
			clone[i].Exclude[j].Tags = append([]string(nil), output.Exclude[j].Tags...)
		}
	}
	return clone
}

func cloneEntries(entries []Entry) []Entry {
	if len(entries) == 0 {
		return nil
	}
	clone := make([]Entry, len(entries))
	for index, entry := range entries {
		clone[index] = Entry{
			Heading:  entry.Heading,
			From:     append([]string(nil), entry.From...),
			Exclude:  append([]string(nil), entry.Exclude...),
			Children: cloneEntries(entry.Children),
		}
	}
	return clone
}

// WriteAtomically serializes and replaces agents.yaml only after validation and
// a successful temporary-file write.
func WriteAtomically(path string, value *Manifest) error {
	output, err := Marshal(value)
	if err != nil {
		return err
	}
	if err := renderfs.WriteAtomically(path, output); err != nil {
		return fmt.Errorf("write manifest: %w", err)
	}
	return nil
}

// Marshal serializes a validated manifest using the same normalized YAML shape
// used for atomic writes.
func Marshal(value *Manifest) ([]byte, error) {
	draft := value.Clone()
	if err := draft.Validate(); err != nil {
		return nil, err
	}
	var output bytes.Buffer
	encoder := yaml.NewEncoder(&output)
	encoder.SetIndent(2)
	if err := encoder.Encode(draft); err != nil {
		return nil, fmt.Errorf("serialize manifest: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return nil, fmt.Errorf("finish manifest serialization: %w", err)
	}
	return output.Bytes(), nil
}

func ensureEOF(decoder *yaml.Decoder) error {
	var extra any
	err := decoder.Decode(&extra)
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return fmt.Errorf("parse manifest: %w", err)
	}
	return fmt.Errorf("parse manifest: multiple YAML documents are not supported")
}

// Validate checks only structural rules. Source references are verified after
// libraries are indexed.
func (m *Manifest) Validate() error {
	if len(m.Sources) == 0 {
		return fmt.Errorf("manifest requires at least one source")
	}
	for name, root := range m.Roots {
		if strings.TrimSpace(name) == "" || strings.Contains(name, ":") || strings.TrimSpace(root) == "" {
			return fmt.Errorf("invalid root %q", name)
		}
	}
	for alias, source := range m.Sources {
		if strings.TrimSpace(alias) == "" || strings.Contains(alias, ":") {
			return fmt.Errorf("invalid source alias %q", alias)
		}
		if source.Root != "" {
			root, found := m.Roots[source.Root]
			if !found {
				return fmt.Errorf("source %q references undeclared root %q", alias, source.Root)
			}
			source.Location = root
			m.Sources[alias] = source
		}
		if strings.TrimSpace(source.Location) == "" {
			return fmt.Errorf("source %q has an empty path", alias)
		}
		if err := sourcepath.ValidateSubdir(source.Subdir); err != nil {
			return fmt.Errorf("source %q: %w", alias, err)
		}
	}
	if strings.TrimSpace(m.Output) == "" && len(m.Doc) > 0 {
		m.Output = "AGENTS.md"
	}
	all := m.EffectiveOutputs()
	for i, output := range all {
		if err := output.validate(); err != nil {
			return fmt.Errorf("output[%d]: %w", i, err)
		}
	}
	for i := range all {
		for j := i + 1; j < len(all); j++ {
			first, second := strings.TrimSuffix(all[i].Path, "/"), strings.TrimSuffix(all[j].Path, "/")
			if first == second || strings.HasPrefix(first, second+"/") || strings.HasPrefix(second, first+"/") {
				return fmt.Errorf("outputs %q and %q collide or overlap", all[i].Path, all[j].Path)
			}
		}
	}
	if len(m.Doc) == 0 && len(m.Outputs) == 0 {
		return fmt.Errorf("manifest requires doc entries or outputs")
	}
	for index, entry := range m.Doc {
		if err := entry.validate(fmt.Sprintf("doc[%d]", index)); err != nil {
			return err
		}
	}
	return nil
}

func (o Output) validate() error {
	path := o.Path
	if strings.TrimSpace(path) == "" || filepath.IsAbs(path) || strings.Contains(path, "\\") {
		return fmt.Errorf("path must be a non-empty relative slash-separated path")
	}
	for _, part := range strings.Split(strings.TrimSuffix(path, "/"), "/") {
		if part == "" || part == "." || part == ".." || part == ".mogent" {
			return fmt.Errorf("unsafe output path %q", path)
		}
	}
	if o.Directory() {
		for _, selector := range append(append([]Selector(nil), o.Include...), o.Exclude...) {
			if err := selector.validate(); err != nil {
				return err
			}
		}
		if o.From != "" {
			if _, _, err := SplitReference(o.From); err != nil {
				return err
			}
		}
		if o.From == "" && len(o.Include) == 0 {
			return fmt.Errorf("directory output %q requires include", path)
		}
		if o.From != "" && len(o.Include) > 0 {
			return fmt.Errorf("directory output %q cannot combine from and include", path)
		}
		return nil
	}
	if strings.HasSuffix(path, ".md") {
		if o.From != "" {
			return fmt.Errorf("Markdown output %q must not use from", path)
		}
		if len(o.Exclude) > 0 && len(o.Include) == 0 {
			return fmt.Errorf("Markdown output %q uses exclude without include", path)
		}
		for _, selector := range append(append([]Selector(nil), o.Include...), o.Exclude...) {
			if err := selector.validate(); err != nil {
				return err
			}
		}
		return nil
	}
	return fmt.Errorf("ambiguous output path %q; use .md or a trailing /", path)
}

func (s Selector) validate() error {
	count := 0
	if s.All != "" {
		count++
		if strings.Contains(s.All, ":") || strings.TrimSpace(s.All) == "" {
			return fmt.Errorf("all selector requires a source alias")
		}
	}
	if s.Source != "" {
		count++
		if _, _, err := SplitReference(s.Source); err != nil {
			return err
		}
	}
	if len(s.Tags) > 0 {
		count++
		for _, tag := range s.Tags {
			if strings.TrimSpace(tag) == "" {
				return fmt.Errorf("tag selector contains an empty tag")
			}
		}
	}
	if count != 1 {
		return fmt.Errorf("selector must contain exactly one of all, source, or tags")
	}
	return nil
}

// EffectiveOutputs returns canonical outputs, or the legacy primary output plus
// legacy additions when the manifest still has a document outline.
func (m *Manifest) EffectiveOutputs() []Output {
	if m.Output == "" {
		return append([]Output(nil), m.Outputs...)
	}
	return append([]Output{{Path: m.Output}}, m.Outputs...)
}

func (e Entry) validate(location string) error {
	if strings.TrimSpace(e.Heading) == "" {
		return fmt.Errorf("%s requires heading", location)
	}
	hasFrom, hasChildren := len(e.From) > 0, len(e.Children) > 0
	if hasFrom == hasChildren {
		return fmt.Errorf("%s must have exactly one of from or children", location)
	}
	if !hasFrom && len(e.Exclude) > 0 {
		return fmt.Errorf("%s uses exclude without from", location)
	}
	for _, reference := range e.From {
		if _, _, err := SplitReference(reference); err != nil {
			return fmt.Errorf("%s: %w", location, err)
		}
	}
	for _, excluded := range e.Exclude {
		if _, _, err := SplitReference(excluded); err != nil {
			return fmt.Errorf("%s: exclusion %q must be source-qualified", location, excluded)
		}
	}
	for index, child := range e.Children {
		if err := child.validate(fmt.Sprintf("%s.children[%d]", location, index)); err != nil {
			return err
		}
	}
	return nil
}

// SplitReference parses explicit source provenance in alias:heading/path form.
func SplitReference(reference string) (string, string, error) {
	alias, path, ok := strings.Cut(strings.TrimSpace(reference), ":")
	if !ok || alias == "" || path == "" || strings.HasPrefix(path, "/") || strings.Contains(path, "//") {
		return "", "", fmt.Errorf("invalid source reference %q; expected alias:heading/path", reference)
	}
	for _, part := range strings.Split(path, "/") {
		if part == "" || part == "." || part == ".." {
			return "", "", fmt.Errorf("invalid source reference %q; path must be a heading path", reference)
		}
	}
	return alias, path, nil
}

func rejectAnchors(node *yaml.Node) error {
	if node.Anchor != "" || node.Kind == yaml.AliasNode || node.Alias != nil {
		return fmt.Errorf("parse manifest: YAML anchors and aliases are not supported")
	}
	for _, child := range node.Content {
		if err := rejectAnchors(child); err != nil {
			return err
		}
	}
	return nil
}
