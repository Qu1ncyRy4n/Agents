package content_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/content"
)

func inventory(t *testing.T, files map[string]string) *content.Inventory {
	t.Helper()
	root := t.TempDir()
	for name, body := range files {
		write(t, root, name, body)
	}
	i, err := content.Discover(root)
	if err != nil {
		t.Fatal(err)
	}
	return i
}

func TestComposeMixedPayloadsPreservesSkillAndOpaqueBytes(t *testing.T) {
	skill := "---\nname: review\ndescription: Review changes.\n---\n# Review\n## Procedure\nInspect.\n"
	i := inventory(t, map[string]string{"review/SKILL.md": skill, "review/check.sh": "#!/bin/sh\nexit 0\n", "settings.toml": "# keep comment\n[x]\ny = 2\n"})
	m, err := content.Compose([]content.Contribution{
		{Inventory: i, Node: "review", Operation: content.Copy, Into: ".agents/skills/review", Origin: "team:review"},
		{Inventory: i, Node: "review/SKILL.md", Heading: []string{"Review", "Procedure"}, Operation: content.RenderMarkdown, Into: "AGENTS.md", Origin: "team:procedure"},
		{Inventory: i, Node: "settings.toml", Operation: content.Copy, Into: "config/tool.toml", Origin: "team:config"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Files) != 4 || string(m.Files[".agents/skills/review/SKILL.md"].Bytes) != skill || string(m.Files["AGENTS.md"].Bytes) != "# Procedure\nInspect.\n" || m.Files[".agents/skills/review/check.sh"].Mode != 0755 {
		t.Fatalf("manifest=%#v", m)
	}
	if string(m.Files["config/tool.toml"].Bytes) != "# keep comment\n[x]\ny = 2\n" {
		t.Fatal("opaque configuration changed")
	}
	for _, event := range m.Events {
		if event.Path == "AGENTS.md" && event.Operation != content.RenderMarkdown {
			t.Fatalf("render provenance=%#v", event)
		}
		if event.Path == "config/tool.toml" && event.Operation != content.Copy {
			t.Fatalf("copy provenance=%#v", event)
		}
	}
}

func TestBundleReplacementRemovesStaleMembersAndReportsHistory(t *testing.T) {
	team := inventory(t, map[string]string{"review/SKILL.md": "# Review\n", "review/old.sh": "old"})
	repo := inventory(t, map[string]string{"review/SKILL.md": "# Repo review\n"})
	base := content.Contribution{Inventory: team, Node: "review", Operation: content.Copy, Into: "skills/review", Origin: "team"}
	override := content.Contribution{Inventory: repo, Node: "review", Operation: content.Copy, Into: "skills/review", Origin: "repo"}
	if _, err := content.Compose([]content.Contribution{base, override}); err == nil || !strings.Contains(err.Error(), "collision") {
		t.Fatalf("collision=%v", err)
	}
	override.Replace = true
	m, err := content.Compose([]content.Contribution{base, override})
	if err != nil {
		t.Fatal(err)
	}
	if _, found := m.Files["skills/review/old.sh"]; found {
		t.Fatal("obsolete member survived")
	}
	removals := 0
	for _, event := range m.Events {
		if event.Action == "replace-remove" && event.Previous == "team" {
			removals++
		}
	}
	if removals != 2 {
		t.Fatalf("events=%v", m.Events)
	}
}

func TestSourceExcludesDoNotRemoveEarlierContributions(t *testing.T) {
	i := inventory(t, map[string]string{"one.md": "one", "one-old.md": "old", "two.md": "two"})
	m, err := content.Compose([]content.Contribution{
		{Inventory: i, Node: "one.md", Operation: content.Copy, Into: "one.md", Origin: "earlier"},
		{Inventory: i, Node: ".", Operation: content.Copy, Into: ".", Origin: "later", Exclude: []string{"one.md"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Files) != 3 || m.Files["one.md"].Origin != "earlier" {
		t.Fatalf("files=%v", m.Files)
	}
}

func TestComposeRejectsPartialSkillsAndBrokenRenderedReferences(t *testing.T) {
	i := inventory(t, map[string]string{"review/SKILL.md": "# Review\n[Checklist](references/check.md)\n", "review/references/check.md": "check"})
	for _, c := range []content.Contribution{
		{Inventory: i, Node: "review/SKILL.md", Operation: content.Copy, Into: "SKILL.md"},
		{Inventory: i, Node: "review", Operation: content.Copy, Into: "review", Exclude: []string{"references"}},
		{Inventory: i, Node: "review/SKILL.md", Operation: content.RenderMarkdown, Into: "AGENTS.md"},
	} {
		if _, err := content.Compose([]content.Contribution{c}); err == nil {
			t.Fatalf("invalid contribution accepted: %#v", c)
		}
	}
	m, err := content.Compose([]content.Contribution{
		{Inventory: i, Node: "review/SKILL.md", Operation: content.RenderMarkdown, Into: "AGENTS.md"},
		{Inventory: inventory(t, map[string]string{"references/check.md": "check"}), Node: ".", Operation: content.Copy, Into: "."},
	})
	if err != nil || len(m.Files) != 2 {
		t.Fatalf("dependency plan=%v %v", m, err)
	}
}

func TestComposeRequiresExplicitAppendAndDiagnosesTypeConflicts(t *testing.T) {
	i := inventory(t, map[string]string{"a.md": "# A\nA.\n", "b.md": "# B\nB.\n"})
	a := content.Contribution{Inventory: i, Node: "a.md", Operation: content.RenderMarkdown, Into: "AGENTS.md", Origin: "a"}
	b := content.Contribution{Inventory: i, Node: "b.md", Operation: content.RenderMarkdown, Into: "AGENTS.md", Origin: "b"}
	if _, err := content.Compose([]content.Contribution{a, b}); err == nil {
		t.Fatal("implicit append accepted")
	}
	b.Append = true
	m, err := content.Compose([]content.Contribution{a, b})
	if err != nil {
		t.Fatal(err)
	}
	if string(m.Files["AGENTS.md"].Bytes) != "# A\nA.\n\n# B\nB.\n" {
		t.Fatalf("rendered=%q", m.Files["AGENTS.md"].Bytes)
	}
	if _, err := content.Compose([]content.Contribution{a, {Inventory: i, Node: "b.md", Operation: content.Copy, Into: "AGENTS.md/b.md"}}); err == nil {
		t.Fatal("file/subtree conflict accepted")
	}
}

func TestMaterializeUsesCapturedBytesAndRefusesOccupiedStaging(t *testing.T) {
	i := inventory(t, map[string]string{"check.sh": "#!/bin/sh\nexit 0\n"})
	m, err := content.Compose([]content.Contribution{{Inventory: i, Node: "check.sh", Operation: content.Copy, Into: "bin/check.sh"}})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(i.Nodes["check.sh"].Source, []byte("changed after plan"), 0600); err != nil {
		t.Fatal(err)
	}
	staging := t.TempDir()
	if err := m.Materialize(staging); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(staging, "bin/check.sh"))
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(filepath.Join(staging, "bin/check.sh"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "#!/bin/sh\nexit 0\n" || info.Mode().Perm() != 0755 {
		t.Fatalf("materialized=%q mode=%v", data, info.Mode())
	}
	if err := m.Materialize(staging); err == nil {
		t.Fatal("occupied staging overwritten")
	}
}
