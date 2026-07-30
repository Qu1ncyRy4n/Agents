package manifest_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/internal/manifest"
)

func TestLoadRejectsDuplicateKeysAndMixedEntryKinds(t *testing.T) {
	temporary := t.TempDir()
	path := filepath.Join(temporary, "agents.yaml")
	if err := os.WriteFile(path, []byte("sources:\n  shared: library\n  shared: other\noutput: AGENTS.md\ndoc:\n  - heading: Rules\n    from: [shared:rules]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, _, err := manifest.Load(path)
	if err == nil || !strings.Contains(err.Error(), "mapping key") {
		t.Fatalf("Load error = %v, want duplicate key error", err)
	}

	if err := os.WriteFile(path, []byte("sources:\n  shared: library\noutput: AGENTS.md\ndoc:\n  - heading: Rules\n    from: [shared:rules]\n    children:\n      - heading: Child\n        from: [shared:rules/child]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, _, err = manifest.Load(path)
	if err == nil || !strings.Contains(err.Error(), "exactly one") {
		t.Fatalf("Load error = %v, want mutually exclusive entry error", err)
	}
}

func TestLoadNormalizesCompactAndExplicitEntries(t *testing.T) {
	temporary := t.TempDir()
	path := filepath.Join(temporary, "agents.yaml")
	content := "sources:\n  shared: library\noutput: AGENTS.md\ndoc:\n  - Identity: shared:identity\n  - heading: Testing\n    from: [shared:testing]\n    exclude: [shared:testing/flaky-retries]\n  - heading: Instructions\n    children:\n      - Format: shared:format\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	value, _, err := manifest.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := value.Doc[0]; got.Heading != "Identity" || len(got.From) != 1 || got.From[0] != "shared:identity" {
		t.Fatalf("compact scalar = %#v", got)
	}
	testing := value.Doc[1]
	if testing.Heading != "Testing" || testing.From[0] != "shared:testing" || testing.Exclude[0] != "shared:testing/flaky-retries" {
		t.Fatalf("explicit entry = %#v", testing)
	}
}

func TestLoadDefaultsOutputAndRejectsAliases(t *testing.T) {
	temporary := t.TempDir()
	path := filepath.Join(temporary, "agents.yaml")
	content := "sources:\n  shared: library\ndoc:\n  - Identity: shared:identity\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	value, _, err := manifest.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if value.Output != "AGENTS.md" {
		t.Fatalf("default output = %q", value.Output)
	}
	if err := os.WriteFile(path, []byte("sources: &sources\n  shared: library\nother: *sources\ndoc: []\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, _, err = manifest.Load(path)
	if err == nil || !strings.Contains(err.Error(), "anchors") {
		t.Fatalf("anchor error = %v", err)
	}
}
