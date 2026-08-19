package sourcepath

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveRejectsUnsafeAndSymlinkedSubdirs(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "libraries", "cdint"), 0o755); err != nil {
		t.Fatal(err)
	}
	resolved, err := Resolve(root, "libraries/cdint")
	if err != nil {
		t.Fatal(err)
	}
	if resolved != filepath.Join(root, "libraries", "cdint") {
		t.Fatalf("resolved = %q", resolved)
	}
	for _, value := range []string{"/absolute", "../escape", "libraries\\cdint", "libraries//cdint"} {
		if _, err := Resolve(root, value); err == nil {
			t.Fatalf("Resolve(%q) succeeded", value)
		}
	}
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(root, "linked")); err != nil {
		t.Fatal(err)
	}
	if _, err := Resolve(root, "linked"); err == nil || !strings.Contains(err.Error(), "symlinked") {
		t.Fatalf("symlink error = %v", err)
	}
	rootLink := filepath.Join(t.TempDir(), "root-link")
	if err := os.Symlink(root, rootLink); err != nil {
		t.Fatal(err)
	}
	if _, err := Resolve(rootLink, ""); err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("root symlink error = %v", err)
	}
}
