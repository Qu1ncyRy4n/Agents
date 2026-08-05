package workspace

import (
	"fmt"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/internal/manifest"
	"github.com/Qu1ncyRy4n/Agents/internal/render"
)

type AddOptions struct {
	Reference string
	Under     string
	Append    bool
	Heading   string
	Rebuild   bool
}

type AddResult struct {
	Heading       string
	Reference     string
	ParentPath    string
	Preview       string
	Section       string
	Tree          string
	ManifestHunk  DiffHunk
	OutputHunk    DiffHunk
	WroteManifest bool
	Rebuilt       bool
}

type DiffHunk struct {
	OldStart int
	OldCount int
	NewStart int
	NewCount int
	Lines    []string
}

// AddSource appends one source reference to the draft, validates the rendered
// preview, and optionally writes the manifest and generated output.
func (s *Session) AddSource(options AddOptions, dryRun bool) (*AddResult, error) {
	oldOutput := s.Output
	oldManifest, err := manifest.Marshal(s.Draft)
	if err != nil {
		return nil, err
	}
	sourceNode, err := s.SourceNode(options.Reference)
	if err != nil {
		return nil, err
	}
	heading := strings.TrimSpace(options.Heading)
	if heading == "" {
		heading = sourceNode.Heading
	}
	entry := manifest.Entry{
		Heading: heading,
		From:    []string{sourceNode.Reference},
	}
	parentPath, siblings, err := s.addTarget(options)
	if err != nil {
		return nil, err
	}
	for _, sibling := range *siblings {
		if sibling.Heading == heading {
			return nil, fmt.Errorf("heading %q already exists under %q; use --heading to choose a different rendered heading", heading, parentLabel(parentPath))
		}
	}
	*siblings = append(*siblings, entry)
	if err := s.MarkDraftChanged(); err != nil {
		return nil, err
	}
	newManifest, err := manifest.Marshal(s.Draft)
	if err != nil {
		return nil, err
	}
	result := &AddResult{
		Heading:      heading,
		Reference:    sourceNode.Reference,
		ParentPath:   parentPath,
		Preview:      s.Output,
		Section:      addedSection(heading, sourceNode.Content, parentPath),
		Tree:         manifestTree(s.Draft.Doc, heading, sourceNode.Reference),
		ManifestHunk: addedHunk(string(oldManifest), string(newManifest)),
		OutputHunk:   addedHunk(oldOutput, s.Output),
	}
	if dryRun {
		return result, nil
	}
	if options.Rebuild {
		if err := s.SaveAndBuild(); err != nil {
			return nil, err
		}
		result.WroteManifest = true
		result.Rebuilt = true
		return result, nil
	}
	if err := manifest.WriteAtomically(s.ManifestPath, s.Draft); err != nil {
		return nil, err
	}
	s.Saved = s.Draft.Clone()
	s.Draft = s.Saved.Clone()
	s.Dirty = false
	rendered, err := render.Build(s.Saved, s.ManifestPath)
	if err != nil {
		return nil, err
	}
	s.Output = rendered.Content
	result.WroteManifest = true
	return result, nil
}

func (s *Session) addTarget(options AddOptions) (string, *[]manifest.Entry, error) {
	if strings.TrimSpace(options.Under) == "" {
		if !options.Append {
			return "", nil, fmt.Errorf("add requires --under <manifest-heading-path> or --append")
		}
		return "", &s.Draft.Doc, nil
	}
	matches := findEntryPaths(s.Draft.Doc, splitManifestPath(options.Under), nil)
	if len(matches) == 0 {
		return "", nil, fmt.Errorf("manifest heading path %q was not found", options.Under)
	}
	if len(matches) > 1 {
		var paths []string
		for _, match := range matches {
			paths = append(paths, strings.Join(match.Headings, "/"))
		}
		return "", nil, fmt.Errorf("manifest heading path %q is ambiguous; matches: %s", options.Under, strings.Join(paths, ", "))
	}
	match := matches[0]
	if len(match.Entry.From) > 0 {
		return "", nil, fmt.Errorf("manifest heading path %q is a source entry and cannot contain children", options.Under)
	}
	return strings.Join(match.Headings, "/"), &match.Entry.Children, nil
}

