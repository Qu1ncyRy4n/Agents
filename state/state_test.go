package state_test

import (
	"encoding/json"
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

func TestPortableStateSurvivesWorkspaceMove(t *testing.T) {
	first := t.TempDir()
	second := t.TempDir()
	firstState := filepath.Join(first, ".mogent", "state.json")
	if err := os.MkdirAll(filepath.Join(first, ".agents", "skills"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(first, "AGENTS.md"), []byte("generated\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(first, ".agents", "skills", "SKILL.md"), []byte("skill\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	hashes, err := state.DirectoryHashes(filepath.Join(first, ".agents", "skills"))
	if err != nil {
		t.Fatal(err)
	}
	if err := state.WriteAll(firstState, map[string]string{filepath.Join(first, "AGENTS.md"): "generated\n"}, map[string]map[string]string{filepath.Join(first, ".agents", "skills"): hashes}); err != nil {
		t.Fatal(err)
	}
	contents, err := os.ReadFile(firstState)
	if err != nil {
		t.Fatal(err)
	}
	var written map[string]any
	if err := json.Unmarshal(contents, &written); err != nil {
		t.Fatal(err)
	}
	if written["version"] != float64(3) || written["output_path"] != nil || written["sha256"] != nil || strings.Contains(string(contents), first) {
		t.Fatalf("state is not clean and portable: %s", contents)
	}
	if err := os.MkdirAll(filepath.Join(second, ".mogent"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(second, ".agents", "skills"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(second, "AGENTS.md"), []byte("generated\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(second, ".agents", "skills", "SKILL.md"), []byte("skill\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	secondState := filepath.Join(second, ".mogent", "state.json")
	if err := os.WriteFile(secondState, contents, 0o600); err != nil {
		t.Fatal(err)
	}
	if got, err := state.Inspect(filepath.Join(second, "AGENTS.md"), secondState); err != nil || got != state.OutputClean {
		t.Fatalf("moved file inspect = %s, %v", got, err)
	}
	if got, err := state.InspectDirectory(filepath.Join(second, ".agents", "skills"), secondState); err != nil || got != state.OutputClean {
		t.Fatalf("moved directory inspect = %s, %v", got, err)
	}
}

func TestPortableStateAcceptsSymlinkWorkspaceInvocation(t *testing.T) {
	workspace := t.TempDir()
	invocation := filepath.Join(t.TempDir(), "workspace")
	if err := os.Symlink(workspace, invocation); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(workspace, "AGENTS.md")
	if err := os.WriteFile(output, []byte("generated\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	invokedState := filepath.Join(invocation, ".mogent", "state.json")
	if err := state.Write(invokedState, filepath.Join(invocation, "AGENTS.md"), "generated\n"); err != nil {
		t.Fatal(err)
	}
	physicalState := filepath.Join(workspace, ".mogent", "state.json")
	if got, err := state.Inspect(output, physicalState); err != nil || got != state.OutputClean {
		t.Fatalf("physical invocation inspect = %s, %v", got, err)
	}
}

func TestLegacyAbsoluteV2StateMigratesOnWrite(t *testing.T) {
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
	legacy := "{\"outputs\":{\"" + output + "\":{\"kind\":\"file\",\"sha256\":\"" + state.Hash([]byte(content)) + "\"}}}\n"
	if err := os.WriteFile(statePath, []byte(legacy), 0o600); err != nil {
		t.Fatal(err)
	}
	if got, err := state.Inspect(output, statePath); err != nil || got != state.OutputClean {
		t.Fatalf("legacy v2 inspect = %s, %v", got, err)
	}
	if err := state.Write(statePath, output, content); err != nil {
		t.Fatal(err)
	}
	migrated, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(migrated), output) || !strings.Contains(string(migrated), "\"AGENTS.md\"") || !strings.Contains(string(migrated), "\"version\": 3") {
		t.Fatalf("migrated state = %s", migrated)
	}
}

func TestWritePreservesOtherOutputRecords(t *testing.T) {
	workspace := t.TempDir()
	statePath := filepath.Join(workspace, ".mogent", "state.json")
	markdown := filepath.Join(workspace, "AGENTS.md")
	skills := filepath.Join(workspace, ".agents", "skills")
	if err := os.MkdirAll(skills, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(markdown, []byte("before\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skills, "SKILL.md"), []byte("skill\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	hashes, err := state.DirectoryHashes(skills)
	if err != nil {
		t.Fatal(err)
	}
	if err := state.WriteAll(statePath, map[string]string{markdown: "before\n"}, map[string]map[string]string{skills: hashes}); err != nil {
		t.Fatal(err)
	}
	if err := state.Write(statePath, markdown, "after\n"); err != nil {
		t.Fatal(err)
	}
	if got, err := state.InspectDirectory(skills, statePath); err != nil || got != state.OutputClean {
		t.Fatalf("skills state after Markdown write = %s, %v", got, err)
	}
}
