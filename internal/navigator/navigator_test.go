package navigator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Qu1ncyRy4n/Agents/internal/state"
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

func TestModelLocalizeSelectedNodeSavesLocalOverride(t *testing.T) {
	temporary, manifestPath := writeNavigatorFixture(t, true)
	sharedPath := filepath.Join(temporary, "library", "core.md")
	originalShared, err := os.ReadFile(sharedPath)
	if err != nil {
		t.Fatal(err)
	}
	model, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	moved, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	moved, _ = moved.(Model).Update(tea.KeyMsg{Type: tea.KeyDown})
	localized, _ := moved.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	draft := localized.(Model)
	if !draft.dirty {
		t.Fatal("expected localization to create a dirty draft")
	}
	if got := entryAt(draft.draft.Doc, []int{1, 0}).From[0]; got != "local:instructions/workflow" {
		t.Fatalf("localized reference = %q", got)
	}
	if draft.draft.Sources["local"] != ".mogent/library" {
		t.Fatalf("local source = %q", draft.draft.Sources["local"])
	}
	if !strings.Contains(draft.View(), "[L] Workflow  <local:instructions/workflow>") {
		t.Fatalf("tree missing local marker:\n%s", draft.View())
	}
	if _, err := os.Stat(filepath.Join(temporary, ".mogent", "library")); !os.IsNotExist(err) {
		t.Fatalf("local library should not exist before save, err=%v", err)
	}

	pending, _ := draft.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	confirmed, _ := pending.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	saved := confirmed.(Model)
	if saved.err != "" {
		t.Fatalf("save failed: %s", saved.err)
	}
	localPath := filepath.Join(temporary, ".mogent", "library", "instructions", "workflow.md")
	localContent, err := os.ReadFile(localPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(localContent), "# Instructions\n\n## Workflow\n\nWork carefully.") {
		t.Fatalf("local override content:\n%s", localContent)
	}
	manifestContent, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(manifestContent), "local: .mogent/library") || !strings.Contains(string(manifestContent), "Workflow: local:instructions/workflow") {
		t.Fatalf("manifest missing local source/reference:\n%s", manifestContent)
	}
	afterShared, err := os.ReadFile(sharedPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(afterShared) != string(originalShared) {
		t.Fatalf("shared source changed:\n%s", afterShared)
	}
}

func TestModelCancelLocalizeDoesNotWriteLocalLibrary(t *testing.T) {
	temporary, manifestPath := writeNavigatorFixture(t, true)
	original, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	model, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	moved, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	moved, _ = moved.(Model).Update(tea.KeyMsg{Type: tea.KeyDown})
	localized, _ := moved.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	pending, _ := localized.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	cancelled, _ := pending.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	draft := cancelled.(Model)
	if draft.dirty {
		t.Fatal("cancel should discard localized draft")
	}
	after, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(original) {
		t.Fatalf("manifest changed after cancel:\n%s", after)
	}
	if _, err := os.Stat(filepath.Join(temporary, ".mogent", "library")); !os.IsNotExist(err) {
		t.Fatalf("local library should not exist after cancel, err=%v", err)
	}
}

func TestModelSaveConfirmationShowsReview(t *testing.T) {
	_, manifestPath := writeNavigatorFixture(t, true)
	model, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	moved, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	moved, _ = moved.(Model).Update(tea.KeyMsg{Type: tea.KeyDown})
	localized, _ := moved.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	pending, _ := localized.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	reviewing := pending.(Model)
	review := reviewing.sourceView(120, 40)
	for _, expected := range []string{
		"Save review",
		"agents.yaml:",
		"AGENTS.md:",
		"local override files:",
		".mogent/library/instructions/workflow.md",
	} {
		if !strings.Contains(review, expected) {
			t.Fatalf("review missing %q:\n%s", expected, review)
		}
	}
}

func TestModelDetectsAndRejectsDirectOutputDrift(t *testing.T) {
	temporary, manifestPath := writeNavigatorFixture(t, true)
	generated := generatedFixtureOutput()
	outputPath := filepath.Join(temporary, "AGENTS.md")
	if err := os.WriteFile(outputPath, []byte(generated), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := state.Write(filepath.Join(temporary, ".mogent", "state.json"), generated); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(outputPath, []byte(strings.Replace(generated, "Work carefully.", "Direct edit.", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	model, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if !model.drift {
		t.Fatal("expected drift to be detected")
	}
	pending, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	confirmed, _ := pending.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	rejected := confirmed.(Model)
	if rejected.err != "" {
		t.Fatalf("reject drift failed: %s", rejected.err)
	}
	if rejected.drift {
		t.Fatal("expected drift to clear after reject")
	}
	output, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(output) != generated {
		t.Fatalf("output was not rebuilt from manifest:\n%s", output)
	}
}

func TestModelImportsDirectOutputDriftIntoSelectedNode(t *testing.T) {
	temporary, manifestPath := writeNavigatorFixture(t, true)
	generated := generatedFixtureOutput()
	edited := strings.Replace(generated, "Work carefully.", "Edited locally.", 1)
	outputPath := filepath.Join(temporary, "AGENTS.md")
	if err := os.WriteFile(outputPath, []byte(generated), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := state.Write(filepath.Join(temporary, ".mogent", "state.json"), generated); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(outputPath, []byte(edited), 0o644); err != nil {
		t.Fatal(err)
	}
	model, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	moved, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	moved, _ = moved.(Model).Update(tea.KeyMsg{Type: tea.KeyDown})
	imported, _ := moved.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
	draft := imported.(Model)
	if !draft.dirty || !draft.importedDrift {
		t.Fatal("expected imported drift to create a dirty import draft")
	}
	if !strings.Contains(draft.output, "Edited locally.") {
		t.Fatalf("preview missing imported edit:\n%s", draft.output)
	}
	pending, _ := draft.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	confirmed, _ := pending.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	saved := confirmed.(Model)
	if saved.err != "" {
		t.Fatalf("save imported drift failed: %s", saved.err)
	}
	localContent, err := os.ReadFile(filepath.Join(temporary, ".mogent", "library", "instructions", "workflow.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(localContent), "Edited locally.") {
		t.Fatalf("local override missing imported edit:\n%s", localContent)
	}
	manifestContent, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(manifestContent), "Workflow: local:instructions/workflow") {
		t.Fatalf("manifest missing imported local reference:\n%s", manifestContent)
	}
	output, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(output), "Edited locally.") {
		t.Fatalf("output missing imported edit:\n%s", output)
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

func generatedFixtureOutput() string {
	return "# Identity\n\nHello.\n\n# Instructions\n\n## Workflow\n\nWork carefully.\n\n## Testing\n\nTest deterministically.\n"
}
