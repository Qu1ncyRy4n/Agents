package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/manifest"
)

func TestBuildOutputsCopiesQMRStyleSkillsAndTracksDrift(t *testing.T) {
	root := t.TempDir()
	skills := filepath.Join(root, "library", "skills", "review")
	if err := os.MkdirAll(filepath.Join(skills, "assets"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skills, "SKILL.md"), []byte("# Review\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skills, "assets", "checklist.txt"), []byte("raw asset\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "library", "rules.md"), []byte("# Rules\nRules.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(root, "agents.yaml")
	value := &manifest.Manifest{Sources: map[string]manifest.Source{"qmr": {Location: "library"}}, Output: "AGENTS.md", Outputs: []manifest.Output{{Path: ".agents/skills/", From: "qmr:skills"}}, Doc: []manifest.Entry{{Heading: "Rules", From: []string{"qmr:rules"}}}}
	if err := manifest.WriteAtomically(manifestPath, value); err != nil {
		t.Fatal(err)
	}
	if err := BuildOutputs(value, manifestPath, "# Rules\n\nRules.\n", false); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, ".agents", "skills", "review", "SKILL.md")); err != nil {
		t.Fatal(err)
	}
	asset, err := os.ReadFile(filepath.Join(root, ".agents", "skills", "review", "assets", "checklist.txt"))
	if err != nil || string(asset) != "raw asset\n" {
		t.Fatalf("asset = %q, %v", asset, err)
	}
	if _, err := os.Stat(filepath.Join(root, ".agents", "skills", "skills")); !os.IsNotExist(err) {
		t.Fatalf("copied selected root itself: %v", err)
	}
	session, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	status, err := session.Status()
	if err != nil {
		t.Fatal(err)
	}
	if len(status.Outputs) != 1 || status.Outputs[0].Status != StatusUpToDate {
		t.Fatalf("status = %#v", status.Outputs)
	}
	if err := os.WriteFile(filepath.Join(root, ".agents", "skills", "review", "SKILL.md"), []byte("edited\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := BuildOutputs(value, manifestPath, "# Rules\n\nRules.\n", false); err == nil || !strings.Contains(err.Error(), "modified") {
		t.Fatalf("overwrite error = %v", err)
	}
	if err := BuildOutputs(value, manifestPath, "# Rules\n\nRules.\n", true); err != nil {
		t.Fatal(err)
	}
}

func TestBuildOutputsRejectsSourceTargetsAndSymlinks(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "library"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "library", "rules.md"), []byte("# Rules\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(root, "agents.yaml")
	value := &manifest.Manifest{Sources: map[string]manifest.Source{"qmr": {Location: "library"}}, Output: "library/generated.md", Doc: []manifest.Entry{{Heading: "Rules", From: []string{"qmr:rules"}}}}
	if err := BuildOutputs(value, manifestPath, "# Rules\n", false); err == nil || !strings.Contains(err.Error(), "source root") {
		t.Fatalf("source target error = %v", err)
	}
	if err := os.Symlink(filepath.Join(root, "outside"), filepath.Join(root, "linked")); err != nil {
		t.Fatal(err)
	}
	value.Output = "linked/AGENTS.md"
	if err := BuildOutputs(value, manifestPath, "# Rules\n", false); err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("symlink error = %v", err)
	}
}
