package assembly_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/assembly"
)

func TestSelectValidatesAndAppendsWithoutBuilding(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "mogent.lib"), []byte("library {format=1;name=\"Team\"}\nmd-section \"intro\" {source=\"intro.md\"}\npreset \"starter\" {include \"self:intro\"}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "intro.md"), []byte("# Intro\n"), 0644); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "mogent.build")
	original := "// keep this\nbuild {format=1;origin=\"authored\"}\nsource \"team" + "\" {local=\".\"}\n"
	if err := os.WriteFile(path, []byte(original), 0600); err != nil {
		t.Fatal(err)
	}
	preview, err := assembly.Select(path, "team:intro", true)
	if err != nil || !strings.HasPrefix(string(preview), original) || !strings.Contains(string(preview), `include "team:intro"`) {
		t.Fatalf("preview=%s err=%v", preview, err)
	}
	stored, err := os.ReadFile(path)
	if err != nil || string(stored) != original {
		t.Fatalf("dry run changed file: %s %v", stored, err)
	}
	if _, err := assembly.Select(path, "team:intro", false); err != nil {
		t.Fatal(err)
	}
	if _, err := assembly.Select(path, "team:intro", false); err == nil {
		t.Fatal("duplicate selection accepted")
	}
	for _, reference := range []string{"team", ":intro", "team:missing", "missing:intro"} {
		if _, err := assembly.Select(path, reference, true); err == nil {
			t.Errorf("invalid reference accepted: %s", reference)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "AGENTS.md")); !os.IsNotExist(err) {
		t.Fatalf("selection installed output: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm() != 0600 {
		t.Fatalf("mode=%v err=%v", info, err)
	}
}

func TestSelectAllowsPlanAndRejectsGeneratedBuild(t *testing.T) {
	root := t.TempDir()
	for name, body := range map[string]string{
		"mogent.lib":   "library {format=1;name=\"Team\"}\nmd-section \"intro\" {source=\"intro.md\"}\n",
		"intro.md":     "# Intro\n",
		"mogent.plan":  "plan {format=1}\nsource \"team\" {local=\".\"}\n",
		"mogent.build": "build {format=1;origin=\"generated\"}\nsource \"team\" {local=\".\"}\n",
	} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := assembly.Select(filepath.Join(root, "mogent.plan"), "team:intro", false); err != nil {
		t.Fatal(err)
	}
	if _, err := assembly.Select(filepath.Join(root, "mogent.build"), "team:intro", false); err == nil || !strings.Contains(err.Error(), "generated") {
		t.Fatalf("generated build select=%v", err)
	}
}

func TestSelectAllAddsOrUpgradesExplicitChildren(t *testing.T) {
	root := t.TempDir()
	for name, body := range map[string]string{
		"mogent.lib":          "library {format=1;name=\"Team\"}\ndir \"skills\" {source=\"skills\";node \"one\" {};node \"two\" {}}\n",
		"skills/one/SKILL.md": "# One\n",
		"skills/two/SKILL.md": "# Two\n",
		"mogent.build":        "build {format=1;origin=\"authored\"}\nsource \"team\" {local=\".\"}\ninclude \"team:skills\"\n",
	} {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
	}
	path := filepath.Join(root, "mogent.build")
	contents, err := assembly.SelectWithOptions(path, "team:skills", assembly.SelectOptions{All: true})
	if err != nil || !strings.Contains(string(contents), `["one", "two"]`) {
		t.Fatalf("all=%s err=%v", contents, err)
	}
	if _, err := assembly.SelectWithOptions(path, "team:skills", assembly.SelectOptions{All: true}); err == nil {
		t.Fatal("all accepted existing explicit selection")
	}
}
