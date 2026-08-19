package manifest_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/manifest"
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

func TestLoadNormalizesCompactAndExplicitSources(t *testing.T) {
	temporary := t.TempDir()
	path := filepath.Join(temporary, "agents.yaml")
	content := "sources:\n  local: ./library\n  shared:\n    location: https://example.com/agents.git\n    subdir: libraries/cdint\ndoc:\n  - Rules: shared:rules\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	value, _, err := manifest.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if value.Sources["local"].Location != "./library" || value.Sources["local"].Subdir != "" {
		t.Fatalf("compact source = %#v", value.Sources["local"])
	}
	if value.Sources["shared"].Location != "https://example.com/agents.git" || value.Sources["shared"].Subdir != "libraries/cdint" {
		t.Fatalf("explicit source = %#v", value.Sources["shared"])
	}
	encoded, err := manifest.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(encoded), "local: ./library") || !strings.Contains(string(encoded), "subdir: libraries/cdint") {
		t.Fatalf("encoded manifest:\n%s", encoded)
	}
}

func TestLoadRejectsUnsafeSourceSubdirs(t *testing.T) {
	temporary := t.TempDir()
	path := filepath.Join(temporary, "agents.yaml")
	for _, subdir := range []string{"/absolute", "../escape", "libraries\\cdint", "libraries//cdint"} {
		content := "sources:\n  shared:\n    location: https://example.com/agents.git\n    subdir: " + subdir + "\ndoc:\n  - Rules: shared:rules\n"
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, _, err := manifest.Load(path); err == nil || !strings.Contains(err.Error(), "subdir") {
			t.Fatalf("subdir %q error = %v", subdir, err)
		}
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
