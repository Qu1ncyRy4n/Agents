package v2

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func artifactFixture(t *testing.T) (*Config, string) {
	t.Helper()
	root := t.TempDir()
	writeLibraryFile(t, root, "library/payloads/review/SKILL.md", "---\nname: review\ndescription: Review a change.\n---\n# Review\n\n## Procedure\nInspect changes.\n\n### Verify\nRun checks.\n")
	writeLibraryFile(t, root, "library/payloads/review/scripts/check.sh", "#!/bin/sh\nexit 0\n")
	if err := os.Chmod(filepath.Join(root, "library/payloads/review/scripts/check.sh"), 0755); err != nil {
		t.Fatal(err)
	}
	writeLibraryFile(t, root, "library/payloads/settings.toml", "# preserve\n[tool]\nx = 1\n")
	writeLibraryFile(t, root, "library/"+LibraryFile, `library {
  format = 2
  id = "typed"
  name = "Typed library"
}
content "physical" {
  directory = "payloads"
}
`)
	writeLibraryFile(t, root, ConfigFile, `mogent { format = 2 }
sources {
  source "shared" { local = "library" }
}
outputs {
  output "package" {
    path = "agent-export"
    kind = "dir-tree"
    source "shared" {
      from = "shared:physical"
      node = "review"
      operation = "copy"
      into = ".agents/skills/review"
    }
    source "shared" {
      from = "shared:physical"
      node = "review/SKILL.md"
      heading = ["Review", "Procedure"]
      operation = "render-markdown"
      into = "AGENTS.md"
    }
    source "shared" {
      from = "shared:physical"
      node = "settings.toml"
      operation = "copy"
      into = "config/tool.toml"
    }
  }
}
`)
	config, err := Load(filepath.Join(root, ConfigFile))
	if err != nil {
		t.Fatal(err)
	}
	return config, root
}

func TestArtifactAPIPlansAndAppliesMixedOutput(t *testing.T) {
	config, root := artifactFixture(t)
	result, err := PlanConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	if result.HasErrors() || len(result.Trees) != 1 || result.Trees[0].FileCount() != 4 {
		t.Fatalf("plan=%#v", result)
	}
	if _, err := os.Stat(filepath.Join(root, "agent-export")); !os.IsNotExist(err) {
		t.Fatalf("plan wrote output: %v", err)
	}
	if _, err := Apply(config, false); err != nil {
		t.Fatal(err)
	}
	guide, err := os.ReadFile(filepath.Join(root, "agent-export/AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(guide) != "# Procedure\nInspect changes.\n\n## Verify\nRun checks.\n" {
		t.Fatalf("guide=%q", guide)
	}
	skill, err := os.ReadFile(filepath.Join(root, "agent-export/.agents/skills/review/SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(skill), "---\nname: review\n") || !strings.Contains(string(skill), "## Procedure") {
		t.Fatalf("skill changed: %q", skill)
	}
	info, err := os.Stat(filepath.Join(root, "agent-export/.agents/skills/review/scripts/check.sh"))
	if err != nil || info.Mode().Perm() != 0755 {
		t.Fatalf("script mode: %v %v", info, err)
	}
	result, err = PlanConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	if result.Trees[0].Status != FileUnchanged {
		t.Fatalf("status=%v", result.Trees[0].Status)
	}
	if err := os.WriteFile(filepath.Join(root, "agent-export/AGENTS.md"), []byte("edited"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := Apply(config, false); err == nil || !strings.Contains(err.Error(), "direct edits") {
		t.Fatalf("drift error=%v", err)
	}
	if _, err := Apply(config, true); err != nil {
		t.Fatal(err)
	}
}

func TestArtifactAPIRejectsImplicitCollisionAndReplacesWholeBundle(t *testing.T) {
	config, root := artifactFixture(t)
	writeLibraryFile(t, root, "library/payloads/overrides/review/SKILL.md", "---\nname: review\ndescription: Repo review.\n---\n# Repo review\n")
	replacement := config.Outputs[0].Sources[0]
	replacement.Node = "overrides/review"
	config.Outputs[0].Sources = append(config.Outputs[0].Sources, replacement)
	if _, err := PlanConfig(config); err == nil || !strings.Contains(err.Error(), "collision") {
		t.Fatalf("collision=%v", err)
	}
	config.Outputs[0].Sources[3].Replace = true
	result, err := Apply(config, false)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "agent-export/.agents/skills/review/scripts/check.sh")); !os.IsNotExist(err) {
		t.Fatalf("stale member=%v", err)
	}
	removals := 0
	for _, event := range result.Trees[0].Manifest.Events {
		if event.Action == "replace-remove" {
			removals++
		}
	}
	if removals != 2 {
		t.Fatalf("events=%v", result.Trees[0].Manifest.Events)
	}
}

func TestArtifactTransactionRestoresOutputAndStateOnLaterFailure(t *testing.T) {
	config, root := artifactFixture(t)
	if _, err := Apply(config, false); err != nil {
		t.Fatal(err)
	}
	statePath := filepath.Join(root, ".mogent/state.json")
	oldState, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatal(err)
	}
	writeLibraryFile(t, root, "library/payloads/settings.toml", "replacement\n")
	result, err := PlanConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := snapshotFiles(map[string]string{}, statePath)
	if err != nil {
		t.Fatal(err)
	}
	backup, err := stageOutputTree(result.Trees[0])
	if err != nil {
		t.Fatal(err)
	}
	snapshot.trees = append(snapshot.trees, backup)
	failure := errors.New("injected failure after staging")
	if err := snapshot.restore(failure); !errors.Is(err, failure) {
		t.Fatalf("restore=%v", err)
	}
	data, err := os.ReadFile(filepath.Join(root, "agent-export/config/tool.toml"))
	if err != nil || string(data) != "# preserve\n[tool]\nx = 1\n" {
		t.Fatalf("restored payload=%q %v", data, err)
	}
	data, err = os.ReadFile(statePath)
	if err != nil || string(data) != string(oldState) {
		t.Fatalf("restored state=%q %v", data, err)
	}
}

func TestLoadLibraryRejectsNamedContentSymlinkEscape(t *testing.T) {
	_, root := artifactFixture(t)
	outside := t.TempDir()
	writeLibraryFile(t, outside, "nested/file.md", "# Outside\n")
	if err := os.Symlink(outside, filepath.Join(root, "library/link")); err != nil {
		t.Fatal(err)
	}
	writeLibraryFile(t, root, "library/"+LibraryFile, "library {\nformat = 2\nid = \"x\"\nname = \"X\"\n}\ncontent \"physical\" {\ndirectory = \"link/nested\"\n}\n")
	if _, err := LoadLibrary(filepath.Join(root, "library")); err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("escape=%v", err)
	}
}

