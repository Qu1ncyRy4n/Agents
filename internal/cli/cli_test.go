package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
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
