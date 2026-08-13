package cli

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/Qu1ncyRy4n/Agents/internal/library"
	"github.com/Qu1ncyRy4n/Agents/internal/workspace"
)

type inventoryOptions struct {
	Coverage   bool
	Tree       bool
	UnusedOnly bool
	ShowTLDR   bool
	Align      bool
	Chars      string
	Width      int
}

type inventoryRow struct {
	node  workspace.SourceNode
	state workspace.CoverageState
}

type treeChars struct {
	branch, last, vertical, blank string
	directory, heading, both      string
}

func displayChars(name string) treeChars {
	if name == "unicode" {
		return treeChars{
			branch: "├── ", last: "└── ", vertical: "│   ", blank: "    ",
			directory: "▸", heading: "¶", both: "◈",
		}
	}
	return treeChars{
		branch: "|-- ", last: "`-- ", vertical: "|   ", blank: "    ",
		directory: "/", heading: "#", both: "/#",
	}
}

func writeAlignedSourceList(stdout io.Writer, nodes []workspace.SourceNode, align, showFile, showLine bool, charSet string) error {
	chars := displayChars(charSet)
	if _, err := fmt.Fprintf(stdout, "Key: %s directory  %s heading  %s directory + heading\n", chars.directory, chars.heading, chars.both); err != nil {
		return fmt.Errorf("write source list: %w", err)
	}
	rows := make([][]string, 0, len(nodes))
	for _, node := range nodes {
		location := ""
		if showFile {
			location = node.File
			if showLine {
				location += fmt.Sprintf(":%d", node.Line)
			}
		} else if showLine {
			location = fmt.Sprintf("line %d", node.Line)
		}
		priority := ""
		if node.Metadata.Priority != nil {
			priority = fmt.Sprintf("p=%.2f", *node.Metadata.Priority)
		}
		rows = append(rows, []string{
			nodeKindMarker(node.Kind, chars), node.Reference, node.Heading,
			strings.Join(node.Metadata.Tags, ","), priority, node.Metadata.TLDR, location,
		})
	}
	return writeColumns(stdout, rows, align)
}

func writeInventory(stdout io.Writer, nodes []workspace.SourceNode, coverage workspace.Coverage, options inventoryOptions) error {
	chars := displayChars(options.Chars)
	if _, err := fmt.Fprintf(stdout, "Key: %s directory  %s heading  %s directory + heading", chars.directory, chars.heading, chars.both); err != nil {
		return fmt.Errorf("write source inventory: %w", err)
	}
	if options.Coverage {
		if _, err := fmt.Fprint(stdout, "  [included|inherited|partial|excluded|unused]"); err != nil {
			return fmt.Errorf("write source inventory: %w", err)
		}
	}
	if _, err := fmt.Fprintln(stdout); err != nil {
		return fmt.Errorf("write source inventory: %w", err)
	}

	stateByReference := make(map[string]workspace.CoverageState)
	visible := make(map[string]bool)
	for _, source := range coverage.Sources {
		for _, node := range source.Nodes {
			reference := source.Alias + ":" + node.Path
			stateByReference[reference] = node.State
			visible[reference] = true
		}
	}
	rowsByAlias := make(map[string][]inventoryRow)
	aliasOrder := make([]string, 0)
	for _, node := range nodes {
		if !visible[node.Reference] {
			continue
		}
		state := stateByReference[node.Reference]
		if options.UnusedOnly && state != workspace.CoverageUnused && !hasUnusedDescendant(node.Reference, stateByReference) {
			continue
		}
		if _, found := rowsByAlias[node.Alias]; !found {
			aliasOrder = append(aliasOrder, node.Alias)
		}
		rowsByAlias[node.Alias] = append(rowsByAlias[node.Alias], inventoryRow{node: node, state: state})
	}
	for _, alias := range aliasOrder {
		rows := rowsByAlias[alias]
		if _, err := fmt.Fprintf(stdout, "\n%s\n", alias); err != nil {
			return fmt.Errorf("write source inventory: %w", err)
		}
		if options.Tree {
			if err := writeInventoryTree(stdout, rows, chars, options); err != nil {
				return err
			}
			continue
		}
		columns := make([][]string, 0, len(rows))
		for _, row := range rows {
			columns = append(columns, inventoryColumns(row, chars, options, ""))
		}
		if err := writeColumns(stdout, columns, options.Align); err != nil {
			return err
		}
	}
	return nil
}

func hasUnusedDescendant(reference string, states map[string]workspace.CoverageState) bool {
	prefix := reference + "/"
	for candidate, state := range states {
		if state == workspace.CoverageUnused && strings.HasPrefix(candidate, prefix) {
			return true
		}
	}
	return false
}

