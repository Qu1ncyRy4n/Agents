package library_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/internal/library"
)

func TestLoadUsesOptionalIDAsAPathSegmentAndOmitsMetadata(t *testing.T) {
	temporary := t.TempDir()
	if err := os.WriteFile(filepath.Join(temporary, "rules.md"), []byte("# Rules\n## Testing <!-- id: strict-testing -->\nUse strict tests.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	index, err := library.Load(temporary)
	if err != nil {
		t.Fatal(err)
	}
	node := index.ByPath["rules/strict-testing"]
	if node == nil {
		t.Fatalf("missing ID-addressed node: %#v", index.ByPath)
	}
	if node.Body != "Use strict tests.\n" {
		t.Fatalf("body retained metadata: %q", node.Body)
	}
	if node.File != filepath.Join(temporary, "rules.md") {
		t.Fatalf("file = %q", node.File)
	}
	if node.Line != 2 {
		t.Fatalf("line = %d", node.Line)
	}
}
