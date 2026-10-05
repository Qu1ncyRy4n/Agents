// Package textdiff renders small, deterministic unified diffs for generated
// text files. It is shared by v1 drift reporting and v2 plan output.
package textdiff

import (
	"fmt"
	"strings"
)

type line struct {
	text    string
	newline bool
}

// Unified returns a unified diff from before to after, or an empty string when
// the two are identical. beforeLabel and afterLabel appear in the header.
func Unified(before, after, beforeLabel, afterLabel string) string {
	left, right := splitLines(before), splitLines(after)
	if sameLines(left, right) {
		return ""
	}
	lengths := make([][]int, len(left)+1)
	for index := range lengths {
		lengths[index] = make([]int, len(right)+1)
	}
	for i := len(left) - 1; i >= 0; i-- {
		for j := len(right) - 1; j >= 0; j-- {
			if left[i] == right[j] {
				lengths[i][j] = lengths[i+1][j+1] + 1
			} else if lengths[i+1][j] >= lengths[i][j+1] {
				lengths[i][j] = lengths[i+1][j]
			} else {
				lengths[i][j] = lengths[i][j+1]
			}
		}
	}
	var output strings.Builder
	fmt.Fprintf(&output, "--- %s\n+++ %s\n@@ -%s +%s @@\n", beforeLabel, afterLabel, hunkRange(1, len(left)), hunkRange(1, len(right)))
	for i, j := 0, 0; i < len(left) || j < len(right); {
		switch {
		case i < len(left) && j < len(right) && left[i] == right[j]:
			writeLine(&output, ' ', left[i])
			i, j = i+1, j+1
		case j < len(right) && (i == len(left) || lengths[i][j+1] > lengths[i+1][j]):
			writeLine(&output, '+', right[j])
			j++
		default:
			writeLine(&output, '-', left[i])
			i++
		}
	}
	return output.String()
}

func splitLines(value string) []line {
	if value == "" {
		return nil
	}
	parts := strings.Split(value, "\n")
	lines := make([]line, 0, len(parts))
	for index, text := range parts {
		if index == len(parts)-1 && text == "" && strings.HasSuffix(value, "\n") {
			continue
		}
		lines = append(lines, line{text: text, newline: index < len(parts)-1})
	}
	return lines
}

func sameLines(first, second []line) bool {
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

func hunkRange(start, count int) string {
	if count == 0 {
		return "0,0"
	}
	if count == 1 {
		return fmt.Sprint(start)
	}
	return fmt.Sprintf("%d,%d", start, count)
}

func writeLine(output *strings.Builder, prefix byte, value line) {
	output.WriteByte(prefix)
	output.WriteString(value.text)
	output.WriteByte('\n')
	if !value.newline {
		output.WriteString("\\ No newline at end of file\n")
	}
}
