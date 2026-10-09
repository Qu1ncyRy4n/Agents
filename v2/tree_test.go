package v2

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func treeFixture(t *testing.T) *Config {
	t.Helper()
	root := t.TempDir()
	writeLibraryFile(t, root, "library/agents/intro/role.md", "# Role\n\nBe helpful.\n")
	writeLibraryFile(t, root, "library/skills/alpha/SKILL.md", "---\nname: alpha\n---\n# Alpha\n")
	writeLibraryFile(t, root, "library/skills/beta/notes.txt", "beta notes\n")
	writeLibraryFile(t, root, "library/"+LibraryFile, `
library {
  format = 2
  id = "test"
  name = "Test"
}
content {
  markdown_root = "agents"
}
sections {
  section "intro" {
    title = "Intro"
    section "role" {
      title = "Role"
      source = "intro/role.md"
    }
  }
}
trees {
  tree "skills" {
    root = "skills"
    entry "alpha" {
      path = "alpha"
      tldr = "Alpha skill."
      tags = ["skill/alpha"]
    }
  }
}
`)
	writeLibraryFile(t, root, ConfigFile, `
mogent { format = 2 }
sources {
  source "shared" { local = "library" }
}
outputs {
  output "agents" {
    path = "AGENTS.md"
    kind = "markdown"
    source "shared" {
      from = "shared:agents"
      select = { all = true }
    }
  }
  output "skills" {
    path = ".agents/skills/"
    kind = "tree"
    source "shared" {
      from = "shared:skills"
      select = { all = true }
    }
  }
}
`)
	config, err := Load(filepath.Join(root, ConfigFile))
	if err != nil {
		t.Fatal(err)
	}
	return config
}

func TestTreeOutputPlansAppliesAndSyncs(t *testing.T) {
	config := treeFixture(t)
	root := filepath.Dir(config.Path)
	if config.Outputs[1].Paths[0] != ".agents/skills" {
		t.Fatalf("tree path = %q", config.Outputs[1].Paths[0])
	}
	result, err := PlanConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Trees) != 1 || result.Trees[0].Status != FileNew || strings.Join(result.Trees[0].Added, " ") != "alpha/SKILL.md beta/notes.txt" {
		t.Fatalf("trees = %#v", result.Trees)
	}
	if _, err := os.Stat(filepath.Join(root, ".agents")); !os.IsNotExist(err) {
		t.Fatalf("plan wrote tree: %v", err)
	}
	if _, err := Apply(config, false); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, ".agents", "skills")
	if content, err := os.ReadFile(filepath.Join(target, "beta", "notes.txt")); err != nil || string(content) != "beta notes\n" {
		t.Fatalf("copied file = %q, %v", content, err)
	}
	entries, err := os.ReadDir(filepath.Join(root, ".agents"))
	if err != nil || len(entries) != 1 {
		t.Fatalf("staging left behind: %v %v", entries, err)
	}
	result, err = PlanConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	if result.Trees[0].Status != FileUnchanged || len(diagnosticCodes(result.Plan)) != 0 {
		t.Fatalf("second plan = %#v %v", result.Trees[0], diagnosticCodes(result.Plan))
	}

	library := filepath.Join(root, "library", "skills")
	writeLibraryFile(t, library, "gamma/SKILL.md", "# Gamma\n")
	if err := os.RemoveAll(filepath.Join(library, "beta")); err != nil {
		t.Fatal(err)
	}
	writeLibraryFile(t, library, "alpha/SKILL.md", "# Alpha v2\n")
	result, err = PlanConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	tree := result.Trees[0]
	if tree.Status != FileChanged || strings.Join(tree.Added, " ") != "gamma/SKILL.md" || strings.Join(tree.Changed, " ") != "alpha/SKILL.md" || strings.Join(tree.Removed, " ") != "beta/notes.txt" {
		t.Fatalf("sync plan = %#v", tree)
	}
	if _, err := Apply(config, false); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(target, "beta")); !os.IsNotExist(err) {
		t.Fatalf("removed source file still present: %v", err)
	}
	if content, err := os.ReadFile(filepath.Join(target, "alpha", "SKILL.md")); err != nil || string(content) != "# Alpha v2\n" {
		t.Fatalf("updated file = %q, %v", content, err)
	}
}

