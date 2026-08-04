package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/internal/cli"
)

func TestRunBuildWritesConfiguredOutput(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(libraryPath, "identity.md"), []byte("# Identity\nHello.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	manifest := "sources:\n  local: library\noutput: generated.md\ndoc:\n  - heading: Agent\n    from: [local:identity]\n"
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"build", "-manifest", manifestPath}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	output, err := os.ReadFile(filepath.Join(temporary, "generated.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(output) != "# Agent\n\nHello.\n" {
		t.Fatalf("output = %q", output)
	}
}

func TestRunStatusReportsWorkspaceState(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(libraryPath, "core.md"), []byte("# Identity\nHello.\n\n# Instructions\nWork.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	manifest := "sources:\n  shared: library\noutput: AGENTS.md\ndoc:\n  - Identity: shared:identity\n"
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"status", "-manifest", manifestPath}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"Manifest: " + manifestPath,
		"Output: " + filepath.Join(temporary, "AGENTS.md"),
		"Output status: missing",
		"Sources: 1",
		"- shared: library (2 nodes)",
	} {
		if !strings.Contains(stdout.String(), expected) {
			t.Fatalf("status missing %q:\n%s", expected, stdout.String())
		}
	}
}
