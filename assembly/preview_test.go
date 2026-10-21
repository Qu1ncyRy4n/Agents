package assembly_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/assembly"
)

func TestPreviewSourceRendersSectionsAndPhysicalPayloadsWithoutWrites(t *testing.T) {
	root := t.TempDir()
	for name, body := range map[string]string{
		"mogent.lib":       "library {format=1;name=\"Team\"}\nmd-section \"guide\" {source=\"guide.md\"}\ndir \"assets\" {source=\"assets\";node \"readme.md\" {};node \"data.bin\" {}}\n",
		"guide.md":         "# Guide\nLine one.\nLine two.\n",
		"assets/readme.md": "# Readme\n\n## Child\nText\n",
		"assets/data.bin":  "\x00\x01",
		"mogent.build":     "build {format=1;origin=\"authored\"}\nsource \"team" + "\" {local=\".\"}\n",
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
	file, err := assembly.LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	guide, err := assembly.PreviewSource(file, path, "team:guide", 2, false)
	if err != nil || !strings.Contains(string(guide.Content), "Line one.") || !strings.Contains(string(guide.Content), "... (1 more lines)") {
		t.Fatalf("guide=%+v err=%v", guide, err)
	}
	directory, err := assembly.PreviewSource(file, path, "team:assets", 0, false)
	if err != nil || len(directory.Children) != 2 {
		t.Fatalf("directory=%+v err=%v", directory, err)
	}
	headings, err := assembly.PreviewSource(file, path, "team:assets/readme.md", 0, true)
	if err != nil || len(headings.Headings) != 2 || headings.Headings[1][1] != "Child" {
		t.Fatalf("headings=%+v err=%v", headings, err)
	}
	opaque, err := assembly.PreviewSource(file, path, "team:assets/data.bin", 0, false)
	if err != nil || !opaque.Opaque || opaque.Bytes != 2 {
		t.Fatalf("opaque=%+v err=%v", opaque, err)
	}
	sectionHeadings, err := assembly.PreviewSource(file, path, "team:guide", 0, true)
	if err != nil || len(sectionHeadings.Headings) != 1 || sectionHeadings.Headings[0][0] != "Guide" {
		t.Fatalf("section headings=%+v err=%v", sectionHeadings, err)
	}
	if _, err := os.Stat(filepath.Join(root, "AGENTS.md")); !os.IsNotExist(err) {
		t.Fatalf("preview installed output: %v", err)
	}
}
