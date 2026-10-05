package state

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteDirectoryPreservesExistingFileRecord(t *testing.T) {
	root := t.TempDir()
	statePath := filepath.Join(root, ".mogent", "state.json")
	file := filepath.Join(root, "AGENTS.md")
	directory := filepath.Join(root, ".agents", "skills")
	if err := os.WriteFile(file, []byte("# Agents\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Write(statePath, file, "# Agents\n"); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, "SKILL.md"), []byte("# Skill\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := WriteDirectory(statePath, directory, map[string]string{"SKILL.md": Hash([]byte("# Skill\n"))}); err != nil {
		t.Fatal(err)
	}
	if status, err := Inspect(file, statePath); err != nil || status != OutputClean {
		t.Fatalf("file status = %q, %v", status, err)
	}
	if status, err := InspectDirectory(directory, statePath); err != nil || status != OutputClean {
		t.Fatalf("directory status = %q, %v", status, err)
	}
}
