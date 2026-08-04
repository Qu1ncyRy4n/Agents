package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/internal/manifest"
)

func TestSessionMarksAndDiscardsDraftChanges(t *testing.T) {
	_, manifestPath := writeSessionFixture(t)
	session, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	session.Draft.Doc = append(session.Draft.Doc, manifest.Entry{
		Heading: "Instructions",
		From:    []string{"shared:instructions"},
	})
	if err := session.MarkDraftChanged(); err != nil {
		t.Fatal(err)
	}
	if !session.Dirty {
		t.Fatal("expected dirty draft")
	}
	if !strings.Contains(session.Output, "# Instructions") {
		t.Fatalf("preview missing draft content:\n%s", session.Output)
	}
	if err := session.DiscardDraft(); err != nil {
		t.Fatal(err)
	}
	if session.Dirty {
		t.Fatal("expected discard to clear dirty state")
	}
	if strings.Contains(session.Output, "# Instructions") {
		t.Fatalf("discarded preview still contains draft content:\n%s", session.Output)
	}
}

func TestSessionSaveAndBuildWritesManifestAndOutput(t *testing.T) {
	temporary, manifestPath := writeSessionFixture(t)
	session, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	session.Draft.Doc = append(session.Draft.Doc, manifest.Entry{
		Heading: "Instructions",
		From:    []string{"shared:instructions"},
	})
	if err := session.MarkDraftChanged(); err != nil {
		t.Fatal(err)
	}
	if err := session.SaveAndBuild(); err != nil {
		t.Fatal(err)
	}
	if session.Dirty {
		t.Fatal("expected save to clear dirty state")
	}
	manifestContents, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(manifestContents), "Instructions: shared:instructions") {
		t.Fatalf("saved manifest missing draft entry:\n%s", manifestContents)
	}
	output, err := os.ReadFile(filepath.Join(temporary, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(output) != "# Identity\n\nHello.\n\n# Instructions\n\nWork carefully.\n" {
		t.Fatalf("output = %q", output)
	}
}

func writeSessionFixture(t *testing.T) (string, string) {
	t.Helper()
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	source := "# Identity\nHello.\n\n# Instructions\nWork carefully.\n"
	if err := os.WriteFile(filepath.Join(libraryPath, "core.md"), []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	config := "sources:\n  shared: library\noutput: AGENTS.md\ndoc:\n  - Identity: shared:identity\n"
	if err := os.WriteFile(manifestPath, []byte(config), 0o644); err != nil {
		t.Fatal(err)
	}
	return temporary, manifestPath
}
