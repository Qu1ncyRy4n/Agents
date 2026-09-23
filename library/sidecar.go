package library

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// SidecarFile is the optional library-root metadata file.
const SidecarFile = "library.mogent.yaml"

// Sidecar describes library-owned identity, source ordering, and relationships.
type Sidecar struct {
	Schema  SidecarSchema           `yaml:"schema"`
	Library SidecarLibrary          `yaml:"library"`
	Tree    []SidecarTreeNode       `yaml:"tree"`
	Groups  map[string]SidecarGroup `yaml:"groups,omitempty"`
}

type SidecarSchema struct {
	ID string `yaml:"id"`
}
type SidecarLibrary struct {
	Name         string `yaml:"name"`
	Version      any    `yaml:"version,omitempty"`
	Description  string `yaml:"description,omitempty"`
	Maintainers  any    `yaml:"maintainers,omitempty"`
	Organization any    `yaml:"organization,omitempty"`
	Repository   any    `yaml:"repository,omitempty"`
	License      any    `yaml:"license,omitempty"`
}
type SidecarGroup struct {
	Mode string `yaml:"mode"`
}
type SidecarTreeNode struct {
	Source         string            `yaml:"source"`
	Title          string            `yaml:"title,omitempty"`
	Children       []SidecarTreeNode `yaml:"children,omitempty"`
	Requires       []string          `yaml:"requires,omitempty"`
	ConflictsWith  []string          `yaml:"conflicts_with,omitempty"`
	ExclusiveGroup string            `yaml:"exclusive_group,omitempty"`
	// These fields are accepted only so check can identify ownership violations.
	Tags []string `yaml:"tags,omitempty"`
	TLDR string   `yaml:"tldr,omitempty"`
}

// CheckReport separates validation errors from compatibility and ownership warnings.
type CheckReport struct{ Errors, Warnings []string }

func (r *CheckReport) addError(format string, args ...any) {
	r.Errors = append(r.Errors, fmt.Sprintf(format, args...))
}
func (r *CheckReport) addWarning(format string, args ...any) {
	r.Warnings = append(r.Warnings, fmt.Sprintf(format, args...))
}