func writeInventoryTree(stdout io.Writer, rows []inventoryRow, chars treeChars, options inventoryOptions) error {
	prefixes := make([]string, len(rows))
	labels := make([]string, len(rows))
	maxLabel := 0
	for index, row := range rows {
		depth := inventoryVisibleDepth(rows, index)
		var prefix strings.Builder
		for ancestorDepth := 0; ancestorDepth < depth; ancestorDepth++ {
			if inventoryAncestorHasLaterSibling(rows, index, ancestorDepth) {
				prefix.WriteString(chars.vertical)
			} else {
				prefix.WriteString(chars.blank)
			}
		}
		if inventoryHasNextAtDepth(rows, index, depth) {
			prefix.WriteString(chars.branch)
		} else {
			prefix.WriteString(chars.last)
		}
		name := row.node.Heading
		if row.node.Kind == library.NodeDirectory {
			name += "/"
		}
		prefixes[index] = prefix.String()
		labels[index] = nodeKindMarker(row.node.Kind, chars) + " " + name
		width := utf8.RuneCountInString(prefixes[index] + labels[index])
		if width > maxLabel {
			maxLabel = width
		}
	}
	seenTLDRFiles := make(map[string]bool)
	outputRows := make([][]string, 0, len(rows))
	for index, row := range rows {
		label := prefixes[index] + labels[index]
		if options.Align {
			label += strings.Repeat(" ", maxLabel-utf8.RuneCountInString(label))
		}
		columns := inventoryColumns(row, chars, options, label)
		if options.ShowTLDR && (row.node.Metadata.TLDR == "" || seenTLDRFiles[row.node.File]) {
			columns[len(columns)-1] = ""
		} else if options.ShowTLDR {
			seenTLDRFiles[row.node.File] = true
		}
		outputRows = append(outputRows, columns)
	}
	if options.Width > 0 {
		return writeFittedTreeRows(stdout, outputRows, options.Width, options.ShowTLDR)
	}
	return writeColumns(stdout, outputRows, options.Align)
}

func writeFittedTreeRows(stdout io.Writer, rows [][]string, width int, hasTLDR bool) error {
	baseWidths := fittedBaseWidths(rows, hasTLDR)
	for _, row := range rows {
		base := row
		tldr := ""
		if hasTLDR && len(row) > 0 {
			base = row[:len(row)-1]
			tldr = row[len(row)-1]
		}
		baseLine := fittedBaseLine(base, baseWidths)
		continuation := "    "
		if len(base) > 0 {
			continuation = treeContinuationPrefix(base[0])
		}
		fullLine := baseLine
		if tldr != "" {
			fullLine += "  " + tldr
		}
		if utf8.RuneCountInString(fullLine) <= width {
			if _, err := fmt.Fprintln(stdout, fullLine); err != nil {
				return fmt.Errorf("write fitted tree: %w", err)
			}
			continue
		}
		if utf8.RuneCountInString(baseLine) <= width {
			if _, err := fmt.Fprintln(stdout, baseLine); err != nil {
				return fmt.Errorf("write fitted tree: %w", err)
			}
		} else if len(base) > 1 {
			if _, err := fmt.Fprintln(stdout, fittedBaseLine(base[:len(base)-1], baseWidths)); err != nil {
				return fmt.Errorf("write fitted tree: %w", err)
			}
			if err := writeWrappedField(stdout, continuation, "Ref: ", base[len(base)-1], width); err != nil {
				return err
			}
		} else if _, err := fmt.Fprintln(stdout, baseLine); err != nil {
			return fmt.Errorf("write fitted tree: %w", err)
		}
		if tldr != "" {
			if err := writeWrappedField(stdout, continuation, "TLDR: ", tldr, width); err != nil {
				return err
			}
		}
	}
	return nil
}

func fittedBaseWidths(rows [][]string, hasTLDR bool) []int {
	var widths []int
	for _, row := range rows {
		limit := len(row)
		if hasTLDR && limit > 0 {
			limit--
		}
		for column := 0; column < limit; column++ {
			if column >= len(widths) {
				widths = append(widths, 0)
			}
			if width := utf8.RuneCountInString(row[column]); width > widths[column] {
				widths[column] = width
			}
		}
	}
	return widths
}

func fittedBaseLine(values []string, widths []int) string {
	var output strings.Builder
	for column, value := range values {
		if column > 0 {
			output.WriteString("  ")
		}
		output.WriteString(value)
		if column < len(values)-1 && column < len(widths) {
			output.WriteString(strings.Repeat(" ", widths[column]-utf8.RuneCountInString(value)))
		}
	}
	return output.String()
}

