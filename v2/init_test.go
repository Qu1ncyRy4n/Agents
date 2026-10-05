package v2

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitWritesMinimalConfigAndRefusesOverwrite(t *testing.T) {
	path := filepath.Join(t.TempDir(), ConfigFile)
	content, err := Init(path, "", "shared", "../library", "agents", "AGENTS.md", false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(content, "source \"shared\"") || !strings.Contains(content, "from = \"shared:agents\"") || !strings.Contains(content, "all = true") {
		t.Fatalf("config = %q", content)
	}
	onDisk, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(onDisk) != content {
		t.Fatalf("on-disk config = %q, want %q", onDisk, content)
	}
	if _, err := Init(path, "", "shared", "../library", "agents", "AGENTS.md", false); err == nil || !strings.Contains(err.Error(), "refusing to overwrite") {
		t.Fatalf("second Init error = %v", err)
	}
}

func TestInitDryRunDoesNotWrite(t *testing.T) {
	path := filepath.Join(t.TempDir(), ConfigFile)
	if _, err := Init(path, "minimal", "shared", "../library", "agents", "AGENTS.md", true); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("config exists after dry run: %v", err)
	}
}

func TestInitTemplatesAndGitSources(t *testing.T) {
	path := filepath.Join(t.TempDir(), ConfigFile)
	content, err := Init(path, "qmr-core", "qmr", "https://example.test/qmr.git", "agents", "AGENTS.md", true)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"git = \"https://example.test/qmr.git\"", "ref = \"HEAD\"", "workflow    = { accept_defaults = true }", "from   = \"qmr:skills\"", "kind = \"tree\""} {
		if !strings.Contains(content, want) {
			t.Fatalf("qmr-core config missing %q:\n%s", want, content)
		}
	}
	config, err := Load(writeConfig(t, content))
	if err != nil {
		t.Fatalf("generated config does not load: %v", err)
	}
	if len(config.Outputs) != 2 || config.Outputs[1].Kind != "tree" {
		t.Fatalf("outputs = %#v", config.Outputs)
	}
	if _, err := Init(path, "ghost", "qmr", "../lib", "agents", "AGENTS.md", true); err == nil || !strings.Contains(err.Error(), "unknown v2 template") {
		t.Fatalf("unknown template error = %v", err)
	}
	if _, err := Init(path, "", "qmr", "https://user:pw@example.test/qmr.git", "agents", "AGENTS.md", true); err == nil || !strings.Contains(err.Error(), "credentials") {
		t.Fatalf("credential URL error = %v", err)
	}
}
