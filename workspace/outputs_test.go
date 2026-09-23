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

func TestSessionStatusDistinguishesRawDirectorySources(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "agents"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "skills", "review"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "agents", "core.md"), []byte("# Core\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "skills", "review", "SKILL.md"), []byte("# Review\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "skills", "review", "checklist.txt"), []byte("check\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(root, "agents.yaml")
	value := &manifest.Manifest{Sources: map[string]manifest.Source{"agents": {Location: "agents"}, "skills": {Location: "skills"}}, Outputs: []manifest.Output{{Path: "AGENTS.md", Include: []manifest.Selector{{All: "agents"}}}, {Path: ".agents/skills/", Include: []manifest.Selector{{All: "skills"}}}}}
	if err := manifest.WriteAtomically(manifestPath, value); err != nil {
		t.Fatal(err)
	}
	session, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	status, err := session.Status()
	if err != nil {
		t.Fatal(err)
	}
	if status.Sources[0].Alias != "agents" || status.Sources[0].DirectoryTree || status.Sources[0].Nodes != 1 {
		t.Fatalf("indexed source = %#v", status.Sources[0])
	}
	if status.Sources[1].Alias != "skills" || !status.Sources[1].DirectoryTree || status.Sources[1].Files != 2 {
		t.Fatalf("raw source = %#v", status.Sources[1])
	}
}

func TestSessionStatusReportsRootSidecarDirectoryTrees(t *testing.T) {
	root := t.TempDir()
	libraryRoot := filepath.Join(root, "library")
	if err := os.MkdirAll(filepath.Join(libraryRoot, "agents"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(libraryRoot, "skills", "review"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(libraryRoot, "agents", "core.md"), []byte("# Core\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(libraryRoot, "skills", "review", "SKILL.md"), []byte("# Review\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	sidecar := "schema: {id: mogent/1}\nlibrary: {name: Test}\ncontent:\n  markdown_roots: [agents]\n  directory_roots: [skills]\ntree:\n  - source: agents/core\n"
	if err := os.WriteFile(filepath.Join(libraryRoot, "library.mogent.yaml"), []byte(sidecar), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(root, "agents.yaml")
	value := &manifest.Manifest{Sources: map[string]manifest.Source{"qmr": {Location: "library"}}, Outputs: []manifest.Output{{Path: "AGENTS.md", Include: []manifest.Selector{{All: "qmr"}}}, {Path: ".agents/skills/", From: "qmr:skills"}}}
	if err := manifest.WriteAtomically(manifestPath, value); err != nil {
		t.Fatal(err)
	}
	session, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	status, err := session.Status()
	if err != nil {
		t.Fatal(err)
	}
	source := status.Sources[0]
	if source.Alias != "qmr" || source.Nodes != 1 || len(source.DirectoryTrees) != 1 || source.DirectoryTrees[0] != (DirectoryTreeStatus{Path: "skills", Files: 1}) {
		t.Fatalf("root sidecar source = %#v", source)
	}
}