func treeContinuationPrefix(label string) string {
	replacements := []struct {
		connector    string
		continuation string
	}{
		{"|-- ", "|   "},
		{"`-- ", "    "},
		{"├── ", "│   "},
		{"└── ", "    "},
	}
	bestIndex := -1
	result := "    "
	for _, replacement := range replacements {
		if index := strings.LastIndex(label, replacement.connector); index > bestIndex {
			bestIndex = index
			result = label[:index] + replacement.continuation
		}
	}
	return result
}

func writeWrappedField(stdout io.Writer, structuralPrefix, label, value string, width int) error {
	prefix := structuralPrefix + label
	continuation := structuralPrefix + strings.Repeat(" ", utf8.RuneCountInString(label))
	lines := wrapWords(value, width-utf8.RuneCountInString(prefix))
	for index, line := range lines {
		linePrefix := continuation
		if index == 0 {
			linePrefix = prefix
		}
		if _, err := fmt.Fprintln(stdout, linePrefix+line); err != nil {
			return fmt.Errorf("write fitted tree: %w", err)
		}
	}
	return nil
}

func wrapWords(value string, width int) []string {
	if width < 8 {
		width = 8
	}
	words := strings.Fields(value)
	if len(words) == 0 {
		return nil
	}
	lines := []string{words[0]}
	for _, word := range words[1:] {
		last := len(lines) - 1
		if utf8.RuneCountInString(lines[last])+1+utf8.RuneCountInString(word) <= width {
			lines[last] += " " + word
			continue
		}
		lines = append(lines, word)
	}
	return lines
}

func inventoryColumns(row inventoryRow, chars treeChars, options inventoryOptions, label string) []string {
	if label == "" {
		label = nodeKindMarker(row.node.Kind, chars) + " " + row.node.Heading
	}
	columns := []string{label}
	if options.Coverage {
		columns = append(columns, "["+string(row.state)+"]")
	}
	columns = append(columns, row.node.Reference)
	if options.ShowTLDR {
		columns = append(columns, row.node.Metadata.TLDR)
	}
	return columns
}

func nodeKindMarker(kind library.NodeKind, chars treeChars) string {
	switch kind {
	case library.NodeDirectory:
		return chars.directory
	case library.NodeDirectoryHeading:
		return chars.both
	default:
		return chars.heading
	}
}

func writeColumns(stdout io.Writer, rows [][]string, align bool) error {
	widths := make([]int, 0)
	if align {
		for _, row := range rows {
			for column, value := range row {
				if column >= len(widths) {
					widths = append(widths, 0)
				}
				if width := utf8.RuneCountInString(value); width > widths[column] {
					widths[column] = width
				}
			}
		}
	}
	for _, row := range rows {
		last := len(row) - 1
		for last >= 0 && row[last] == "" {
			last--
		}
		for column := 0; column <= last; column++ {
			if column > 0 {
				if _, err := fmt.Fprint(stdout, "  "); err != nil {
					return fmt.Errorf("write aligned columns: %w", err)
				}
			}
			value := row[column]
			if _, err := fmt.Fprint(stdout, value); err != nil {
				return fmt.Errorf("write aligned columns: %w", err)
			}
			if align && column < last {
				if _, err := fmt.Fprint(stdout, strings.Repeat(" ", widths[column]-utf8.RuneCountInString(value))); err != nil {
					return fmt.Errorf("write aligned columns: %w", err)
				}
			}
		}
		if _, err := fmt.Fprintln(stdout); err != nil {
			return fmt.Errorf("write aligned columns: %w", err)
		}
	}
	return nil
}

func inventoryVisibleDepth(rows []inventoryRow, index int) int {
	depth := 0
	for previous := 0; previous < index; previous++ {
		if rows[previous].node.Alias == rows[index].node.Alias && strings.HasPrefix(rows[index].node.Path, rows[previous].node.Path+"/") {
			depth++
		}
	}
	return depth
}

func inventoryHasNextAtDepth(rows []inventoryRow, index, depth int) bool {
	for next := index + 1; next < len(rows); next++ {
		nextDepth := inventoryVisibleDepth(rows, next)
		if nextDepth == depth {
			return true
		}
		if nextDepth < depth {
			return false
		}
	}
	return false
}

func inventoryAncestorHasLaterSibling(rows []inventoryRow, index, ancestorDepth int) bool {
	for previous := index - 1; previous >= 0; previous-- {
		if inventoryVisibleDepth(rows, previous) != ancestorDepth || !strings.HasPrefix(rows[index].node.Path, rows[previous].node.Path+"/") {
			continue
		}
		return inventoryHasNextAtDepth(rows, previous, ancestorDepth)
	}
	return false
}
