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

type operation struct {
	prefix byte
	line   line
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
	operations := make([]operation, 0, len(left)+len(right))
	for i, j := 0, 0; i < len(left) || j < len(right); {
		switch {
		case i < len(left) && j < len(right) && left[i] == right[j]:
			operations = append(operations, operation{prefix: ' ', line: left[i]})
			i, j = i+1, j+1
		case j < len(right) && (i == len(left) || lengths[i][j+1] > lengths[i+1][j]):
			operations = append(operations, operation{prefix: '+', line: right[j]})
			j++
		default:
			operations = append(operations, operation{prefix: '-', line: left[i]})
			i++
		}
	}
	var output strings.Builder
	fmt.Fprintf(&output, "--- %s\n+++ %s\n", beforeLabel, afterLabel)
	writeHunks(&output, operations, 3)
	return output.String()
}

// writeHunks writes conventional unified-diff hunks: changed lines plus a
// small amount of surrounding context. Keeping distant edits separate makes
// generated-output drift readable without concealing either change.
func writeHunks(output *strings.Builder, operations []operation, context int) {
	include := make([]bool, len(operations))
	for index, operation := range operations {
		if operation.prefix == ' ' {
			continue
		}
		start, end := index-context, index+context
		if start < 0 {
			start = 0
		}
		if end >= len(operations) {
			end = len(operations) - 1
		}
		for candidate := start; candidate <= end; candidate++ {
			include[candidate] = true
		}
	}
	oldBefore, newBefore := 0, 0
	for index := 0; index < len(operations); {
		if !include[index] {
			advancePosition(operations[index], &oldBefore, &newBefore)
			index++
			continue
		}
		start, oldStart, newStart := index, oldBefore+1, newBefore+1
		oldCount, newCount := 0, 0
		for index < len(operations) && include[index] {
			operation := operations[index]
			if operation.prefix != '+' {
				oldCount++
			}
			if operation.prefix != '-' {
				newCount++
			}
			advancePosition(operation, &oldBefore, &newBefore)
			index++
		}
		if oldCount == 0 {
			oldStart--
		}
		if newCount == 0 {
			newStart--
		}
		fmt.Fprintf(output, "@@ -%s +%s @@\n", hunkRange(oldStart, oldCount), hunkRange(newStart, newCount))
		for _, operation := range operations[start:index] {
			writeLine(output, operation.prefix, operation.line)
		}
	}
}

func advancePosition(operation operation, oldPosition, newPosition *int) {
	if operation.prefix != '+' {
		*oldPosition++
	}
	if operation.prefix != '-' {
		*newPosition++
	}
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