// LoadSidecar reads a root sidecar. A missing sidecar is reported as os.ErrNotExist.
func LoadSidecar(root string) (*Sidecar, []string, error) {
	path := filepath.Join(root, SidecarFile)
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	var sidecar Sidecar
	decoder := yaml.NewDecoder(bytes.NewReader(contents))
	decoder.KnownFields(true)
	if err := decoder.Decode(&sidecar); err != nil {
		return nil, nil, fmt.Errorf("parse %s: %w", path, err)
	}
	var extra any
	if err := decoder.Decode(&extra); err == nil {
		return nil, nil, fmt.Errorf("parse %s: sidecar must contain one YAML document", path)
	} else if !errors.Is(err, io.EOF) {
		return nil, nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if sidecar.Schema.ID == "" {
		return &sidecar, []string{"schema.id is omitted; assuming mogent/1"}, nil
	}
	if sidecar.Schema.ID != "mogent/1" {
		return nil, nil, fmt.Errorf("parse %s: unsupported schema.id %q", path, sidecar.Schema.ID)
	}
	return &sidecar, nil, nil
}

// LoadWithSidecar indexes a source and applies its optional sidecar's
// presentation metadata. Load remains the raw source inventory API for scan
// and validation callers.
func LoadWithSidecar(root string) (*Index, []string, error) {
	index, err := Load(root)
	if err != nil {
		return nil, nil, err
	}
	sidecar, warnings, err := LoadSidecar(root)
	if errors.Is(err, os.ErrNotExist) {
		return index, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}
	report := validateSidecar(sidecar, index)
	report.Warnings = append(warnings, report.Warnings...)
	if len(report.Errors) > 0 {
		return nil, report.Warnings, fmt.Errorf("library sidecar validation failed: %s", strings.Join(report.Errors, "; "))
	}
	applySidecar(index, sidecar)
	return index, report.Warnings, nil
}

// Check validates a sidecar against the Markdown source inventory.
func Check(root string) (*CheckReport, error) {
	sidecar, warnings, err := LoadSidecar(root)
	if err != nil {
		return nil, err
	}
	index, err := Load(root)
	if err != nil {
		return nil, err
	}
	report := validateSidecar(sidecar, index)
	report.Warnings = append(warnings, report.Warnings...)
	sort.Strings(report.Warnings)
	return report, nil
}

func validateSidecar(sidecar *Sidecar, index *Index) *CheckReport {
	report := &CheckReport{}
	if strings.TrimSpace(sidecar.Library.Name) == "" {
		report.addError("library.name is required")
	}
	declared := map[string]SidecarTreeNode{}
	parents := map[string]string{}
	var visit func([]SidecarTreeNode, string)
	visit = func(nodes []SidecarTreeNode, parent string) {
		for _, node := range nodes {
			field := "tree source " + fmt.Sprintf("%q", node.Source)
			if !validSidecarPath(node.Source) {
				report.addError("%s is not a normalized library-relative source path", field)
			}
			if _, found := declared[node.Source]; found {
				report.addError("duplicate tree source %q", node.Source)
			}
			declared[node.Source], parents[node.Source] = node, parent
			if len(node.Tags) > 0 || node.TLDR != "" {
				report.addWarning("tree source %q contains tags or tldr; Markdown frontmatter owns those fields", node.Source)
			}
			visit(node.Children, node.Source)
		}
	}
	visit(sidecar.Tree, "")
	for path, node := range index.ByPath {
		if _, found := declared[path]; !found {
			report.addError("tree is missing discovered source %q", path)
		}
		if len(node.Metadata.Requires) > 0 || len(node.Metadata.ConflictsWith) > 0 {
			report.addWarning("source %q has requires or conflicts_with in Markdown frontmatter; sidecar owns relationship fields", path)
		}
	}
	for path, node := range declared {
		actual, found := index.ByPath[path]
		for _, target := range append(append([]string(nil), node.Requires...), node.ConflictsWith...) {
			if _, ok := declared[target]; !ok {
				report.addError("tree source %q references unknown relationship target %q", path, target)
			}
		}
		if !found {
			report.addError("tree source %q is not a discovered directory or Markdown heading", path)
			continue
		}
		wantParent := ""
		for candidate, possible := range index.ByPath {
			for _, child := range possible.Children {
				if child == actual {
					wantParent = candidate
					break
				}
			}
			if wantParent != "" {
				break
			}
		}
		if parents[path] != wantParent {
			report.addError("tree source %q is nested under %q, but source hierarchy requires %q", path, parents[path], wantParent)
		}
		if node.ExclusiveGroup != "" {
			if _, ok := sidecar.Groups[node.ExclusiveGroup]; !ok {
				report.addError("tree source %q references unknown exclusive group %q", path, node.ExclusiveGroup)
			}
		}
	}
	members := map[string]int{}
	for _, node := range declared {
		if node.ExclusiveGroup != "" {
			members[node.ExclusiveGroup]++
		}
	}
	for name, group := range sidecar.Groups {
		if group.Mode != "at-most-one" {
			report.addError("group %q has invalid mode %q (expected at-most-one)", name, group.Mode)
		}
		if members[name] == 0 {
			report.addError("group %q has no member tree nodes", name)
		}
	}
	checkRequirementCycles(declared, report)
	sort.Strings(report.Errors)
	sort.Strings(report.Warnings)
	return report
}

func applySidecar(index *Index, sidecar *Sidecar) {
	var apply func([]SidecarTreeNode) []*Node
	apply = func(nodes []SidecarTreeNode) []*Node {
		ordered := make([]*Node, 0, len(nodes))
		for _, declared := range nodes {
			node := index.ByPath[declared.Source]
			if declared.Title != "" {
				node.Heading = declared.Title
			}
			node.Children = apply(declared.Children)
			ordered = append(ordered, node)
		}
		return ordered
	}
	index.Roots = apply(sidecar.Tree)
}

func validSidecarPath(path string) bool {
	return path != "" && !strings.Contains(path, "\\") && filepath.ToSlash(filepath.Clean(path)) == path && !strings.HasPrefix(path, "../") && path != "." && !strings.HasPrefix(path, "/")
}
func checkRequirementCycles(nodes map[string]SidecarTreeNode, report *CheckReport) {
	state := map[string]int{}
	var visit func(string)
	visit = func(path string) {
		if state[path] == 1 {
			report.addError("requires cycle includes %q", path)
			return
		}
		if state[path] == 2 {
			return
		}
		state[path] = 1
		for _, next := range nodes[path].Requires {
			if _, ok := nodes[next]; ok {
				visit(next)
			}
		}
		state[path] = 2
	}
	paths := make([]string, 0, len(nodes))
	for path := range nodes {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		visit(path)
	}
}

// Scan inventories a library as a deterministic, exhaustive candidate sidecar.
func Scan(root string) (*Sidecar, error) {
	index, err := Load(root)
	if err != nil {
		return nil, err
	}
	var convert func([]*Node) []SidecarTreeNode
	convert = func(nodes []*Node) []SidecarTreeNode {
		result := make([]SidecarTreeNode, 0, len(nodes))
		for _, node := range nodes {
			result = append(result, SidecarTreeNode{Source: node.Path, Children: convert(node.Children)})
		}
		return result
	}
	return &Sidecar{Schema: SidecarSchema{ID: "mogent/1"}, Library: SidecarLibrary{Name: filepath.Base(filepath.Clean(root))}, Tree: convert(index.Roots)}, nil
}

// MarshalSidecar returns the canonical YAML representation used by init and scan.
func MarshalSidecar(sidecar *Sidecar) ([]byte, error) {
	contents, err := yaml.Marshal(sidecar)
	if err != nil {
		return nil, err
	}
	return contents, nil
}
