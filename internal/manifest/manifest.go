// Package manifest loads and validates the YAML document outline.
package manifest

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Manifest is the complete, repository-local build specification.
type Manifest struct {
	Sources map[string]string `yaml:"sources"`
	Vars    map[string]any    `yaml:"vars,omitempty"`
	Output  string            `yaml:"output,omitempty"`
	Doc     []Entry           `yaml:"doc"`
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
		Sources: make(map[string]string, len(m.Sources)),
		Vars:    make(map[string]any, len(m.Vars)),
		Output:  m.Output,
		Doc:     cloneEntries(m.Doc),
	}
	for key, value := range m.Sources {
		clone.Sources[key] = value
	}
	for key, value := range m.Vars {
		clone.Vars[key] = value
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
	draft := value.Clone()
	if err := draft.Validate(); err != nil {
		return err
	}
	var output bytes.Buffer
	encoder := yaml.NewEncoder(&output)
	encoder.SetIndent(2)
	if err := encoder.Encode(draft); err != nil {
		return fmt.Errorf("serialize manifest: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return fmt.Errorf("finish manifest serialization: %w", err)
	}
	directory := filepath.Dir(path)
	temporary, err := os.CreateTemp(directory, ".mogent-manifest-*")
	if err != nil {
		return fmt.Errorf("create temporary manifest: %w", err)
	}
	temporaryName := temporary.Name()
	if _, err := temporary.Write(output.Bytes()); err != nil {
		return cleanupTemporary(fmt.Errorf("write temporary manifest: %w", err), temporary, temporaryName)
	}
	if err := temporary.Chmod(0o644); err != nil {
		return cleanupTemporary(fmt.Errorf("set manifest permissions: %w", err), temporary, temporaryName)
	}
	if err := temporary.Close(); err != nil {
		if removeErr := os.Remove(temporaryName); removeErr != nil {
			return errors.Join(fmt.Errorf("close temporary manifest: %w", err), fmt.Errorf("remove temporary manifest: %w", removeErr))
		}
		return fmt.Errorf("close temporary manifest: %w", err)
	}
	if err := os.Rename(temporaryName, path); err != nil {
		if removeErr := os.Remove(temporaryName); removeErr != nil {
			return errors.Join(fmt.Errorf("replace manifest: %w", err), fmt.Errorf("remove temporary manifest: %w", removeErr))
		}
		return fmt.Errorf("replace manifest: %w", err)
	}
	return nil
}

func cleanupTemporary(buildErr error, temporary *os.File, path string) error {
	closeErr := temporary.Close()
	removeErr := os.Remove(path)
	if closeErr != nil || removeErr != nil {
		return errors.Join(buildErr, closeErr, removeErr)
	}
	return buildErr
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
	for alias, sourcePath := range m.Sources {
		if strings.TrimSpace(alias) == "" || strings.Contains(alias, ":") {
			return fmt.Errorf("invalid source alias %q", alias)
		}
		if strings.TrimSpace(sourcePath) == "" {
			return fmt.Errorf("source %q has an empty path", alias)
		}
	}
	if strings.TrimSpace(m.Output) == "" {
		m.Output = "AGENTS.md"
	}
	if len(m.Doc) == 0 {
		return fmt.Errorf("manifest requires at least one doc entry")
	}
	for index, entry := range m.Doc {
		if err := entry.validate(fmt.Sprintf("doc[%d]", index)); err != nil {
			return err
		}
	}
	return nil
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