func TestTreeOutputRefusesEditedTargetUntilForced(t *testing.T) {
	config := treeFixture(t)
	root := filepath.Dir(config.Path)
	if _, err := Apply(config, false); err != nil {
		t.Fatal(err)
	}
	edited := filepath.Join(root, ".agents", "skills", "alpha", "SKILL.md")
	if err := os.WriteFile(edited, []byte("hand edit\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	result, err := PlanConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	if result.Trees[0].Status != FileChanged || strings.Join(diagnosticCodes(result.Plan), " ") != "MOGENT208" {
		t.Fatalf("plan = %#v %v", result.Trees[0], diagnosticCodes(result.Plan))
	}
	if _, err := Apply(config, false); err == nil || !strings.Contains(err.Error(), "direct edits in tree output") {
		t.Fatalf("apply error = %v", err)
	}
	if _, err := Apply(config, true); err != nil {
		t.Fatal(err)
	}
	if content, err := os.ReadFile(edited); err != nil || strings.Contains(string(content), "hand edit") {
		t.Fatalf("forced apply left edit: %q %v", content, err)
	}
}

func TestTreeOutputRefusesUnmanagedTargetAndSymlinkSource(t *testing.T) {
	config := treeFixture(t)
	root := filepath.Dir(config.Path)
	writeLibraryFile(t, root, ".agents/skills/mine/SKILL.md", "# Mine\n")
	result, err := PlanConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(result.Trees[0].Removed, " ") != "mine/SKILL.md" || !strings.Contains(result.Plan.Diagnostics[0].Message, "not managed") {
		t.Fatalf("plan = %#v %#v", result.Trees[0], result.Plan.Diagnostics)
	}
	if _, err := Apply(config, false); err == nil || !strings.Contains(err.Error(), "untracked tree output") {
		t.Fatalf("apply error = %v", err)
	}
	if _, err := Apply(config, true); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(root, "library", "agents"), filepath.Join(root, "library", "skills", "link")); err != nil {
		t.Skip("symlinks unsupported:", err)
	}
	if _, err := PlanConfig(config); err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("symlink source error = %v", err)
	}
}

func TestTreeOutputExcludesPaths(t *testing.T) {
	config := treeFixture(t)
	config.Outputs[1].Sources[0].Exclude = []string{"beta"}
	root := filepath.Dir(config.Path)

	result, err := PlanConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(result.Trees[0].Added, " "); got != "alpha/SKILL.md" {
		t.Fatalf("added = %q", got)
	}
	if _, err := Apply(config, false); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, ".agents", "skills")
	if _, err := os.Stat(filepath.Join(target, "beta")); !os.IsNotExist(err) {
		t.Fatalf("excluded directory exists: %v", err)
	}

	config.Outputs[1].Sources[0].Exclude = nil
	result, err = PlanConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(result.Trees[0].Added, " "); got != "beta/notes.txt" {
		t.Fatalf("added after removing exclude = %q", got)
	}
}

func TestLoadLibraryValidatesTreeEntries(t *testing.T) {
	root := t.TempDir()
	writeLibraryFile(t, root, "agents/a.md", "# A\n")
	writeLibraryFile(t, root, "skills/x/SKILL.md", "# X\n")
	writeLibraryFile(t, root, LibraryFile, "library {\n  format = 2\n  id = \"t\"\n  name = \"T\"\n}\ncontent {\n  markdown_root = \"agents\"\n}\nsections {\n  section \"a\" {\n    title = \"A\"\n    source = \"a.md\"\n  }\n}\ntrees {\n  tree \"skills\" {\n    root = \"skills\"\n    entry \"x\" {\n      path = \"missing\"\n    }\n  }\n}\n")
	if _, err := LoadLibrary(root); err == nil || !strings.Contains(err.Error(), "entry \"x\" path \"missing\"") {
		t.Fatalf("LoadLibrary error = %v", err)
	}
}
