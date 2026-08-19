package library_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/library"
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

func TestLoadUsesDirectoryPathAsHeadingPrefix(t *testing.T) {
	temporary := t.TempDir()
	sourceDir := filepath.Join(temporary, "lang", "go")
	if err := os.MkdirAll(sourceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "testing.md"), []byte("# Testing\nUse Go tests.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	index, err := library.Load(temporary)
	if err != nil {
		t.Fatal(err)
	}
	node := index.ByPath["lang/go/testing"]
	if node == nil {
		t.Fatalf("missing directory-prefixed node: %#v", index.ByPath)
	}
	if node.Heading != "Testing" {
		t.Fatalf("heading = %q", node.Heading)
	}
	if _, found := index.ByPath["testing"]; found {
		t.Fatalf("unexpected unprefixed node path: %#v", index.ByPath)
	}
	if index.ByPath["lang"].Kind != library.NodeDirectory || index.ByPath["lang/go"].Kind != library.NodeDirectory {
		t.Fatalf("directory nodes = %#v, %#v", index.ByPath["lang"], index.ByPath["lang/go"])
	}
	if len(index.Roots) != 1 || index.Roots[0].Path != "lang" || len(index.Roots[0].Children) != 1 {
		t.Fatalf("directory tree = %#v", index.Roots)
	}
}

func TestLoadMergesCollidingDirectoryAndHeadingIdentity(t *testing.T) {
	temporary := t.TempDir()
	if err := os.Mkdir(filepath.Join(temporary, "lang"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(temporary, "lang.md"), []byte("# Lang\nLanguage policy.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(temporary, "lang", "go.md"), []byte("# Go\nGo policy.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	index, err := library.Load(temporary)
	if err != nil {
		t.Fatal(err)
	}
	lang := index.ByPath["lang"]
	if lang.Kind != library.NodeDirectoryHeading || lang.Body != "Language policy.\n" {
		t.Fatalf("lang = %#v", lang)
	}
	if len(lang.Children) != 1 || lang.Children[0].Path != "lang/go" {
		t.Fatalf("lang children = %#v", lang.Children)
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

func TestLoadRejectsSymlinks(t *testing.T) {
	tests := []struct {
		name       string
		createLink func(t *testing.T, temporary string) string
		want       []string
	}{
		{
			name: "source root",
			createLink: func(t *testing.T, temporary string) string {
				t.Helper()
				target := filepath.Join(temporary, "real-library")
				if err := os.Mkdir(target, 0o755); err != nil {
					t.Fatal(err)
				}
				writeTestFile(t, filepath.Join(target, "rules.md"), "# Rules\nBody.\n")
				root := filepath.Join(temporary, "linked-library")
				createTestSymlink(t, target, root)
				return root
			},
			want: []string{"is an unsupported symlink", "declare the target directory directly"},
		},
		{
			name: "Markdown file",
			createLink: func(t *testing.T, temporary string) string {
				t.Helper()
				root := filepath.Join(temporary, "library")
				if err := os.Mkdir(root, 0o755); err != nil {
					t.Fatal(err)
				}
				target := filepath.Join(temporary, "private.md")
				writeTestFile(t, target, "# Private\nDo not import.\n")
				createTestSymlink(t, target, filepath.Join(root, "linked.md"))
				return root
			},
			want: []string{"contains unsupported symlink", "linked.md", "declare the target directory as a separate source"},
		},
		{
			name: "directory",
			createLink: func(t *testing.T, temporary string) string {
				t.Helper()
				root := filepath.Join(temporary, "library")
				if err := os.Mkdir(root, 0o755); err != nil {
					t.Fatal(err)
				}
				target := filepath.Join(temporary, "shared-directory")
				if err := os.Mkdir(target, 0o755); err != nil {
					t.Fatal(err)
				}
				writeTestFile(t, filepath.Join(target, "rules.md"), "# Rules\nBody.\n")
				createTestSymlink(t, target, filepath.Join(root, "linked-directory"))
				return root
			},
			want: []string{"contains unsupported symlink", "linked-directory", "ordinary source content"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			temporary := t.TempDir()
			root := test.createLink(t, temporary)
			_, err := library.Load(root)
			if err == nil {
				t.Fatal("expected symlink error")
			}
			for _, want := range test.want {
				if !strings.Contains(err.Error(), want) {
					t.Fatalf("error %q does not contain %q", err, want)
				}
			}
		})
	}
}

func writeTestFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func createTestSymlink(t *testing.T, target, path string) {
	t.Helper()
	if err := os.Symlink(target, path); err != nil {
		t.Skipf("symlink creation is unavailable: %v", err)
	}
}
