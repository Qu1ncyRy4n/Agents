package workspace

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/internal/library"
	"github.com/Qu1ncyRy4n/Agents/internal/manifest"
	"github.com/Qu1ncyRy4n/Agents/internal/render"
)

type AddOptions struct {
	Reference string
	Under     string
	Append    bool
	First     bool
	Last      bool
	Before    string
	After     string
	Heading   string
	Rebuild   bool
}

type AddResult struct {
	Heading       string
	Reference     string
	ParentPath    string
	Placement     string
	Preview       string
	Section       string
	Tree          string
	SourceTree    string
	Relations     []AddRelation
	ManifestHunk  DiffHunk
	OutputHunk    DiffHunk
	WroteManifest bool
	Rebuilt       bool
}

type AddRelation struct {
	Reference string
	Reason    string
	Exact     bool
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
	relations := s.addRelations(sourceNode.Reference)
	heading := strings.TrimSpace(options.Heading)
	if heading == "" {
		heading = sourceNode.Heading
	}
	entry := manifest.Entry{
		Heading: heading,
		From:    []string{sourceNode.Reference},
	}
	parentPath, siblings, insertion, err := s.addTarget(options)
	if err != nil {
		return nil, err
	}
	for _, sibling := range *siblings {
		if sibling.Heading == heading {
			return nil, fmt.Errorf("heading %q already exists under %q; use --heading to choose a different rendered heading", heading, parentLabel(parentPath))
		}
	}
	*siblings = insertManifestEntry(*siblings, insertion, entry)
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
		Placement:    addPlacementDescription(options, parentPath),
		Preview:      s.Output,
		Section:      addedSection(heading, sourceNode.Content, parentPath),
		Tree:         manifestTree(s.Draft.Doc, heading, sourceNode.Reference),
		SourceTree:   s.sourceTree(sourceNode.Reference),
		Relations:    relations,
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

func addPlacementDescription(options AddOptions, parentPath string) string {
	if target := strings.TrimSpace(options.Before); target != "" {
		return "before " + target
	}
	if target := strings.TrimSpace(options.After); target != "" {
		return "after " + target
	}
	if options.Append {
		return "last at document root"
	}
	position := "last"
	if options.First {
		position = "first"
	}
	return position + " under " + parentLabel(parentPath)
}

func (s *Session) addTarget(options AddOptions) (string, *[]manifest.Entry, int, error) {
	under := strings.TrimSpace(options.Under)
	before := strings.TrimSpace(options.Before)
	after := strings.TrimSpace(options.After)
	locationCount := 0
	for _, selected := range []bool{under != "", options.Append, before != "", after != ""} {
		if selected {
			locationCount++
		}
	}
	if locationCount != 1 {
		return "", nil, 0, fmt.Errorf("add requires exactly one placement: --under, --append, --before, or --after")
	}
	if options.First && options.Last {
		return "", nil, 0, fmt.Errorf("add cannot use --first and --last together")
	}
	if (options.First || options.Last) && under == "" {
		return "", nil, 0, fmt.Errorf("--first and --last require --under")
	}
	if before != "" || after != "" {
		target := before
		if target == "" {
			target = after
		}
		locations := findEntryLocations(&s.Draft.Doc, splitManifestPath(target), nil)
		if len(locations) == 0 {
			if suggestion := closestManifestHeadingPath(target, s.ManifestHeadingPaths()); suggestion != "" {
				return "", nil, 0, fmt.Errorf("manifest heading path %q was not found; did you mean %q?", target, suggestion)
			}
			return "", nil, 0, fmt.Errorf("manifest heading path %q was not found", target)
		}
		if len(locations) > 1 {
			return "", nil, 0, fmt.Errorf("manifest heading path %q is ambiguous", target)
		}
		location := locations[0]
		index := location.Index
		if after != "" {
			index++
		}
		return strings.Join(location.Parents, "/"), location.Entries, index, nil
	}
	if options.Append {
		return "", &s.Draft.Doc, len(s.Draft.Doc), nil
	}
	matches := findEntryPaths(s.Draft.Doc, splitManifestPath(under), nil)
	if len(matches) == 0 {
		if suggestion := closestManifestHeadingPath(under, s.ManifestHeadingPaths()); suggestion != "" {
			return "", nil, 0, fmt.Errorf("manifest heading path %q was not found; did you mean %q?", under, suggestion)
		}
		return "", nil, 0, fmt.Errorf("manifest heading path %q was not found", under)
	}
	if len(matches) > 1 {
		var paths []string
		for _, match := range matches {
			paths = append(paths, strings.Join(match.Headings, "/"))
		}
		return "", nil, 0, fmt.Errorf("manifest heading path %q is ambiguous; matches: %s", under, strings.Join(paths, ", "))
	}
	match := matches[0]
	if len(match.Entry.From) > 0 {
		return "", nil, 0, fmt.Errorf("manifest heading path %q is a source entry and cannot contain children", under)
	}
	children := &match.Entry.Children
	insertion := len(*children)
	if options.First {
		insertion = 0
	}
	return strings.Join(match.Headings, "/"), children, insertion, nil
}

func insertManifestEntry(entries []manifest.Entry, index int, entry manifest.Entry) []manifest.Entry {
	entries = append(entries, manifest.Entry{})
	copy(entries[index+1:], entries[index:])
	entries[index] = entry
	return entries
}

type entryLocation struct {
	Entries *[]manifest.Entry
	Index   int
	Parents []string
}

func findEntryLocations(entries *[]manifest.Entry, target []string, parents []string) []entryLocation {
	var matches []entryLocation
	for index := range *entries {
		entry := &(*entries)[index]
		headings := append(append([]string(nil), parents...), entry.Heading)
		if samePath(headings, target) {
			matches = append(matches, entryLocation{Entries: entries, Index: index, Parents: append([]string(nil), parents...)})
		}
		matches = append(matches, findEntryLocations(&entry.Children, target, headings)...)
	}
	return matches
}

func (s *Session) sourceTree(reference string) string {
	alias, path, err := manifest.SplitReference(reference)
	if err != nil {
		return ""
	}
	node := s.Sources[alias].ByPath[path]
	var output strings.Builder
	writeSourceSelectionTree(&output, node, "", true)
	return strings.TrimRight(output.String(), "\n")
}

func writeSourceSelectionTree(output *strings.Builder, node *library.Node, prefix string, last bool) {
	output.WriteString(prefix)
	if last {
		output.WriteString("`-- ")
	} else {
		output.WriteString("|-- ")
	}
	switch node.Kind {
	case library.NodeDirectory:
		output.WriteString("/ ")
	case library.NodeDirectoryHeading:
		output.WriteString("/# ")
	default:
		output.WriteString("# ")
	}
	output.WriteString(node.Heading)
	if node.Kind == library.NodeDirectory {
		output.WriteByte('/')
	}
	output.WriteByte('\n')
	nextPrefix := prefix
	if last {
		nextPrefix += "    "
	} else {
		nextPrefix += "|   "
	}
	for index, child := range node.Children {
		writeSourceSelectionTree(output, child, nextPrefix, index == len(node.Children)-1)
	}
}

func (s *Session) addRelations(reference string) []AddRelation {
	alias, path, err := manifest.SplitReference(reference)
	if err != nil {
		return nil
	}
	var relations []AddRelation
	for _, existing := range manifestReferences(s.Draft.Doc) {
		existingAlias, existingPath, splitErr := manifest.SplitReference(existing)
		if splitErr != nil {
			continue
		}
		if existingAlias == alias && (existingPath == path || strings.HasPrefix(existingPath, path+"/") || strings.HasPrefix(path, existingPath+"/")) {
			relations = append(relations, AddRelation{Reference: existing, Reason: "overlapping source subtree", Exact: true})
			continue
		}
		if segment := sharedPathSegment(path, existingPath); segment != "" {
			relations = append(relations, AddRelation{Reference: existing, Reason: "shared source group " + segment})
		}
	}
	sort.SliceStable(relations, func(first, second int) bool {
		if relations[first].Exact != relations[second].Exact {
			return relations[first].Exact
		}
		return relations[first].Reference < relations[second].Reference
	})
	return relations
}

func manifestReferences(entries []manifest.Entry) []string {
	var references []string
	for _, entry := range entries {
		references = append(references, entry.From...)
		references = append(references, manifestReferences(entry.Children)...)
	}
	return references
}

func sharedPathSegment(first, second string) string {
	ignored := map[string]bool{"shared-baseline": true, "instructions": true, "constraints": true, "format": true, "identity": true, "lang": true}
	secondSegments := make(map[string]bool)
	for _, segment := range strings.Split(second, "/") {
		secondSegments[segment] = true
	}
	for _, segment := range strings.Split(first, "/") {
		if !ignored[segment] && secondSegments[segment] {
			return segment
		}
	}
	return ""
}

func closestManifestHeadingPath(path string, candidates []string) string {
	best := ""
	bestDistance := 8
	target := strings.ToLower(path)
	for _, candidate := range candidates {
		distance := levenshtein(target, strings.ToLower(candidate))
		if strings.Contains(strings.ToLower(candidate), target) || strings.Contains(target, strings.ToLower(candidate)) {
			distance--
		}
		if distance < bestDistance {
			best = candidate
			bestDistance = distance
		}
	}
	return best
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
