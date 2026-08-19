package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/starter"
)

func TestInitializeDryRunAndBuild(t *testing.T) {
	temporary := t.TempDir()
	writeStarterLibrary(t, filepath.Join(temporary, "library"))
	template, err := starter.Get("minimal")
	if err != nil {
		t.Fatal(err)
	}
	value, err := template.Manifest(map[string]string{"shared": "library"}, "AGENTS.md")
	if err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	result, err := Initialize(InitOptions{ManifestPath: manifestPath, Manifest: value, DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result.ManifestYAML, "shared: library") || !strings.Contains(result.Preview, "# Identity") {
		t.Fatalf("dry-run result = %#v", result)
	}
	if _, err := os.Stat(manifestPath); !os.IsNotExist(err) {
		t.Fatalf("dry run wrote manifest: %v", err)
	}
	result, err = Initialize(InitOptions{ManifestPath: manifestPath, Manifest: value, Build: true})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Wrote || !result.Built {
		t.Fatalf("result = %#v", result)
	}
	if _, err := os.Stat(filepath.Join(temporary, "AGENTS.md")); err != nil {
		t.Fatal(err)
	}
	if _, err := Initialize(InitOptions{ManifestPath: manifestPath, Manifest: value}); err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("existing manifest error = %v", err)
	}
}

func writeStarterLibrary(t *testing.T, path string) {
	t.Helper()
	if err := os.Mkdir(path, 0o755); err != nil {
		t.Fatal(err)
	}
	content := "# Shared Baseline\n\n## Identity\n\n### Role\nAct carefully.\n\n### Source Of Truth\nRead the design.\n\n## Instructions\n\n### Focused Change Loop\nKeep changes focused.\n\n## Constraints\n\n### Safe Defaults\nStay safe.\n\n## Format\n\n### Clear Handoff\nReport checks.\n"
	if err := os.WriteFile(filepath.Join(path, "baseline.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
