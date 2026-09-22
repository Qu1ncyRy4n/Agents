package state_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/state"
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
	if err := state.Write(statePath, outputPath, "generated\n"); err != nil {
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

func TestInspectReportsOutputState(t *testing.T) {
	temporary := t.TempDir()
	outputPath := filepath.Join(temporary, "AGENTS.md")
	statePath := filepath.Join(temporary, ".mogent", "state.json")

	got, err := state.Inspect(outputPath, statePath)
	if err != nil {
		t.Fatal(err)
	}
	if got != state.OutputMissing {
		t.Fatalf("missing state = %s", got)
	}

	if err := os.WriteFile(outputPath, []byte("handwritten\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err = state.Inspect(outputPath, statePath)
	if err != nil {
		t.Fatal(err)
	}
	if got != state.OutputUntracked {
		t.Fatalf("untracked state = %s", got)
	}

	if err := state.Write(statePath, outputPath, "handwritten\n"); err != nil {
		t.Fatal(err)
	}
	got, err = state.Inspect(outputPath, statePath)
	if err != nil {
		t.Fatal(err)
	}
	if got != state.OutputClean {
		t.Fatalf("clean state = %s", got)
	}

	if err := os.WriteFile(outputPath, []byte("edited\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err = state.Inspect(outputPath, statePath)
	if err != nil {
		t.Fatal(err)
	}
	if got != state.OutputModified {
		t.Fatalf("modified state = %s", got)
	}
}

func TestCheckOverwriteRejectsAChangedOutputPath(t *testing.T) {
	temporary := t.TempDir()
	firstOutput := filepath.Join(temporary, "AGENTS.md")
	secondOutput := filepath.Join(temporary, "CLAUDE.md")
	statePath := filepath.Join(temporary, ".mogent", "state.json")
	if err := os.WriteFile(firstOutput, []byte("generated\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := state.Write(statePath, firstOutput, "generated\n"); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(secondOutput, []byte("generated\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := state.CheckOverwrite(secondOutput, statePath, false); err == nil || !strings.Contains(err.Error(), "output path does not match") {
		t.Fatalf("changed output path error = %v", err)
	}
	got, err := state.Inspect(secondOutput, statePath)
	if err != nil {
		t.Fatal(err)
	}
	if got != state.OutputUntracked {
		t.Fatalf("changed output path state = %s", got)
	}
}

func TestWriteCreatesPrivateStateFile(t *testing.T) {
	temporary := t.TempDir()
	statePath := filepath.Join(temporary, ".mogent", "state.json")
	if err := state.Write(statePath, filepath.Join(temporary, "AGENTS.md"), "generated\n"); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(statePath)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("mode = %o, want 600", got)
	}
}

func TestInspectReadsLegacySingleOutputState(t *testing.T) {
	temporary := t.TempDir()
	output := filepath.Join(temporary, "AGENTS.md")
	content := "generated\n"
	if err := os.WriteFile(output, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	statePath := filepath.Join(temporary, ".mogent", "state.json")
	if err := os.MkdirAll(filepath.Dir(statePath), 0o755); err != nil {
		t.Fatal(err)
	}
	legacy := "{\"output_path\":\"" + output + "\",\"sha256\":\"" + state.Hash([]byte(content)) + "\"}\n"
	if err := os.WriteFile(statePath, []byte(legacy), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := state.Inspect(output, statePath)
	if err != nil || got != state.OutputClean {
		t.Fatalf("legacy inspect = %s, %v", got, err)
	}
}
