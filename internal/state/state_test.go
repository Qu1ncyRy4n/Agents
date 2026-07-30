package state_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/internal/state"
)

func TestCheckOverwriteProtectsUntrackedAndEditedOutput(t *testing.T) {
	temporary := t.TempDir()
	outputPath := filepath.Join(temporary, "AGENTS.md")
	statePath := filepath.Join(temporary, ".mogent", "state.json")
	if err := os.WriteFile(outputPath, []byte("handwritten\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := state.CheckOverwrite(outputPath, statePath, false); err == nil || !strings.Contains(err.Error(), "untracked") {
		t.Fatalf("untracked error = %v", err)
	}
	if err := state.Write(statePath, "generated\n"); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(outputPath, []byte("direct edit\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := state.CheckOverwrite(outputPath, statePath, false); err == nil || !strings.Contains(err.Error(), "direct edits") {
		t.Fatalf("direct edit error = %v", err)
	}
	if err := state.CheckOverwrite(outputPath, statePath, true); err != nil {
		t.Fatalf("force error = %v", err)
	}
}
