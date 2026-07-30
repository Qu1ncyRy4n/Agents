package navigator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestModelShowsTreeFinalAndSourceContext(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	source := "# Identity\nHello {{ .project_name }}.\n\n# Instructions\n\n## Workflow\nWork carefully.\n"
	if err := os.WriteFile(filepath.Join(libraryPath, "core.md"), []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	config := "sources:\n  shared: library\nvars:\n  project_name: Demo\ndoc:\n  - Identity: shared:identity\n  - heading: Instructions\n    children:\n      - Workflow: shared:instructions/workflow\n"
	if err := os.WriteFile(manifestPath, []byte(config), 0o644); err != nil {
		t.Fatal(err)
	}

	model, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	model.width = 150
	model.height = 30
	view := model.View()

	for _, expected := range []string{
		"Manifest tree",
		"Final AGENTS.md",
		"Selected-source context",
		"[x] Identity  <shared:identity>",
		"Hello Demo.",
		"Source: shared:identity",
	} {
		if !strings.Contains(view, expected) {
			t.Fatalf("view missing %q:\n%s", expected, view)
		}
	}
}

func TestModelNarrowViewSwitchesDetailPane(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(libraryPath, "core.md"), []byte("# Identity\nHello.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	config := "sources:\n  shared: library\ndoc:\n  - Identity: shared:identity\n"
	if err := os.WriteFile(manifestPath, []byte(config), 0o644); err != nil {
		t.Fatal(err)
	}

	model, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	model.width = 100
	model.height = 20
	if !strings.Contains(model.View(), "Detail: Final") {
		t.Fatalf("initial narrow view should show final detail:\n%s", model.View())
	}
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})
	switched := updated.(Model)
	if !strings.Contains(switched.View(), "Detail: Source") {
		t.Fatalf("narrow view should switch to source detail:\n%s", switched.View())
	}
}

func TestModelConfirmedSaveBuildWritesOutput(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(libraryPath, "core.md"), []byte("# Identity\nHello.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	config := "sources:\n  shared: library\noutput: AGENTS.md\ndoc:\n  - Identity: shared:identity\n"
	if err := os.WriteFile(manifestPath, []byte(config), 0o644); err != nil {
		t.Fatal(err)
	}

	model, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	pending, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	confirmed, _ := pending.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	saved := confirmed.(Model)
	if saved.err != "" {
		t.Fatalf("save failed: %s", saved.err)
	}
	output, err := os.ReadFile(filepath.Join(temporary, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(output) != "# Identity\n\nHello.\n" {
		t.Fatalf("output = %q", output)
	}
}

func TestModelDraftAddSaveWritesManifestAndOutput(t *testing.T) {
	temporary, manifestPath := writeNavigatorFixture(t, false)
	model, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	model.width = 120
	model.height = 24
	moved, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	added, _ := moved.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	draft := added.(Model)
	if !draft.dirty {
		t.Fatal("expected draft to be dirty after add")
	}
	if !strings.Contains(draft.View(), "state: ~ draft") {
		t.Fatalf("view missing dirty state:\n%s", draft.View())
	}
	if !strings.Contains(draft.output, "## Testing") {
		t.Fatalf("draft output missing Testing:\n%s", draft.output)
	}

	pending, _ := draft.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	confirmed, _ := pending.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	saved := confirmed.(Model)
	if saved.err != "" {
		t.Fatalf("save failed: %s", saved.err)
	}
	if saved.dirty {
		t.Fatal("expected clean state after save")
	}
	manifestContents, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(manifestContents), "Testing: shared:instructions/testing") {
		t.Fatalf("saved manifest missing Testing:\n%s", manifestContents)
	}
	output, err := os.ReadFile(filepath.Join(temporary, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(output), "## Testing\n\nTest deterministically.") {
		t.Fatalf("output missing Testing:\n%s", output)
	}
}

func TestModelCancelledDraftDoesNotWriteFiles(t *testing.T) {
	_, manifestPath := writeNavigatorFixture(t, false)
	original, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	model, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	moved, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	added, _ := moved.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	pending, _ := added.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	cancelled, _ := pending.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	draft := cancelled.(Model)
	if draft.dirty {
		t.Fatal("cancel should discard the unsaved draft")
	}
	if strings.Contains(draft.output, "## Testing") {
		t.Fatalf("cancelled draft output still contains Testing:\n%s", draft.output)
	}
	after, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(original) {
		t.Fatalf("manifest changed after cancel:\n%s", after)
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(manifestPath), "AGENTS.md")); !os.IsNotExist(err) {
		t.Fatalf("AGENTS.md should not exist after cancel, err=%v", err)
	}
}

func TestModelFailedSaveRestoresManifestAndOutput(t *testing.T) {
	temporary, manifestPath := writeNavigatorFixture(t, false)
	original, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(temporary, ".mogent"), []byte("not a directory\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	model, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	moved, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	added, _ := moved.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	pending, _ := added.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	confirmed, _ := pending.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	failed := confirmed.(Model)
	if failed.err == "" {
		t.Fatal("expected failed save to report an error")
	}
	after, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(original) {
		t.Fatalf("manifest changed after failed save:\n%s", after)
	}
	if _, err := os.Stat(filepath.Join(temporary, "AGENTS.md")); !os.IsNotExist(err) {
		t.Fatalf("AGENTS.md should not remain after failed save, err=%v", err)
	}
}

func TestModelDraftRemoveAndReorder(t *testing.T) {
	_, manifestPath := writeNavigatorFixture(t, true)
	model, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	moved, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	moved, _ = moved.(Model).Update(tea.KeyMsg{Type: tea.KeyDown})
	reordered, _ := moved.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'D'}})
	draft := reordered.(Model)
	parent := entryAt(draft.draft.Doc, []int{1})
	if got := parent.Children[1].Heading; got != "Workflow" {
		t.Fatalf("expected Workflow to move after Testing, got %q", got)
	}
	removed, _ := draft.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	draft = removed.(Model)
	parent = entryAt(draft.draft.Doc, []int{1})
	if len(parent.Children) != 1 || parent.Children[0].Heading != "Testing" {
		t.Fatalf("expected Workflow removed, children=%v", parent.Children)
	}
	if !draft.dirty {
		t.Fatal("expected dirty draft after remove")
	}
}

func writeNavigatorFixture(t *testing.T, includeTesting bool) (string, string) {
	t.Helper()
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	source := "# Identity\nHello.\n\n# Instructions\n\n## Workflow\nWork carefully.\n\n## Testing\nTest deterministically.\n"
	if err := os.WriteFile(filepath.Join(libraryPath, "core.md"), []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	testingEntry := ""
	if includeTesting {
		testingEntry = "      - Testing: shared:instructions/testing\n"
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	config := "sources:\n  shared: library\noutput: AGENTS.md\ndoc:\n  - Identity: shared:identity\n  - heading: Instructions\n    children:\n      - Workflow: shared:instructions/workflow\n" + testingEntry
	if err := os.WriteFile(manifestPath, []byte(config), 0o644); err != nil {
		t.Fatal(err)
	}
	return temporary, manifestPath
}
