package cli

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

type treeChars struct {
	branch, last, vertical, blank string
}

func displayChars(name string) treeChars {
	if name == "unicode" {
		return treeChars{branch: "├── ", last: "└── ", vertical: "│   ", blank: "    "}
	}
	return treeChars{branch: "|-- ", last: "`-- ", vertical: "|   ", blank: "    "}
}

func writeColumns(stdout io.Writer, rows [][]string, align bool) error {
	widths := make([]int, 0)
	if align {
		for _, row := range rows {
			for column, value := range row {
				if column >= len(widths) {
					widths = append(widths, 0)
				}
				if width := displayWidth(value); width > widths[column] {
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
				if _, err := fmt.Fprint(stdout, strings.Repeat(" ", widths[column]-displayWidth(value))); err != nil {
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

// displayWidth ignores ANSI control sequences so styled terminal cells retain
// their column alignment while redirected output can remain plain text.
func displayWidth(value string) int {
	width := 0
	for len(value) > 0 {
		if strings.HasPrefix(value, "\x1b[") {
			consumed := false
			for index := 2; index < len(value); index++ {
				if value[index] >= '@' && value[index] <= '~' {
					value = value[index+1:]
					consumed = true
					break
				}
			}
			if consumed {
				continue
			}
		}
		_, size := utf8.DecodeRuneInString(value)
		width++
		value = value[size:]
	}
	return width
}
