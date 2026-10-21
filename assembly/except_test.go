package assembly_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/assembly"
)

func TestExceptAddsDescendantToNarrowestSelectedAncestor(t *testing.T) {
	root := t.TempDir()
	for name, body := range map[string]string{
		"mogent.lib":   "library {format=1;name=\"Team\"}\nmd-section \"workflow\" {md-section \"first\" {source=\"first.md" + "\"};md-section \"second\" {source=\"second.md\"}}\n",
		"first.md":     "# First\n",
		"second.md":    "# Second\n",
		"mogent.build": "// keep comment\nbuild {format=1;origin=\"authored\"}\nsource \"team\" {local=\".\"}\ninclude \"team:workflow\"\n",
	} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
	}
	path := filepath.Join(root, "mogent.build")
	preview, err := assembly.Except(path, "team:workflow/first", true)
	if err != nil || preview.Relative != "first" || !strings.Contains(string(preview.Content), `exclude = ["first"]`) {
		t.Fatalf("preview=%+v err=%v", preview, err)
	}
	stored, err := os.ReadFile(path)
	if err != nil || strings.Contains(string(stored), "exclude") {
		t.Fatalf("dry run modified file: %s %v", stored, err)
	}
	if _, err := assembly.Except(path, "team:workflow/first", false); err != nil {
		t.Fatal(err)
	}
	updated, err := os.ReadFile(path)
	if err != nil || !strings.Contains(string(updated), "// keep comment") || !strings.Contains(string(updated), `exclude = ["first"]`) {
		t.Fatalf("updated=%s err=%v", updated, err)
	}
	if _, err := assembly.Except(path, "team:workflow/first", false); err == nil {
		t.Fatal("duplicate exception accepted")
	}
	for _, reference := range []string{"team:workflow", "team:missing", "missing:workflow/first"} {
		if _, err := assembly.Except(path, reference, true); err == nil {
			t.Errorf("invalid exception accepted: %s", reference)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "AGENTS.md")); !os.IsNotExist(err) {
		t.Fatalf("exception installed output: %v", err)
	}
}
