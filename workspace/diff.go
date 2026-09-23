package workspace

import (
	"fmt"
	"strings"
)

type diffLine struct {
	text    string
	newline bool
}

// UnifiedDiff returns a deterministic unified diff from generated to on-disk Markdown.
func UnifiedDiff(expected, actual string) string {
	before, after := splitDiffLines(expected), splitDiffLines(actual)
	if sameDiffLines(before, after) {
		return ""
	}
	lengths := make([][]int, len(before)+1)
	for index := range lengths {
		lengths[index] = make([]int, len(after)+1)
	}
	for i := len(before) - 1; i >= 0; i-- {
		for j := len(after) - 1; j >= 0; j-- {
			if before[i] == after[j] {
				lengths[i][j] = lengths[i+1][j+1] + 1
			} else if lengths[i+1][j] >= lengths[i][j+1] {
				lengths[i][j] = lengths[i+1][j]
			} else {
				lengths[i][j] = lengths[i][j+1]
			}
		}
	}
	var output strings.Builder
	fmt.Fprintf(&output, "--- expected\n+++ actual\n@@ -%s +%s @@\n", diffRange(1, len(before)), diffRange(1, len(after)))
	for i, j := 0, 0; i < len(before) || j < len(after); {
		switch {
		case i < len(before) && j < len(after) && before[i] == after[j]:
			writeDiffLine(&output, ' ', before[i])
			i, j = i+1, j+1
		case j < len(after) && (i == len(before) || lengths[i][j+1] > lengths[i+1][j]):
			writeDiffLine(&output, '+', after[j])
			j++
		default:
			writeDiffLine(&output, '-', before[i])
			i++
		}
	}
	return output.String()
}

func splitDiffLines(value string) []diffLine {
	if value == "" {
		return nil
	}
	parts := strings.Split(value, "\n")
	lines := make([]diffLine, 0, len(parts))
	for index, text := range parts {
		if index == len(parts)-1 && text == "" && strings.HasSuffix(value, "\n") {
			continue
		}
		lines = append(lines, diffLine{text: text, newline: index < len(parts)-1})
	}
	return lines
}

func sameDiffLines(first, second []diffLine) bool {
	if len(first) != len(second) {
		return false
	}
	for index := range first {
		if first[index] != second[index] {
			return false
		}
	}
	return true
}

func diffRange(start, count int) string {
	if count == 0 {
		return "0,0"
	}
	if count == 1 {
		return fmt.Sprint(start)
	}
	return fmt.Sprintf("%d,%d", start, count)
}

func writeDiffLine(output *strings.Builder, prefix byte, line diffLine) {
	output.WriteByte(prefix)
	output.WriteString(line.text)
	output.WriteByte('\n')
	if !line.newline {
		output.WriteString("\\ No newline at end of file\n")
	}
}
