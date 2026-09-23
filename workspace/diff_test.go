package workspace

import "testing"

func TestUnifiedDiffReportsChangedMarkdown(t *testing.T) {
	got := UnifiedDiff("# Title\nold\n", "# Title\nnew\n")
	want := "--- expected\n+++ actual\n@@ -1,2 +1,2 @@\n # Title\n-old\n+new\n"
	if got != want {
		t.Fatalf("diff = %q, want %q", got, want)
	}
}

func TestUnifiedDiffIsEmptyForMatchingMarkdown(t *testing.T) {
	if got := UnifiedDiff("# Title\n", "# Title\n"); got != "" {
		t.Fatalf("diff = %q", got)
	}
}
