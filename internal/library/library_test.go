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

func TestLoadParsesFileFrontmatterMetadata(t *testing.T) {
	temporary := t.TempDir()
	source := `---
tags: [go, testing]
tldr: Prefer deterministic Go tests.
priority: 0.8
scope: project
requires: [shared:instructions/workflow]
conflicts_with: [shared:testing/fast-only]
---
# Go

## Testing
Use table tests.
`
	if err := os.WriteFile(filepath.Join(temporary, "go.md"), []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	index, err := library.Load(temporary)
	if err != nil {
		t.Fatal(err)
	}
	node := index.ByPath["go/testing"]
	if node == nil {
		t.Fatalf("missing node: %#v", index.ByPath)
	}
	if node.Line != 11 {
		t.Fatalf("line = %d", node.Line)
	}
	if node.Body != "Use table tests.\n" {
		t.Fatalf("body retained metadata: %q", node.Body)
	}
	if node.Metadata.TLDR != "Prefer deterministic Go tests." {
		t.Fatalf("metadata = %#v", node.Metadata)
	}
	if len(node.Metadata.Tags) != 2 || node.Metadata.Tags[0] != "go" || node.Metadata.Tags[1] != "testing" {
		t.Fatalf("tags = %#v", node.Metadata.Tags)
	}
	if node.Metadata.Priority == nil || *node.Metadata.Priority != 0.8 {
		t.Fatalf("priority = %#v", node.Metadata.Priority)
	}
	if node.Metadata.Scope != "project" {
		t.Fatalf("scope = %q", node.Metadata.Scope)
	}
	if len(node.Metadata.Requires) != 1 || node.Metadata.Requires[0] != "shared:instructions/workflow" {
		t.Fatalf("requires = %#v", node.Metadata.Requires)
	}
	if len(node.Metadata.ConflictsWith) != 1 || node.Metadata.ConflictsWith[0] != "shared:testing/fast-only" {
		t.Fatalf("conflicts_with = %#v", node.Metadata.ConflictsWith)
	}
}

func TestLoadRejectsInvalidFrontmatterMetadata(t *testing.T) {
	temporary := t.TempDir()
	source := "---\npriority: 1.5\n---\n# Go\nBody.\n"
	if err := os.WriteFile(filepath.Join(temporary, "go.md"), []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := library.Load(temporary); err == nil {
		t.Fatal("expected invalid priority error")
	}
}