func TestArtifactCanRenderExistingAuthoredSelections(t *testing.T) {
	config, root := artifactFixture(t)
	sidecar := filepath.Join(root, "library", LibraryFile)
	data, err := os.ReadFile(sidecar)
	if err != nil {
		t.Fatal(err)
	}
	writeLibraryFile(t, root, "library/"+LibraryFile, string(data)+`content {
  markdown_root = "agents"
}
sections {
  section "intro" {
    title = "Intro"
    section "role" {
      title = "Role"
      source = "role.md"
    }
  }
}
`)
	writeLibraryFile(t, root, "library/agents/role.md", "# Role\nBe helpful.\n")
	data, err = os.ReadFile(config.Path)
	if err != nil {
		t.Fatal(err)
	}
	updated := strings.Replace(string(data), "    source \"shared\" {", `    source "shared" {
      from = "shared:agents"
      select = { intro = true }
      operation = "render-markdown"
      into = "GUIDE.md"
    }
    source "shared" {`, 1)
	writeLibraryFile(t, root, ConfigFile, updated)
	config, err = Load(config.Path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Apply(config, false); err != nil {
		t.Fatal(err)
	}
	guide, err := os.ReadFile(filepath.Join(root, "agent-export/GUIDE.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(guide) != "# Intro\n\n## Role\n\nBe helpful.\n" {
		t.Fatalf("authored guide=%q", guide)
	}
}

func TestLoadRejectsOverlappingArtifactOutputs(t *testing.T) {
	config, root := artifactFixture(t)
	data, err := os.ReadFile(config.Path)
	if err != nil {
		t.Fatal(err)
	}
	withSecond := strings.Replace(string(data), "outputs {", "outputs {\n output \"nested\" {\n path = \"agent-export/AGENTS.md\"\n kind = \"markdown\"\n source \"shared\" {\n from = \"shared:agents\"\n select = { all = true }\n }\n }", 1)
	writeLibraryFile(t, root, ConfigFile, withSecond)
	if _, err := Load(config.Path); err == nil || !strings.Contains(err.Error(), "overlap") {
		t.Fatalf("overlap=%v", err)
	}
}