type entryPath struct {
	Entry    *manifest.Entry
	Headings []string
}

func findEntryPaths(entries []manifest.Entry, target []string, parents []string) []entryPath {
	var matches []entryPath
	for index := range entries {
		entry := &entries[index]
		headings := append(append([]string(nil), parents...), entry.Heading)
		if samePath(headings, target) {
			matches = append(matches, entryPath{Entry: entry, Headings: headings})
		}
		matches = append(matches, findEntryPaths(entry.Children, target, headings)...)
	}
	return matches
}

func splitManifestPath(path string) []string {
	parts := strings.Split(path, "/")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func samePath(headings []string, target []string) bool {
	if len(headings) != len(target) {
		return false
	}
	for index := range headings {
		if headings[index] != target[index] {
			return false
		}
	}
	return true
}

func parentLabel(path string) string {
	if path == "" {
		return "document root"
	}
	return path
}

func addedSection(heading string, content string, parentPath string) string {
	var output strings.Builder
	level := 1
	if parentPath != "" {
		level = len(splitManifestPath(parentPath)) + 1
	}
	output.WriteString(strings.Repeat("#", level))
	output.WriteByte(' ')
	output.WriteString(heading)
	output.WriteString("\n\n")
	if body := strings.TrimSpace(content); body != "" {
		output.WriteString(body)
		output.WriteString("\n")
	}
	return output.String()
}

func manifestTree(entries []manifest.Entry, addedHeading string, addedReference string) string {
	var output strings.Builder
	writeManifestTree(&output, entries, "", addedHeading, addedReference)
	return strings.TrimRight(output.String(), "\n")
}

func writeManifestTree(output *strings.Builder, entries []manifest.Entry, prefix string, addedHeading string, addedReference string) {
	for index, entry := range entries {
		last := index == len(entries)-1
		marker := "  "
		if entry.Heading == addedHeading && len(entry.From) == 1 && entry.From[0] == addedReference {
			marker = "+ "
		}
		output.WriteString(marker)
		output.WriteString(prefix)
		if last {
			output.WriteString("`- ")
		} else {
			output.WriteString("|- ")
		}
		output.WriteString(entry.Heading)
		if len(entry.From) > 0 {
			output.WriteString("  <")
			output.WriteString(strings.Join(entry.From, ", "))
			output.WriteString(">")
		}
		output.WriteByte('\n')
		nextPrefix := prefix
		if last {
			nextPrefix += "   "
		} else {
			nextPrefix += "|  "
		}
		writeManifestTree(output, entry.Children, nextPrefix, addedHeading, addedReference)
	}
}

func addedHunk(oldText string, newText string) DiffHunk {
	oldLines := splitPatchLines(oldText)
	newLines := splitPatchLines(newText)
	prefix := 0
	for prefix < len(oldLines) && prefix < len(newLines) && oldLines[prefix] == newLines[prefix] {
		prefix++
	}
	suffix := 0
	for suffix < len(oldLines)-prefix && suffix < len(newLines)-prefix &&
		oldLines[len(oldLines)-1-suffix] == newLines[len(newLines)-1-suffix] {
		suffix++
	}
	oldCount := len(oldLines) - prefix - suffix
	newCount := len(newLines) - prefix - suffix
	return DiffHunk{
		OldStart: hunkOldStart(prefix, oldCount),
		OldCount: oldCount,
		NewStart: prefix + 1,
		NewCount: newCount,
		Lines:    append([]string(nil), newLines[prefix:prefix+newCount]...),
	}
}

func splitPatchLines(text string) []string {
	trimmed := strings.TrimRight(text, "\n")
	if trimmed == "" {
		return nil
	}
	return strings.Split(trimmed, "\n")
}

func hunkOldStart(prefix int, oldCount int) int {
	if oldCount == 0 {
		return prefix
	}
	return prefix + 1
}
