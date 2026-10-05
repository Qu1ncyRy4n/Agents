package workspace

import "github.com/Qu1ncyRy4n/Agents/textdiff"

// UnifiedDiff returns a deterministic unified diff from generated to on-disk Markdown.
func UnifiedDiff(expected, actual string) string {
	return textdiff.Unified(expected, actual, "expected", "actual")
}
