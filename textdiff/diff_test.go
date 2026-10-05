package textdiff

import "testing"

func TestUnifiedReportsChangedText(t *testing.T) {
	got := Unified("# Title\nold\n", "# Title\nnew\n", "expected", "actual")
	want := "--- expected\n+++ actual\n@@ -1,2 +1,2 @@\n # Title\n-old\n+new\n"
	if got != want {
		t.Fatalf("diff = %q, want %q", got, want)
	}
}

func TestUnifiedIsEmptyForMatchingText(t *testing.T) {
	if got := Unified("# Title\n", "# Title\n", "a", "b"); got != "" {
		t.Fatalf("diff = %q", got)
	}
}

func TestUnifiedRendersNewFileAsAdditions(t *testing.T) {
	got := Unified("", "one\ntwo\n", "/dev/null", "b/AGENTS.md")
	want := "--- /dev/null\n+++ b/AGENTS.md\n@@ -0,0 +1,2 @@\n+one\n+two\n"
	if got != want {
		t.Fatalf("diff = %q, want %q", got, want)
	}
}
