package v2

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitWritesMinimalConfigAndRefusesOverwrite(t *testing.T) {
	path := filepath.Join(t.TempDir(), ConfigFile)
	content, err := Init(path, "shared", "../library", "agents", "AGENTS.md", false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(content, "source \"shared\"") || !strings.Contains(content, "from = \"shared:agents\"") {
		t.Fatalf("config = %q", content)
	}
	onDisk, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(onDisk) != content {
		t.Fatalf("on-disk config = %q, want %q", onDisk, content)
	}
	if _, err := Init(path, "shared", "../library", "agents", "AGENTS.md", false); err == nil || !strings.Contains(err.Error(), "refusing to overwrite") {
		t.Fatalf("second Init error = %v", err)
	}
}

func TestInitDryRunDoesNotWrite(t *testing.T) {
	path := filepath.Join(t.TempDir(), ConfigFile)
	if _, err := Init(path, "shared", "../library", "agents", "AGENTS.md", true); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("config exists after dry run: %v", err)
	}
}
