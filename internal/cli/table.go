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
